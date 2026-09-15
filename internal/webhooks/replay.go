package webhooks

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ziyan-junaideen/edge-cli/internal/jsonapi"
)

const maxResponseBodyBytes = 1024 * 1024

type DeliveryVersion string

const (
	DeliveryVersionV1 DeliveryVersion = "v1"
	DeliveryVersionV2 DeliveryVersion = "v2"
	DeliveryVersionV3 DeliveryVersion = "v3"
)

type ReplayMaterial struct {
	Event     EventResource
	SecretKey string
}

type EventResource struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Attributes    EventAttributes    `json:"attributes"`
	Relationships EventRelationships `json:"relationships"`
}

type EventAttributes struct {
	Mode         string          `json:"mode"`
	ResourceType string          `json:"resource_type"`
	ResourceID   string          `json:"resource_id"`
	Slug         string          `json:"slug"`
	Data         json.RawMessage `json:"data"`
}

type EventRelationships struct {
	Merchant ResourceRelationship `json:"merchant"`
}

type ResourceRelationship struct {
	Data jsonapi.ResourceIdentifier `json:"data"`
}

type PreparedReplay struct {
	Body    []byte
	Headers http.Header
}

type ReplayResponse struct {
	StatusCode int
	Status     string
	Body       []byte
}

func ParseDeliveryVersion(value string) (DeliveryVersion, error) {
	version := DeliveryVersion(strings.ToLower(strings.TrimSpace(value)))
	switch version {
	case DeliveryVersionV1, DeliveryVersionV2, DeliveryVersionV3:
		return version, nil
	default:
		return "", fmt.Errorf("unsupported webhook delivery version %q; expected v1, v2, or v3", value)
	}
}

func ValidateTargetURL(targetURL string) error {
	parsedTargetURL, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return fmt.Errorf("parse target URL: %w", err)
	}
	if parsedTargetURL.Scheme != "http" && parsedTargetURL.Scheme != "https" {
		return fmt.Errorf("target URL must use http or https")
	}
	if parsedTargetURL.Host == "" {
		return fmt.Errorf("target URL must include a host")
	}
	return nil
}

func ExtractReplayMaterial(delivery jsonapi.Resource, document jsonapi.Document) (ReplayMaterial, error) {
	eventIdentifier, err := jsonapi.RelationshipIdentifier(delivery, "event")
	if err != nil {
		return ReplayMaterial{}, err
	}
	event, err := jsonapi.FindIncluded(document, eventIdentifier)
	if err != nil {
		return ReplayMaterial{}, err
	}

	subscriptionIdentifier, err := jsonapi.RelationshipIdentifier(delivery, "webhook_subscription")
	if err != nil {
		return ReplayMaterial{}, err
	}
	subscription, err := jsonapi.FindIncluded(document, subscriptionIdentifier)
	if err != nil {
		return ReplayMaterial{}, err
	}

	merchantIdentifier, err := jsonapi.RelationshipIdentifier(event, "merchant")
	if err != nil {
		return ReplayMaterial{}, err
	}

	mode, err := stringAttribute(event, "mode")
	if err != nil {
		return ReplayMaterial{}, err
	}
	resourceType, err := stringAttribute(event, "resource_type")
	if err != nil {
		return ReplayMaterial{}, err
	}
	resourceID, err := stringAttribute(event, "resource_id")
	if err != nil {
		return ReplayMaterial{}, err
	}
	slug, err := stringAttribute(event, "slug")
	if err != nil {
		return ReplayMaterial{}, err
	}
	secretKey, err := stringAttribute(subscription, "secret_key")
	if err != nil {
		return ReplayMaterial{}, err
	}

	eventData := event.Attributes["data"]
	if len(eventData) == 0 {
		eventData = json.RawMessage("null")
	}

	return ReplayMaterial{
		SecretKey: secretKey,
		Event: EventResource{
			ID:   event.ID,
			Type: "events",
			Attributes: EventAttributes{
				Mode:         mode,
				ResourceType: resourceType,
				ResourceID:   resourceID,
				Slug:         slug,
				Data:         eventData,
			},
			Relationships: EventRelationships{
				Merchant: ResourceRelationship{Data: merchantIdentifier},
			},
		},
	}, nil
}

func PrepareReplay(material ReplayMaterial, version DeliveryVersion, timestamp time.Time) (PreparedReplay, error) {
	var payload any = material.Event
	if version == DeliveryVersionV2 || version == DeliveryVersionV3 {
		payload = struct {
			Data EventResource `json:"data"`
		}{Data: material.Event}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return PreparedReplay{}, fmt.Errorf("encode webhook payload: %w", err)
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "Edge Technologies, Inc. Webhook Client/1.0; edge-cli")

	switch version {
	case DeliveryVersionV1, DeliveryVersionV2:
		digest := sha1.Sum([]byte(material.SecretKey))
		headers.Set("X-Hub-Signature", base64.StdEncoding.EncodeToString(digest[:]))
	case DeliveryVersionV3:
		unixTimestamp := timestamp.Unix()
		signedPayload := fmt.Sprintf("%d.%s", unixTimestamp, body)
		messageAuthenticationCode := hmac.New(sha256.New, []byte(material.SecretKey))
		_, _ = messageAuthenticationCode.Write([]byte(signedPayload))
		headers.Set(
			"Edge-Signature",
			fmt.Sprintf("t=%d,v3=%s", unixTimestamp, hex.EncodeToString(messageAuthenticationCode.Sum(nil))),
		)
	default:
		return PreparedReplay{}, fmt.Errorf("unsupported webhook delivery version %q", version)
	}

	return PreparedReplay{Body: body, Headers: headers}, nil
}

func Play(ctx context.Context, httpClient *http.Client, targetURL string, preparedReplay PreparedReplay) (ReplayResponse, error) {
	if err := ValidateTargetURL(targetURL); err != nil {
		return ReplayResponse{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(preparedReplay.Body))
	if err != nil {
		return ReplayResponse{}, fmt.Errorf("create webhook request: %w", err)
	}
	request.Header = preparedReplay.Headers.Clone()

	response, err := httpClient.Do(request)
	if err != nil {
		return ReplayResponse{}, fmt.Errorf("deliver webhook: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBodyBytes+1))
	if err != nil {
		return ReplayResponse{}, fmt.Errorf("read webhook response: %w", err)
	}
	if len(body) > maxResponseBodyBytes {
		return ReplayResponse{}, fmt.Errorf("webhook response exceeds %d bytes", maxResponseBodyBytes)
	}

	return ReplayResponse{
		StatusCode: response.StatusCode,
		Status:     response.Status,
		Body:       body,
	}, nil
}

func stringAttribute(resource jsonapi.Resource, attributeName string) (string, error) {
	rawValue, ok := resource.Attributes[attributeName]
	if !ok {
		return "", fmt.Errorf("attribute %q is missing from %s %s", attributeName, resource.Type, resource.ID)
	}

	var value string
	if err := json.Unmarshal(rawValue, &value); err != nil {
		return "", fmt.Errorf("decode attribute %q from %s %s: %w", attributeName, resource.Type, resource.ID, err)
	}
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("attribute %q is empty on %s %s", attributeName, resource.Type, resource.ID)
	}
	return value, nil
}
