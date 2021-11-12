package collector

import (
	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
	"prometheus_libvirt_exporter/collector/qga"
	"prometheus_libvirt_exporter/internal"
)

const (
	cpuCollectorSubsystem = "cpu"
)

type cpuCollector struct {
	cpuCores         *prometheus.Desc
	cpuSystemTime    *prometheus.Desc
	cpuCpuTime       *prometheus.Desc
	cpuUserTime      *prometheus.Desc
	cpuload1         *prometheus.Desc
	cpuload5         *prometheus.Desc
	cpuload15        *prometheus.Desc
	qgaCpuSystemTime *prometheus.Desc
	qgaCpuCpuTime    *prometheus.Desc
	qgaCpuUserTime   *prometheus.Desc
	qgaCpuStealTime  *prometheus.Desc
}

func init() {
	registerCollector("cpu", newCPUCollector)
}

func newCPUCollector() (Collector, error) {
	c := &cpuCollector{
		cpuCores: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "cores"),
			"",
			[]string{"uuid"}, nil),

		cpuSystemTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "system_time"),
			"",
			[]string{"uuid"}, nil),
		cpuCpuTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "cpu_time"),
			"",
			[]string{"uuid"}, nil),
		cpuload1: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "load1"),
			"",
			[]string{"uuid"}, nil),
		cpuload5: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "load5"),
			"",
			[]string{"uuid"}, nil),
		cpuload15: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "load15"),
			"",
			[]string{"uuid"}, nil),
		cpuUserTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "user_time"),
			"",
			[]string{"uuid"}, nil),

		qgaCpuSystemTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "qga_system_time"),
			"",
			[]string{"uuid"}, nil),
		qgaCpuCpuTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "qga_cpu_time"),
			"",
			[]string{"uuid"}, nil),
		qgaCpuUserTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "qga_user_time"),
			"",
			[]string{"uuid"}, nil),
		qgaCpuStealTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "qga_steal_time"),
			"",
			[]string{"uuid"}, nil),
	}

	return c, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric, stats *libvirt.DomainStats, uuid string, rs rpcSet) error {
	cpuStats := stats.Cpu
	info, err := stats.Domain.GetInfo()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.cpuSystemTime,
		prometheus.GaugeValue,
		float64(cpuStats.System),
		uuid)
	ch <- prometheus.MustNewConstMetric(c.cpuCpuTime,
		prometheus.GaugeValue,
		float64(cpuStats.Time),
		uuid)
	ch <- prometheus.MustNewConstMetric(c.cpuUserTime,
		prometheus.GaugeValue,
		float64(cpuStats.User),
		uuid)

	ch <- prometheus.MustNewConstMetric(c.cpuCores,
		prometheus.GaugeValue,
		float64(info.NrVirtCpu),
		uuid)

	if rs.GuestFileRead {
		dataStat, err := qga.ReadFile(stats.Domain, "/proc/stat")
		if err != nil {
			return err
		}
		s, _ := internal.GetStat(dataStat)
		//fmt.Printf("%+v\n", s)

		ch <- prometheus.MustNewConstMetric(c.qgaCpuSystemTime,
			prometheus.GaugeValue,
			float64(s.CPUTotal.System),
			uuid)
		ch <- prometheus.MustNewConstMetric(c.qgaCpuStealTime,
			prometheus.GaugeValue,
			float64(s.CPUTotal.Steal),
			uuid)
		ch <- prometheus.MustNewConstMetric(c.qgaCpuUserTime,
			prometheus.GaugeValue,
			float64(s.CPUTotal.User),
			uuid)

		dataLoad, err := qga.ReadFile(stats.Domain, "/proc/loadavg")
		if err != nil {
			return err
		}
		l, err := internal.GetLoad(dataLoad)
		if err != nil {
			return err
		}
		//fmt.Printf("%+v\n", l)
		ch <- prometheus.MustNewConstMetric(c.cpuload1,
			prometheus.GaugeValue,
			float64(l[0]),
			uuid)
		ch <- prometheus.MustNewConstMetric(c.cpuload5,
			prometheus.GaugeValue,
			float64(l[1]),
			uuid)
		ch <- prometheus.MustNewConstMetric(c.cpuload15,
			prometheus.GaugeValue,
			float64(l[2]),
			uuid)
	}
	return nil
}
