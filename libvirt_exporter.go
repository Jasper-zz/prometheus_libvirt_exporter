package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"path"
	"prometheus_libvirt_exporter/collector"
	"runtime"
	"sort"
)

var (
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address on which to expose metrics and web interface.",
	).Default(":9108").String()
	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	maxRequests = kingpin.Flag(
		"web.max-requests",
		"Maximum number of parallel scrape requests. Use 0 to disable.",
	).Default("40").Int()
	includeExporterMetrics = kingpin.Flag(
		"web.disable-exporter-metrics",
		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
	).Bool()
	libvirtURI = kingpin.Flag(
		"libvirt.uri",
		"Libvirt URI from which to extract metrics.",
	).Default("qemu:///system").String()
	logLevel = kingpin.Flag(
		"log.level",
		"log level Debug|Info|Warn|Error",
	).Default("Info").String()
)

type handler struct {
	exporterMetricsRegistry *prometheus.Registry
	unfilteredHandler       http.Handler
	logger                  logrus.Logger
}

func newHandler(logger logrus.Logger) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		logger:                  logger,
	}

	if *includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
		)
	}

	if innerHandler, err := h.innerHandler(); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		h.unfilteredHandler = innerHandler
	}

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	h.logger.Debug("collect query: ", filters)

	if len(filters) == 0 {
		h.unfilteredHandler.ServeHTTP(w, r)
		return
	}

	filteredHandler, err := h.innerHandler(filters...)
	if err != nil {
		h.logger.Infof("Couldn't create filtered metrics handler: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Couldn't create filtered metrics handler: %s", err)))
		return
	}

	filteredHandler.ServeHTTP(w, r)
}

func (h *handler) innerHandler(filters ...string) (http.Handler, error) {
	lc, err := collector.NewLibvirtCollector(*libvirtURI, h.logger, filters...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	if len(filters) == 0 {
		h.logger.Info("msg=Enabled_collectors")
		collectors := []string{}
		for n := range lc.Collectors {
			collectors = append(collectors, n)
		}
		sort.Strings(collectors)
		for _, c := range collectors {
			h.logger.Info("collector=", c)
		}
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("libvirt_exporter"))
	if err := r.Register(lc); err != nil {
		return nil, fmt.Errorf("couldn't register libvirt collector: %s", err)
	}

	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, r},
		promhttp.HandlerOpts{
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: *maxRequests,
			Registry:            h.exporterMetricsRegistry,
		},
	)

	return handler, nil
}

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" caller=%s:%d -", filename, f.Line)
		},
	})
	logger.SetReportCaller(true)
	logrusLevel, err := logrus.ParseLevel(*logLevel)
	if err == nil {
		logger.SetLevel(logrusLevel)
	}

	http.Handle(*metricsPath, newHandler(*logger))
	// http.Handle(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Libvirt Exporter</title></head>
			<body>
			<h1>Libvirt Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	logger.Info("ListeningOn=", *listenAddress)
	logger.Fatal(http.ListenAndServe(*listenAddress, nil))
	//
	//log.Info("Listening on ", *listenAddress)
	//log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
