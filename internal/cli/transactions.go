package cli

import (
	"context"

	"github.com/edgepayments/ept-cli/internal/edgeapi"
	"github.com/edgepayments/ept-cli/internal/jsonapi"
	"github.com/edgepayments/ept-cli/internal/output"
	"github.com/spf13/cobra"
)

func newPaymentDemandsCommand(options *globalOptions) *cobra.Command {
	return newReadOnlyResourceCommand(options, resourceCommandDefinition{
		Use:              "payment-demands",
		Aliases:          []string{"payment_demands", "demands"},
		Short:            "Inspect payment demands",
		ListShort:        "List payment demands",
		ShowShort:        "Show a payment demand",
		ShowArgumentName: "payment_demand_id",
		AllowedIncludes: includeSet(
			"merchant",
			"buyer",
			"payer",
			"receiver",
			"payment_method",
			"billing_address",
			"shipping_address",
		),
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
