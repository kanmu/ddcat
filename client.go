package ddcat

import (
	"context"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type Client struct {
	apiKey string
	appKey string
	api    *datadogV2.LogsApi
}

func NewClient(apiKey string, appKey string) *Client {
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewLogsApi(apiClient)

	client := &Client{
		apiKey: apiKey,
		appKey: appKey,
		api:    api,
	}

	return client
}

func (client *Client) withAPIKey(ctx context.Context) context.Context {
	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: client.apiKey,
			},
			"appKeyAuth": {
				Key: client.appKey,
			},
		},
	)

	return ctx
}

func (client *Client) ListLogs(req datadogV2.LogsListRequest, callback func(datadogV2.LogsListResponse)) error {
	ctx := client.withAPIKey(context.Background())

	if req.Page == nil {
		req.SetPage(*datadogV2.NewLogsListRequestPage())
	}

	req.Page.SetLimit(5000)

	for {
		resp, _, err := client.api.ListLogs(ctx, *datadogV2.NewListLogsOptionalParameters().WithBody(req))

		if err != nil {
			return err
		}

		callback(resp)

		if resp.Meta.Page == nil || resp.Meta.Page.After == nil {
			break
		}

		req.Page.SetCursor(*resp.Meta.Page.After)
	}

	return nil
}
