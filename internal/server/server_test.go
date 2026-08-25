package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"boussinesq-z/internal/circular"
)

func testRoutes() http.Handler {
	store := circular.NewStore("../../data/circular-influence.csv")
	return Routes(store)
}

func postJSON(t *testing.T, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRoutes().ServeHTTP(rec, req)
	return rec
}

func TestStressEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/stress", map[string]interface{}{
		"P": 100000,
		"z": 2,
		"r": 0,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["sigma_z"] == nil {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestSuperposeEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/superpose", map[string]interface{}{
		"forces": []map[string]interface{}{
			{"P": 100000, "z": 2, "r": 0},
			{"P": 50000, "z": 3, "r": 1},
		},
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCircularEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/circular", map[string]interface{}{
		"q": 100000,
		"a": 1,
		"z": 2,
		"r": 0,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["influence_i"] == nil {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestInvalidStressReturns400(t *testing.T) {
	rec := postJSON(t, "/api/stress", map[string]interface{}{
		"P": 100,
		"z": 0,
		"r": 0,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestTableSnapshotEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/table-snapshot", nil)
	rec := httptest.NewRecorder()
	testRoutes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	testRoutes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/stress", nil)
	rec := httptest.NewRecorder()
	testRoutes().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestZeroLoadEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/stress", map[string]interface{}{
		"P": 0,
		"z": 2,
		"r": 0,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["sigma_z"] != 0.0 {
		t.Fatalf("body = %s", rec.Body.String())
	}
}
