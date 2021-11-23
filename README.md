# Libvirt exporter

Prometheus exporter for vm metrics written in Go with pluggable metric collectors.

## Installation and Usage

If you are new to Prometheus there is a [simple step-by-step guide](https://prometheus.io/docs/guides/node-exporter/).

The `libvirt_exporter` listens on HTTP port 9108 by default. See the `--help` output for more options.

### Docker

For situations where Docker deployment is needed, some extra flags must be used to allow the `libvirt_exporter` access to the host libvirt-sock file.

```bash
docker build -t prom/prometheus_libvirt_exporter:v2.0.1 .
```

```bash
docker run -d \
  --restart="always" \
  --net="host" \
  --name libvirt_exporter \
  -v "/var/run/libvirt/libvirt-sock:/var/run/libvirt/libvirt-sock" \
  prom/prometheus-libvirt-exporter:2.0.1
```

For Docker compose, similar flag changes are needed.

```yaml
---
version: '3.8'

services:
  libvirt_exporter:
    image: prom/prometheus-libvirt-exporter:2.0.1
    container_name: libvirt_exporter
    network_mode: host
    restart: always
    volumes:
      - '/var/run/libvirt/libvirt-sock:/var/run/libvirt/libvirt-sock'
```

## Collectors

The tablesbelow list all existing collectors.

| Name          | Description                                                         |
| ------------- | ------------------------------------------------------------------- |
| cpu           | Exposes VM CPU statistics                                           |
| meminfo       | Exposes memory statistics.                                          |
| diskstats     | Exposes disk I/O statistics.                                        |
| netdev        | Exposes network interface statistics such as bytes transferred.     |
| netstat       | Exposes network statistics.                |

### Extend collectors

| Name          | Description                                                         |
| ------------- | ------------------------------------------------------------------- |
| filesystem    | Exposes filesystem statistics, such as disk space used.             |
| loadavg       | Exposes load average.                                               |

### Filtering enabled collectors

The `libvirt_exporter` will expose all metrics from enabled collectors by default.  This is the recommended way to collect metrics to avoid errors when comparing metrics of different families.

For advanced use the `libvirt_exporter` can be passed an optional list of collectors to filter metrics. The `collect[]` parameter may be used multiple times.  In Prometheus configuration you can use this syntax under the [scrape config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#<scrape_config>).

```
  params:
    collect[]:
      - foo
      - bar
```

This can be useful for having different Prometheus servers collect specific metrics from nodes.

## Development building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)
* RHEL/CentOS: `glibc-static` package.

Building:

```bash
    git clone https://github.com/Jasper-zz/prometheus_libvirt_exporter.git
    cd prometheus_libvirt_exporter
    make build
    ./.build/prometheus_libvirt_exporter <flags>
```

To see all available configuration flags:

```bash
    ./.build/prometheus_libvirt_exporter -h
```
