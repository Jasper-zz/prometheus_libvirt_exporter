package main

import (
	"fmt"
	"path"
	"net/http"
	"prometheus_libvirt_exporter/collector"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
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
	logger logrus.Logger
}

func newHandler(logger logrus.Logger) *handler {
	h := &handler{
		logger: logger,
	}

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	if len(filters) != 0 {
		h.logger.Info("collect filters:", filters)
	}

	le, err := collector.NewLibvirtExporter(*libvirtURI, h.logger, filters...)
	if err != nil {
		h.logger.Errorf("could't create collcetor: %s", err)
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(le)
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{registry},
		promhttp.HandlerOpts{
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: *maxRequests,
		},
	)

	if *includeExporterMetrics {
		registry.MustRegister(
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
		)
	}

	handler.ServeHTTP(w, r)
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
