package cli

import (
	"context"

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

			merchants, rawData, err := client.ListMerchants(context.Background())
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), rawData)
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

			merchant, rawData, err := client.ShowMerchant(context.Background(), args[0])
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), rawData)
			}
			return output.Merchant(command.OutOrStdout(), merchant)
		},
	}

	merchantsCommand.AddCommand(listCommand, showCommand)
	return merchantsCommand
}
