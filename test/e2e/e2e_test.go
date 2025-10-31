package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const baseURL = "http://localhost:8081"

func mustHTTP(t *testing.T) *http.Client {
	t.Helper()
	return &http.Client{Timeout: 10 * time.Second}
}

func requireServerUp(t *testing.T) {
	t.Helper()
	httpClient := mustHTTP(t)
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := httpClient.Get(baseURL + "/metrics")
		if err == nil && resp.StatusCode < 500 {
			_ = resp.Body.Close()
			return
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("server is not up at %s", baseURL)
}

func TestAcceptOrderDeliveryAndGetOrders(t *testing.T) {
	if os.Getenv("E2E") == "" {
		t.Skip("E2E env var not set; skip E2E tests")
	}
	requireServerUp(t)

	httpClient := mustHTTP(t)

	orderID := "ord-e2e-1"
	userID := "user-e2e-1"

	// 1) Accept order delivery
	payload := map[string]any{
		"orderId":        orderID,
		"recipientId":    userID,
		"storageTime":    "86400s", // 1 day
		"cost":           1500,
		"weight":         700,
		"packaging":      "BOX",
		"additionalFilm": false,
	}
	body, _ := json.Marshal(payload)
	resp, err := httpClient.Post(baseURL+"/v1/pvz-service/accept-order-delivery", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("accept-order-delivery request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("accept-order-delivery unexpected status: %d, body: %s", resp.StatusCode, string(b))
	}

	// 2) Get orders for the user and ensure our order is present
	resp, err = httpClient.Get(baseURL + "/v1/pvz-service/get-orders?userId=" + userID + "&samePVZ=true")
	if err != nil {
		t.Fatalf("get-orders request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("get-orders unexpected status: %d, body: %s", resp.StatusCode, string(b))
	}

	var got struct {
		Orders []struct {
			OrderID string `json:"orderId"`
		} `json:"orders"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode get-orders response: %v", err)
	}
	found := false
	for _, o := range got.Orders {
		if o.OrderID == orderID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("order %s not found in get-orders response", orderID)
	}
}

func TestAcceptOrderDelivery_ValidationError(t *testing.T) {
	if os.Getenv("E2E") == "" {
		t.Skip("E2E env var not set; skip E2E tests")
	}
	requireServerUp(t)

	httpClient := mustHTTP(t)

	// Business validation: additionalFilm cannot be true when packaging is FILM
	payload := map[string]any{
		"orderId":        "ord-e2e-invalid",
		"recipientId":    "user-e2e-invalid",
		"storageTime":    "3600s",
		"cost":           100,
		"weight":         100,
		"packaging":      "FILM",
		"additionalFilm": true,
	}
	body, _ := json.Marshal(payload)
	resp, err := httpClient.Post(baseURL+"/v1/pvz-service/accept-order-delivery", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("accept-order-delivery request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest { // mapped from gRPC InvalidArgument
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 400 Bad Request, got %d, body: %s", resp.StatusCode, string(b))
	}
}
