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
	"testing"

	"github.com/stretchr/testify/assert"
)

// Success test.
func TestMetric(t *testing.T) {
	testData := `{"name":"m1","tags":{"tag":"val"},"datapoints":[[123456,1]]}`
	j, err := NewMetric("m1").AddTag("tag", "val").AddDataPoint(123456, 1).Build()

	assert.Nil(t, err, "Dont' expect error")
	assert.Equal(t, string(j), testData, "Metric build output must be same")
}

// Failure test.
func TestEmptyMetricName(t *testing.T) {
	j, err := NewMetric("").Build()

	assert.Nil(t, j, "Metric object must be nil")
	assert.Equal(t, ErrorMetricNameInvalid, err, "Invalid metric name error expected")
}

// Failure test.
func TestEmptyTagName(t *testing.T) {
	j, err := NewMetric("m1").AddTag("", "abc").Build()
	assert.Nil(t, j, "Metric object must be nil")
	assert.Equal(t, ErrorTagNameInvalid, err, "Invalid tag name error expected")
}

// Failure test.
func TestEmptyTagValue(t *testing.T) {
	m, err := NewMetric("m1").AddTag("xyz", "").Build()
	assert.Nil(t, m, "Metric object must be nil")
	assert.Equal(t, ErrorTagValueInvalid, err, "Invalid tag value error expected")
}

// Failure test.
func TestTTLValueNegative(t *testing.T) {
	m, err := NewMetric("m1").AddTTL(-1).Build()
	assert.Nil(t, m, "Metric object must be nil")
	assert.Equal(t, ErrorTTLInvalid, err, "Invalid TTL error expected")
}
