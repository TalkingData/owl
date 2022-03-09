// Copyright 2016 Ajit Yagaty
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builder

import "encoding/json"

type MetricBuilder interface {
	// Add a new metric to the builder.
	AddMetric(name string) Metric

	// Get a list of all the metrics that are part of the builder.
	GetMetrics() []Metric

	// Encode the Metrics list into JSON.
	Build() ([]byte, error)
}

// Type that implements the MetricBuilder interface.
type mBuilder struct {
	Metrics []Metric `json:"metrics"`
}

func NewMetricBuilder() MetricBuilder {
	return &mBuilder{}
}

func (mb *mBuilder) AddMetric(name string) Metric {
	m := NewMetric(name)
	mb.Metrics = append(mb.Metrics, m)
	return m
}

func (mb *mBuilder) GetMetrics() []Metric {
	return mb.Metrics
}

func (mb *mBuilder) Build() ([]byte, error) {
	// Make sure the contents of each metric object are correct.
	for _, m := range mb.Metrics {
		err := m.validate()
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(mb.Metrics)
}
