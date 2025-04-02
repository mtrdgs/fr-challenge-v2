package main

import (
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
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
	payload := jsonResponse{}

	// decode request
	err := app.readJSON(w, r, &requestQuote)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// verify how many invalid arguments request has
	// example of invalid argument: missing zipcode
	invalidArgs := app.checkRequest(requestQuote)
	if len(invalidArgs) > 0 {
		payload.Error = true
		payload.Message = "Missing arguments!"
		payload.Data = invalidArgs

		app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	// build request (needed for external api)
	requestAPI := app.buildRequestAPI(requestQuote)
	//app.writeJSON(w, http.StatusOK, requestAPI)

	// simulate request
	responseAPI, err := app.postSimulateAPI(requestAPI)
	if err != nil {
		payload.Error = true
		payload.Message = "Failed to connect to API"
		payload.Data = err.Error()

		app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}
	app.writeJSON(w, http.StatusOK, responseAPI)

}
