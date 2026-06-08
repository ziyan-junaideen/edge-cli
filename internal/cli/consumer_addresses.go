package cli

import (
	"context"

	"github.com/ziyan-junaideen/edge-cli/internal/edgeapi"
	"github.com/ziyan-junaideen/edge-cli/internal/output"
	"github.com/spf13/cobra"
)

func newConsumerAddressesCommand(options *globalOptions) *cobra.Command {
	consumerAddressesCommand := &cobra.Command{
		Use:     "consumer-addresses",
		Aliases: []string{"addresses", "consumer_addresses"},
		Short:   "Inspect consumer addresses",
	}

	allowedIncludes := includeSet("merchant", "customer")

	var listIncludeValues []string
	var listPreloadValues []string
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List consumer addresses",
		RunE: func(command *cobra.Command, args []string) error {
			includes, err := parseIncludes(listIncludeValues, listPreloadValues, allowedIncludes)
			if err != nil {
				return err
			}

			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			consumerAddresses, document, err := client.ListConsumerAddresses(context.Background(), edgeapi.QueryOptions{Include: includes})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return output.ConsumerAddressCollection(command.OutOrStdout(), consumerAddresses)
		},
	}

	var showIncludeValues []string
	var showPreloadValues []string
	showCommand := &cobra.Command{
		Use:   "show <consumer_address_id>",
		Short: "Show a consumer address",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			includes, err := parseIncludes(showIncludeValues, showPreloadValues, allowedIncludes)
			if err != nil {
				return err
			}

			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			consumerAddress, document, err := client.ShowConsumerAddress(context.Background(), args[0], edgeapi.QueryOptions{Include: includes})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return output.ShowResource(command.OutOrStdout(), consumerAddress, document, includes)
		},
	}

	listCommand.Flags().StringArrayVar(&listIncludeValues, "include", nil, "JSON:API relationship to include; repeat or comma-separate values")
	listCommand.Flags().StringArrayVar(&listPreloadValues, "preload", nil, "alias for --include")
	showCommand.Flags().StringArrayVar(&showIncludeValues, "include", nil, "JSON:API relationship to include; repeat or comma-separate values")
	showCommand.Flags().StringArrayVar(&showPreloadValues, "preload", nil, "alias for --include")

	consumerAddressesCommand.AddCommand(listCommand, showCommand)
	return consumerAddressesCommand
}
