package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ziyan-junaideen/edge-cli/internal/edgeapi"
	"github.com/ziyan-junaideen/edge-cli/internal/jsonapi"
	"github.com/ziyan-junaideen/edge-cli/internal/output"
	"github.com/ziyan-junaideen/edge-cli/internal/webhooks"
)

func newWebhookDeliveriesCommand(options *globalOptions) *cobra.Command {
	webhookDeliveriesCommand := newReadOnlyResourceCommand(options, resourceCommandDefinition{
		Use:              "webhook-deliveries",
		Aliases:          []string{"webhook_deliveries", "webhooks"},
		Short:            "Inspect and replay webhook deliveries",
		ListShort:        "List webhook deliveries",
		ShowShort:        "Show a webhook delivery",
		ShowArgumentName: "webhook_delivery_id",
		AllowedIncludes:  includeSet("event", "merchant", "webhook_subscription"),
		List: func(ctx context.Context, client *edgeapi.Client, queryOptions edgeapi.QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
			return client.ListWebhookDeliveries(ctx, queryOptions)
		},
		Show: func(ctx context.Context, client *edgeapi.Client, id string, queryOptions edgeapi.QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
			return client.ShowWebhookDelivery(ctx, id, queryOptions)
		},
		RenderCollection: func(command *cobra.Command, resources []jsonapi.Resource) error {
			return output.WebhookDeliveryCollection(command.OutOrStdout(), resources)
		},
		RenderMember: func(command *cobra.Command, resource jsonapi.Resource, document jsonapi.Document, includes []string) error {
			return output.ShowResource(command.OutOrStdout(), resource, document, includes)
		},
	})

	webhookDeliveriesCommand.AddCommand(newReplayWebhookDeliveryCommand(options))
	return webhookDeliveriesCommand
}

func newReplayWebhookDeliveryCommand(options *globalOptions) *cobra.Command {
	var targetURL string
	var deliveryVersionValue string
	var timeout time.Duration
	var dryRun bool

	command := &cobra.Command{
		Use:     "replay <webhook_delivery_id>",
		Aliases: []string{"play"},
		Short:   "Replay a webhook delivery to another endpoint",
		Args:    cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			if err := webhooks.ValidateTargetURL(targetURL); err != nil {
				return err
			}
			if timeout <= 0 {
				return fmt.Errorf("timeout must be greater than zero")
			}

			deliveryVersion, err := webhooks.ParseDeliveryVersion(deliveryVersionValue)
			if err != nil {
				return err
			}

			client, _, err := newAPIClient(options)
			if err != nil {
				return err
			}

			delivery, document, err := client.ShowWebhookDelivery(
				command.Context(),
				args[0],
				edgeapi.QueryOptions{Include: []string{"event", "webhook_subscription"}},
			)
			if err != nil {
				return err
			}

			material, err := webhooks.ExtractReplayMaterial(delivery, document)
			if err != nil {
				return fmt.Errorf("prepare webhook delivery %s: %w", args[0], err)
			}
			preparedReplay, err := webhooks.PrepareReplay(material, deliveryVersion, time.Now())
			if err != nil {
				return err
			}

			if dryRun {
				return renderPreparedReplay(command, targetURL, deliveryVersion, preparedReplay, options.jsonOutput)
			}

			httpClient := &http.Client{
				Timeout: timeout,
				CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			response, err := webhooks.Play(command.Context(), httpClient, targetURL, preparedReplay)
			if err != nil {
				return err
			}

			if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
				responseBody := strings.TrimSpace(string(response.Body))
				if responseBody == "" {
					return fmt.Errorf("webhook endpoint returned %s", response.Status)
				}
				return fmt.Errorf("webhook endpoint returned %s: %s", response.Status, responseBody)
			}

			if options.jsonOutput {
				return output.JSON(command.OutOrStdout(), map[string]any{
					"delivery_id": args[0],
					"target_url":  targetURL,
					"version":     deliveryVersion,
					"status":      response.Status,
					"response":    strings.TrimSpace(string(response.Body)),
				})
			}

			_, err = fmt.Fprintf(command.OutOrStdout(), "Replayed webhook delivery %s to %s: %s\n", args[0], targetURL, response.Status)
			if err != nil {
				return err
			}
			if responseBody := strings.TrimSpace(string(response.Body)); responseBody != "" {
				_, err = fmt.Fprintln(command.OutOrStdout(), responseBody)
			}
			return err
		},
	}

	command.Flags().StringVar(&targetURL, "to", "", "HTTP endpoint that should receive the replayed webhook")
	command.Flags().StringVar(&deliveryVersionValue, "delivery-version", "v3", "webhook payload and signature version: v1, v2, or v3")
	command.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "maximum time to wait for the webhook endpoint")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "print the signed request without sending it")
	_ = command.MarkFlagRequired("to")

	return command
}

func renderPreparedReplay(command *cobra.Command, targetURL string, deliveryVersion webhooks.DeliveryVersion, preparedReplay webhooks.PreparedReplay, jsonOutput bool) error {
	if jsonOutput {
		return output.JSON(command.OutOrStdout(), map[string]any{
			"target_url": targetURL,
			"version":    deliveryVersion,
			"headers":    preparedReplay.Headers,
			"body":       json.RawMessage(preparedReplay.Body),
		})
	}

	if _, err := fmt.Fprintf(command.OutOrStdout(), "POST %s\n", targetURL); err != nil {
		return err
	}
	for headerName, headerValues := range preparedReplay.Headers {
		for _, headerValue := range headerValues {
			if _, err := fmt.Fprintf(command.OutOrStdout(), "%s: %s\n", headerName, headerValue); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintf(command.OutOrStdout(), "\n%s\n", preparedReplay.Body)
	return err
}
