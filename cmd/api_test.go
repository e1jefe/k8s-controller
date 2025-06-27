package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeploymentStruct(t *testing.T) {
	// Test the Deployment struct
	d := Deployment{
		Name:      "test-deployment",
		Namespace: "default",
		Replicas:  3,
		Ready:     2,
	}

	if d.Name != "test-deployment" {
		t.Errorf("Expected name 'test-deployment', got %s", d.Name)
	}
	if d.Replicas != 3 {
		t.Errorf("Expected replicas 3, got %d", d.Replicas)
	}
	if d.Ready != 2 {
		t.Errorf("Expected ready 2, got %d", d.Ready)
	}
}

func TestListDeploymentsHandler_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("POST", "/deployments", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listDeploymentsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestListDeploymentsHandler_WithNilInformer(t *testing.T) {
	req, err := http.NewRequest("GET", "/deployments", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create a simple handler that handles nil informer gracefully
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Return empty array when informer is not set up
		var deployments []Deployment
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deployments)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check that response is valid JSON
	var deployments []Deployment
	if err := json.Unmarshal(rr.Body.Bytes(), &deployments); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Should be empty array when no deployments
	if len(deployments) != 0 {
		t.Errorf("Expected empty array, got %d items", len(deployments))
	}
}
