# ddcat

[![test](https://github.com/winebarrel/ddcat/actions/workflows/build.yml/badge.svg)](https://github.com/winebarrel/ddcat/actions/workflows/build.yml)

CLI to display Datadog Logs.

## Usage

```
Usage: ddcat --api-key=STRING --app-key=STRING <query>

Arguments:
  <query>    Search query. see https://docs.datadoghq.com/logs/explorer/search_syntax/

Flags:
  -h, --help                   Show context-sensitive help.
      --api-key=STRING         Datadog API key ($DD_API_KEY).
      --app-key=STRING         Datadog APP key ($DD_APP_KEY).
      --indexes=INDEXES,...    Indexes to search.
      --from=STRING            The minimum time for the requested logs.
      --to=STRING              The maximum time for the requested logs.
      --sort="timestamp"       Sort parameters when querying logs ('timestamp', '-timestamp').
      --with-attrs             Include attributes in displayed logs.
      --with-tags              Include tags in displayed logs
      --version
```

```
$ ddcat --from 'now-1m' 'service:web OK'
{"timestamp":"2023-07-26T20:10:30+09:00","status":"info","service":"web","message":"GET /user/info 200 OK"}
{"timestamp":"2023-07-26T20:10:31+09:00","status":"info","service":"web","message":"POST /entry 200 OK"}
{"timestamp":"2023-07-26T20:10:32+09:00","status":"info","service":"web","message":"GET /entry 200 OK"}
...
```

## Related Links

* https://docs.datadoghq.com/logs/explorer/search_syntax/
* https://docs.datadoghq.com/logs/guide/access-your-log-data-programmatically/
* https://docs.datadoghq.com/api/latest/logs/#search-logs

## Installation

```
brew install winebarrel/ddcat/ddcat
```
