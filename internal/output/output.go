package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ziyan-junaideen/edge-cli/internal/jsonapi"
)

func JSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func MerchantCollection(writer io.Writer, merchants []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tBUSINESS\tEMAIL\tWEBSITE\tENTITY"); err != nil {
		return err
	}
	for _, merchant := range merchants {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\n",
			merchant.ID,
			attributeString(merchant, "business_name"),
			attributeString(merchant, "business_email"),
			attributeString(merchant, "business_website"),
			attributeString(merchant, "entity_type"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func CustomerCollection(writer io.Writer, customers []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tNAME\tEMAIL\tPHONE\tCREATED"); err != nil {
		return err
	}
	for _, customer := range customers {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\n",
			customer.ID,
			attributeString(customer, "name"),
			attributeString(customer, "email"),
			attributeString(customer, "phone_number"),
			attributeString(customer, "created_at"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func ConsumerAddressCollection(writer io.Writer, consumerAddresses []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tLINE 1\tCITY\tSTATE\tZIP\tCOUNTRY"); err != nil {
		return err
	}
	for _, consumerAddress := range consumerAddresses {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			consumerAddress.ID,
			attributeString(consumerAddress, "line_1"),
			attributeString(consumerAddress, "city"),
			attributeString(consumerAddress, "state"),
			attributeString(consumerAddress, "zip"),
			attributeString(consumerAddress, "country"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func PaymentDemandCollection(writer io.Writer, paymentDemands []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tAMOUNT\tSTATE\tPROCESSOR\tSUCCEEDED\tCREATED"); err != nil {
		return err
	}
	for _, paymentDemand := range paymentDemands {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			paymentDemand.ID,
			moneyString(paymentDemand, "amount_cents", "amount_currency"),
			firstAttributeString(paymentDemand, "state", "status"),
			attributeString(paymentDemand, "processor_state"),
			attributeString(paymentDemand, "succeeded_at"),
			attributeString(paymentDemand, "created_at"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func PaymentSubscriptionCollection(writer io.Writer, paymentSubscriptions []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tAMOUNT\tSTATUS\tINTERVAL\tCREATED"); err != nil {
		return err
	}
	for _, paymentSubscription := range paymentSubscriptions {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\n",
			paymentSubscription.ID,
			moneyString(paymentSubscription, "amount_cents", "amount_currency"),
			attributeString(paymentSubscription, "status"),
			firstAttributeString(paymentSubscription, "billing_interval", "interval"),
			attributeString(paymentSubscription, "created_at"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func PaymentMethodCollection(writer io.Writer, paymentMethods []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tKIND\tNICKNAME\tLAST FOUR\tEXPIRY\tSTATE"); err != nil {
		return err
	}
	for _, paymentMethod := range paymentMethods {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s/%s\t%s\n",
			paymentMethod.ID,
			attributeString(paymentMethod, "kind"),
			attributeString(paymentMethod, "nickname"),
			attributeString(paymentMethod, "last_four"),
			attributeString(paymentMethod, "expiry_month"),
			attributeString(paymentMethod, "expiry_year"),
			attributeString(paymentMethod, "external_state"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func RefundDemandCollection(writer io.Writer, refundDemands []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tAMOUNT\tSTATE\tREASON\tCREATED"); err != nil {
		return err
	}
	for _, refundDemand := range refundDemands {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\n",
			refundDemand.ID,
			moneyString(refundDemand, "amount_cents", "amount_currency"),
			attributeString(refundDemand, "state"),
			attributeString(refundDemand, "reason"),
			attributeString(refundDemand, "created_at"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func UserResourceCollection(writer io.Writer, resources []jsonapi.Resource) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tTYPE\tSUMMARY\tSTATUS\tCREATED"); err != nil {
		return err
	}
	for _, resource := range resources {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\n",
			resource.ID,
			resource.Type,
			userResourceSummary(resource),
			firstAttributeString(resource, "status", "account_status", "action"),
			attributeString(resource, "created_at"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
}

func ShowResource(writer io.Writer, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
	if err := resourceDetails(writer, resource, document, includes); err != nil {
		return err
	}
	if err := RelationshipIdentifiers(writer, resource); err != nil {
		return err
	}
	return IncludedResources(writer, resource, document, relationshipNames(includes, defaultRelationshipNames(resource.Type)))
}

func resourceDetails(writer io.Writer, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
	switch resource.Type {
	case "merchants":
		return Merchant(writer, resource)
	case "customers":
		return Customer(writer, resource)
	case "consumer_addresses":
		return ConsumerAddress(writer, resource)
	case "payment_demands":
		return PaymentDemand(writer, resource, jsonapi.Document{}, nil)
	case "payment_subscriptions":
		return PaymentSubscription(writer, resource, jsonapi.Document{}, nil)
	case "payment_methods":
		return PaymentMethod(writer, resource)
	case "refund_demands":
		return RefundDemand(writer, resource)
	case "account_alerts", "accounts", "memberships", "merchant_punitive_actions", "permissions", "red_flags":
		return UserResource(writer, resource, jsonapi.Document{}, nil)
	default:
		_, err := fmt.Fprintf(writer, "ID: %s\nType: %s\n", resource.ID, resource.Type)
		return err
	}
}

func Merchant(writer io.Writer, merchant jsonapi.Resource) error {
	_, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Business: %s\n"+
			"Description: %s\n"+
			"Email: %s\n"+
			"Phone: %s\n"+
			"Website: %s\n"+
			"Entity: %s\n"+
			"Category Code: %s\n"+
			"Timezone: %s\n"+
			"Zip Code: %s\n"+
			"Support Email: %s\n"+
			"Support URL: %s\n"+
			"Terms URL: %s\n"+
			"Privacy Policy URL: %s\n"+
			"Prior Bankruptcies: %s\n"+
			"Active At: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		merchant.ID,
		merchant.Type,
		attributeString(merchant, "business_name"),
		attributeString(merchant, "business_description"),
		attributeString(merchant, "business_email"),
		attributeString(merchant, "phone_number"),
		attributeString(merchant, "business_website"),
		attributeString(merchant, "entity_type"),
		attributeString(merchant, "category_code"),
		attributeString(merchant, "business_timezone"),
		attributeString(merchant, "business_zip_code"),
		attributeString(merchant, "business_support_email"),
		attributeString(merchant, "business_support_url"),
		attributeString(merchant, "business_terms_url"),
		attributeString(merchant, "business_privacy_policy_url"),
		attributeString(merchant, "prior_bankruptcies"),
		attributeString(merchant, "active_at"),
		attributeString(merchant, "created_at"),
		attributeString(merchant, "updated_at"),
	)
	return err
}

func Customer(writer io.Writer, customer jsonapi.Resource) error {
	_, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Name: %s\n"+
			"Email: %s\n"+
			"Phone: %s\n"+
			"Description: %s\n"+
			"Blocked At: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		customer.ID,
		customer.Type,
		attributeString(customer, "name"),
		attributeString(customer, "email"),
		attributeString(customer, "phone_number"),
		attributeString(customer, "description"),
		attributeString(customer, "blocked_at"),
		attributeString(customer, "created_at"),
		attributeString(customer, "updated_at"),
	)
	return err
}

func ConsumerAddress(writer io.Writer, consumerAddress jsonapi.Resource) error {
	_, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Line 1: %s\n"+
			"Line 2: %s\n"+
			"City: %s\n"+
			"State: %s\n"+
			"Zip: %s\n"+
			"Country: %s\n"+
			"Discarded At: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		consumerAddress.ID,
		consumerAddress.Type,
		attributeString(consumerAddress, "line_1"),
		attributeString(consumerAddress, "line_2"),
		attributeString(consumerAddress, "city"),
		attributeString(consumerAddress, "state"),
		attributeString(consumerAddress, "zip"),
		attributeString(consumerAddress, "country"),
		attributeString(consumerAddress, "discarded_at"),
		attributeString(consumerAddress, "created_at"),
		attributeString(consumerAddress, "updated_at"),
	)
	return err
}

func PaymentDemand(writer io.Writer, paymentDemand jsonapi.Resource, document jsonapi.Document, includes []string) error {
	if _, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Description: %s\n"+
			"Amount: %s\n"+
			"Discount: %s\n"+
			"Fee: %s\n"+
			"State: %s\n"+
			"Processor State: %s\n"+
			"Capture Method: %s\n"+
			"Purchase Reference: %s\n"+
			"Purchase Kind: %s\n"+
			"Payer Timezone: %s\n"+
			"Idempotency Key: %s\n"+
			"Email Receipt: %s\n"+
			"CVC2 Check: %s\n"+
			"Address Line 1 Verification: %s\n"+
			"Postal Code Verification: %s\n"+
			"3DS Version: %s\n"+
			"3DS Status: %s\n"+
			"3DS Cryptogram: %s\n"+
			"ECI: %s\n"+
			"Directory Transaction EID: %s\n"+
			"ACS Transaction EID: %s\n"+
			"Succeeded At: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		paymentDemand.ID,
		paymentDemand.Type,
		attributeString(paymentDemand, "description"),
		moneyString(paymentDemand, "amount_cents", "amount_currency"),
		moneyAttributeString(paymentDemand, "discount_cents", "amount_currency"),
		moneyAttributeString(paymentDemand, "fee_cents", "amount_currency"),
		firstAttributeString(paymentDemand, "state", "status"),
		attributeString(paymentDemand, "processor_state"),
		attributeString(paymentDemand, "capture_method"),
		attributeString(paymentDemand, "purchase_reference"),
		attributeString(paymentDemand, "purchase_kind"),
		attributeString(paymentDemand, "payer_timezone"),
		attributeString(paymentDemand, "idempotency_key"),
		attributeString(paymentDemand, "email_receipt"),
		attributeString(paymentDemand, "cvc2_check"),
		attributeString(paymentDemand, "address_line1_verification"),
		attributeString(paymentDemand, "postal_code_verification"),
		attributeString(paymentDemand, "threeds_version"),
		attributeString(paymentDemand, "threeds_status"),
		attributeString(paymentDemand, "threeds_cryptogram"),
		attributeString(paymentDemand, "eci"),
		attributeString(paymentDemand, "directory_transaction_eid"),
		attributeString(paymentDemand, "acs_transaction_eid"),
		attributeString(paymentDemand, "succeeded_at"),
		attributeString(paymentDemand, "created_at"),
		attributeString(paymentDemand, "updated_at"),
	); err != nil {
		return err
	}

	if err := LineItems(writer, paymentDemand); err != nil {
		return err
	}

	return Relationships(writer, paymentDemand, document, relationshipNames(includes, []string{
		"payer",
		"buyer",
		"receiver",
		"payment_method",
		"billing_address",
		"shipping_address",
		"merchant",
	}))
}

func PaymentSubscription(writer io.Writer, paymentSubscription jsonapi.Resource, document jsonapi.Document, includes []string) error {
	if _, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Description: %s\n"+
			"Slug: %s\n"+
			"Amount: %s\n"+
			"Discount: %s\n"+
			"Fee: %s\n"+
			"Status: %s\n"+
			"Billing Period: %s\n"+
			"Billing Cycle Anchor: %s\n"+
			"Proration Behavior: %s\n"+
			"Purchase Reference: %s\n"+
			"Purchase Kind: %s\n"+
			"Payer Timezone: %s\n"+
			"Idempotency Key: %s\n"+
			"Email Receipt: %s\n"+
			"CVC2 Check: %s\n"+
			"Address Line 1 Verification: %s\n"+
			"Postal Code Verification: %s\n"+
			"3DS Version: %s\n"+
			"3DS Status: %s\n"+
			"3DS Cryptogram: %s\n"+
			"ECI: %s\n"+
			"Directory Transaction EID: %s\n"+
			"ACS Transaction EID: %s\n"+
			"Canceled At: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		paymentSubscription.ID,
		paymentSubscription.Type,
		attributeString(paymentSubscription, "description"),
		attributeString(paymentSubscription, "slug"),
		moneyString(paymentSubscription, "amount_cents", "amount_currency"),
		moneyAttributeString(paymentSubscription, "discount_cents", "amount_currency"),
		moneyAttributeString(paymentSubscription, "fee_cents", "amount_currency"),
		attributeString(paymentSubscription, "status"),
		firstAttributeString(paymentSubscription, "billing_period", "billing_interval", "interval"),
		firstAttributeString(paymentSubscription, "billing_cycle_anchor_at", "billing_anchor_at"),
		attributeString(paymentSubscription, "proration_behavior"),
		attributeString(paymentSubscription, "purchase_reference"),
		attributeString(paymentSubscription, "purchase_kind"),
		attributeString(paymentSubscription, "payer_timezone"),
		attributeString(paymentSubscription, "idempotency_key"),
		attributeString(paymentSubscription, "email_receipt"),
		attributeString(paymentSubscription, "cvc2_check"),
		attributeString(paymentSubscription, "address_line1_verification"),
		attributeString(paymentSubscription, "postal_code_verification"),
		attributeString(paymentSubscription, "threeds_version"),
		attributeString(paymentSubscription, "threeds_status"),
		attributeString(paymentSubscription, "threeds_cryptogram"),
		attributeString(paymentSubscription, "eci"),
		attributeString(paymentSubscription, "directory_transaction_eid"),
		attributeString(paymentSubscription, "acs_transaction_eid"),
		attributeString(paymentSubscription, "canceled_at"),
		attributeString(paymentSubscription, "created_at"),
		attributeString(paymentSubscription, "updated_at"),
	); err != nil {
		return err
	}

	if err := LineItems(writer, paymentSubscription); err != nil {
		return err
	}

	return Relationships(writer, paymentSubscription, document, relationshipNames(includes, []string{
		"payer",
		"buyer",
		"receiver",
		"payment_method",
		"billing_address",
		"shipping_address",
		"merchant",
	}))
}

func PaymentMethod(writer io.Writer, paymentMethod jsonapi.Resource) error {
	expiry := expiryString(paymentMethod)
	_, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Kind: %s\n"+
			"Nickname: %s\n"+
			"Description: %s\n"+
			"BIN: %s\n"+
			"Last Four: %s\n"+
			"Card PAN Token: %s\n"+
			"Card CVV Token: %s\n"+
			"Expiry: %s\n"+
			"External State: %s\n"+
			"Discarded At: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		paymentMethod.ID,
		paymentMethod.Type,
		attributeString(paymentMethod, "kind"),
		attributeString(paymentMethod, "nickname"),
		attributeString(paymentMethod, "description"),
		attributeString(paymentMethod, "card_bin"),
		attributeString(paymentMethod, "last_four"),
		attributeString(paymentMethod, "card_pan_token"),
		attributeString(paymentMethod, "card_cvv_token"),
		expiry,
		attributeString(paymentMethod, "external_state"),
		attributeString(paymentMethod, "discarded_at"),
		attributeString(paymentMethod, "created_at"),
		attributeString(paymentMethod, "updated_at"),
	)
	return err
}

func RefundDemand(writer io.Writer, refundDemand jsonapi.Resource) error {
	_, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Amount: %s\n"+
			"State: %s\n"+
			"Reason: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		refundDemand.ID,
		refundDemand.Type,
		moneyString(refundDemand, "amount_cents", "amount_currency"),
		attributeString(refundDemand, "state"),
		attributeString(refundDemand, "reason"),
		attributeString(refundDemand, "created_at"),
		attributeString(refundDemand, "updated_at"),
	)
	return err
}

func UserResource(writer io.Writer, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
	if _, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Summary: %s\n"+
			"Status: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		resource.ID,
		resource.Type,
		userResourceSummary(resource),
		firstAttributeString(resource, "status", "account_status", "action"),
		attributeString(resource, "created_at"),
		attributeString(resource, "updated_at"),
	); err != nil {
		return err
	}

	return Relationships(writer, resource, document, relationshipNames(includes, defaultRelationshipNames(resource.Type)))
}

func attributeString(resource jsonapi.Resource, name string) string {
	rawValue, ok := resource.Attributes[name]
	if !ok {
		return ""
	}

	var stringValue string
	if err := json.Unmarshal(rawValue, &stringValue); err == nil {
		return stringValue
	}

	var anyValue any
	if err := json.Unmarshal(rawValue, &anyValue); err == nil && anyValue != nil {
		return fmt.Sprint(anyValue)
	}

	return ""
}

func moneyString(resource jsonapi.Resource, centsAttribute string, currencyAttribute string) string {
	return moneyAttributeString(resource, centsAttribute, currencyAttribute)
}

func moneyAttributeString(resource jsonapi.Resource, centsAttribute string, currencyAttribute string) string {
	currency := attributeString(resource, currencyAttribute)
	cents := attributeString(resource, centsAttribute)
	if cents == "" {
		return strings.TrimSpace(currency)
	}

	minorUnits, err := strconv.ParseInt(cents, 10, 64)
	if err != nil {
		return strings.TrimSpace(cents + " " + currency)
	}

	decimals := currencyDecimals(currency)
	if decimals == 0 {
		return strings.TrimSpace(fmt.Sprintf("%d %s", minorUnits, currency))
	}

	divisor := int64(1)
	for index := 0; index < decimals; index++ {
		divisor *= 10
	}

	sign := ""
	if minorUnits < 0 {
		sign = "-"
		minorUnits = -minorUnits
	}

	major := minorUnits / divisor
	minor := minorUnits % divisor
	amount := fmt.Sprintf("%s%d.%0*d", sign, major, decimals, minor)
	return strings.TrimSpace(amount + " " + currency)
}

func currencyDecimals(currency string) int {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "BIF", "CLP", "DJF", "GNF", "JPY", "KMF", "KRW", "MGA", "PYG", "RWF", "UGX", "VND", "VUV", "XAF", "XOF", "XPF":
		return 0
	default:
		return 2
	}
}

func firstAttributeString(resource jsonapi.Resource, names ...string) string {
	for _, name := range names {
		value := attributeString(resource, name)
		if value != "" {
			return value
		}
	}
	return ""
}

func expiryString(resource jsonapi.Resource) string {
	expiryMonth := attributeString(resource, "expiry_month")
	expiryYear := attributeString(resource, "expiry_year")
	if expiryMonth == "" && expiryYear == "" {
		return ""
	}
	if expiryMonth == "" {
		return expiryYear
	}
	if expiryYear == "" {
		return expiryMonth
	}
	return expiryMonth + "/" + expiryYear
}

func LineItems(writer io.Writer, resource jsonapi.Resource) error {
	rawLineItems, ok := resource.Attributes["line_items"]
	if !ok || len(rawLineItems) == 0 || string(rawLineItems) == "null" {
		return nil
	}

	var lineItems []map[string]json.RawMessage
	if err := json.Unmarshal(rawLineItems, &lineItems); err != nil {
		return err
	}
	if len(lineItems) == 0 {
		return nil
	}

	if _, err := fmt.Fprintln(writer, "\nLine Items:"); err != nil {
		return err
	}
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "NAME\tQTY\tAMOUNT\tTAX\tDISCOUNT\tSKU\tDESCRIPTION"); err != nil {
		return err
	}

	for _, lineItem := range lineItems {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			rawString(lineItem, "name"),
			rawString(lineItem, "quantity"),
			rawMoney(lineItem, "amount_cents", "amount_currency"),
			rawMoney(lineItem, "tax_cents", "tax_currency"),
			rawMoney(lineItem, "discount_cents", "discount_currency"),
			rawString(lineItem, "sku"),
			rawString(lineItem, "description"),
		); err != nil {
			return err
		}
	}

	return table.Flush()
}

func rawString(values map[string]json.RawMessage, name string) string {
	rawValue, ok := values[name]
	if !ok || len(rawValue) == 0 || string(rawValue) == "null" {
		return ""
	}

	var stringValue string
	if err := json.Unmarshal(rawValue, &stringValue); err == nil {
		return stringValue
	}

	var anyValue any
	if err := json.Unmarshal(rawValue, &anyValue); err == nil && anyValue != nil {
		return fmt.Sprint(anyValue)
	}

	return ""
}

func rawMoney(values map[string]json.RawMessage, centsAttribute string, currencyAttribute string) string {
	return moneyAttributeString(jsonapi.Resource{
		Attributes: map[string]json.RawMessage{
			centsAttribute:    values[centsAttribute],
			currencyAttribute: values[currencyAttribute],
		},
	}, centsAttribute, currencyAttribute)
}

func relationshipNames(includes []string, defaultNames []string) []string {
	if len(includes) == 0 {
		return defaultNames
	}
	return includes
}

func Relationships(writer io.Writer, resource jsonapi.Resource, document jsonapi.Document, relationshipNames []string) error {
	return IncludedResources(writer, resource, document, relationshipNames)
}

func RelationshipIdentifiers(writer io.Writer, resource jsonapi.Resource) error {
	if len(resource.Relationships) == 0 {
		return nil
	}

	relationshipNames := make([]string, 0, len(resource.Relationships))
	for relationshipName := range resource.Relationships {
		relationshipNames = append(relationshipNames, relationshipName)
	}
	sort.Strings(relationshipNames)

	rows := [][2]string{}
	for _, relationshipName := range relationshipNames {
		rawRelationship := resource.Relationships[relationshipName]
		relationship, err := jsonapi.DecodeRelationship(rawRelationship)
		if err != nil {
			return err
		}
		identifiers, err := jsonapi.DecodeResourceIdentifiers(relationship.Data)
		if err != nil {
			return err
		}
		if len(identifiers) == 0 {
			continue
		}

		for index, identifier := range identifiers {
			relationLabel := ""
			if index == 0 {
				relationLabel = titleize(relationshipName)
			}
			rows = append(rows, [2]string{relationLabel, identifier.Type + " " + identifier.ID})
		}
	}

	if len(rows) == 0 {
		return nil
	}

	if _, err := fmt.Fprintln(writer, "\nRelationships:"); err != nil {
		return err
	}
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "RELATION\tRESOURCE"); err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := fmt.Fprintf(table, "%s\t%s\n", row[0], row[1]); err != nil {
			return err
		}
	}
	return table.Flush()
}

func IncludedResources(writer io.Writer, resource jsonapi.Resource, document jsonapi.Document, relationshipNames []string) error {
	includedResources, err := jsonapi.DecodeIncluded(document.Included)
	if err != nil {
		return err
	}
	if len(includedResources) == 0 {
		return nil
	}

	includedIndex := map[string]jsonapi.Resource{}
	for _, includedResource := range includedResources {
		includedIndex[resourceKey(includedResource.Type, includedResource.ID)] = includedResource
	}

	printedHeader := false
	for _, relationshipName := range relationshipNames {
		rawRelationship, ok := resource.Relationships[relationshipName]
		if !ok {
			continue
		}

		relationship, err := jsonapi.DecodeRelationship(rawRelationship)
		if err != nil {
			return err
		}
		identifiers, err := jsonapi.DecodeResourceIdentifiers(relationship.Data)
		if err != nil {
			return err
		}
		if len(identifiers) == 0 {
			continue
		}

		for _, identifier := range identifiers {
			includedResource, ok := includedIndex[resourceKey(identifier.Type, identifier.ID)]
			if !ok {
				continue
			}
			if !printedHeader {
				if _, err := fmt.Fprintln(writer, "\nIncluded:"); err != nil {
					return err
				}
				printedHeader = true
			}
			if _, err := fmt.Fprintf(writer, "\n%s:\n", titleize(relationshipName)); err != nil {
				return err
			}
			if err := resourceDetails(writer, includedResource, jsonapi.Document{}, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceKey(resourceType string, resourceID string) string {
	return resourceType + ":" + resourceID
}

func summarizeResource(resource jsonapi.Resource) string {
	switch resource.Type {
	case "customers":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "name"),
			attributeString(resource, "email"),
			attributeString(resource, "phone_number"),
		})
	case "consumer_addresses":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "line_1"),
			attributeString(resource, "line_2"),
			attributeString(resource, "city"),
			attributeString(resource, "state"),
			attributeString(resource, "zip"),
			attributeString(resource, "country"),
		})
	case "payment_methods":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "kind"),
			"last four " + attributeString(resource, "last_four"),
			attributeString(resource, "external_state"),
		})
	case "merchants":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "business_name"),
			attributeString(resource, "business_email"),
		})
	case "payment_demands":
		return compactSummary(resource.ID, []string{
			moneyString(resource, "amount_cents", "amount_currency"),
			firstAttributeString(resource, "state", "status"),
			attributeString(resource, "processor_state"),
		})
	case "accounts":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "name"),
			attributeString(resource, "email"),
			attributeString(resource, "account_status"),
		})
	case "memberships":
		return compactSummary(resource.ID, []string{
			"opener " + attributeString(resource, "opener"),
		})
	case "permissions":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "slug"),
			attributeString(resource, "description"),
		})
	case "red_flags":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "violation"),
			attributeString(resource, "source"),
			attributeString(resource, "status"),
		})
	case "account_alerts":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "status"),
			"acknowledged " + attributeString(resource, "acknowledged_at"),
		})
	case "merchant_punitive_actions":
		return compactSummary(resource.ID, []string{
			attributeString(resource, "action"),
		})
	default:
		return resource.Type + " " + resource.ID
	}
}

func userResourceSummary(resource jsonapi.Resource) string {
	switch resource.Type {
	case "account_alerts":
		return compactSummary("", []string{
			attributeString(resource, "status"),
			"acknowledged " + attributeString(resource, "acknowledged_at"),
		})
	case "accounts":
		return compactSummary("", []string{
			attributeString(resource, "name"),
			attributeString(resource, "email"),
			attributeString(resource, "account_status"),
		})
	case "memberships":
		return compactSummary("", []string{
			"opener " + attributeString(resource, "opener"),
		})
	case "merchant_punitive_actions":
		return attributeString(resource, "action")
	case "permissions":
		return compactSummary("", []string{
			attributeString(resource, "slug"),
			attributeString(resource, "description"),
		})
	case "red_flags":
		return compactSummary("", []string{
			attributeString(resource, "violation"),
			attributeString(resource, "source"),
			attributeString(resource, "comment"),
		})
	default:
		return summarizeResource(resource)
	}
}

func defaultRelationshipNames(resourceType string) []string {
	switch resourceType {
	case "account_alerts":
		return []string{"merchant", "red_flag"}
	case "accounts":
		return []string{"memberships", "personal_identification"}
	case "memberships":
		return []string{"account", "merchant", "permissions"}
	case "merchant_punitive_actions":
		return []string{"merchant", "red_flag"}
	case "permissions":
		return []string{"merchant_tokens"}
	case "red_flags":
		return []string{"merchant"}
	default:
		return nil
	}
}

func compactSummary(resourceID string, parts []string) string {
	nonEmptyParts := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "last four" {
			continue
		}
		nonEmptyParts = append(nonEmptyParts, part)
	}
	if len(nonEmptyParts) == 0 {
		return resourceID
	}
	if resourceID == "" {
		return strings.Join(nonEmptyParts, ", ")
	}
	return resourceID + " (" + strings.Join(nonEmptyParts, ", ") + ")"
}

func titleize(value string) string {
	words := strings.Split(value, "_")
	for index, word := range words {
		if word == "" {
			continue
		}
		words[index] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}
