package cli

import (
	"context"

	"github.com/edgepayments/ept-cli/internal/edgeapi"
	"github.com/edgepayments/ept-cli/internal/output"
	"github.com/spf13/cobra"
)

func newMerchantsCommand(options *globalOptions) *cobra.Command {
	merchantsCommand := &cobra.Command{
		Use:   "merchants",
		Short: "Inspect merchant accounts",
	}

	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List merchant accounts accessible by the token",
		RunE: func(command *cobra.Command, args []string) error {
			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			merchants, document, err := client.ListMerchants(context.Background(), edgeapi.QueryOptions{})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return output.MerchantCollection(command.OutOrStdout(), merchants)
		},
	}

	showCommand := &cobra.Command{
		Use:   "show <merchant_id>",
		Short: "Show a merchant account",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			merchant, document, err := client.ShowMerchant(context.Background(), args[0], edgeapi.QueryOptions{})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return output.Merchant(command.OutOrStdout(), merchant)
		},
	}

	merchantsCommand.AddCommand(listCommand, showCommand)
	return merchantsCommand
}
