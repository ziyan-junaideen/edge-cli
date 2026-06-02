package edgeapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCustomersSendsIncludeQuery(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/customers" {
			t.Fatalf("expected /v2/customers path, got %q", request.URL.Path)
		}
		if request.URL.Query().Get("include") != "addresses,merchant" {
			t.Fatalf("expected include query, got %q", request.URL.RawQuery)
		}
		if request.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("expected bearer token header")
		}

		responseWriter.Header().Set("Content-Type", "application/vnd.api+json")
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":   "customer-id",
					"type": "customers",
					"attributes": map[string]any{
						"name": "Jane Doe",
					},
				},
			},
			"included": []map[string]any{
				{
					"id":   "address-id",
					"type": "consumer_addresses",
				},
			},
		})
	}))
	defer server.Close()

	client, err := New(Config{
		APIURL:             server.URL + "/v2",
		Token:              "test-token",
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	client.httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	customers, document, err := client.ListCustomers(context.Background(), QueryOptions{Include: []string{"addresses", "merchant"}})
	if err != nil {
		t.Fatalf("ListCustomers returned error: %v", err)
	}

	if len(customers) != 1 {
		t.Fatalf("expected one customer, got %d", len(customers))
	}
	if len(document.Included) == 0 {
		t.Fatal("expected included data to be preserved")
	}
}

func TestShowPaymentDemandSendsIncludeQuery(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/payment_demands/payment-demand-id" {
			t.Fatalf("expected payment demand path, got %q", request.URL.Path)
		}
		if request.URL.Query().Get("include") != "payer,billing_address,payment_method" {
			t.Fatalf("expected include query, got %q", request.URL.RawQuery)
		}

		responseWriter.Header().Set("Content-Type", "application/vnd.api+json")
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"data": map[string]any{
				"id":   "payment-demand-id",
				"type": "payment_demands",
				"attributes": map[string]any{
					"amount_cents":    1000,
					"amount_currency": "USD",
				},
			},
		})
	}))
	defer server.Close()

	client, err := New(Config{
		APIURL:             server.URL + "/v2",
		Token:              "test-token",
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	paymentDemand, _, err := client.ShowPaymentDemand(
		context.Background(),
		"payment-demand-id",
		QueryOptions{Include: []string{"payer", "billing_address", "payment_method"}},
	)
	if err != nil {
		t.Fatalf("ShowPaymentDemand returned error: %v", err)
	}

	if paymentDemand.ID != "payment-demand-id" {
		t.Fatalf("expected payment demand id, got %q", paymentDemand.ID)
	}
}
