package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

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

func (app *Config) checkRequest(req requestQuote) (args []string) {
	args = make([]string, 0)

	// contains zipcode?
	if strings.EqualFold(req.Recipient.Address.Zipcode, "") {
		args = append(args, "Zipcode is required")
	}

	// contains volume?
	if len(req.Volumes) == 0 {
		args = append(args, "Volumes is required")
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

func (app *Config) buildRequestAPI(reqQuote requestQuote) (reqAPI requestAPI) {
	// shipper
	reqAPI.Shipper.RegisteredNumber = "25438296000158"        // hardcoded for now
	reqAPI.Shipper.Token = "1d52a9b6b78cf07b08586152459a5c90" // hardcoded for now
	reqAPI.Shipper.PlatformCode = "5AKVkHqCn"                 // hardcoded for now

	// recipient
	reqAPI.Recipient.Type = 0        // fixed
	reqAPI.Recipient.Country = "BRA" // fixed
	reqAPI.Recipient.Zipcode, _ = strconv.Atoi(reqQuote.Recipient.Address.Zipcode)

	// dispatchers
	var dispatcher dispatcher
	dispatcher.RegisteredNumber = "25438296000158" // hardcoded for now
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
