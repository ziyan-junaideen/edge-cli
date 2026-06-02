package cli

import (
	"context"

	"github.com/edgepayments/ept-cli/internal/edgeapi"
	"github.com/edgepayments/ept-cli/internal/jsonapi"
	"github.com/edgepayments/ept-cli/internal/output"
	"github.com/spf13/cobra"
)

type resourceCommandDefinition struct {
	Use              string
	Aliases          []string
	Short            string
	ListShort        string
	ShowShort        string
	ShowArgumentName string
	AllowedIncludes  map[string]struct{}
	List             func(context.Context, *edgeapi.Client, edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error)
	Show             func(context.Context, *edgeapi.Client, string, edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error)
	RenderCollection func(command *cobra.Command, resources []jsonapi.Resource) error
	RenderMember     func(command *cobra.Command, resource jsonapi.Resource, document jsonapi.Document, includes []string) error
}

func newReadOnlyResourceCommand(options *globalOptions, definition resourceCommandDefinition) *cobra.Command {
	resourceCommand := &cobra.Command{
		Use:     definition.Use,
		Aliases: definition.Aliases,
		Short:   definition.Short,
	}

	var listIncludeValues []string
	var listPreloadValues []string
	listCommand := &cobra.Command{
		Use:   "list",
		Short: definition.ListShort,
		RunE: func(command *cobra.Command, args []string) error {
			includes, err := parseIncludes(listIncludeValues, listPreloadValues, definition.AllowedIncludes)
			if err != nil {
				return err
			}

			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			resources, document, err := definition.List(context.Background(), client, edgeapi.QueryOptions{Include: includes})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return definition.RenderCollection(command, resources)
		},
	}

	var showIncludeValues []string
	var showPreloadValues []string
	showCommand := &cobra.Command{
		Use:   "show <" + definition.ShowArgumentName + ">",
		Short: definition.ShowShort,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			includes, err := parseIncludes(showIncludeValues, showPreloadValues, definition.AllowedIncludes)
			if err != nil {
				return err
			}

			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			resource, document, err := definition.Show(context.Background(), client, args[0], edgeapi.QueryOptions{Include: includes})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), document)
			}
			return definition.RenderMember(command, resource, document, includes)
		},
	}

	listCommand.Flags().StringArrayVar(&listIncludeValues, "include", nil, "JSON:API relationship to include; repeat or comma-separate values")
	listCommand.Flags().StringArrayVar(&listPreloadValues, "preload", nil, "alias for --include")
	showCommand.Flags().StringArrayVar(&showIncludeValues, "include", nil, "JSON:API relationship to include; repeat or comma-separate values")
	showCommand.Flags().StringArrayVar(&showPreloadValues, "preload", nil, "alias for --include")

	resourceCommand.AddCommand(listCommand, showCommand)
	return resourceCommand
}
