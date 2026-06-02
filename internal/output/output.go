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

func Merchant(writer io.Writer, merchant jsonapi.Resource) error {
	_, err := fmt.Fprintf(
		writer,
		"ID: %s\nType: %s\nBusiness: %s\nEmail: %s\nWebsite: %s\nEntity: %s\nCreated: %s\nUpdated: %s\n",
		merchant.ID,
		merchant.Type,
		attributeString(merchant, "business_name"),
		attributeString(merchant, "business_email"),
		attributeString(merchant, "business_website"),
		attributeString(merchant, "entity_type"),
		attributeString(merchant, "created_at"),
		attributeString(merchant, "updated_at"),
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
