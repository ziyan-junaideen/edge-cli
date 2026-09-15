package edgeapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestCreatePaymentDemandSendsUpdatedPaymentSchema(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", request.Method)
		}
		if request.URL.Path != "/v2/payment_demands" {
			t.Fatalf("expected payment demand path, got %q", request.URL.Path)
		}
		if request.URL.Query().Get("include") != "payer,billing_address" {
			t.Fatalf("expected include query, got %q", request.URL.RawQuery)
		}
		if request.Header.Get("Content-Type") != "application/vnd.api+json" {
			t.Fatalf("expected JSON:API content type, got %q", request.Header.Get("Content-Type"))
		}

		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		data := body["data"].(map[string]any)
		attributes := data["attributes"].(map[string]any)
		if attributes["tax_detail"].(map[string]any)["tax_cents"] != float64(825) {
			t.Fatalf("expected tax detail, got %#v", attributes["tax_detail"])
		}
		if attributes["shipping_detail"].(map[string]any)["shipping_cents"] != float64(500) {
			t.Fatalf("expected shipping detail, got %#v", attributes["shipping_detail"])
		}
		lineItems := attributes["line_items"].([]any)
		if len(lineItems) != 1 || lineItems[0].(map[string]any)["commodity_code"] != "12345678" {
			t.Fatalf("expected complete line item, got %#v", lineItems)
		}
		relationships := data["relationships"].(map[string]any)
		payer := relationships["payer"].(map[string]any)["data"].(map[string]any)
		if payer["id"] != "customer-id" || payer["type"] != "customers" {
			t.Fatalf("expected payer relationship, got %#v", payer)
		}

		responseWriter.Header().Set("Content-Type", "application/vnd.api+json")
		responseWriter.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"data": map[string]any{
				"id":         "payment-demand-id",
				"type":       "payment_demands",
				"attributes": attributes,
			},
		})
	}))
	defer server.Close()

	client, err := New(Config{APIURL: server.URL + "/v2", Token: "test-token", InsecureSkipVerify: true})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	amountCents, quantity := 10000, 1
	resource, _, err := client.CreatePaymentDemand(context.Background(), CreatePaymentDemandParams{
		Attributes: CreatePaymentDemandAttributes{
			AmountCents:    amountCents,
			AmountCurrency: "USD",
			IdempotencyKey: "idempotency-key",
			TaxDetail:      &PaymentTaxDetail{TaxCents: 825, TaxCurrency: "USD"},
			ShippingDetail: &PaymentShippingDetail{ShippingCents: 500, ShippingCurrency: "USD"},
			LineItems: []PaymentLineItem{{
				Name: "Widget", CommodityCode: "12345678", AmountCents: &amountCents,
				AmountCurrency: "USD", Quantity: &quantity,
			}},
		},
		PayerID: "customer-id",
	}, QueryOptions{Include: []string{"payer", "billing_address"}})
	if err != nil {
		t.Fatalf("CreatePaymentDemand returned error: %v", err)
	}
	if resource.ID != "payment-demand-id" {
		t.Fatalf("expected created payment demand, got %q", resource.ID)
	}
}

func TestShowWebhookDeliveryRequestsReplayRelationships(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/webhook_deliveries/delivery-id" {
			t.Fatalf("expected webhook delivery path, got %q", request.URL.Path)
		}
		if request.URL.Query().Get("include") != "event,webhook_subscription" {
			t.Fatalf("expected replay includes, got %q", request.URL.RawQuery)
		}

		responseWriter.Header().Set("Content-Type", "application/vnd.api+json")
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"data": map[string]any{
				"id":   "delivery-id",
				"type": "webhook_deliveries",
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

	delivery, _, err := client.ShowWebhookDelivery(
		context.Background(),
		"delivery-id",
		QueryOptions{Include: []string{"event", "webhook_subscription"}},
	)
	if err != nil {
		t.Fatalf("ShowWebhookDelivery returned error: %v", err)
	}
	if delivery.ID != "delivery-id" {
		t.Fatalf("expected delivery id, got %q", delivery.ID)
	}
}

func TestAPIErrorFormatsForbiddenJSONAPIError(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(responseWriter).Encode(map[string]any{
			"errors": []map[string]any{
				{
					"status": "403",
					"code":   "forbidden",
					"title":  "Forbidden",
					"detail": "Token does not have permission to access payment_demands.",
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

	_, _, err = client.ListPaymentDemands(context.Background(), QueryOptions{})
	if err == nil {
		t.Fatal("expected API error")
	}

	var apiError APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiError.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", apiError.StatusCode)
	}

	errorMessage := err.Error()
	for _, expected := range []string{"403 Forbidden", "forbidden", "Token does not have permission"} {
		if !strings.Contains(errorMessage, expected) {
			t.Fatalf("expected error message to contain %q, got %q", expected, errorMessage)
		}
	}
}

func TestAPIErrorIncludesPlainTextBody(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusBadGateway)
		_, _ = responseWriter.Write([]byte("upstream unavailable"))
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

	_, _, err = client.ListPaymentDemands(context.Background(), QueryOptions{})
	if err == nil {
		t.Fatal("expected API error")
	}

	errorMessage := err.Error()
	for _, expected := range []string{"502 Bad Gateway", "upstream unavailable"} {
		if !strings.Contains(errorMessage, expected) {
			t.Fatalf("expected error message to contain %q, got %q", expected, errorMessage)
		}
	}
}
