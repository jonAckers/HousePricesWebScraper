package main

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsses"
	"github.com/aws/aws-cdk-go/awscdk/v2/triggers"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

const (
	DATABASE_NAME = "HousePriceDatabase"
)

func NewHousePricesStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create a VPC
	vpc := awsec2.NewVpc(stack, jsii.String("HousePriceVpc"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				SubnetType: awsec2.SubnetType_PUBLIC,
				Name:       jsii.String("Public"),
			},
			{
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
				Name:       jsii.String("PrivateWithEgress"),
			},
			{
				SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
				Name:       jsii.String("Private"),
			},
		},
	})

	// Create a Security Group that allows inbound traffic on port 5432 for PostgreSQL
	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("HousePriceSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
	})
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(5432)), jsii.String("Allow PostgreSQL access"), jsii.Bool(false))

	// Create an RDS instance
	rdsInstance := awsrds.NewDatabaseInstance(stack, jsii.String("HousePriceRDS"), &awsrds.DatabaseInstanceProps{
		Engine: awsrds.DatabaseInstanceEngine_Postgres(&awsrds.PostgresInstanceEngineProps{
			Version: awsrds.PostgresEngineVersion_VER_13(),
		}),
		Vpc:            vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
		},
		SecurityGroups: &[]awsec2.ISecurityGroup{securityGroup},
		InstanceType:   awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
		Credentials:    awsrds.Credentials_FromGeneratedSecret(jsii.String("dbadmin"), &awsrds.CredentialsBaseOptions{}),
		DatabaseName:   jsii.String(DATABASE_NAME),
		MultiAz:       jsii.Bool(false),
		PubliclyAccessible: jsii.Bool(false),
	})

	// Create a trigger to update the database in case of any new migrations
	deployDbTrigger := triggers.NewTriggerFunction(stack, jsii.String("DeployDatabaseTrigger"), &triggers.TriggerFunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String(filepath.Join("..", "./database/")), &awss3assets.AssetOptions{}),
		Vpc: vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
		},
		Environment: &map[string]*string{
			"SECRET_NAME": rdsInstance.Secret().SecretArn(),
		},
	})
	deployDbTrigger.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("secretsmanager:GetSecretValue"),
			Resources: &[]*string{rdsInstance.Secret().SecretArn()},
		}),
	)
	deployDbTrigger.Node().AddDependency(rdsInstance)

	// Create Lambda function and grant read access to the RDS secret
	scraperLambda := awslambda.NewFunction(stack, jsii.String("HousePriceScraper"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String(filepath.Join("..", "./scraper/")), &awss3assets.AssetOptions{}),
		Vpc: vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
		Environment: &map[string]*string{
			"SECRET_NAME": rdsInstance.Secret().SecretArn(),
		},
		SecurityGroups: &[]awsec2.ISecurityGroup{securityGroup},
		Timeout: awscdk.Duration_Seconds(jsii.Number(30)),
	})
	scraperLambda.AddToRolePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("secretsmanager:GetSecretValue"),
			Resources: &[]*string{rdsInstance.Secret().SecretArn()},
		}),
	)

	// Create VPC endpoints for
	vpc.AddInterfaceEndpoint(jsii.String("SecretsManagerVpcEndpoint"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_SECRETS_MANAGER(),
	})
	vpc.AddInterfaceEndpoint(jsii.String("RdsVpcEndpoint"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_RDS(),
	})

	// Create an EventBridge rule to schedule Lambda execution every hour
	rule := awsevents.NewRule(stack, jsii.String("ScrapeHousePricesEveryHour"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Rate(awscdk.Duration_Hours(jsii.Number(1))),
	})
	rule.AddTarget(awseventstargets.NewLambdaFunction(scraperLambda, &awseventstargets.LambdaFunctionProps{}))

	// Configure Simple Email Service
	// Verify outbound email address
	awsses.NewEmailIdentity(stack, jsii.String("VerifiedEmailIdentity"), &awsses.EmailIdentityProps{
		Identity: awsses.Identity_Email(jsii.String("jonathon_ackers@hotmail.co.uk")),
	})
	// Grant SES permissions to Lambda to send emails
	scraperLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("ses:SendEmail", "ses:SendRawEmail"),
		Resources: jsii.Strings("*"),
	}))

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewHousePricesStack(app, "HousePricesScraper", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) to deploy the stack in
func env() *awscdk.Environment {
	return &awscdk.Environment{
	 Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
