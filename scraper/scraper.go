package main

import (
	"encoding/json"
	"log/slog"
	"regexp"

	"github.com/jonackers/HousePricesWebScraper/scraper/internal/database"

	"github.com/gocolly/colly"
)

const RIGHTMOVE_URL = "https://www.rightmove.co.uk/property-for-sale/find.html?locationIdentifier=OUTCODE%5E1369&maxPrice=450000&radius=10.0&sortType=6&propertyTypes=&includeSSTC=false&mustHave=&dontShow=&furnishTypes=&keywords="

func parseHouseDetails() ([]Property, error) {
	var properties []Property

	// Initialise the collector
	c := colly.NewCollector()

	// Regex to extract JSON data from the script tag
	jsonRegex := regexp.MustCompile(`window\.jsonModel = \{"properties":(\[.*\]),"resultCount".*`)

	// OnHTML callback to target the <script> tag
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		// Extract JSON using regex
		jsonData := jsonRegex.FindStringSubmatch(scriptContent)

		if len(jsonData) > 1 {
			// Parse JSON
			err := json.Unmarshal([]byte(jsonData[1]), &properties)
			if err != nil {
				slog.Error("Error parsing JSON:", err)
				return
			}

			// Loop through the properties and print the required data
			slog.Info("Scrape complete!", "property_count", len(properties))
		}
	})

	// Visit website
	err := c.Visit(RIGHTMOVE_URL)
	if err != nil {
		slog.Error("Failed to visit rightmove", "err", err)
	}

	return properties, err
}

func (cfg *dbConfig) getNewProperties(properties []Property) (propertyLists, error) {
	var newProperties []Property
	var updatedProperties []Property

	// Iterate through the properties to see if there are any new ones
	for _, property := range properties {
		dbProperty, err := cfg.db.GetPropertyById(cfg.ctx, int32(property.Id))

		if err == nil {
			// Property found in the database, so must not be new
			// Check if property has been updated
			if (property.ListingUpdate.Date.After(dbProperty.ListingUpdateDate)) {
				slog.Info("Existing property updated!", "property", property)
				err := cfg.updateDbProperty(property)
				if err != nil {
					return propertyLists{}, err
				}
				updatedProperties = append(updatedProperties, property)
			}
		} else {
			// Property not found, so must be new
			slog.Info("New property found!", "property", property)
			// Add new property to the database and add it to list to return
			err := cfg.addNewPropertyToDb(property)
			if err != nil {
				return propertyLists{}, err
			}
			newProperties = append(newProperties, property)
		}
	}

	return propertyLists{
		new: newProperties,
		updated: updatedProperties,
	}, nil
}

func (cfg *dbConfig) addNewPropertyToDb(property Property) error {
	_, err := cfg.db.CreateProperty(cfg.ctx, database.CreatePropertyParams{
		ID: int32(property.Id),
		Bedrooms: int32(property.Bedrooms),
		Bathrooms: int32(property.Bathrooms),
		Description: property.Description,
		Address: property.Address,
		Latitude: property.Location.Latitude,
		Longitude: property.Location.Longitude,
		Type: property.Type,
		ListingUpdateReason: property.ListingUpdate.Reason,
		ListingUpdateDate: property.ListingUpdate.Date,
		PriceAmount: int32(property.Price.Amount),
		PriceCurrencyCode: property.Price.CurrencyCode,
		EstateAgentTelephone: property.EstateAgent.Telephone,
		EstateAgentName: property.EstateAgent.Name,
	})

	return err
}

func (cfg *dbConfig) updateDbProperty(property Property) error {
	err := cfg.db.UpdatePropertyById(cfg.ctx, database.UpdatePropertyByIdParams{
		ID: int32(property.Id),
		Bedrooms: int32(property.Bedrooms),
		Bathrooms: int32(property.Bathrooms),
		Description: property.Description,
		Address: property.Address,
		Latitude: property.Location.Latitude,
		Longitude: property.Location.Longitude,
		Type: property.Type,
		ListingUpdateReason: property.ListingUpdate.Reason,
		ListingUpdateDate: property.ListingUpdate.Date,
		PriceAmount: int32(property.Price.Amount),
		PriceCurrencyCode: property.Price.CurrencyCode,
		EstateAgentTelephone: property.EstateAgent.Telephone,
		EstateAgentName: property.EstateAgent.Name,
	})

	return err
}
