package cli

import (
	"context"

	"github.com/edgepayments/ept-cli/internal/edgeapi"
	"github.com/edgepayments/ept-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCustomersCommand(options *globalOptions) *cobra.Command {
	customersCommand := &cobra.Command{
		Use:   "customers",
		Short: "Inspect customers",
	}

	allowedIncludes := includeSet("merchant", "addresses", "payment_demands")

	var listIncludeValues []string
	var listPreloadValues []string
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List customers",
		RunE: func(command *cobra.Command, args []string) error {
			includes, err := parseIncludes(listIncludeValues, listPreloadValues, allowedIncludes)
			if err != nil {
				return err
			}

			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			customers, document, err := client.ListCustomers(context.Background(), edgeapi.QueryOptions{Include: includes})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return output.CustomerCollection(command.OutOrStdout(), customers)
		},
	}

	var showIncludeValues []string
	var showPreloadValues []string
	showCommand := &cobra.Command{
		Use:   "show <customer_id>",
		Short: "Show a customer",
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

			customer, document, err := client.ShowCustomer(context.Background(), args[0], edgeapi.QueryOptions{Include: includes})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return output.ShowResource(command.OutOrStdout(), customer, document, includes)
		},
	}

	listCommand.Flags().StringArrayVar(&listIncludeValues, "include", nil, "JSON:API relationship to include; repeat or comma-separate values")
	listCommand.Flags().StringArrayVar(&listPreloadValues, "preload", nil, "alias for --include")
	showCommand.Flags().StringArrayVar(&showIncludeValues, "include", nil, "JSON:API relationship to include; repeat or comma-separate values")
	showCommand.Flags().StringArrayVar(&showPreloadValues, "preload", nil, "alias for --include")

	customersCommand.AddCommand(listCommand, showCommand)
	return customersCommand
}
