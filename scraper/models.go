package main

type Property struct {
	Id             int    `json:"id"`
	Bedrooms       int    `json:"bedrooms"`
	Bathrooms      int    `json:"bathrooms"`
	Description    string `json:"summary"`
	Address        string `json:"displayAddress"`
	Location       struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Type            string `json:"propertySubType"`
	ListingUpdate   struct {
		Reason string `json:"listingUpdateReason"`
		Date   string `json:"listingUpdateDate"`
	} `json:"listingUpdate"`
	Price struct {
		Amount       int    `json:"amount"`
		CurrencyCode string `json:"currencyCode"`
	} `json:"price"`
	EstateAgent struct {
		Telephone  string `json:"contactTelephone"`
		Name       string `json:"branchDisplayName"`
	} `json:"customer"`
}
