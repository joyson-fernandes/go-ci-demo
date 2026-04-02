package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("health check returned %d, want 200", w.Code)
	}
}

func TestWeatherHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/weather/london", nil)
	w := httptest.NewRecorder()
	weatherHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("weather returned %d, want 200", w.Code)
	}
	var resp WeatherResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Location != "London, UK" {
		t.Errorf("location = %q, want London, UK", resp.Location)
	}
}

func TestWeatherNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/weather/atlantis", nil)
	w := httptest.NewRecorder()
	weatherHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("unknown city returned %d, want 404", w.Code)
	}
}

func TestCitiesHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/cities", nil)
	w := httptest.NewRecorder()
	citiesHandler(w, req)
	var cities []string
	json.NewDecoder(w.Body).Decode(&cities)
	if len(cities) != 8 {
		t.Errorf("got %d cities, want 8", len(cities))
	}
}
