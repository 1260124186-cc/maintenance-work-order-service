package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/store"
)

func TestHealthAndWorkOrderLifecycle(t *testing.T) {
	server := NewServer(store.NewMemoryRepository())

	health := httptest.NewRecorder()
	server.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health status = %d", health.Code)
	}

	body := bytes.NewBufferString(`{"asset_id":"FAN-01","title":"Inspect belt","priority":"urgent","labels":["safety"]}`)
	create := httptest.NewRecorder()
	server.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/work-orders", body))
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", create.Code, create.Body.String())
	}

	var order struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(create.Body).Decode(&order); err != nil {
		t.Fatalf("decode created work order: %v", err)
	}

	assign := httptest.NewRecorder()
	server.ServeHTTP(assign, httptest.NewRequest(http.MethodPost, "/work-orders/"+order.ID+"/assign", bytes.NewBufferString(`{"technician":"Jordan"}`)))
	if assign.Code != http.StatusOK {
		t.Fatalf("assign status = %d, body = %s", assign.Code, assign.Body.String())
	}

	complete := httptest.NewRecorder()
	server.ServeHTTP(complete, httptest.NewRequest(http.MethodPost, "/work-orders/"+order.ID+"/complete", nil))
	if complete.Code != http.StatusOK {
		t.Fatalf("complete status = %d, body = %s", complete.Code, complete.Body.String())
	}
}

func TestUnknownAssetReturnsNotFound(t *testing.T) {
	server := NewServer(store.NewMemoryRepository())
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/work-orders",
		bytes.NewBufferString(`{"asset_id":"NOPE","title":"Inspect","priority":"normal"}`)))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
