package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly"
)

const RIGHTMOVE_URL = "https://www.rightmove.co.uk/property-for-sale/find.html?locationIdentifier=OUTCODE%5E1369&maxPrice=450000&radius=10.0&sortType=6&propertyTypes=&includeSSTC=false&mustHave=&dontShow=&furnishTypes=&keywords="

func parseHouseDetails() []Property {
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
				fmt.Println("Error parsing JSON:", err)
				return
			}

			// Loop through the properties and print the required data
			log.Printf("Scrape complete! Found %d properties.\n", len(properties))
		}
	})

	// Visit website
	err := c.Visit(RIGHTMOVE_URL)
	if err != nil {
		log.Fatal(err)
	}

	return properties
}
