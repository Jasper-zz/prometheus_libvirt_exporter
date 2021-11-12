package collector

import (
	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	memCollectorSubsystem = "mem"
)

type memCollector struct {
	memTotal     *prometheus.Desc
	memUsed      *prometheus.Desc
	memAvailable *prometheus.Desc
}

func init() {
	registerCollector("mem", newMemCollector)
}

func newMemCollector() (Collector, error) {
	c := &memCollector{
		memTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, memCollectorSubsystem, "total",
			),
			"",
			[]string{"uuid"}, nil),
		memUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, memCollectorSubsystem, "used"),
			"used with buff/cache",
			[]string{"uuid"}, nil),
		memAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, memCollectorSubsystem, "available"),
			"available without buff/cache",
			[]string{"uuid"}, nil),
	}

	return c, nil
}

func (c *memCollector) Update(ch chan<- prometheus.Metric, stats *libvirt.DomainStats, uuid string, rs rpcSet) error {
	memStats := stats.Balloon
	ch <- prometheus.MustNewConstMetric(c.memTotal,
		prometheus.GaugeValue,
		float64(memStats.Available),
		uuid)
	ch <- prometheus.MustNewConstMetric(c.memUsed,
		prometheus.GaugeValue,
		float64(memStats.Available-memStats.Unused),
		uuid)
	ch <- prometheus.MustNewConstMetric(c.memAvailable,
		prometheus.GaugeValue,
		float64(memStats.Usable),
		uuid)
	return nil
}
