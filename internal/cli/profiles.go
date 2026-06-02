package cli

import (
	"fmt"
	"sort"

	"github.com/edgepayments/ept-cli/internal/config"
	"github.com/edgepayments/ept-cli/internal/output"
	"github.com/spf13/cobra"
)

func newProfilesCommand(options *globalOptions) *cobra.Command {
	profilesCommand := &cobra.Command{
		Use:   "profiles",
		Short: "Manage API profiles",
	}

	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List configured profiles",
		RunE: func(command *cobra.Command, args []string) error {
			_, configFile, err := config.Load(config.Overrides{})
			if err != nil {
				return err
			}
			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), configFile)
			}
			profileNames := make([]string, 0, len(configFile.Profiles))
			for name := range configFile.Profiles {
				profileNames = append(profileNames, name)
			}
			sort.Strings(profileNames)
			for _, name := range profileNames {
				profile := configFile.Profiles[name]
				activeMarker := " "
				if name == configFile.ActiveProfile {
					activeMarker = "*"
				}
				fmt.Fprintf(command.OutOrStdout(), "%s %s\t%s\n", activeMarker, name, profile.APIURL)
			}
			return nil
		},
	}

	useCommand := &cobra.Command{
		Use:   "use <name>",
		Short: "Set the active profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			_, configFile, err := config.Load(config.Overrides{})
			if err != nil {
				return err
			}
			if _, ok := configFile.Profiles[args[0]]; !ok {
				return fmt.Errorf("profile %q is not configured", args[0])
			}
			configFile.ActiveProfile = args[0]
			if err := config.Save(configFile); err != nil {
				return err
			}
			fmt.Fprintf(command.OutOrStdout(), "Active profile is now %q\n", args[0])
			return nil
		},
	}

	var apiURL string
	var caCert string
	setCommand := &cobra.Command{
		Use:   "set <name>",
		Short: "Create or update a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			_, configFile, err := config.Load(config.Overrides{})
			if err != nil {
				return err
			}
			if configFile.Profiles == nil {
				configFile.Profiles = map[string]config.Profile{}
			}

			profile := configFile.Profiles[args[0]]
			if apiURL != "" {
				normalizedURL, err := config.NormalizeAPIURL(apiURL)
				if err != nil {
					return err
				}
				profile.APIURL = normalizedURL
			}
			if caCert != "" {
				profile.CACert = caCert
			}
			if options.insecureSet {
				profile.InsecureSkipVerify = options.insecureSkipVerify
			}
			if profile.APIURL == "" {
				return fmt.Errorf("profile %q requires --api-url", args[0])
			}

			configFile.Profiles[args[0]] = profile
			if err := config.Save(configFile); err != nil {
				return err
			}
			fmt.Fprintf(command.OutOrStdout(), "Saved profile %q\n", args[0])
			return nil
		},
	}
	setCommand.Flags().StringVar(&apiURL, "api-url", "", "Edge API base URL")
	setCommand.Flags().StringVar(&caCert, "ca-cert", "", "CA certificate path")

	profilesCommand.AddCommand(listCommand, useCommand, setCommand)
	return profilesCommand
}
