package main

import (
	"log"
	"os"
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
	Query     string   `arg:"" name:"query" required:"" help:"Search query. see https://docs.datadoghq.com/logs/explorer/search_syntax/"`
	Indexes   []string `help:"Indexes to search."`
	From      string   `help:"The minimum time for the requested logs."`
	To        string   `help:"The maximum time for the requested logs."`
	Sort      string   `enum:"timestamp,-timestamp" default:"timestamp" help:"Sort parameters when querying logs ('timestamp', '-timestamp')."`
	WithAttrs bool     `help:"Include attributes in displayed logs."`
	WithTags  bool     `help:"Include tags in displayed logs	"`
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

	err = client.ListLogs(*req, print(os.Stdout, cli.WithAttrs, cli.WithTags))

	if err != nil {
		log.Fatal(err)
	}
}
