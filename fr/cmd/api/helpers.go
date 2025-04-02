package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
