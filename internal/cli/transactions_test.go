package cli

import (
	"strings"
	"testing"
)

func TestParsePaymentLineItemsUsesCompleteSchema(t *testing.T) {
	lineItems, err := parsePaymentLineItems([]string{`{
		"name":"Widget",
		"sku":"widget-001",
		"unit_of_measure":"unit",
		"description":"A standard widget",
		"commodity_code":"12345678",
		"amount_cents":10000,
		"amount_currency":"usd",
		"tax_cents":825,
		"tax_currency":"usd",
		"discount_cents":250,
		"discount_currency":"usd",
		"quantity":2
	}`})
	if err != nil {
		t.Fatalf("parsePaymentLineItems returned error: %v", err)
	}
	if len(lineItems) != 1 {
		t.Fatalf("expected one line item, got %d", len(lineItems))
	}
	lineItem := lineItems[0]
	if lineItem.CommodityCode != "12345678" || lineItem.AmountCurrency != "USD" || lineItem.TaxCurrency != "USD" {
		t.Fatalf("unexpected parsed line item: %#v", lineItem)
	}
}

func TestParsePaymentLineItemsRequiresSchemaFields(t *testing.T) {
	_, err := parsePaymentLineItems([]string{`{"name":"Widget"}`})
	if err == nil || !strings.Contains(err.Error(), "requires amount_cents, amount_currency, and quantity") {
		t.Fatalf("expected required field error, got %v", err)
	}
}

func TestPaymentDemandCreateCommandExposesUpdatedAttributes(t *testing.T) {
	command := newPaymentDemandsCommand(&globalOptions{})
	createCommand, _, err := command.Find([]string{"create"})
	if err != nil {
		t.Fatalf("find create command: %v", err)
	}
	for _, flagName := range []string{
		"tax-cents", "tax-currency", "shipping-cents", "shipping-currency", "line-item",
		"payer-timezone", "threeds-version", "threeds-status", "threeds-cryptogram", "eci",
		"directory-transaction-eid", "acs-transaction-eid", "payer", "buyer", "receiver",
	} {
		if createCommand.Flags().Lookup(flagName) == nil {
			t.Errorf("expected --%s flag", flagName)
		}
	}
}
