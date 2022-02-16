package collector

import (
	"encoding/json"
	"fmt"
	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const namespace = "libvirt"

var (
	factories          = make(map[string]Collector)
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"node_exporter: Duration of a collector scrape.",
		[]string{"domain"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"node_exporter: Whether a collector succeeded.",
		[]string{"domain"},
		nil,
	)
)

type Commands struct {
	Name            string `json:"name"`
	Enabled         bool   `json:"enabled"`
	SuccessResponse bool   `json:"success-response"`
}

type SupportedCommands struct {
	Return struct {
		Version           string     `json:"version"`
		SupportedCommands []Commands `json:"supported_commands"`
	}
}

type rpcSet struct {
	GuestFileRead bool
	GuestExec     bool
}

type Collector interface {
	Update(ch chan<- prometheus.Metric, dom *libvirt.DomainStats, uuid string, rs rpcSet) error
}

type LibvirtCollector struct {
	Uri        string
	Collectors map[string]Collector
	logger     logrus.Logger
}

func registerCollector(collector string, factory func() (Collector, error)) {
	c, err := factory()
	if err != nil {
		logrus.Debug("failed to init collector ", collector)
		return
	}
	factories[collector] = c
}

func NewLibvirtCollector(uri string, logger logrus.Logger, filters ...string) (*LibvirtCollector, error) {
	collectors := make(map[string]Collector)
	if len(filters) == 0 {
		collectors = factories
	} else {
		for _, filter := range filters {
			_, exist := factories[filter]
			if !exist {
				return nil, fmt.Errorf("missing collector: %s", filter)
			}
			collectors[filter] = factories[filter]
		}
	}

	return &LibvirtCollector{
		Uri:        uri,
		Collectors: collectors,
		logger:     logger,
	}, nil
}

func (l *LibvirtCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

func (l *LibvirtCollector) excute(ch chan<- prometheus.Metric, domStats libvirt.DomainStats) {
	var success float64
	begin := time.Now()
	uuid, _ := domStats.Domain.GetUUIDString()

	frSet, _ := GetRpcSet(domStats.Domain)

	wg := sync.WaitGroup{}
	for name, c := range l.Collectors {
		wg.Add(1)
		go func(n string, c Collector) {
			err := c.Update(ch, &domStats, uuid, frSet)
			if err != nil {
				l.logger.Debug("uuid=", uuid, " collector=", n, " error=", err)
				success = 0
			} else {
				success = 1
			}
			wg.Done()
		}(name, c)
	}
	wg.Wait()

	duration := time.Since(begin)
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), uuid)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, uuid)
}

func (l *LibvirtCollector) Collect(ch chan<- prometheus.Metric) {
	conn, err := libvirt.NewConnect(l.Uri)
	if err != nil {
		l.logger.Fatalf("failed to connect to %s.", l.Uri)
	}
	defer conn.Close()

	// get all doms stats
	statsAll, err := conn.GetAllDomainStats([]*libvirt.Domain{},
		libvirt.DOMAIN_STATS_STATE|
			libvirt.DOMAIN_STATS_CPU_TOTAL|
			libvirt.DOMAIN_STATS_BALLOON|
			libvirt.DOMAIN_STATS_BLOCK|
			libvirt.DOMAIN_STATS_INTERFACE,
		//libvirt.CONNECT_GET_ALL_DOMAINS_STATS_NOWAIT, // maybe in future
		libvirt.CONNECT_GET_ALL_DOMAINS_STATS_ACTIVE)

	defer func(statsAll []libvirt.DomainStats) {
		for _, domStat := range statsAll {
			domStat.Domain.Free()
		}
	}(statsAll)

	if err != nil {
		l.logger.Warn("failed to get stats.")
	}

	wg := sync.WaitGroup{}
	for _, stats := range statsAll {
		wg.Add(1)
		go func(s libvirt.DomainStats) {
			l.excute(ch, s)
			wg.Done()
		}(stats)
	}
	wg.Wait()
}

func GetRpcSet(dom *libvirt.Domain) (rpcSet, error) {
	rs := rpcSet{false, false}
	cmdSet, err := dom.QemuAgentCommand("{\"execute\":\"guest-info\"}", 1, 0)
	if err != nil {
		return rs, err
	}
	sc := SupportedCommands{}
	err = json.Unmarshal([]byte(cmdSet), &sc)
	if err != nil {
		return rs, err
	}
	rs = rpcSet{true, true}
	for _, cmd := range sc.Return.SupportedCommands {
		if cmd.Name == "guest-file-read" || cmd.Name == "guest-file-open" || cmd.Name == "guest-file-close" {
			if !cmd.Enabled {
				rs.GuestFileRead = cmd.Enabled
			}
		}
		if cmd.Name == "guest-exec" || cmd.Name == "guest-exec-status" {
			if !cmd.Enabled {
				rs.GuestExec = cmd.Enabled
			}
		}
	}
	return rs, err
}
