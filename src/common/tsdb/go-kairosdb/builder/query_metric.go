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

// Query request for a metric. If a metric is queried by name only then all
// data points for all tags are returned. You can narrow down the query by
// adding tags so only data points associated with those tags are returned.
//
// Aggregators may be added to the metric. An aggregator performs an operation
// on the data such as summing or averaging. If multiple aggregators are added,
// the output of the first is sent to the input of the next, and so forth until
// all aggregators have been processed, These are processed in the order they
// were added.
//
// The results of the query can be grouped in various ways using a grouper.
// For example, if you had a metric with a customer tag, the resulting data
// points could be grouped by the different customers. Multiple groupers can be
// used so you could, for example, group by tag and value.
//
// Note that aggregation is very fast but grouping can slow down the query.

type OrderType string

const (
	ASCENDING  OrderType = "asc"
	DESCENDING OrderType = "desc"
)

type QueryMetric interface {
	// Add a map of tags. This narrows the query to only show data points
	// associated with the tags' values.
	AddTags(tags map[string][]string) QueryMetric

	// Adds a tag with multiple values. This narrows the query to only show
	// data points associated with the tag's values.
	AddTag(name string, val string) QueryMetric

	// Adds an aggregator to the metric.
	AddAggregator(aggr Aggregator) QueryMetric

	// Adds a grouper to the metric.
	AddGrouper(tagks []string) QueryMetric

	// Limits the number of data point returned from the query.
	// The limit is done before aggregators are executed.
	SetLimit(limit int) QueryMetric

	// Orders the data points. The server default is ascending.
	SetOrder(order OrderType) QueryMetric

	// Validates the contents of the QueryMetric instance.
	Validate() error
}

type qMetric struct {
	Tags        map[string][]string `json:"tags,omitempty"`
	Name        string              `json:"name,omitempty"`
	Limit       int                 `json:"limit,omitempty"`
	Aggregators []Aggregator        `json:"aggregators,omitempty"`
	Order       OrderType           `json:"order,omitempty"`
	GroupBys    []GroupBy           `json:"group_by",omitempty`
}

type GroupBy struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func NewQueryMetric(name string) QueryMetric {
	return &qMetric{
		Name:        name,
		Tags:        make(map[string][]string),
		Aggregators: make([]Aggregator, 0),
		GroupBys:    make([]GroupBy, 0),
	}
}

func (qm *qMetric) AddTags(tags map[string][]string) QueryMetric {
	for k, v := range tags {
		qm.Tags[k] = append(qm.Tags[k], v...)
	}

	return qm
}

func (qm *qMetric) AddTag(name string, value string) QueryMetric {
	qm.Tags[name] = append(qm.Tags[name], value)
	return qm
}

func (qm *qMetric) AddAggregator(aggr Aggregator) QueryMetric {
	qm.Aggregators = append(qm.Aggregators, aggr)
	return qm
}

// TODO: This is just a placeholder. Need to define the Grouper type.
func (qm *qMetric) AddGrouper(tagks []string) QueryMetric {
	qm.GroupBys = append(qm.GroupBys, GroupBy{
		Name: "tag",
		Tags: tagks,
	})
	return qm
}

func (qm *qMetric) SetLimit(limit int) QueryMetric {
	qm.Limit = limit
	return qm
}

func (qm *qMetric) SetOrder(order OrderType) QueryMetric {
	qm.Order = order
	return qm
}

func (qm *qMetric) Validate() error {
	if qm.Name == "" {
		return ErrorQMetricNameInvalid
	}

	for k, v := range qm.Tags {
		if k == "" {
			return ErrorQMetricTagNameInvalid
		} else if len(v) == 0 {
			return ErrorQMetricTagValueInvalid
		}
	}

	if qm.Limit < 0 {
		return ErrorQMetricLimitInvalid
	}

	for _, aggr := range qm.Aggregators {
		err := aggr.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
