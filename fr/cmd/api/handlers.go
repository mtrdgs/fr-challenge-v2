package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
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

	// call simulate api
	responseAPI, err := app.postSimulateAPI(requestAPI)
	if err != nil {
		payload.Error = true
		payload.Message = "Failed to connect to API"
		payload.Data = err.Error()

		app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}
	//app.writeJSON(w, http.StatusOK, responseAPI)

	// format response from api, to be used in mongo
	quoteResult := app.formatResponseAPI(responseAPI)

	// save result in mongo
	err = app.Models.QuoteEntry.Insert(quoteResult)
	if err != nil {
		payload.Error = true
		payload.Message = "Failed to insert into Mongo"
		payload.Data = err.Error()

		app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	// done correctly!
	app.writeJSON(w, http.StatusOK, quoteResult)
}

func (app *Config) Metrics(w http.ResponseWriter, r *http.Request) {
	var lastQuotes int64
	var err error

	// check if quertystring is set
	queryString := r.URL.Query().Get("last_quotes")
	if !strings.EqualFold(queryString, "") {
		// it is! let's see if it has a valid entry
		lastQuotes, _ = strconv.ParseInt(queryString, 10, 64)
		if lastQuotes == 0 {
			app.errorJSON(w, errors.New("invalid 'last_quotes' value"), http.StatusBadRequest)
			return
		}
	} else {
		app.errorJSON(w, errors.New("parameter 'last_quotes' is required"), http.StatusBadRequest)
		return
	}

	// retrieve quotes from db
	quotes, err := app.Models.QuoteEntry.FindSpecific(lastQuotes)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusOK, quotes)
}
