package output

import (
	"encoding/json"
	"fmt"
	"io"
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
