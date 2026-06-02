package edgeapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/edgepayments/ept-cli/internal/jsonapi"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	token      string
}

type Config struct {
	APIURL             string
	Token              string
	CACert             string
	InsecureSkipVerify bool
}

type APIError struct {
	StatusCode int
	Errors     []jsonapi.Error
	Body       string
}

type QueryOptions struct {
	Include []string
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.Token) == "" {
		return nil, errors.New("api token is required")
	}

	baseURL, err := url.Parse(config.APIURL)
	if err != nil {
		return nil, fmt.Errorf("parse api URL: %w", err)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	tlsConfig, err := tlsClientConfig(config)
	if err != nil {
		return nil, err
	}
	transport.TLSClientConfig = tlsConfig

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		token: config.Token,
	}, nil
}

func (client *Client) ListMerchants(ctx context.Context, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	return client.listResources(ctx, "merchants", options)
}

func (client *Client) ShowMerchant(ctx context.Context, merchantID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	return client.showResource(ctx, "merchants", merchantID, options)
}

func (client *Client) ListCustomers(ctx context.Context, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	return client.listResources(ctx, "customers", options)
}

func (client *Client) ShowCustomer(ctx context.Context, customerID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	return client.showResource(ctx, "customers", customerID, options)
}

func (client *Client) ListConsumerAddresses(ctx context.Context, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	return client.listResources(ctx, "consumer_addresses", options)
}

func (client *Client) ShowConsumerAddress(ctx context.Context, consumerAddressID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	return client.showResource(ctx, "consumer_addresses", consumerAddressID, options)
}

func (client *Client) ListPaymentDemands(ctx context.Context, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	return client.listResources(ctx, "payment_demands", options)
}

func (client *Client) ShowPaymentDemand(ctx context.Context, paymentDemandID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	return client.showResource(ctx, "payment_demands", paymentDemandID, options)
}

func (client *Client) ListPaymentSubscriptions(ctx context.Context, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	return client.listResources(ctx, "payment_subscriptions", options)
}

func (client *Client) ShowPaymentSubscription(ctx context.Context, paymentSubscriptionID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	return client.showResource(ctx, "payment_subscriptions", paymentSubscriptionID, options)
}

func (client *Client) ListPaymentMethods(ctx context.Context, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	return client.listResources(ctx, "payment_methods", options)
}

func (client *Client) ShowPaymentMethod(ctx context.Context, paymentMethodID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	return client.showResource(ctx, "payment_methods", paymentMethodID, options)
}

func (client *Client) ListRefundDemands(ctx context.Context, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	return client.listResources(ctx, "refund_demands", options)
}

func (client *Client) ShowRefundDemand(ctx context.Context, refundDemandID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	return client.showResource(ctx, "refund_demands", refundDemandID, options)
}

func (client *Client) listResources(ctx context.Context, path string, options QueryOptions) ([]jsonapi.Resource, jsonapi.Document, error) {
	document, err := client.get(ctx, path, options)
	if err != nil {
		return nil, jsonapi.Document{}, err
	}

	resources, err := jsonapi.DecodeResourceCollection(document.Data)
	return resources, document, err
}

func (client *Client) showResource(ctx context.Context, path string, resourceID string, options QueryOptions) (jsonapi.Resource, jsonapi.Document, error) {
	document, err := client.get(ctx, path+"/"+url.PathEscape(resourceID), options)
	if err != nil {
		return jsonapi.Resource{}, jsonapi.Document{}, err
	}

	resource, err := jsonapi.DecodeResource(document.Data)
	return resource, document, err
}

func (client *Client) get(ctx context.Context, path string, options QueryOptions) (jsonapi.Document, error) {
	requestURL := *client.baseURL
	requestURL.Path = strings.TrimRight(client.baseURL.Path, "/") + "/" + strings.TrimLeft(path, "/")
	query := requestURL.Query()
	if len(options.Include) > 0 {
		query.Set("include", strings.Join(options.Include, ","))
	}
	requestURL.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return jsonapi.Document{}, err
	}
	request.Header.Set("Authorization", "Bearer "+client.token)
	request.Header.Set("Accept", "application/vnd.api+json")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return jsonapi.Document{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return jsonapi.Document{}, err
	}

	var document jsonapi.Document
	if len(body) > 0 {
		if err := json.Unmarshal(body, &document); err != nil && response.StatusCode < 400 {
			return jsonapi.Document{}, fmt.Errorf("decode response: %w", err)
		}
	}

	if response.StatusCode >= 400 {
		return jsonapi.Document{}, APIError{StatusCode: response.StatusCode, Errors: document.Errors, Body: string(body)}
	}

	return document, nil
}

func (apiError APIError) Error() string {
	if len(apiError.Errors) == 0 {
		return fmt.Sprintf("api request failed with status %d", apiError.StatusCode)
	}

	firstError := apiError.Errors[0]
	message := firstError.Detail
	if message == "" {
		message = firstError.Title
	}
	if message == "" {
		message = firstError.Code
	}
	if message == "" {
		message = fmt.Sprintf("api request failed with status %d", apiError.StatusCode)
	}
	return message
}

func tlsClientConfig(config Config) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: config.InsecureSkipVerify,
	}

	if strings.TrimSpace(config.CACert) == "" {
		return tlsConfig, nil
	}

	certificateBytes, err := os.ReadFile(config.CACert)
	if err != nil {
		return nil, fmt.Errorf("read CA certificate: %w", err)
	}

	certificatePool, err := x509.SystemCertPool()
	if err != nil {
		certificatePool = x509.NewCertPool()
	}
	if !certificatePool.AppendCertsFromPEM(certificateBytes) {
		return nil, errors.New("CA certificate did not contain a valid PEM certificate")
	}

	tlsConfig.RootCAs = certificatePool
	return tlsConfig, nil
}
