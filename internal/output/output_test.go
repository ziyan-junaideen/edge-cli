package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/edgepayments/ept-cli/internal/jsonapi"
)

func TestMoneyStringFormatsCentsAsMajorUnits(t *testing.T) {
	resource := resourceWithAttributes(t, map[string]any{
		"amount_cents":    10000,
		"amount_currency": "USD",
	})

	if got := moneyString(resource, "amount_cents", "amount_currency"); got != "100.00 USD" {
		t.Fatalf("expected 100.00 USD, got %q", got)
	}
}

func TestPaymentDemandCollectionFormatsAmount(t *testing.T) {
	var buffer bytes.Buffer
	paymentDemand := resourceWithAttributes(t, map[string]any{
		"amount_cents":    10000,
		"amount_currency": "USD",
		"processor_state": "succeeded",
	})
	paymentDemand.ID = "payment-demand-id"

	if err := PaymentDemandCollection(&buffer, []jsonapi.Resource{paymentDemand}); err != nil {
		t.Fatalf("PaymentDemandCollection returned error: %v", err)
	}

	if !strings.Contains(buffer.String(), "100.00 USD") {
		t.Fatalf("expected formatted amount in output, got:\n%s", buffer.String())
	}
	if strings.Contains(buffer.String(), "10000 USD") {
		t.Fatalf("expected raw cent amount to be hidden, got:\n%s", buffer.String())
	}
}

func TestShowResourcePrintsRelationshipIDsAndIncludedDetails(t *testing.T) {
	var buffer bytes.Buffer
	paymentDemand := resourceWithAttributes(t, map[string]any{
		"amount_cents":    10000,
		"amount_currency": "USD",
		"processor_state": "succeeded",
		"capture_method":  "automatic",
		"email_receipt":   true,
		"succeeded_at":    "2026-06-01T17:27:43Z",
		"created_at":      "2026-06-01T17:27:18Z",
		"updated_at":      "2026-06-01T17:27:43Z",
	})
	paymentDemand.ID = "payment-demand-id"
	paymentDemand.Type = "payment_demands"
	paymentDemand.Relationships = map[string]json.RawMessage{
		"payer": mustMarshal(t, map[string]any{
			"data": map[string]any{"type": "customers", "id": "customer-id"},
		}),
	}

	includedCustomer := resourceWithAttributes(t, map[string]any{
		"name":         "Jane Doe",
		"email":        "jane@example.com",
		"phone_number": "+15555550100",
		"description":  "VIP",
		"created_at":   "2026-06-01T17:20:00Z",
		"updated_at":   "2026-06-01T17:21:00Z",
	})
	includedCustomer.ID = "customer-id"
	includedCustomer.Type = "customers"

	document := jsonapi.Document{
		Included: mustMarshal(t, []jsonapi.Resource{includedCustomer}),
	}

	if err := ShowResource(&buffer, paymentDemand, document, []string{"payer"}); err != nil {
		t.Fatalf("ShowResource returned error: %v", err)
	}

	output := buffer.String()
	for _, expected := range []string{
		"Relationships:",
		"Payer",
		"customers customer-id",
		"Included:",
		"ID: customer-id",
		"Type: customers",
		"Name: Jane Doe",
		"Email: jane@example.com",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestPaymentSubscriptionShowsExpandedAttributesAndLineItems(t *testing.T) {
	var buffer bytes.Buffer
	paymentSubscription := resourceWithAttributes(t, map[string]any{
		"status":                     "active",
		"description":                "Test",
		"amount_cents":               10000,
		"amount_currency":            "USD",
		"slug":                       "paylink_slug",
		"purchase_reference":         "00000001",
		"email_receipt":              true,
		"payer_timezone":             "Asia/Colombo",
		"created_at":                 "2026-06-04T03:59:09Z",
		"updated_at":                 "2026-06-04T03:59:25Z",
		"discount_cents":             0,
		"billing_period":             "one_month",
		"billing_cycle_anchor_at":    "2026-06-04T03:58:37Z",
		"idempotency_key":            "idempotency-key",
		"purchase_kind":              "order",
		"cvc2_check":                 "unprocessed",
		"address_line1_verification": "unverified",
		"postal_code_verification":   "unverified",
		"fee_cents":                  0,
		"proration_behavior":         "none",
		"acs_transaction_eid":        "acs-id",
		"directory_transaction_eid":  "directory-id",
		"eci":                        "05",
		"threeds_cryptogram":         "cryptogram",
		"threeds_status":             "Y",
		"threeds_version":            "2.2.0",
		"line_items": []map[string]any{
			{
				"name":              "Test Subscription",
				"description":       "Test",
				"amount_cents":      10000,
				"amount_currency":   "USD",
				"quantity":          1,
				"tax_cents":         0,
				"tax_currency":      "USD",
				"discount_cents":    0,
				"discount_currency": "USD",
			},
		},
	})
	paymentSubscription.ID = "payment-subscription-id"
	paymentSubscription.Type = "payment_subscriptions"

	if err := PaymentSubscription(&buffer, paymentSubscription, jsonapi.Document{}, nil); err != nil {
		t.Fatalf("PaymentSubscription returned error: %v", err)
	}

	output := buffer.String()
	for _, expected := range []string{
		"Description: Test",
		"Slug: paylink_slug",
		"Amount: 100.00 USD",
		"Discount: 0.00 USD",
		"Fee: 0.00 USD",
		"Billing Period: one_month",
		"Proration Behavior: none",
		"Purchase Reference: 00000001",
		"Idempotency Key: idempotency-key",
		"CVC2 Check: unprocessed",
		"3DS Version: 2.2.0",
		"Line Items:",
		"Test Subscription",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}

func resourceWithAttributes(t *testing.T, attributes map[string]any) jsonapi.Resource {
	t.Helper()

	rawAttributes := map[string]json.RawMessage{}
	for key, value := range attributes {
		encodedValue, err := json.Marshal(value)
		if err != nil {
			t.Fatalf("marshal attribute %q: %v", key, err)
		}
		rawAttributes[key] = encodedValue
	}
	return jsonapi.Resource{Attributes: rawAttributes}
}

func mustMarshal(t *testing.T, value any) json.RawMessage {
	t.Helper()

	encodedValue, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal value: %v", err)
	}
	return encodedValue
}
