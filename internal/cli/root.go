package cli

import (
	"fmt"
	"os"

	"github.com/edgepayments/ept-cli/internal/config"
	"github.com/spf13/cobra"
)

type globalOptions struct {
	profileName        string
	apiURL             string
	caCert             string
	insecureSkipVerify bool
	insecureSet        bool
	jsonOutput         bool
}

func NewRootCommand() *cobra.Command {
	options := &globalOptions{}

	rootCommand := &cobra.Command{
		Use:           "edge",
		Short:         "Edge Payment Technologies command line client",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCommand.PersistentFlags().StringVar(&options.profileName, "profile", "", "profile name")
	rootCommand.PersistentFlags().StringVar(&options.apiURL, "api-url", "", "Edge API base URL")
	rootCommand.PersistentFlags().StringVar(&options.caCert, "ca-cert", "", "CA certificate path for local development")
	rootCommand.PersistentFlags().BoolVar(&options.insecureSkipVerify, "insecure-skip-verify", false, "skip TLS verification for non-production development endpoints")
	rootCommand.PersistentFlags().BoolVar(&options.jsonOutput, "json", false, "print JSON output")

	rootCommand.PersistentPreRunE = func(command *cobra.Command, args []string) error {
		options.insecureSet = command.Flags().Changed("insecure-skip-verify") || command.InheritedFlags().Changed("insecure-skip-verify")
		return nil
	}

	rootCommand.AddCommand(newAuthCommand(options))
	rootCommand.AddCommand(newConsumerAddressesCommand(options))
	rootCommand.AddCommand(newCustomersCommand(options))
	rootCommand.AddCommand(newMerchantsCommand(options))
	rootCommand.AddCommand(newPaymentDemandsCommand(options))
	rootCommand.AddCommand(newPaymentMethodsCommand(options))
	rootCommand.AddCommand(newPaymentSubscriptionsCommand(options))
	rootCommand.AddCommand(newProfilesCommand(options))
	rootCommand.AddCommand(newRefundDemandsCommand(options))

	rootCommand.SetOut(os.Stdout)
	rootCommand.SetErr(os.Stderr)
	return rootCommand
}

func loadRuntime(options *globalOptions) (config.Runtime, error) {
	overrides := config.Overrides{
		ProfileName: options.profileName,
		APIURL:      options.apiURL,
		CACert:      options.caCert,
	}
	if options.insecureSet {
		overrides.InsecureSkipVerify = &options.insecureSkipVerify
	}
	runtime, _, err := config.Load(overrides)
	if err != nil {
		return config.Runtime{}, err
	}
	if runtime.InsecureSkipVerify {
		fmt.Fprintf(os.Stderr, "warning: TLS verification is disabled for profile %q\n", runtime.ProfileName)
	}
	return runtime, nil
}
