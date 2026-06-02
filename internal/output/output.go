package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/edgepayments/ept-cli/internal/jsonapi"
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
			"%s\t%s %s\t%s\t%s\t%s\t%s\n",
			paymentDemand.ID,
			attributeString(paymentDemand, "amount_cents"),
			attributeString(paymentDemand, "amount_currency"),
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
			"%s\t%s %s\t%s\t%s\t%s\n",
			paymentSubscription.ID,
			attributeString(paymentSubscription, "amount_cents"),
			attributeString(paymentSubscription, "amount_currency"),
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
			"%s\t%s %s\t%s\t%s\t%s\n",
			refundDemand.ID,
			attributeString(refundDemand, "amount_cents"),
			attributeString(refundDemand, "amount_currency"),
			attributeString(refundDemand, "state"),
			attributeString(refundDemand, "reason"),
			attributeString(refundDemand, "created_at"),
		); err != nil {
			return err
		}
	}
	return table.Flush()
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
			"Amount: %s %s\n"+
			"State: %s\n"+
			"Processor State: %s\n"+
			"Capture Method: %s\n"+
			"Email Receipt: %s\n"+
			"Succeeded At: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		paymentDemand.ID,
		paymentDemand.Type,
		attributeString(paymentDemand, "amount_cents"),
		attributeString(paymentDemand, "amount_currency"),
		firstAttributeString(paymentDemand, "state", "status"),
		attributeString(paymentDemand, "processor_state"),
		attributeString(paymentDemand, "capture_method"),
		attributeString(paymentDemand, "email_receipt"),
		attributeString(paymentDemand, "succeeded_at"),
		attributeString(paymentDemand, "created_at"),
		attributeString(paymentDemand, "updated_at"),
	); err != nil {
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
			"Amount: %s %s\n"+
			"Status: %s\n"+
			"Interval: %s\n"+
			"Billing Cycle Anchor: %s\n"+
			"Email Receipt: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		paymentSubscription.ID,
		paymentSubscription.Type,
		attributeString(paymentSubscription, "amount_cents"),
		attributeString(paymentSubscription, "amount_currency"),
		attributeString(paymentSubscription, "status"),
		firstAttributeString(paymentSubscription, "billing_interval", "interval"),
		firstAttributeString(paymentSubscription, "billing_cycle_anchor_at", "billing_anchor_at"),
		attributeString(paymentSubscription, "email_receipt"),
		attributeString(paymentSubscription, "created_at"),
		attributeString(paymentSubscription, "updated_at"),
	); err != nil {
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
	_, err := fmt.Fprintf(
		writer,
		"ID: %s\n"+
			"Type: %s\n"+
			"Kind: %s\n"+
			"Nickname: %s\n"+
			"Description: %s\n"+
			"BIN: %s\n"+
			"Last Four: %s\n"+
			"Expiry: %s/%s\n"+
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
		attributeString(paymentMethod, "expiry_month"),
		attributeString(paymentMethod, "expiry_year"),
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
			"Amount: %s %s\n"+
			"State: %s\n"+
			"Reason: %s\n"+
			"Created: %s\n"+
			"Updated: %s\n",
		refundDemand.ID,
		refundDemand.Type,
		attributeString(refundDemand, "amount_cents"),
		attributeString(refundDemand, "amount_currency"),
		attributeString(refundDemand, "state"),
		attributeString(refundDemand, "reason"),
		attributeString(refundDemand, "created_at"),
		attributeString(refundDemand, "updated_at"),
	)
	return err
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

func firstAttributeString(resource jsonapi.Resource, names ...string) string {
	for _, name := range names {
		value := attributeString(resource, name)
		if value != "" {
			return value
		}
	}
	return ""
}

func relationshipNames(includes []string, defaultNames []string) []string {
	if len(includes) == 0 {
		return defaultNames
	}
	return includes
}

func Relationships(writer io.Writer, resource jsonapi.Resource, document jsonapi.Document, relationshipNames []string) error {
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
			if _, err := fmt.Fprintf(writer, "%s: %s\n", titleize(relationshipName), summarizeResource(includedResource)); err != nil {
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
			attributeString(resource, "amount_cents") + " " + attributeString(resource, "amount_currency"),
			firstAttributeString(resource, "state", "status"),
			attributeString(resource, "processor_state"),
		})
	default:
		return resource.Type + " " + resource.ID
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
