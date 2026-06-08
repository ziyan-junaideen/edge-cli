package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ziyan-junaideen/edge-cli/internal/secrets"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newAuthCommand(options *globalOptions) *cobra.Command {
	authCommand := &cobra.Command{
		Use:   "auth",
		Short: "Manage API authentication",
	}

	var token string
	loginCommand := &cobra.Command{
		Use:   "login",
		Short: "Store an API token for the active profile",
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := loadRuntime(options)
			if err != nil {
				return err
			}

			apiToken := strings.TrimSpace(token)
			if apiToken == "" {
				fmt.Fprint(command.ErrOrStderr(), "Token: ")
				tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Fprintln(command.ErrOrStderr())
				if err != nil {
					return err
				}
				apiToken = strings.TrimSpace(string(tokenBytes))
			}
			if apiToken == "" {
				return errors.New("token is required")
			}

			secretStore, err := secrets.Open()
			if err != nil {
				return err
			}
			if err := secretStore.SetToken(runtime.ProfileName, apiToken); err != nil {
				return err
			}

			fmt.Fprintf(command.OutOrStdout(), "Stored token for profile %q\n", runtime.ProfileName)
			return nil
		},
	}
	loginCommand.Flags().StringVar(&token, "token", "", "API token")

	statusCommand := &cobra.Command{
		Use:   "status",
		Short: "Show whether the active profile has a token",
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := loadRuntime(options)
			if err != nil {
				return err
			}

			secretStore, err := secrets.Open()
			if err != nil {
				return err
			}
			_, err = secretStore.Token(runtime.ProfileName)
			if errors.Is(err, secrets.ErrNotFound) {
				fmt.Fprintf(command.OutOrStdout(), "No token configured for profile %q\n", runtime.ProfileName)
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Fprintf(command.OutOrStdout(), "Token configured for profile %q\n", runtime.ProfileName)
			return nil
		},
	}

	logoutCommand := &cobra.Command{
		Use:   "logout",
		Short: "Remove the API token for the active profile",
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := loadRuntime(options)
			if err != nil {
				return err
			}

			secretStore, err := secrets.Open()
			if err != nil {
				return err
			}
			if err := secretStore.DeleteToken(runtime.ProfileName); err != nil {
				return err
			}

			fmt.Fprintf(command.OutOrStdout(), "Removed token for profile %q\n", runtime.ProfileName)
			return nil
		},
	}

	authCommand.AddCommand(loginCommand, statusCommand, logoutCommand)
	return authCommand
}
