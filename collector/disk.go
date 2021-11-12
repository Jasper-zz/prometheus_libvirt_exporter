package collector

import (
	"prometheus_libvirt_exporter/collector/qga"
	"prometheus_libvirt_exporter/internal"

	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	diskCollectorSubsystem = "disk"
)

type diskCollector struct {
	readRequests  *prometheus.Desc
	writeRequests *prometheus.Desc
	readBytes     *prometheus.Desc
	writeBytes    *prometheus.Desc
	availBytes    *prometheus.Desc
	sizeBytes     *prometheus.Desc
	inodes        *prometheus.Desc
	availInodes   *prometheus.Desc
}

func init() {
	registerCollector("disk", newDiskCollector)
}

func newDiskCollector() (Collector, error) {
	c := &diskCollector{
		readRequests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "read_requests"),
			"",
			[]string{"uuid", "target_device"}, nil),
		writeRequests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "write_requests"),
			"",
			[]string{"uuid", "target_device"}, nil),
		readBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "read_bytes"),
			"",
			[]string{"uuid", "target_device"}, nil),
		writeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "write_bytes"),
			"",
			[]string{"uuid", "target_device"}, nil),

		sizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "size_bytes"),
			"",
			[]string{"uuid", "target_device", "fstype", "mountpoint"}, nil),
		availBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "avail_bytes"),
			"",
			[]string{"uuid", "target_device", "fstype", "mountpoint"}, nil),

		inodes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "inodes"),
			"",
			[]string{"uuid", "target_device", "fstype", "mountpoint"}, nil),
		availInodes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskCollectorSubsystem, "avail_inodes"),
			"",
			[]string{"uuid", "target_device", "fstype", "mountpoint"}, nil),
	}

	return c, nil
}

func (c *diskCollector) Update(ch chan<- prometheus.Metric, stats *libvirt.DomainStats, uuid string, rs rpcSet) error {
	diskStats := stats.Block
	for _, v := range diskStats {
		if v.Name != "hda" {
			ch <- prometheus.MustNewConstMetric(c.readRequests,
				prometheus.GaugeValue,
				float64(v.RdReqs),
				uuid,
				v.Name)
			ch <- prometheus.MustNewConstMetric(c.writeRequests,
				prometheus.GaugeValue,
				float64(v.WrReqs),
				uuid,
				v.Name)
			ch <- prometheus.MustNewConstMetric(c.readBytes,
				prometheus.GaugeValue,
				float64(v.RdBytes),
				uuid,
				v.Name)
			ch <- prometheus.MustNewConstMetric(c.writeBytes,
				prometheus.GaugeValue,
				float64(v.WrBytes),
				uuid,
				v.Name)
		}
	}
	if rs.GuestExec {
		execArg := qga.GuestExecArg{
			Path:          "/usr/bin/df",
			Arg:           []string{"--output=source,fstype,target,itotal,iavail,size,avail"},
			CaptureOutput: true,
		}
		execRet, err := qga.Exec(stats.Domain, execArg)
		if err != nil {
			return err
		}
		fsStats, err := internal.GetFilesystem(execRet)
		if err != nil {
			return err
		}
		//fmt.Printf("%+v\n", fsStats)

		for _, s := range fsStats {
			ch <- prometheus.MustNewConstMetric(c.sizeBytes,
				prometheus.GaugeValue,
				float64(s.Size),
				uuid,
				s.Labels.Device,
				s.Labels.FsType,
				s.Labels.MountPoint)
			ch <- prometheus.MustNewConstMetric(c.availBytes,
				prometheus.GaugeValue,
				float64(s.Avail),
				uuid,
				s.Labels.Device,
				s.Labels.FsType,
				s.Labels.MountPoint)
			ch <- prometheus.MustNewConstMetric(c.inodes,
				prometheus.GaugeValue,
				float64(s.Inodes),
				uuid,
				s.Labels.Device,
				s.Labels.FsType,
				s.Labels.MountPoint)
			ch <- prometheus.MustNewConstMetric(c.availInodes,
				prometheus.GaugeValue,
				float64(s.IAvail),
				uuid,
				s.Labels.Device,
				s.Labels.FsType,
				s.Labels.MountPoint)
		}
	}
	return nil
}
