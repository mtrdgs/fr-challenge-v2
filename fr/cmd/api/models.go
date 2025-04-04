package main

import "time"

type requestQuote struct {
	Recipient recipientQuote `json:"recipient"`
	Volumes   []volume       `json:"volumes"`
}

type recipientQuote struct {
	Address address `json:"address"`
}

type address struct {
	Zipcode string `json:"zipcode"`
}

type volume struct {
	Category      string  `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight int     `json:"unitary_weight"`
	Price         int     `json:"price"`
	UnitaryPrice  int     `json:"unitary_price"`
	Sku           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}

// RequestAPI
type requestAPI struct {
	Shipper        shipper      `json:"shipper"`
	Recipient      recipientApi `json:"recipient"`
	Dispatchers    []dispatcher `json:"dispatchers"`
	SimulationType []int        `json:"simulation_type"`
	Returns        returns      `json:"returns"`
}

// Shipper
type shipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

// RecipientApi
type recipientApi struct {
	Type    int    `json:"type"`
	Country string `json:"country"`
	Zipcode int    `json:"zipcode"`
}

// Returns
type returns struct {
	Composition  bool `json:"composition"`
	Volumes      bool `json:"volumes"`
	AppliedRules bool `json:"applied_rules"`
}

// Dispatcher -
type dispatcher struct {
	RegisteredNumber string   `json:"registered_number"`
	Zipcode          int      `json:"zipcode"`
	Volumes          []volume `json:"volumes"`
}

// ResponseAPI -
type responseAPI struct {
	Dispatchers []dispatcherAPI `json:"dispatchers"`
}

// DispatcherAPI -
type dispatcherAPI struct {
	ID                         string  `json:"id"`
	RequestID                  string  `json:"request_id"`
	RegisteredNumberShipper    string  `json:"registered_number_shipper"`
	RegisteredNumberDispatcher string  `json:"registered_number_dispatcher"`
	ZipcodeOrigin              int     `json:"zipcode_origin"`
	Offers                     []offer `json:"offers"`
}

// Offer -
type offer struct {
	Offer                       int                         `json:"offer"`
	TableReference              string                      `json:"table_reference"`
	SimulationType              int                         `json:"simulation_type"`
	Carrier                     carrier                     `json:"carrier"`
	Service                     string                      `json:"service"`
	DeliveryTime                deliveryTime                `json:"delivery_time,omitempty"`
	Expiration                  time.Time                   `json:"expiration"`
	CostPrice                   float64                     `json:"cost_price"`
	FinalPrice                  float64                     `json:"final_price"`
	Weights                     weights                     `json:"weights"`
	OriginalDeliveryTime        originalDeliveryTime        `json:"original_delivery_time,omitempty"`
	HomeDelivery                bool                        `json:"home_delivery"`
	CarrierOriginalDeliveryTime carrierOriginalDeliveryTime `json:"carrier_original_delivery_time,omitempty"`
	Modal                       string                      `json:"modal"`
}

// Carrier -
type carrier struct {
	Name             string `json:"name"`
	RegisteredNumber string `json:"registered_number"`
	StateInscription string `json:"state_inscription"`
	Logo             string `json:"logo"`
	Reference        int    `json:"reference"`
	CompanyName      string `json:"company_name"`
}

// DeliveryTime -
type deliveryTime struct {
	Days          int    `json:"days"`
	EstimatedDate string `json:"estimated_date"`
}

// Weights -
type weights struct {
	Real  int     `json:"real"`
	Cubed float64 `json:"cubed"`
	Used  float64 `json:"used"`
}

// OroriginalDeliveryTime -
type originalDeliveryTime struct {
	Days          int    `json:"days"`
	EstimatedDate string `json:"estimated_date"`
}

// CarrierOriginalDeliveryTime -
type carrierOriginalDeliveryTime struct {
	Days          int    `json:"days"`
	EstimatedDate string `json:"estimated_date"`
}
