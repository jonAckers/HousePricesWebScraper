package main

import (
	"fmt"
	"strings"
)

func buildEmailBody(properties []Property) string {
	var builder strings.Builder

	// Add the header
	builder.WriteString(fmt.Sprintf("There are %d new properties on Rightmove!\n\n", len(properties)))

	// Add each property's details
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
