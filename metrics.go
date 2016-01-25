package metrics

import (
	cm "google.golang.org/api/cloudmonitoring/v2beta2"
	"net/http"
)

func newMetricDescriptor(name string) *cm.MetricDescriptor {
	return &cm.MetricDescriptor{
		Name: name,
		TypeDescriptor: &cm.MetricDescriptorTypeDescriptor{
			MetricType: "gauge",
			ValueType:  "int64",
		},
	}
}

func CreateMetric(client *http.Client, project, metricName string) error {
	service, err := cm.New(client)
	if err != nil {
		return err
	}
	_, err = service.MetricDescriptors.Create(project, newMetricDescriptor(metricName)).Do()
	if err != nil {
		return err
	}

	return nil
}

func ListMetrics(client *http.Client, project string) (*cm.ListMetricDescriptorsResponse, error) {
	service, err := cm.New(client)
	if err != nil {
		return nil, err
	}
	return service.MetricDescriptors.List(project, &cm.ListMetricDescriptorsRequest{}).Do()
}
