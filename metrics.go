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

func (r *reporter) reportMeter(name string, val metrics.Meter) {
	r.ensureMetric(fmt.Sprintf("%s.count", name), Int)
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("error retrieving hostname, can't report:", err.Error())
		return
	}
	labels := make(map[string]string, 1)
	labels[hostnameLabelKey] = hostname

	timeseries := []*Timeseries{
		&Timeseries{
			MetricName: NameInDomain(fmt.Sprintf("%s.count", name)),
			Now:        time.Now(),
			Int64Value: val.Count(),
			Labels:     labels,
		},
	}
	err = WriteTimeseries(r.client, r.project, timeseries)
	if err != nil {
		fmt.Println("error writing timeseries:", err.Error())
	}
}

type reporter struct {
	client  *http.Client
	project string
}

func newReporter(client *http.Client, project string) *reporter {
	return &reporter{client, project}
}

func (r *reporter) report(name string, val interface{}) {
	fmt.Println("reporting", name)
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
