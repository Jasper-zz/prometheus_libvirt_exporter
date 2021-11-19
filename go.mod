module prometheus_libvirt_exporter

go 1.15

require (
	github.com/libvirt/libvirt-go v7.0.0+incompatible
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
