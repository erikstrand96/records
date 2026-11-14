package models

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewProblemDetails(t *testing.T) {
	p := NewProblemDetails(http.StatusNotFound, "The requested record was not found")

	if p.Status != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, p.Status)
	}
	if p.Title != "Not Found" {
		t.Errorf("expected title 'Not Found', got '%s'", p.Title)
	}
	if p.Type != "about:blank" {
		t.Errorf("expected type 'about:blank', got '%s'", p.Type)
	}
	if p.Detail != "The requested record was not found" {
		t.Errorf("expected detail 'The requested record was not found', got '%s'", p.Detail)
	}
}

func TestNewProblemDetailsWithTitle(t *testing.T) {
	p := NewProblemDetailsWithTitle(http.StatusBadRequest, "Custom Bad Request", "Invalid input syntax")

	if p.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, p.Status)
	}
	if p.Title != "Custom Bad Request" {
		t.Errorf("expected title 'Custom Bad Request', got '%s'", p.Title)
	}
	if p.Detail != "Invalid input syntax" {
		t.Errorf("expected detail 'Invalid input syntax', got '%s'", p.Detail)
	}

	// Empty title should fallback to status text
	p2 := NewProblemDetailsWithTitle(http.StatusInternalServerError, "", "DB error")
	if p2.Title != "Internal Server Error" {
		t.Errorf("expected title 'Internal Server Error', got '%s'", p2.Title)
	}
}

func TestProblemDetailsBuilders(t *testing.T) {
	p := NewProblemDetails(http.StatusBadRequest, "Invalid request body").
		WithType("https://example.com/probs/invalid-body").
		WithTitle("Validation Failed").
		WithStatus(http.StatusUnprocessableEntity).
		WithDetail("Field 'title' is required").
		WithInstance("/records/123").
		WithExtension("field", "title").
		WithExtensions(map[string]any{
			"errorCode": "REQ_001",
		})

	if p.Type != "https://example.com/probs/invalid-body" {
		t.Errorf("unexpected type: %s", p.Type)
	}
	if p.Title != "Validation Failed" {
		t.Errorf("unexpected title: %s", p.Title)
	}
	if p.Status != http.StatusUnprocessableEntity {
		t.Errorf("unexpected status: %d", p.Status)
	}
	if p.Detail != "Field 'title' is required" {
		t.Errorf("unexpected detail: %s", p.Detail)
	}
	if p.Instance != "/records/123" {
		t.Errorf("unexpected instance: %s", p.Instance)
	}
	if p.Extensions["field"] != "title" {
		t.Errorf("expected extension field to be 'title', got %v", p.Extensions["field"])
	}
	if p.Extensions["errorCode"] != "REQ_001" {
		t.Errorf("expected extension errorCode to be 'REQ_001', got %v", p.Extensions["errorCode"])
	}
}

func TestProblemDetailsError(t *testing.T) {
	p1 := &ProblemDetails{Status: 404, Title: "Not Found", Detail: "Item missing"}
	if p1.Error() != "404 Not Found: Item missing" {
		t.Errorf("unexpected error string: %s", p1.Error())
	}

	p2 := &ProblemDetails{Detail: "Something went wrong"}
	if p2.Error() != "Something went wrong" {
		t.Errorf("unexpected error string: %s", p2.Error())
	}

	p3 := &ProblemDetails{Status: 500, Title: "Internal Server Error"}
	if p3.Error() != "500 Internal Server Error" {
		t.Errorf("unexpected error string: %s", p3.Error())
	}

	p4 := &ProblemDetails{Status: 502}
	if p4.Error() != "HTTP 502" {
		t.Errorf("unexpected error string: %s", p4.Error())
	}

	p5 := &ProblemDetails{}
	if p5.Error() != "problem details error" {
		t.Errorf("unexpected error string: %s", p5.Error())
	}

	var pNil *ProblemDetails
	if pNil.Error() != "problem details error" {
		t.Errorf("unexpected error string for nil problem details: %s", pNil.Error())
	}
}

func TestProblemDetailsMarshalJSON(t *testing.T) {
	p := &ProblemDetails{
		Type:     "https://example.com/probs/out-of-credit",
		Title:    "You do not have enough credit.",
		Status:   403,
		Detail:   "Your current balance is 30, but that costs 50.",
		Instance: "/account/12345/msgs/abc",
		Extensions: map[string]any{
			"balance": 30,
		},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("failed to marshal problem details: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if raw["type"] != "https://example.com/probs/out-of-credit" {
		t.Errorf("unexpected type in JSON: %v", raw["type"])
	}
	if raw["title"] != "You do not have enough credit." {
		t.Errorf("unexpected title in JSON: %v", raw["title"])
	}
	if raw["status"] != float64(403) {
		t.Errorf("unexpected status in JSON: %v", raw["status"])
	}
	if raw["detail"] != "Your current balance is 30, but that costs 50." {
		t.Errorf("unexpected detail in JSON: %v", raw["detail"])
	}
	if raw["instance"] != "/account/12345/msgs/abc" {
		t.Errorf("unexpected instance in JSON: %v", raw["instance"])
	}
	if raw["balance"] != float64(30) {
		t.Errorf("unexpected extension member in JSON: %v", raw["balance"])
	}
}

func TestProblemDetailsUnmarshalJSON(t *testing.T) {
	jsonPayload := `{
		"type": "https://example.com/probs/out-of-credit",
		"title": "You do not have enough credit.",
		"status": 403,
		"detail": "Your current balance is 30, but that costs 50.",
		"instance": "/account/12345/msgs/abc",
		"balance": 30,
		"accounts": ["/account/12345", "/account/67890"]
	}`

	var p ProblemDetails
	if err := json.Unmarshal([]byte(jsonPayload), &p); err != nil {
		t.Fatalf("failed to unmarshal problem details: %v", err)
	}

	if p.Type != "https://example.com/probs/out-of-credit" {
		t.Errorf("unexpected type: %s", p.Type)
	}
	if p.Title != "You do not have enough credit." {
		t.Errorf("unexpected title: %s", p.Title)
	}
	if p.Status != 403 {
		t.Errorf("unexpected status: %d", p.Status)
	}
	if p.Detail != "Your current balance is 30, but that costs 50." {
		t.Errorf("unexpected detail: %s", p.Detail)
	}
	if p.Instance != "/account/12345/msgs/abc" {
		t.Errorf("unexpected instance: %s", p.Instance)
	}
	if p.Extensions == nil {
		t.Fatalf("expected extensions map to be populated")
	}
	if p.Extensions["balance"] != float64(30) {
		t.Errorf("unexpected extension balance: %v", p.Extensions["balance"])
	}
	accounts, ok := p.Extensions["accounts"].([]any)
	if !ok || len(accounts) != 2 {
		t.Errorf("unexpected extension accounts: %v", p.Extensions["accounts"])
	}
}

func TestProblemDetailsWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	prob := NewProblemDetails(http.StatusNotFound, "Resource not found").
		WithInstance("/records/999")

	err := prob.Write(rec)
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != ContentTypeProblemJSON {
		t.Errorf("expected Content-Type '%s', got '%s'", ContentTypeProblemJSON, contentType)
	}

	var decoded ProblemDetails
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if decoded.Status != http.StatusNotFound {
		t.Errorf("expected decoded status %d, got %d", http.StatusNotFound, decoded.Status)
	}
	if decoded.Title != "Not Found" {
		t.Errorf("expected decoded title 'Not Found', got '%s'", decoded.Title)
	}
	if decoded.Detail != "Resource not found" {
		t.Errorf("expected decoded detail 'Resource not found', got '%s'", decoded.Detail)
	}
	if decoded.Instance != "/records/999" {
		t.Errorf("expected decoded instance '/records/999', got '%s'", decoded.Instance)
	}
}

func TestWriteProblem(t *testing.T) {
	rec := httptest.NewRecorder()
	prob := NewProblemDetails(http.StatusInternalServerError, "Internal server error occurred")

	err := WriteProblem(rec, prob)
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if rec.Header().Get("Content-Type") != ContentTypeProblemJSON {
		t.Errorf("expected Content-Type '%s', got '%s'", ContentTypeProblemJSON, rec.Header().Get("Content-Type"))
	}
}
