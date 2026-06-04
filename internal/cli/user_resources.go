package cli

import (
	"context"

	"github.com/edgepayments/ept-cli/internal/edgeapi"
	"github.com/edgepayments/ept-cli/internal/jsonapi"
	"github.com/edgepayments/ept-cli/internal/output"
	"github.com/spf13/cobra"
)

func newAccountAlertsCommand(options *globalOptions) *cobra.Command {
	return newUserResourceCommand(options, resourceCommandDefinition{
		Use:              "account-alerts",
		Aliases:          []string{"account_alerts", "alerts"},
		Short:            "Inspect account alerts",
		ListShort:        "List account alerts",
		ShowShort:        "Show an account alert",
		ShowArgumentName: "account_alert_id",
		AllowedIncludes:  includeSet("merchant", "red_flag"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListAccountAlerts(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowAccountAlert(ctx, id, options)
		},
	})
}

func newAccountsCommand(options *globalOptions) *cobra.Command {
	return newUserResourceCommand(options, resourceCommandDefinition{
		Use:              "accounts",
		Short:            "Inspect accounts",
		ListShort:        "List accounts",
		ShowShort:        "Show an account",
		ShowArgumentName: "account_id",
		AllowedIncludes:  includeSet("memberships", "personal_identification"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListAccounts(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowAccount(ctx, id, options)
		},
	})
}

func newMembershipsCommand(options *globalOptions) *cobra.Command {
	return newUserResourceCommand(options, resourceCommandDefinition{
		Use:              "memberships",
		Short:            "Inspect memberships",
		ListShort:        "List memberships",
		ShowShort:        "Show a membership",
		ShowArgumentName: "membership_id",
		AllowedIncludes:  includeSet("account", "merchant", "permissions"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListMemberships(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowMembership(ctx, id, options)
		},
	})
}

func newMerchantPunitiveActionsCommand(options *globalOptions) *cobra.Command {
	return newUserResourceCommand(options, resourceCommandDefinition{
		Use:              "merchant-punitive-actions",
		Aliases:          []string{"merchant_punitive_actions", "punitive-actions"},
		Short:            "Inspect merchant punitive actions",
		ListShort:        "List merchant punitive actions",
		ShowShort:        "Show a merchant punitive action",
		ShowArgumentName: "merchant_punitive_action_id",
		AllowedIncludes:  includeSet("merchant", "red_flag"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListMerchantPunitiveActions(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowMerchantPunitiveAction(ctx, id, options)
		},
	})
}

func newPermissionsCommand(options *globalOptions) *cobra.Command {
	return newUserResourceCommand(options, resourceCommandDefinition{
		Use:              "permissions",
		Short:            "Inspect permissions",
		ListShort:        "List permissions",
		ShowShort:        "Show a permission",
		ShowArgumentName: "permission_id",
		AllowedIncludes:  includeSet("merchant_tokens"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListPermissions(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowPermission(ctx, id, options)
		},
	})
}

func newRedFlagsCommand(options *globalOptions) *cobra.Command {
	return newUserResourceCommand(options, resourceCommandDefinition{
		Use:              "red-flags",
		Aliases:          []string{"red_flags"},
		Short:            "Inspect red flags",
		ListShort:        "List red flags",
		ShowShort:        "Show a red flag",
		ShowArgumentName: "red_flag_id",
		AllowedIncludes:  includeSet("merchant"),
		List: func(ctx context.Context, client *edgeapi.Client, options edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListRedFlags(ctx, options)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, options edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowRedFlag(ctx, id, options)
		},
	})
}

func newUserResourceCommand(options *globalOptions, definition resourceCommandDefinition) *cobra.Command {
	definition.RenderCollection = func(command *cobra.Command, resources []jsonapi.Resource) error {
		return output.UserResourceCollection(command.OutOrStdout(), resources)
	}
	definition.RenderMember = func(command *cobra.Command, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
		return output.ShowResource(command.OutOrStdout(), resource, document, includes)
	}
	return newReadOnlyResourceCommand(options, definition)
}
