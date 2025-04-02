package main

import (
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type requestQuote struct {
	Recipient struct {
		Address struct {
			Zipcode string `json:"zipcode"`
		} `json:"address"`
	} `json:"recipient"`
	Volumes []volume `json:"volumes"`
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

func (app *Config) Fr(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "connected!",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) Quote(w http.ResponseWriter, r *http.Request) {
	requestQuote := requestQuote{}

	// decode request
	err := app.readJSON(w, r, &requestQuote)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// check if its returning
	app.writeJSON(w, http.StatusOK, requestQuote)
}
