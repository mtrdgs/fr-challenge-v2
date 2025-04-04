package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mtrdgs/fr/data"
)

// ResponseMetrics -
type responseMetrics struct {
	Metrics []metric `json:"metrics"`
}

// Metric -
type metric struct {
	ResultsPerCarrier    map[string]int     `json:"results_per_carrier"`
	TotalPricePerCarrier map[string]float64 `json:"total_price_per_carrier"`
	AvgPricePerCarrier   map[string]float64 `json:"avg_price_per_carrier"`
	CheapestFreight      map[string]float64 `json:"cheapest_freight"`
	PriciestFreight      map[string]float64 `json:"priciest_freight"`
}

// writeJSON - writes a json response to client (from a response's status code and data)
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// readJSON - reads a request's body and converts it into json
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1Mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

// errorJSON - sends a json error response to client
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

// checkRequest - verifies if user's request has all needed arguments
func (app *Config) checkRequest(req requestQuote) (args []string) {
	args = make([]string, 0)

	// contains zipcode?
	if strings.EqualFold(req.Recipient.Address.Zipcode, "") {
		args = append(args, "Zipcode is required")
	}

	// contains volume?
	if len(req.Volumes) == 0 {
		args = append(args, "Volumes are required")
	}

	// contains specific variables?
	for key, value := range req.Volumes {
		// category
		if strings.EqualFold(value.Category, "") {
			args = append(args, fmt.Sprintf("Category is required for Volume[%d]", key))
		}

		// amount
		if value.Amount == 0 {
			args = append(args, fmt.Sprintf("Amount is required for Volume[%d]", key))
		}

		// price
		if value.Price == 0 {
			args = append(args, fmt.Sprintf("Price is required for Volume[%d]", key))
		}

		// sku
		if strings.EqualFold(value.Sku, "") {
			args = append(args, fmt.Sprintf("SKU is required for Volume[%d]", key))
		}

		// height
		if value.Height == 0 {
			args = append(args, fmt.Sprintf("Height is required for Volume[%d]", key))
		}

		// width
		if value.Width == 0 {
			args = append(args, fmt.Sprintf("Width is required for Volume[%d]", key))
		}

		// length
		if value.Length == 0 {
			args = append(args, fmt.Sprintf("Length is required for Volume[%d]", key))
		}
	}

	return args
}

// buildRequestAPI - creates request to be used at freterapido api from user's input
func (app *Config) buildRequestAPI(reqQuote requestQuote) (reqAPI requestAPI) {
	// shipper
	reqAPI.Shipper.RegisteredNumber = os.Getenv("REGISTERED_NUMBER")
	reqAPI.Shipper.Token = os.Getenv("TOKEN")
	reqAPI.Shipper.PlatformCode = os.Getenv("PLATFORM_CODE")

	// recipient
	reqAPI.Recipient.Type = 0        // fixed
	reqAPI.Recipient.Country = "BRA" // fixed
	reqAPI.Recipient.Zipcode, _ = strconv.Atoi(reqQuote.Recipient.Address.Zipcode)

	// dispatchers
	var dispatcher dispatcher
	dispatcher.RegisteredNumber = os.Getenv("REGISTERED_NUMBER")
	dispatcher.Zipcode = reqAPI.Recipient.Zipcode
	for _, volume := range reqQuote.Volumes {
		volume.UnitaryPrice = volume.Price / volume.Amount
		dispatcher.Volumes = append(dispatcher.Volumes, volume)
	}
	reqAPI.Dispatchers = append(reqAPI.Dispatchers, dispatcher)

	// simulation type
	reqAPI.SimulationType = append(reqAPI.SimulationType, 0) // fixed

	// returns
	reqAPI.Returns.Composition = false
	reqAPI.Returns.Volumes = false
	reqAPI.Returns.AppliedRules = false

	return reqAPI
}

// postSimulateAPI - calls freterapido api
func (app *Config) postSimulateAPI(reqAPI requestAPI) (resAPI responseAPI, err error) {
	// build request
	payload, err := json.Marshal(reqAPI)
	if err != nil {
		return resAPI, err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(payload))
	if err != nil {
		return resAPI, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// send request
	res, err := app.Client.Do(req)
	if err != nil {
		return resAPI, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return resAPI, errors.New("appi returned a non-200 status code")
	}

	// decode response
	err = json.NewDecoder(res.Body).Decode(&resAPI)
	if err != nil {
		return resAPI, fmt.Errorf("failed to decode response: %w", err)
	}

	return resAPI, nil
}

// formatResponseAPI - returns a json to be used in mongo insert operation, using info from freterapido api response
func (app *Config) formatResponseAPI(entry responseAPI) (result data.QuoteEntry) {
	// has dispatchers?
	if len(entry.Dispatchers) == 0 {
		return result
	}

	// format response from api
	for _, value := range entry.Dispatchers[0].Offers {
		result.Carrier = append(result.Carrier, data.Carrier{
			Name:     value.Carrier.Name,
			Service:  value.Modal,
			Deadline: value.CarrierOriginalDeliveryTime.Days,
			Price:    value.FinalPrice,
		})
	}

	return result
}

// prepareMetricsResponse - calcs the info from quotes (stored in the db) and formats it as readable json to send to client"
func (app *Config) prepareMetricsResponse(quotes []data.QuoteEntry) responseMetrics {
	// create map to store metrics
	mapMetrics := metric{
		ResultsPerCarrier:    make(map[string]int),
		TotalPricePerCarrier: make(map[string]float64),
		AvgPricePerCarrier:   make(map[string]float64),
		CheapestFreight:      make(map[string]float64),
		PriciestFreight:      make(map[string]float64),
	}

	// iterate over quotes
	for _, quote := range quotes {
		for _, carrier := range quote.Carrier {
			// total results
			mapMetrics.ResultsPerCarrier[carrier.Name]++

			// total price
			mapMetrics.TotalPricePerCarrier[carrier.Name] = (mapMetrics.TotalPricePerCarrier[carrier.Name] + carrier.Price)
			mapMetrics.TotalPricePerCarrier[carrier.Name] = float64(int(mapMetrics.TotalPricePerCarrier[carrier.Name]*100)) / 100

			// avg price
			mapMetrics.AvgPricePerCarrier[carrier.Name] = mapMetrics.TotalPricePerCarrier[carrier.Name] / float64(mapMetrics.ResultsPerCarrier[carrier.Name])
			mapMetrics.AvgPricePerCarrier[carrier.Name] = float64(int(mapMetrics.AvgPricePerCarrier[carrier.Name]*100)) / 100

			// cheapest freight
			if mapMetrics.CheapestFreight[carrier.Name] == 0 || carrier.Price < mapMetrics.CheapestFreight[carrier.Name] {
				mapMetrics.CheapestFreight[carrier.Name] = carrier.Price
			}

			// priciest freight
			if mapMetrics.PriciestFreight[carrier.Name] == 0 || carrier.Price > mapMetrics.PriciestFreight[carrier.Name] {
				mapMetrics.PriciestFreight[carrier.Name] = carrier.Price
			}
		}
	}

	// append metrics to response
	return responseMetrics{Metrics: []metric{mapMetrics}}
}
