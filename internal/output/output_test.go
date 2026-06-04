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

func TestPaymentDemandShowsExpandedAttributesAndLineItems(t *testing.T) {
	var buffer bytes.Buffer
	paymentDemand := resourceWithAttributes(t, map[string]any{
		"description":                "Test charge",
		"amount_cents":               10000,
		"amount_currency":            "USD",
		"discount_cents":             250,
		"fee_cents":                  53,
		"processor_state":            "succeeded",
		"capture_method":             "automatic",
		"purchase_reference":         "00000001",
		"purchase_kind":              "order",
		"payer_timezone":             "Asia/Colombo",
		"idempotency_key":            "idempotency-key",
		"email_receipt":              true,
		"cvc2_check":                 "match",
		"address_line1_verification": "match",
		"postal_code_verification":   "match",
		"threeds_version":            "2.2.0",
		"threeds_status":             "Y",
		"threeds_cryptogram":         "cryptogram",
		"eci":                        "05",
		"directory_transaction_eid":  "directory-id",
		"acs_transaction_eid":        "acs-id",
		"succeeded_at":               "2026-06-01T17:27:43Z",
		"created_at":                 "2026-06-01T17:27:18Z",
		"updated_at":                 "2026-06-01T17:27:43Z",
		"line_items": []map[string]any{
			{
				"name":              "Widget",
				"description":       "A standard widget",
				"amount_cents":      10000,
				"amount_currency":   "USD",
				"quantity":          1,
				"tax_cents":         825,
				"tax_currency":      "USD",
				"discount_cents":    250,
				"discount_currency": "USD",
			},
		},
	})
	paymentDemand.ID = "payment-demand-id"
	paymentDemand.Type = "payment_demands"

	if err := PaymentDemand(&buffer, paymentDemand, jsonapi.Document{}, nil); err != nil {
		t.Fatalf("PaymentDemand returned error: %v", err)
	}

	output := buffer.String()
	for _, expected := range []string{
		"Description: Test charge",
		"Amount: 100.00 USD",
		"Discount: 2.50 USD",
		"Fee: 0.53 USD",
		"Purchase Reference: 00000001",
		"Payer Timezone: Asia/Colombo",
		"Idempotency Key: idempotency-key",
		"CVC2 Check: match",
		"3DS Version: 2.2.0",
		"Line Items:",
		"Widget",
		"8.25 USD",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
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

func TestShowResourceSkipsRelationshipSectionWhenNoRelationshipIDs(t *testing.T) {
	var buffer bytes.Buffer
	merchant := resourceWithAttributes(t, map[string]any{
		"business_name": "Edge",
	})
	merchant.ID = "merchant-id"
	merchant.Type = "merchants"
	merchant.Relationships = map[string]json.RawMessage{
		"customers": mustMarshal(t, map[string]any{
			"links": map[string]any{"self": "https://example.test/relationships/customers"},
		}),
	}

	if err := ShowResource(&buffer, merchant, jsonapi.Document{}, nil); err != nil {
		t.Fatalf("ShowResource returned error: %v", err)
	}

	if strings.Contains(buffer.String(), "Relationships:") {
		t.Fatalf("expected no empty relationship section, got:\n%s", buffer.String())
	}
}

func TestPaymentMethodDoesNotRenderEmptyExpirySlash(t *testing.T) {
	var buffer bytes.Buffer
	paymentMethod := resourceWithAttributes(t, map[string]any{
		"kind":           "visa",
		"last_four":      "0004",
		"card_pan_token": "pan-token",
		"card_cvv_token": "cvv-token",
		"external_state": "confirmed",
	})
	paymentMethod.ID = "payment-method-id"
	paymentMethod.Type = "payment_methods"

	if err := PaymentMethod(&buffer, paymentMethod); err != nil {
		t.Fatalf("PaymentMethod returned error: %v", err)
	}

	if strings.Contains(buffer.String(), "Expiry: /") {
		t.Fatalf("expected empty expiry not to render as slash, got:\n%s", buffer.String())
	}
	for _, expected := range []string{"Card PAN Token: pan-token", "Card CVV Token: cvv-token"} {
		if !strings.Contains(buffer.String(), expected) {
			t.Fatalf("expected token field %q in output, got:\n%s", expected, buffer.String())
		}
	}
}

func TestPaymentMethodRendersExpiryWhenPresent(t *testing.T) {
	var buffer bytes.Buffer
	paymentMethod := resourceWithAttributes(t, map[string]any{
		"expiry_month": 12,
		"expiry_year":  2030,
	})
	paymentMethod.ID = "payment-method-id"
	paymentMethod.Type = "payment_methods"

	if err := PaymentMethod(&buffer, paymentMethod); err != nil {
		t.Fatalf("PaymentMethod returned error: %v", err)
	}

	if !strings.Contains(buffer.String(), "Expiry: 12/2030") {
		t.Fatalf("expected expiry in output, got:\n%s", buffer.String())
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
