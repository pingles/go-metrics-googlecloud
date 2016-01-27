package metrics

import (
	"fmt"
	metrics "github.com/rcrowley/go-metrics"
	"net/http"
	"os"
	"time"
)

const hostnameLabelKey = "hostname"

func (r *reporter) ensureMetric(name string, t string) error {
	fullMetricName := NameInDomain(name)
	metric := &Metric{
		Name: fullMetricName,
		Type: t,
		Labels: map[string]string{
			NameInDomain(hostnameLabelKey): "Hostname of machine sending metric.",
		},
	}
	fmt.Println("creating metric:", metric.Name)
	err := CreateMetric(r.client, r.project, metric)
	if err != nil {
		fmt.Println("error creating metric:", err.Error())
	}
	return nil
}

func (r *reporter) reportMetric(metricName string, metricType string, val interface{}) {
	r.ensureMetric(metricName, metricType)
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("error retrieving hostname, can't report:", err.Error())
		return
	}
	labels := make(map[string]string, 1)
	labels[hostnameLabelKey] = hostname

	timeseries := &Timeseries{
		MetricName: NameInDomain(metricName),
		Now:        time.Now(),
		Labels:     labels,
	}

	if metricType == Int {
		timeseries.Int64Value = val.(int64)
	} else if metricType == Double {
		timeseries.DoubleValue = val.(float64)
	}

	err = WriteTimeseries(r.client, r.project, []*Timeseries{timeseries})
	if err != nil {
		fmt.Println("error writing timeseries:", err.Error())
	}
}

func (r *reporter) reportMeter(name string, val metrics.Meter) {
	r.reportMetric(fmt.Sprintf("%s.count", name), Int, val.Count())
}

type reporter struct {
	client  *http.Client
	project string
}

func newReporter(client *http.Client, project string) *reporter {
	return &reporter{client, project}
}

func (r *reporter) report(name string, val interface{}) {
	switch metric := val.(type) {
	case metrics.Meter:
		r.reportMeter(name, metric)
	}
}

func GoogleCloudMonitoring(r metrics.Registry, d time.Duration, client *http.Client, project string) {
	reporter := newReporter(client, project)
	for _ = range time.Tick(d) {
		r.Each(reporter.report)
	}
}
