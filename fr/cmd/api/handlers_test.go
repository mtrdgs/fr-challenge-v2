package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mtrdgs/fr/data"
)

var testApp Config

func TestConfig_Quote(t *testing.T) {
	// call mocked repository
	repo := data.NewMongoTestRepository(nil)
	testApp.Repo = repo

	postBody := map[string]interface{}{
		"message": "test message",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest(http.MethodPost, "/quote", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Quote)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected http.StatusBadRequest but got %d", rr.Code)
	}
}

func TestConfig_Metrics(t *testing.T) {
	// call mocked repository
	repo := data.NewMongoTestRepository(nil)
	testApp.Repo = repo

	req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Metrics)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected http.StatusOK but got %d", rr.Code)
	}
}
