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
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Success test.
func TestMetricBuilder(t *testing.T) {
	// Read in the exmaple JSON ouput from file.
	data, _ := ioutil.ReadFile("../test_resources/multiple_metrics.json")
	// Trim the whitespace/newline.
	str1 := strings.TrimSpace(string(data))

	// Instantiate a builder.
	b := NewMetricBuilder()
	b.AddMetric("metric1").
		AddDataPoint(1, int64(10)).
		AddDataPoint(2, int64(30)).
		AddTag("tag1", "tab1value").
		AddTag("tag2", "tab2value")

	b.AddMetric("metric2").
		AddDataPoint(2, int64(30)).
		AddDataPoint(3, 2.3).
		AddTag("tag3", "tab3value")

	// Get the JSON output from the builder instance.
	s, _ := b.Build()
	str2 := string(s)

	assert.Equal(t, str1, str2, "Builder output & file contents must be equal")
}

// Failure test.
func TestMetricBuilderEmptyMetricName(t *testing.T) {
	b := NewMetricBuilder()
	b.AddMetric("")

	s, err := b.Build()
	assert.Equal(t, ErrorMetricNameInvalid, err, "Invalid metric name error expected")
	assert.Nil(t, s, "Build output must be nil")
}

// Failure test.
func TestMetricBuilderEmptyTagName(t *testing.T) {
	b := NewMetricBuilder()
	b.AddMetric("metric1").AddTag("", "val")

	s, err := b.Build()
	assert.Equal(t, ErrorTagNameInvalid, err, "Invalid tag name error expected")
	assert.Nil(t, s, "Build output must be nil")
}

// Failure test.
func TestMetricBuilderEmptyTagValue(t *testing.T) {
	b := NewMetricBuilder()
	b.AddMetric("metric1").AddTag("tag", "")

	s, err := b.Build()
	assert.Equal(t, ErrorTagValueInvalid, err, "Invalid tag value error expected")
	assert.Nil(t, s, "Build output must be nil")
}

// Success test.
func TestMetricBuilderTSNegative(t *testing.T) {
	b := NewMetricBuilder()
	b.AddMetric("metric1").
		AddTag("tag", "val").
		AddDataPoint(-1, 10)

	// Build should succeed.
	s, err := b.Build()
	assert.Nil(t, err, "Don't expect an error")
	assert.NotNil(t, s, "Should have a non-nil output")
}
