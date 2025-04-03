package main

import "time"

// RequestAPI
type requestAPI struct {
	Shipper struct {
		RegisteredNumber string `json:"registered_number"`
		Token            string `json:"token"`
		PlatformCode     string `json:"platform_code"`
	} `json:"shipper"`
	Recipient struct {
		Type    int    `json:"type"`
		Country string `json:"country"`
		Zipcode int    `json:"zipcode"`
	} `json:"recipient"`
	Dispatchers    []dispatcher `json:"dispatchers"`
	SimulationType []int        `json:"simulation_type"`
	Returns        struct {
		Composition  bool `json:"composition"`
		Volumes      bool `json:"volumes"`
		AppliedRules bool `json:"applied_rules"`
	} `json:"returns"`
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
	Offer          int    `json:"offer"`
	TableReference string `json:"table_reference"`
	SimulationType int    `json:"simulation_type"`
	Carrier        struct {
		Name             string `json:"name"`
		RegisteredNumber string `json:"registered_number"`
		StateInscription string `json:"state_inscription"`
		Logo             string `json:"logo"`
		Reference        int    `json:"reference"`
		CompanyName      string `json:"company_name"`
	} `json:"carrier"`
	Service      string `json:"service"`
	DeliveryTime struct {
		Days          int    `json:"days"`
		EstimatedDate string `json:"estimated_date"`
	} `json:"delivery_time,omitempty"`
	Expiration time.Time `json:"expiration"`
	CostPrice  float64   `json:"cost_price"`
	FinalPrice float64   `json:"final_price"`
	Weights    struct {
		Real  int     `json:"real"`
		Cubed float64 `json:"cubed"`
		Used  float64 `json:"used"`
	} `json:"weights"`
	OriginalDeliveryTime struct {
		Days          int    `json:"days"`
		EstimatedDate string `json:"estimated_date"`
	} `json:"original_delivery_time,omitempty"`
	HomeDelivery                bool `json:"home_delivery"`
	CarrierOriginalDeliveryTime struct {
		Days          int    `json:"days"`
		EstimatedDate string `json:"estimated_date"`
	} `json:"carrier_original_delivery_time,omitempty"`
	Modal string `json:"modal"`
}
