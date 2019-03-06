# Circonus Agent

The circonus-agent is intended to be a drop-in replacement for NAD. There is, however, one specific caveat -- native plugins (.js) do not work. Unless modified to run `node` independently and follow [plugin output guidelines](#output)


# Features

1. Replacement for NAD, written in go
1. Builtin metric [collectors](#builtin-collectors) -- the default Linux builtins emit the common metrics needed for cosi visuals (graphs, worksheets, & dashboards)
1. [Plugin](#plugins) architecture for local metric collection
1. Local HTTP [Receiver](#receiver) for POST/PUT metric collection
1. Local [StatsD](#statsd) listener for application metrics
1. Prometheus format support
    1. Receive HTTP `PUT|POST` to `/prom` endpoint (e.g. `PUT http://127.0.0.1:2609/prom`)
    1. Fetch (see [Prometheus collector](https://github.com/circonus-labs/circonus-agent/blob/master/etc/README.md#prometheus-collector) for details)
    1. Extract HTTP `GET` of `/prom` endpoint will emit metrics in Prometheus format (e.g. `GET http://127.0.0.1:2609/prom`)


# Install via COSI

```
curl -sSL https://setup.circonus.com/install | bash \
    -s -- \
    --cosiurl https://setup.circonus.com/ \
    --key <insert api key> \
    --app <insert api app>
```

Features of the COSI installed circonus-agent on Linux systems:

  * includes (if OS supports) [protocol_observer](https://github.com/circonus-labs/wirelatency), no longer needs to be built/installed manually
  * includes (if OS supports) [circonus-logwatch](https://github.com/circonus-labs/circonus-logwatch), no longer needs to be installed manually
  * includes OS/version/architecture-specific NAD plugins (non-javascript only) -- **Note:** the circonus-agent is **not** capable of using NAD _native plugins_ since they require NodeJS

Supported Operating Systems (x86_64 and/or amd64):

  * RHEL7 (CentOS, RedHat, Oracle)
  * RHEL6 (CentOS, RedHat, amzn)
  * Ubuntu18
  * Ubuntu16
  * Ubuntu14
  * Debian9
  * Debian8
  * FreeBSD 12
  * FreeBSD 11

Please continue to use the original cosi(w/NAD) for OmniOS and Raspian - cosi v2 support for these is TBD. Note: after installing NAD a binary circonus-agent can be used as a drop-in replacement (configure circonus-agent _plugins directory_ to be NAD plugins directory -- javascript plugins will not function). Binaries for OmniOS (`solaris_x86_64`) and Raspian (`linux_arm`) are available in the [circonus-agent repository](https://github.com/circonus-labs/circonus-agent/releases/latest).


# Quick Start

**Manually** installing on a system where the original [cosi](https://github.com/circonus-labs/circonus-one-step-install) installed and configured NAD.

1. `mkdir -p /opt/circonus/agent`
1. Download [latest release](../../releases/latest) from repository
1. Extract archive into `/opt/circonus/agent`
1. If planning to use `--check-enable-new-metrics`, ensure the `state` directory is owned by the user `circonus-agentd` will run as -- note, this is no longer required if a **new** check is created by cosi or the circonus-agent
1. Create a [config](https://github.com/circonus-labs/circonus-agent/blob/master/etc/README.md#main-configuration) (see minimal example below) or use command line parameters
1. Copy, edit, and install one of the service configurations in `service/`
1. Stop NAD (e.g. `systemctl stop nad` or `/etc/init.d/nad stop`)
1. Start the circonus-agent (e.g. `systemctl start circonus-agent` or `/etc/init.d/circonus-agent start`)
1. Disable NAD service so that it will not start at next reboot

Example, minimal, configuration using existing cosi install, configuration would be placed into `/opt/circonus/agent/etc/circonus-agent.toml`:

```toml
# enable debug for more verbose messages
#debug = true

# set the plugin directory to NAD's plugins
plugin-dir = "/opt/circonus/nad/etc/node-agent.d"

[api]
key = "cosi" # use cosi api configuration

[check]
bundle_id = "cosi" # use cosi system check bundle

[reverse]
enabled = true
```


# Options

```
$ /opt/circonus/agent/sbin/circonus-agentd -h
```

```
Flags:
      --api-app string                    [ENV: CA_API_APP] Circonus API Token app (default "circonus-agent")
      --api-ca-file string                [ENV: CA_API_CA_FILE] Circonus API CA certificate file
      --api-key string                    [ENV: CA_API_KEY] Circonus API Token key
      --api-url string                    [ENV: CA_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --check-broker string               [ENV: CA_CHECK_BROKER] ID of Broker to use or 'select' for random selection of valid broker, if creating a check bundle (default "select")
  -C, --check-create                      [ENV: CA_CHECK_CREATE] Create check bundle (for reverse and auto enable new metrics)
  -I, --check-id string                   [ENV: CA_CHECK_ID] Check Bundle ID or 'cosi' for cosi system check (for reverse)
      --check-metric-filters string       [ENV: CA_CHECK_METRIC_FILTERS] List of filters used to manage which metrics are collected
      --check-metric-refresh-ttl string   [ENV: CA_CHECK_METRIC_REFRESH_TTL] Refresh check metrics TTL (default "5m")
  -S, --check-metric-streamtags           [ENV: CA_CHECK_METRIC_STREAMTAGS] Add check tags to metrics as stream tags
      --check-tags string                 [ENV: CA_CHECK_TAGS] Tags [comma separated list] to use, if creating a check bundle
  -T, --check-target string               [ENV: CA_CHECK_TARGET] Check target host (for creating a new check) (default "cosi-tool-c7")
      --check-title string                [ENV: CA_CHECK_TITLE] Title [display name] to use, if creating a check bundle (default "<check-target> /agent")
      --collectors strings                [ENV: CA_COLLECTORS] List of builtin collectors to enable (default [procfs/cpu,procfs/diskstats,procfs/if,procfs/loadavg,procfs/vm])
  -c, --config string                     config file (default is /opt/circonus/agent/etc/circonus-agent.(json|toml|yaml)
  -d, --debug                             [ENV: CA_DEBUG] Enable debug messages
      --debug-api                         [ENV: CA_DEBUG_API] Enable Circonus API debug messages
      --debug-cgm                         [ENV: CA_DEBUG_CGM] Enable CGM debug messages
      --debug-dump-metrics string         [ENV: CA_DEBUG_DUMP_METRICS] Directory to dump sent metrics
  -h, --help                              help for circonus-agent
  -l, --listen strings                    [ENV: CA_LISTEN] Listen spec e.g. :2609, [::1], [::1]:2609, 127.0.0.1, 127.0.0.1:2609, foo.bar.baz, foo.bar.baz:2609 (default ":2609")
  -L, --listen-socket strings             [ENV: CA_LISTEN_SOCKET] Unix socket to create
      --log-level string                  [ENV: CA_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty                        [ENV: CA_LOG_PRETTY] Output formatted/colored log lines [ignored on windows]
      --no-gzip                           Disable gzip HTTP responses
      --no-statsd                         [ENV: CA_NO_STATSD] Disable StatsD listener
  -p, --plugin-dir string                 [ENV: CA_PLUGIN_DIR] Plugin directory
      --plugin-list strings               [ENV: CA_PLUGIN_LIST] List of explicit plugin commands to run
      --plugin-ttl-units string           [ENV: CA_PLUGIN_TTL_UNITS] Default plugin TTL units (default "s")
  -r, --reverse                           [ENV: CA_REVERSE] Enable reverse connection
      --reverse-broker-ca-file string     [ENV: CA_REVERSE_BROKER_CA_FILE] Broker CA certificate file
      --reverse-max-conn-retry int        [ENV: CA_REVERSE_MAX_CONN_RETRY] Max attempts to retry persistently failing reverse connection to broker [-1=indefinitely] (default -1)
      --show-config string                Show config (json|toml|yaml) and exit
      --ssl-cert-file string              [ENV: CA_SSL_CERT_FILE] SSL Certificate file (PEM cert and CAs concatenated together) (default "/opt/circonus/agent/etc/circonus-agent.pem")
      --ssl-key-file string               [ENV: CA_SSL_KEY_FILE] SSL Key file (default "/opt/circonus/agent/etc/circonus-agent.key")
      --ssl-listen string                 [ENV: CA_SSL_LISTEN] SSL listen address and port [IP]:[PORT] - setting enables SSL
      --ssl-verify                        [ENV: CA_SSL_VERIFY] Enable SSL verification (default true)
      --statsd-group-cid string           [ENV: CA_STATSD_GROUP_CID] StatsD group check bundle ID
      --statsd-group-counters string      [ENV: CA_STATSD_GROUP_COUNTERS] StatsD group metric counter handling (average|sum) (default "sum")
      --statsd-group-gauges string        [ENV: CA_STATSD_GROUP_GAUGES] StatsD group gauge operator (default "average")
      --statsd-group-prefix string        [ENV: CA_STATSD_GROUP_PREFIX] StatsD group metric prefix (default "group.")
      --statsd-group-sets string          [ENV: CA_STATSD_GROPUP_SETS] StatsD group set operator (default "sum")
      --statsd-host-category string       [ENV: CA_STATSD_HOST_CATEGORY] StatsD host metric category (default "statsd")
      --statsd-host-prefix string         [ENV: CA_STATSD_HOST_PREFIX] StatsD host metric prefix
      --statsd-port string                [ENV: CA_STATSD_PORT] StatsD port (default "8125")
  -V, --version                           Show version and exit
```


# Configuration

The Circonus agent can be configured via the command line, environment variables, and/or a configuration file. For details on using configuration files, see the configuration section of [etc/README.md](etc/README.md#main-configuration)



# Plugins

For documentation on plugins please refer to [plugins/README.md](plugins/README.md).



# Receiver

The Circonus agent provides a special handler for the endpoint `/write` which will accept HTTP POST and HTTP PUT requests containing structured JSON.

The structure of the JSON expected by the receiver is the same as the JSON format accepted from plugins. See the [JSON](plugins/README.md#json) section of the [plugin documentation](plugins/README.md) for details on the structure.

The URL syntax for sending metrics is `/write/ID` where `ID` is a prefix for all of the metrics being sent in the request.

For example:

HTTP POST `http://127.0.0.1:2609/write/test` with a payload of:

```json
{
    "t1": {
        "_type": "i",
        "_value": 32
    },
    "t2": {
        "_tags": ["abc:123"],
        "_type": "s",
        "_value": "foo"
    }
}
```

Would result in metrics in the Circonus UI of:

```
test`t1 numeric 32
test`t2|ST[abc:123] text "foo"
```



# StatsD

The Circonus  agent provides a StatsD listener by default (disable: `--no-statsd`, configure port: `--statsd-port`). It accepts the basic [StatsD metric types](https://github.com/etsy/statsd/blob/master/docs/metric_types.md#statsd-metric-types) as well as, Circonus specific metric types `h` and `t`. In addition, the StatsD listener support adding stream tags to metrics via `|#tag_list` added to a metric (where *tag_list* is a comma separated list of key:value pairs).

Syntax: `name:value|type[|@rate][|#tag_list]`

| Type | Note                            |
| ---- | ------------------------------- |
| `c`  | Counter                         |
| `g`  | Gauge                           |
| `h`  | Histogram - Circonus specific   |
| `ms` | Timing - treated as a Histogram |
| `s`  | Sets - treated as a Counter     |
| `t`  | Text - Circonus specific        |

>NOTE: the derivative metrics automatically generated with some StatsD types are not created by Circonus, as the data is already available within the Circonus UI.



# Builtin collectors

The circonus-agent has builtin collectors offering a higher level of efficiency over executing plugins. The circonus-agent `--collectors` command line option controls which collectors are enabled. Builtin collectors take precedence over plugins - if a builtin collector exists with the same ID as a plugin, the plugin will not be activated. Configuration files for builtins are located in the circonus-agent `etc` directory (e.g. `/opt/circonus/agent/etc` or `C:\circonus-agent\etc`).

Configuration:

* Command line `--collectors` (space delimited list)
* Environment `CA_COLLECTORS` (space delimited list)
* Config file `collectors` (array of strings)

* Windows default WMI collectors: `['wmi/cache', 'wmi/disk', 'wmi/ip', 'wmi/interface', 'wmi/memory', 'wmi/object', 'wmi/paging_file' 'wmi/processor', 'wmi/tcp', 'wmi/udp']`
* Linux default ProcFS collectors: `['procfs/cpu','procfs/diskstats','procfs/if','procfs/loadavg','procfs/vm']`
* Common `prometheus` (disabled if no configuration file exists)

For complete list of collectors and details on collector specific configuration see [etc/README.md](etc/README.md#collector-configurations).

To **disable** all default builtin collectors pass `--collectors=""` on the command line or configure `collectors` attribute in a configuration file.

# Manual build

1. Clone repo `git clone https://github.com/circonus-labs/circonus-agent.git`
  * circonus-agent uses go modules, go1.11+ is required
  * clone **oustide** of `GOPATH` (or use `GO111MODULE=on`)
1. Build `go build -o circonus-agentd`
1. Install `cp circonus-agentd /opt/circonus/agent/sbin`

Unless otherwise noted, the source files are distributed under the BSD-style license found in the [LICENSE](LICENSE) file.
