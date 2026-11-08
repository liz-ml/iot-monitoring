package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	router := SetupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"Esperado status 200, recebido %d",
			recorder.Code,
		)
	}

	var response struct {
		Status  string `json:"status"`
		Service string `json:"service"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Resposta JSON inválida: %v", err)
	}

	if response.Status != "running" {
		t.Errorf("Status esperado 'running', recebido %q", response.Status)
	}

	if response.Service != "user-service" {
		t.Errorf("Serviço inesperado: %q", response.Service)
	}
}
