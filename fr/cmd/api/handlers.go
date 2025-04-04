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

// Fr - a test page to see if there's connectivity/response
func (app *Config) Fr(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "connected!",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// Quote -
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

	// call freterapido api (simulate module)
	responseAPI, err := app.postSimulateAPI(requestAPI)
	if err != nil {
		payload.Error = true
		payload.Message = "Failed to connect to freterapido API"
		payload.Data = err.Error()

		app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	// format response from api, to be used in mongo
	quoteResult := app.formatResponseAPI(responseAPI)

	// save result in mongo
	err = app.Repo.Insert(quoteResult)
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

// Metrics -
func (app *Config) Metrics(w http.ResponseWriter, r *http.Request) {
	var lastQuotes int64
	var err error

	// check if quertystring is set
	queryString := r.URL.Query().Get("last_quotes")
	if !strings.EqualFold(queryString, "") {
		// it is! let's see if it has a valid entry (any int number)
		lastQuotes, err = strconv.ParseInt(queryString, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("invalid 'last_quotes' value"), http.StatusBadRequest)
			return
		}
	}

	// retrieve quotes from db
	quotes, err := app.Repo.FindSpecific(lastQuotes)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// format quotes
	responseMetrics := app.prepareMetricsResponse(quotes)

	// done correctly!
	app.writeJSON(w, http.StatusOK, responseMetrics)
}
