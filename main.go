package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/dasa-health/azure_metrics_exporter/azure"
	"github.com/dasa-health/azure_metrics_exporter/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/subosito/gotenv"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests.").Default(":9276").String()
)

func init() {
	prometheus.MustRegister(version.NewCollector("azure_exporter"))
	gotenv.Load()
}

// Collector generic collector type
type Collector struct {
	tagValue string
}

// Describe implemented with dummy data to satisfy interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

// Collect - collect results from Azure Montior API and create Prometheus metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	resourceAggregation := os.Getenv("metric_aggregation")

	if c.tagValue == "" {
		logger.Error("Tag value is empty", nil)
	}

	logger.Info("Get all resources", nil)

	ac, err := azure.GetAccessToken()

	if err != nil {
		logger.Error("Failed to get access token: %v", err)
	}
	resources, err := ac.GetResources(c.tagValue)

	if err != nil {
		logger.Error("Failed to get all resources: %v", err)
	}

	for _, resource := range resources.Value {

		if !azure.ValidateTypeMetric(resource.Type) {
			continue
		}

		logger.Info(fmt.Sprintf("Retrieves all metric definitions of resource [ %s ]", resource.Name), nil)

		typeMetrics, err := ac.GetMetricTypes(resource.ID, resource.Type)

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get metrics types from resources %s: %v", resource.Name, err), nil)
		}

		logger.Info(fmt.Sprintf("Treats metric definitions found from resource [ %s ]", resource.Name), nil)

		typeMetricsTreated := azure.TreatTypeMetric(typeMetrics)

		for _, typeMetric := range typeMetricsTreated {

			metricValueData, err := ac.GetMetric(resource.ID, typeMetric, resourceAggregation)

			if err != nil {
				logger.Error(fmt.Sprintf("Failed to get metrics for target %s: %v", resource.ID, err), nil)
				continue
			}

			if metricValueData.Value == nil {
				logger.Error(fmt.Sprintf("Metric %v not found at target %v\n", typeMetric, resource.ID), nil)
				continue
			}
			if len(metricValueData.Value) <= 0 || len(metricValueData.Value[0].Timeseries) <= 0 || len(metricValueData.Value[0].Timeseries[0].Data) == 0 {
				logger.Error(fmt.Sprintf("No metric data returned for metric %v at target %v\n", typeMetric, resource.ID), nil)
				continue
			}
			for _, value := range metricValueData.Value {

				metricName, err := SanitizeMetricName(value.Name.Value, value.Unit)

				if err != nil {
					logger.Error(fmt.Sprintf("Failed to get metrics types from resources %s: %v", resource.Name, err), nil)
				}
				defer recoverMetric(resource.Name, value.Name.Value)

				if len(value.Timeseries) <= 0 || len(value.Timeseries[0].Data) <= 0 {
					continue
				}

				metricValue := value.Timeseries[0].Data[len(value.Timeseries[0].Data)-1]

				labels := CreateResourceLabels(value.ID, resource.Name, resource.Type, IdentifyEnvironmentResource(resource.Name))

				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(metricName+"_total", metricName+"_total", nil, labels),
					prometheus.GaugeValue,
					metricValue.Total,
				)

				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(metricName+"_average", metricName+"_average", nil, labels),
					prometheus.GaugeValue,
					metricValue.Average,
				)

				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(metricName+"_min", metricName+"_min", nil, labels),
					prometheus.GaugeValue,
					metricValue.Minimum,
				)

				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(metricName+"_max", metricName+"_max", nil, labels),
					prometheus.GaugeValue,
					metricValue.Maximum,
				)
			}
		}
	}

	logger.Info("Finally Get all resources", nil)

}

func recoverMetric(resource, metric string) {
	if r := recover(); r != nil {
		logger.Info(fmt.Sprintf("Recovered error from metric %s from resource %s : %v", metric, resource, r), nil)
		debug.PrintStack()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()
	collector := &Collector{tagValue: r.URL.Query().Get("tagValue")}
	registry.MustRegister(collector)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head>
            <title>Azure Exporter</title>
            </head>
            <body>
            <h1>Azure Exporter</h1>
						<p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})

	http.HandleFunc("/metrics", handler)
	log.Printf("azure_metrics_exporter listening on port %v", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
		os.Exit(1)
	}

}
