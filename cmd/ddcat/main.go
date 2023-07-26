package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/alecthomas/kong"
	"github.com/araddon/dateparse"
	"github.com/winebarrel/ddcat"
)

func init() {
	log.SetFlags(0)
}

var version string

var cli struct {
	APIKey    string   `env:"DD_API_KEY" required:"" help:"Datadog API key."`
	APPKey    string   `env:"DD_APP_KEY" required:"" help:"Datadog APP key."`
	Query     string   `arg:"" name:"query" required:""`
	Indexes   []string `help:""`
	From      string   `help:""`
	To        string   `help:""`
	Sort      string   `enum:"timestamp,-timestamp" default:"timestamp"`
	WithAttrs bool     `help:""`
	WithTags  bool     `help:""`
	Version   kong.VersionFlag
}

func init() {
	log.SetFlags(0)
}

func buildReq() (*datadogV2.LogsListRequest, error) {
	req := datadogV2.NewLogsListRequest()

	filter := datadogV2.NewLogsQueryFilter()
	filter.SetQuery(cli.Query)

	if len(cli.Indexes) >= 1 {
		filter.SetIndexes(cli.Indexes)
	}

	if cli.From != "" {
		from := cli.From

		if t, err := dateparse.ParseLocal(from); err == nil {
			from = t.Format(time.RFC3339)
		}

		filter.SetFrom(from)
	}

	if cli.To != "" {
		to := cli.To

		if t, err := dateparse.ParseLocal(to); err == nil {
			to = t.Format(time.RFC3339)
		}

		filter.SetTo(to)
	}

	req.SetFilter(*filter)

	options := datadogV2.NewLogsQueryOptions()
	zone, err := time.LoadLocation("Local")

	if err != nil {
		return nil, err
	}

	options.SetTimezone(zone.String())

	req.SetSort(datadogV2.LogsSort(cli.Sort))
	return req, nil
}

type Row struct {
	Timestamp  string         `json:"timestamp"`
	Status     string         `json:"status,omitempty"`
	Service    string         `json:"service,omitempty"`
	Host       string         `json:"host,omitempty"`
	Message    string         `json:"message,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
}

func buildRow(log *datadogV2.LogAttributes) *Row {
	row := &Row{
		Timestamp: log.GetTimestamp().Local().Format(time.RFC3339),
		Status:    log.GetStatus(),
		Service:   log.GetService(),
		Host:      log.GetHost(),
		Message:   log.GetMessage(),
	}

	if cli.WithAttrs {
		if attrs, ok := log.GetAttributesOk(); ok {
			row.Attributes = *attrs
		}
	}

	if cli.WithTags {
		if tags, ok := log.GetTagsOk(); ok {
			row.Tags = *tags
		}
	}

	return row
}

func main() {
	kong.Parse(
		&cli,
		kong.Vars{"version": version},
	)

	client := ddcat.NewClient(cli.APIKey, cli.APPKey)
	req, err := buildReq()

	if err != nil {
		log.Fatal(err)
	}

	err = client.ListLogs(*req, func(resp datadogV2.LogsListResponse) {
		for _, data := range resp.Data {
			attrs, ok := data.GetAttributesOk()

			if !ok {
				continue
			}

			row := buildRow(attrs)
			line, err := json.Marshal(row)

			if err != nil {
				panic(err)
			}

			fmt.Println(string(line))

			// var msg string

			// if attrs.Message != nil {
			// 	msg = *attrs.Message
			// }

			// fmt.Printf("%v\t%s\t%s\n", attrs.Timestamp.Local(), *attrs.Status, msg)
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}
