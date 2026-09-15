package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ziyan-junaideen/edge-cli/internal/jsonapi"
)

func TestExtractReplayMaterialAndPrepareV3(t *testing.T) {
	delivery, document := replayDocument(t)

	material, err := ExtractReplayMaterial(delivery, document)
	if err != nil {
		t.Fatalf("ExtractReplayMaterial returned error: %v", err)
	}

	timestamp := time.Unix(1_755_230_400, 0)
	preparedReplay, err := PrepareReplay(material, DeliveryVersionV3, timestamp)
	if err != nil {
		t.Fatalf("PrepareReplay returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(preparedReplay.Body, &payload); err != nil {
		t.Fatalf("decode prepared body: %v", err)
	}
	event := payload["data"].(map[string]any)
	attributes := event["attributes"].(map[string]any)
	if event["id"] != "event-id" || attributes["slug"] != "succeeded" {
		t.Fatalf("unexpected webhook event: %#v", event)
	}
	if _, ok := attributes["created_at"]; ok {
		t.Fatal("webhook payload must not include the API-only created_at attribute")
	}

	messageAuthenticationCode := hmac.New(sha256.New, []byte("subscription-secret"))
	_, _ = messageAuthenticationCode.Write([]byte("1755230400." + string(preparedReplay.Body)))
	expectedHeader := "t=1755230400,v3=" + hex.EncodeToString(messageAuthenticationCode.Sum(nil))
	if preparedReplay.Headers.Get("Edge-Signature") != expectedHeader {
		t.Fatalf("expected %q, got %q", expectedHeader, preparedReplay.Headers.Get("Edge-Signature"))
	}
}

func TestPrepareLegacyPayloadShapes(t *testing.T) {
	delivery, document := replayDocument(t)
	material, err := ExtractReplayMaterial(delivery, document)
	if err != nil {
		t.Fatalf("ExtractReplayMaterial returned error: %v", err)
	}

	v1Replay, err := PrepareReplay(material, DeliveryVersionV1, time.Time{})
	if err != nil {
		t.Fatalf("prepare v1 replay: %v", err)
	}
	v2Replay, err := PrepareReplay(material, DeliveryVersionV2, time.Time{})
	if err != nil {
		t.Fatalf("prepare v2 replay: %v", err)
	}

	if strings.Contains(string(v1Replay.Body), `"data":{"id":"event-id"`) {
		t.Fatalf("v1 payload should be a bare event: %s", v1Replay.Body)
	}
	if !strings.Contains(string(v2Replay.Body), `"data":{"id":"event-id"`) {
		t.Fatalf("v2 payload should wrap the event in data: %s", v2Replay.Body)
	}
	if v1Replay.Headers.Get("X-Hub-Signature") == "" {
		t.Fatal("expected legacy signature header")
	}
	if v1Replay.Headers.Get("Edge-Signature") != "" {
		t.Fatal("legacy replay should not contain Edge-Signature")
	}
}

func TestPlayPostsExactPreparedBodyAndHeaders(t *testing.T) {
	preparedReplay := PreparedReplay{
		Body: []byte(`{"data":{"id":"event-id"}}`),
		Headers: http.Header{
			"Content-Type":   []string{"application/json"},
			"Edge-Signature": []string{"t=1,v3=abc"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		if request.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", request.Method)
		}
		if string(body) != string(preparedReplay.Body) {
			t.Errorf("expected body %q, got %q", preparedReplay.Body, body)
		}
		if request.Header.Get("Edge-Signature") != "t=1,v3=abc" {
			t.Errorf("signature header was not forwarded")
		}
		responseWriter.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	response, err := Play(context.Background(), server.Client(), server.URL, preparedReplay)
	if err != nil {
		t.Fatalf("Play returned error: %v", err)
	}
	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", response.StatusCode)
	}
}

func replayDocument(t *testing.T) (jsonapi.Resource, jsonapi.Document) {
	t.Helper()

	rawDocument := []byte(`{
  "data": {
    "id": "delivery-id",
    "type": "webhook_deliveries",
    "relationships": {
      "event": {"data": {"id": "event-id", "type": "events"}},
      "webhook_subscription": {"data": {"id": "subscription-id", "type": "webhook_subscriptions"}}
    }
  },
  "included": [
    {
      "id": "event-id",
      "type": "events",
      "attributes": {
        "mode": "live",
        "resource_type": "transaction.payment_demands",
        "resource_id": "payment-demand-id",
        "slug": "succeeded",
        "data": {"attributes": {"amount_cents": 1200}},
        "created_at": "2026-09-10T00:00:00Z"
      },
      "relationships": {
        "merchant": {"data": {"id": "merchant-id", "type": "merchants"}}
      }
    },
    {
      "id": "subscription-id",
      "type": "webhook_subscriptions",
      "attributes": {"secret_key": "subscription-secret"}
    }
  ]
}`)

	var document jsonapi.Document
	if err := json.Unmarshal(rawDocument, &document); err != nil {
		t.Fatalf("decode fixture document: %v", err)
	}
	delivery, err := jsonapi.DecodeResource(document.Data)
	if err != nil {
		t.Fatalf("decode fixture delivery: %v", err)
	}
	return delivery, document
}
