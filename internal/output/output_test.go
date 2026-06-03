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
