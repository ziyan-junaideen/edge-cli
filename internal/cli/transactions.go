package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ziyan-junaideen/edge-cli/internal/edgeapi"
	"github.com/ziyan-junaideen/edge-cli/internal/jsonapi"
	"github.com/ziyan-junaideen/edge-cli/internal/output"
)

func newPaymentDemandsCommand(options *globalOptions) *cobra.Command {
	allowedIncludes := includeSet(
		"merchant",
		"buyer",
		"payer",
		"receiver",
		"payment_method",
		"billing_address",
		"shipping_address",
	)
	command := newReadOnlyResourceCommand(options, resourceCommandDefinition{
		Use:              "payment-demands",
		Aliases:          []string{"payment_demands", "demands"},
		Short:            "Inspect payment demands",
		ListShort:        "List payment demands",
		ShowShort:        "Show a payment demand",
		ShowArgumentName: "payment_demand_id",
		AllowedIncludes:  allowedIncludes,
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListPaymentDemands(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowPaymentDemand(ctx, id, options)
		},
		RenderCollection: func(command *cobra.Command, resources []jsonapi.Resource) error {
			return output.PaymentDemandCollection(command.OutOrStdout(), resources)
		},
		RenderMember: func(command *cobra.Command, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
			return output.ShowResource(command.OutOrStdout(), resource, document, includes)
		},
	})
	command.AddCommand(newCreatePaymentDemandCommand(options, allowedIncludes))
	return command
}

type createPaymentDemandOptions struct {
	amountCents             int
	amountCurrency          string
	idempotencyKey          string
	confirmed               bool
	description             string
	emailReceipt            bool
	captureMethod           string
	discountCents           int
	purchaseReference       string
	purchaseKind            string
	payerTimezone           string
	taxCents                int
	taxCurrency             string
	shippingCents           int
	shippingCurrency        string
	lineItems               []string
	threedsVersion          string
	threedsStatus           string
	threedsCryptogram       string
	eci                     string
	directoryTransactionEID string
	acsTransactionEID       string
	paymentMethodID         string
	billingAddressID        string
	shippingAddressID       string
	payerID                 string
	buyerID                 string
	receiverID              string
	includeValues           []string
	preloadValues           []string
}

func newCreatePaymentDemandCommand(options *globalOptions, allowedIncludes map[string]struct{}) *cobra.Command {
	createOptions := &createPaymentDemandOptions{}
	command := &cobra.Command{
		Use:   "create",
		Short: "Create a payment intent or confirmed payment demand",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, args []string) error {
			includes, err := parseIncludes(createOptions.includeValues, createOptions.preloadValues, allowedIncludes)
			if err != nil {
				return err
			}
			if err := validateCreatePaymentDemandOptions(createOptions); err != nil {
				return err
			}

			lineItems, err := parsePaymentLineItems(createOptions.lineItems)
			if err != nil {
				return err
			}

			attributes := edgeapi.CreatePaymentDemandAttributes{
				AmountCents:             createOptions.amountCents,
				AmountCurrency:          createOptions.amountCurrency,
				IdempotencyKey:          createOptions.idempotencyKey,
				Confirmed:               createOptions.confirmed,
				Description:             createOptions.description,
				CaptureMethod:           createOptions.captureMethod,
				PurchaseReference:       createOptions.purchaseReference,
				PurchaseKind:            createOptions.purchaseKind,
				PayerTimezone:           createOptions.payerTimezone,
				LineItems:               lineItems,
				ThreedsVersion:          createOptions.threedsVersion,
				ThreedsStatus:           createOptions.threedsStatus,
				ThreedsCryptogram:       createOptions.threedsCryptogram,
				ECI:                     createOptions.eci,
				DirectoryTransactionEID: createOptions.directoryTransactionEID,
				ACSTransactionEID:       createOptions.acsTransactionEID,
			}
			if command.Flags().Changed("email-receipt") {
				attributes.EmailReceipt = &createOptions.emailReceipt
			}
			if command.Flags().Changed("discount-cents") {
				attributes.DiscountCents = &createOptions.discountCents
			}
			if command.Flags().Changed("tax-cents") {
				attributes.TaxDetail = &edgeapi.PaymentTaxDetail{
					TaxCents: createOptions.taxCents, TaxCurrency: createOptions.taxCurrency,
				}
			}
			if command.Flags().Changed("shipping-cents") {
				attributes.ShippingDetail = &edgeapi.PaymentShippingDetail{
					ShippingCents: createOptions.shippingCents, ShippingCurrency: createOptions.shippingCurrency,
				}
			}

			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}
			resource, document, err := client.CreatePaymentDemand(context.Background(), edgeapi.CreatePaymentDemandParams{
				Attributes:        attributes,
				PaymentMethodID:   createOptions.paymentMethodID,
				BillingAddressID:  createOptions.billingAddressID,
				ShippingAddressID: createOptions.shippingAddressID,
				PayerID:           createOptions.payerID,
				BuyerID:           createOptions.buyerID,
				ReceiverID:        createOptions.receiverID,
			}, edgeapi.QueryOptions{Include: includes})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return output.ShowResource(command.OutOrStdout(), resource, document, includes)
		},
	}

	flags := command.Flags()
	flags.IntVar(&createOptions.amountCents, "amount-cents", 0, "payment amount in cents")
	flags.StringVar(&createOptions.amountCurrency, "amount-currency", "USD", "ISO 4217 payment currency")
	flags.StringVar(&createOptions.idempotencyKey, "idempotency-key", "", "unique key used to identify a repeated payment request")
	flags.BoolVar(&createOptions.confirmed, "confirmed", false, "create and queue a payment demand instead of an unconfirmed intent")
	flags.StringVar(&createOptions.description, "description", "", "payment description")
	flags.BoolVar(&createOptions.emailReceipt, "email-receipt", true, "email a receipt for the payment")
	flags.StringVar(&createOptions.captureMethod, "capture-method", "automatic", "capture method (automatic or manual)")
	flags.IntVar(&createOptions.discountCents, "discount-cents", 0, "total discount in cents")
	flags.StringVar(&createOptions.purchaseReference, "purchase-reference", "", "merchant purchase reference")
	flags.StringVar(&createOptions.purchaseKind, "purchase-kind", "", "purchase kind (order or invoice)")
	flags.StringVar(&createOptions.payerTimezone, "payer-timezone", "", "IANA timezone of the payer")
	flags.IntVar(&createOptions.taxCents, "tax-cents", 0, "total tax in cents")
	flags.StringVar(&createOptions.taxCurrency, "tax-currency", "USD", "ISO 4217 tax currency")
	flags.IntVar(&createOptions.shippingCents, "shipping-cents", 0, "total shipping charge in cents")
	flags.StringVar(&createOptions.shippingCurrency, "shipping-currency", "USD", "ISO 4217 shipping currency")
	flags.StringArrayVar(&createOptions.lineItems, "line-item", nil, "line item as JSON; repeat for multiple items")
	flags.StringVar(&createOptions.threedsVersion, "threeds-version", "", "3DS version returned by authentication")
	flags.StringVar(&createOptions.threedsStatus, "threeds-status", "", "3DS status returned by authentication")
	flags.StringVar(&createOptions.threedsCryptogram, "threeds-cryptogram", "", "3DS cryptogram returned by authentication")
	flags.StringVar(&createOptions.eci, "eci", "", "3DS electronic commerce indicator")
	flags.StringVar(&createOptions.directoryTransactionEID, "directory-transaction-eid", "", "3DS directory transaction ID")
	flags.StringVar(&createOptions.acsTransactionEID, "acs-transaction-eid", "", "3DS ACS transaction ID")
	flags.StringVar(&createOptions.paymentMethodID, "payment-method", "", "payment method ID")
	flags.StringVar(&createOptions.billingAddressID, "billing-address", "", "billing address ID")
	flags.StringVar(&createOptions.shippingAddressID, "shipping-address", "", "shipping address ID")
	flags.StringVar(&createOptions.payerID, "payer", "", "payer customer ID")
	flags.StringVar(&createOptions.buyerID, "buyer", "", "buyer customer ID")
	flags.StringVar(&createOptions.receiverID, "receiver", "", "receiver customer ID")
	flags.StringArrayVar(&createOptions.includeValues, "include", nil, "JSON:API relationship to include; repeat or comma-separate values")
	flags.StringArrayVar(&createOptions.preloadValues, "preload", nil, "alias for --include")

	_ = command.MarkFlagRequired("amount-cents")
	_ = command.MarkFlagRequired("idempotency-key")
	return command
}

func validateCreatePaymentDemandOptions(options *createPaymentDemandOptions) error {
	if strings.ToUpper(options.amountCurrency) != "USD" {
		return fmt.Errorf("--amount-currency must be USD")
	}
	if strings.ToUpper(options.taxCurrency) != "USD" {
		return fmt.Errorf("--tax-currency must be USD")
	}
	if strings.ToUpper(options.shippingCurrency) != "USD" {
		return fmt.Errorf("--shipping-currency must be USD")
	}
	if options.captureMethod != "automatic" && options.captureMethod != "manual" {
		return fmt.Errorf("--capture-method must be automatic or manual")
	}
	if options.purchaseKind != "" && options.purchaseKind != "order" && options.purchaseKind != "invoice" {
		return fmt.Errorf("--purchase-kind must be order or invoice")
	}
	options.amountCurrency = strings.ToUpper(options.amountCurrency)
	options.taxCurrency = strings.ToUpper(options.taxCurrency)
	options.shippingCurrency = strings.ToUpper(options.shippingCurrency)
	return nil
}

func parsePaymentLineItems(values []string) ([]edgeapi.PaymentLineItem, error) {
	lineItems := make([]edgeapi.PaymentLineItem, 0, len(values))
	for index, value := range values {
		decoder := json.NewDecoder(strings.NewReader(value))
		decoder.DisallowUnknownFields()
		var lineItem edgeapi.PaymentLineItem
		if err := decoder.Decode(&lineItem); err != nil {
			return nil, fmt.Errorf("decode --line-item %d: %w", index+1, err)
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			return nil, fmt.Errorf("decode --line-item %d: unexpected trailing JSON", index+1)
		}
		if lineItem.AmountCents == nil || lineItem.AmountCurrency == "" || lineItem.Quantity == nil {
			return nil, fmt.Errorf("--line-item %d requires amount_cents, amount_currency, and quantity", index+1)
		}
		if strings.ToUpper(lineItem.AmountCurrency) != "USD" ||
			(lineItem.TaxCurrency != "" && strings.ToUpper(lineItem.TaxCurrency) != "USD") ||
			(lineItem.DiscountCurrency != "" && strings.ToUpper(lineItem.DiscountCurrency) != "USD") {
			return nil, fmt.Errorf("--line-item %d currencies must be USD", index+1)
		}
		if *lineItem.Quantity < 1 {
			return nil, fmt.Errorf("--line-item %d quantity must be at least 1", index+1)
		}
		lineItem.AmountCurrency = strings.ToUpper(lineItem.AmountCurrency)
		lineItem.TaxCurrency = strings.ToUpper(lineItem.TaxCurrency)
		lineItem.DiscountCurrency = strings.ToUpper(lineItem.DiscountCurrency)
		lineItems = append(lineItems, lineItem)
	}
	return lineItems, nil
}

func newPaymentSubscriptionsCommand(options *globalOptions) *cobra.Command {
	return newReadOnlyResourceCommand(options, resourceCommandDefinition{
		Use:              "payment-subscriptions",
		Aliases:          []string{"payment_subscriptions", "subscriptions"},
		Short:            "Inspect payment subscriptions",
		ListShort:        "List payment subscriptions",
		ShowShort:        "Show a payment subscription",
		ShowArgumentName: "payment_subscription_id",
		AllowedIncludes: includeSet(
			"merchant",
			"billing_address",
			"shipping_address",
			"buyer",
			"receiver",
			"payer",
			"payment_method",
		),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListPaymentSubscriptions(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowPaymentSubscription(ctx, id, options)
		},
		RenderCollection: func(command *cobra.Command, resources []jsonapi.Resource) error {
			return output.PaymentSubscriptionCollection(command.OutOrStdout(), resources)
		},
		RenderMember: func(command *cobra.Command, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
			return output.ShowResource(command.OutOrStdout(), resource, document, includes)
		},
	})
}

func newPaymentMethodsCommand(options *globalOptions) *cobra.Command {
	return newReadOnlyResourceCommand(options, resourceCommandDefinition{
		Use:              "payment-methods",
		Aliases:          []string{"payment_methods", "methods"},
		Short:            "Inspect payment methods",
		ListShort:        "List payment methods",
		ShowShort:        "Show a payment method",
		ShowArgumentName: "payment_method_id",
		AllowedIncludes:  includeSet("address", "customer", "payment_demands", "merchant"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListPaymentMethods(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowPaymentMethod(ctx, id, options)
		},
		RenderCollection: func(command *cobra.Command, resources []jsonapi.Resource) error {
			return output.PaymentMethodCollection(command.OutOrStdout(), resources)
		},
		RenderMember: func(command *cobra.Command, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
			return output.ShowResource(command.OutOrStdout(), resource, document, includes)
		},
	})
}

func newRefundDemandsCommand(options *globalOptions) *cobra.Command {
	return newReadOnlyResourceCommand(options, resourceCommandDefinition{
		Use:              "refund-demands",
		Aliases:          []string{"refund_demands", "refunds"},
		Short:            "Inspect refund demands",
		ListShort:        "List refund demands",
		ShowShort:        "Show a refund demand",
		ShowArgumentName: "refund_demand_id",
		AllowedIncludes:  includeSet("merchant", "payment_demand"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListRefundDemands(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowRefundDemand(ctx, id, options)
		},
		RenderCollection: func(command *cobra.Command, resources []jsonapi.Resource) error {
			return output.RefundDemandCollection(command.OutOrStdout(), resources)
		},
		RenderMember: func(command *cobra.Command, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
			return output.ShowResource(command.OutOrStdout(), resource, document, includes)
		},
	})
}
