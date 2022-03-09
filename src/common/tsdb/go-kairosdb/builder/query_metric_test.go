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

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Success test.
func TestQMetric(t *testing.T) {
	testData := `{"tags":{"tag1":"val1"},"name":"qm1","limit":100,"order":"desc"}`
	qm := NewQueryMetric("qm1").AddTag("tag1", "val1").SetLimit(100).SetOrder(DESCENDING)
	err := qm.Validate()

	assert.Nil(t, err, "No error expected")
	j, _ := json.Marshal(qm)
	assert.Equal(t, testData, string(j), "Query Metric json output must match")
}

// Failure test.
func TestQMetricNameEmpty(t *testing.T) {
	err := NewQueryMetric("").Validate()

	assert.Equal(t, ErrorQMetricNameInvalid, err, "Query Metric name cannot be empty")
}

// Failure test.
func TestQMetricTagNameEmpty(t *testing.T) {
	err := NewQueryMetric("qm1").AddTag("", "val").Validate()

	assert.Equal(t, ErrorQMetricTagNameInvalid, err, "Query Metric tag name cannot be empty")
}

// Failure test.
func TestQMetricTagValueEmpty(t *testing.T) {
	err := NewQueryMetric("qm1").AddTag("tag", "").Validate()

	assert.Equal(t, ErrorQMetricTagValueInvalid, err, "Query Metric tag value cannot be empty")
}

// Failure test.
func TestQMetricMapTagNameEmpty(t *testing.T) {
	tm := map[string]string{
		"tag1": "val1",
		"":     "val2",
	}
	err := NewQueryMetric("qm1").AddTags(tm).Validate()

	assert.Equal(t, ErrorQMetricTagNameInvalid, err, "Query Metric tag name cannot be empty")
}

// Failure test.
func TestQMetricMapTagValueEmpty(t *testing.T) {
	tm := map[string]string{
		"tag1": "val1",
		"tag2": "",
	}
	err := NewQueryMetric("qm1").AddTags(tm).Validate()

	assert.Equal(t, ErrorQMetricTagValueInvalid, err, "Query Metric tag value cannot be empty")
}

// Failure test.
func TestQMetricLimitNegative(t *testing.T) {
	err := NewQueryMetric("qm1").AddTag("tag", "val").SetLimit(-1).Validate()

	assert.Equal(t, ErrorQMetricLimitInvalid, err, "Query Metric limit cannot be negative")
}
