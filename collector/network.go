package collector

import (
	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	networkCollectorSubsystem = "network"
)

type networkCollector struct {
	receiveBytes    *prometheus.Desc
	receivePackets  *prometheus.Desc
	receiveErrors   *prometheus.Desc
	receiveDrops    *prometheus.Desc
	transmitBytes   *prometheus.Desc
	transmitPackets *prometheus.Desc
	transmitErrors  *prometheus.Desc
	transmitDrops   *prometheus.Desc
}

func init() {
	registerCollector("network", newNetworkCollector)
}

func newNetworkCollector() (Collector, error) {
	c := &networkCollector{
		receiveBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "receive_bytes",
			),
			"",
			[]string{"uuid", "target_device"}, nil),
		receivePackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "receive_packets"),
			"",
			[]string{"uuid", "target_device"}, nil),
		receiveErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "receive_errors"),
			"",
			[]string{"uuid", "target_device"}, nil),
		receiveDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "receive_drops"),
			"",
			[]string{"uuid", "target_device"}, nil),
		transmitBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "transmit_bytes",
			),
			"",
			[]string{"uuid", "target_device"}, nil),
		transmitPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "transmit_packets"),
			"",
			[]string{"uuid", "target_device"}, nil),
		transmitErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "transmit_errors"),
			"",
			[]string{"uuid", "target_device"}, nil),
		transmitDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, networkCollectorSubsystem, "transmit_drops"),
			"",
			[]string{"uuid", "target_device"}, nil),
	}

	return c, nil
}

func (c *networkCollector) Update(ch chan<- prometheus.Metric, stats *libvirt.DomainStats, uuid string, rs rpcSet) error {
	netStats := stats.Net
	for _, v := range netStats {
		ch <- prometheus.MustNewConstMetric(c.receiveBytes,
			prometheus.GaugeValue,
			float64(v.RxBytes),
			uuid,
			v.Name)
		ch <- prometheus.MustNewConstMetric(c.receivePackets,
			prometheus.GaugeValue,
			float64(v.RxPkts),
			uuid,
			v.Name)
		ch <- prometheus.MustNewConstMetric(c.receiveErrors,
			prometheus.GaugeValue,
			float64(v.RxErrs),
			uuid,
			v.Name)
		ch <- prometheus.MustNewConstMetric(c.receiveDrops,
			prometheus.GaugeValue,
			float64(v.RxDrop),
			uuid,
			v.Name)
		ch <- prometheus.MustNewConstMetric(c.transmitBytes,
			prometheus.GaugeValue,
			float64(v.TxBytes),
			uuid,
			v.Name)
		ch <- prometheus.MustNewConstMetric(c.transmitPackets,
			prometheus.GaugeValue,
			float64(v.TxPkts),
			uuid,
			v.Name)
		ch <- prometheus.MustNewConstMetric(c.transmitErrors,
			prometheus.GaugeValue,
			float64(v.TxErrs),
			uuid,
			v.Name)
		ch <- prometheus.MustNewConstMetric(c.transmitDrops,
			prometheus.GaugeValue,
			float64(v.TxDrop),
			uuid,
			v.Name)
	}
	return nil
}
