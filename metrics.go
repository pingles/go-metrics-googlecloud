package metrics

import (
	"fmt"
	cm "google.golang.org/api/cloudmonitoring/v2beta2"
	"net/http"
	"strings"
)

const (
	Int    = "int64"
	Double = "double"
)

type Metric struct {
	Name        string
	Description string
	Labels      map[string]string
	Type        string
}

func newMetricDescriptor(metric *Metric) *cm.MetricDescriptor {
	labelDescriptors := make([]*cm.MetricDescriptorLabelDescriptor, len(metric.Labels))
	index := 0
	for key, description := range metric.Labels {
		labelDescriptors[index] = &cm.MetricDescriptorLabelDescriptor{
			Key:         key,
			Description: description,
		}
		index += 1
	}

	return &cm.MetricDescriptor{
		Name:        metric.Name,
		Description: metric.Description,
		Labels:      labelDescriptors,
		TypeDescriptor: &cm.MetricDescriptorTypeDescriptor{
			MetricType: "gauge",
			ValueType:  metric.Type,
		},
	}
}

// ensures names are prefixed with the cloud monitoring domain. name would
// normally be hierarchically structured, separated with /.
func NameInDomain(name string) string {
	return fmt.Sprintf("custom.cloudmonitoring.googleapis.com/%s", name)
}

func CreateMetric(client *http.Client, project string, m *Metric) error {
	service, err := cm.New(client)
	if err != nil {
		return err
	}
	_, err = service.MetricDescriptors.Create(project, newMetricDescriptor(m)).Do()
	if err != nil {
		return err
	}

	return nil
}

func DeleteMetric(client *http.Client, project, name string) error {
	service, err := cm.New(client)
	if err != nil {
		return err
	}
	_, err = service.MetricDescriptors.Delete(project, name).Do()

	return err
}

func IsCustom(desc *cm.MetricDescriptor) bool {
	return strings.HasPrefix(desc.Name, NameInDomain(""))
}

func ListMetrics(client *http.Client, project string) ([]*cm.MetricDescriptor, error) {
	service, err := cm.New(client)
	if err != nil {
		return nil, err
	}
	resp, err := service.MetricDescriptors.List(project, &cm.ListMetricDescriptorsRequest{}).Do()
	if err != nil {
		return nil, err
	}
	return resp.Metrics, nil
}
