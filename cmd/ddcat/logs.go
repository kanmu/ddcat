package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type Row struct {
	Timestamp  string         `json:"timestamp"`
	Status     string         `json:"status,omitempty"`
	Service    string         `json:"service,omitempty"`
	Host       string         `json:"host,omitempty"`
	Message    string         `json:"message,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
}

func buildRow(log *datadogV2.LogAttributes, withAttrs bool, withTags bool) *Row {
	row := &Row{
		Timestamp: log.GetTimestamp().Local().Format(time.RFC3339),
		Status:    log.GetStatus(),
		Service:   log.GetService(),
		Host:      log.GetHost(),
		Message:   log.GetMessage(),
	}

	if withAttrs {
		if attrs, ok := log.GetAttributesOk(); ok {
			row.Attributes = *attrs
		}
	}

	if withTags {
		if tags, ok := log.GetTagsOk(); ok {
			row.Tags = *tags
		}
	}

	return row
}

func print(out io.Writer, withAttrs bool, withTags bool) func(resp datadogV2.LogsListResponse) {
	return func(resp datadogV2.LogsListResponse) {
		for _, data := range resp.Data {
			attrs, ok := data.GetAttributesOk()

			if !ok {
				continue
			}

			row := buildRow(attrs, withAttrs, withTags)
			line, err := json.Marshal(row)

			if err != nil {
				panic(err)
			}

			fmt.Fprintln(out, string(line))
		}
	}
}
