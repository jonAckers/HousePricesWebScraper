package main

import (
	"fmt"
	"strings"
)

func buildEmailBody(properties propertyLists) string {
	var builder strings.Builder

	// Add the header
	builder.WriteString(fmt.Sprintf("There are %d new properties and %d updated properties on Rightmove!\n\n",
						len(properties.new), len(properties.updated)))

	// Add new property details
	if len(properties.new) > 0 {
		builder.WriteString("\nNew Properties\n")
		builder.WriteString(buildBodyForList(properties.new))
	}

	// Add updated property details
	if len(properties.updated) > 0 {
		builder.WriteString("\nUpdated Properties\n")
		builder.WriteString(buildBodyForList(properties.updated))
	}

	return builder.String()
}

func buildBodyForList(properties []Property) string {
	var builder strings.Builder

	for _, property := range properties {
		builder.WriteString("----------\n")
		builder.WriteString(fmt.Sprintf("%s\n", property.Address))
		builder.WriteString(fmt.Sprintf("https://www.rightmove.co.uk/properties/%d\n", property.Id))
		builder.WriteString(fmt.Sprintf("- Price: %s%d\n", getCurrencySymbol(property.Price.CurrencyCode), property.Price.Amount))
		builder.WriteString(fmt.Sprintf("- Bedrooms: %d\n", property.Bedrooms))
		builder.WriteString(fmt.Sprintf("- Bathrooms: %d\n", property.Bathrooms))
		builder.WriteString(fmt.Sprintf("- Type: %s\n", property.Type))
		builder.WriteString(fmt.Sprintf("- Estate Agent: %s - %s\n\n", property.EstateAgent.Name, property.EstateAgent.Telephone))
		builder.WriteString(fmt.Sprintf("%s\n", property.Description))
		builder.WriteString("----------\n\n")
	}

	return builder.String()
}

func getCurrencySymbol(code string) string {
	if code == "GBP" {
		return "£"
	}
	if code == "EUR" {
		return "€"
	}
	if code == "USD" {
		return "$"
	}
	return ""
}
