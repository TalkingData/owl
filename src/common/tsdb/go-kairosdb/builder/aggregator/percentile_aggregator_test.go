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

package aggregator

import (
	"testing"

	"owl/common/tsdb/go-kairosdb/builder/utils"

	"github.com/stretchr/testify/assert"
)

// Success test.
func TestPercentileAggr(t *testing.T) {
	pa := NewPercentileAggregator(0.5, 20, utils.MINUTES)
	err := pa.Validate()

	assert.Nil(t, err, "No error expected")
	assert.EqualValues(t, "percentile", pa.Name(), "Percentile aggregator's name must be set to 'percentile'")
	assert.EqualValues(t, 0.5, pa.Percentile(), "Percentile aggregator percentile value must be set to 0.5")
	assert.EqualValues(t, 20, pa.Value(), "Percentile aggregator time value must be set 20")
	assert.EqualValues(t, utils.MINUTES, pa.Unit(), "Percentile aggregator time unit must be set 'minutes'")
}

// Failure test.
func TestPercentileAggrZeroPercentile(t *testing.T) {
	pa := NewPercentileAggregator(0.0, 100, utils.MINUTES)
	err := pa.Validate()

	assert.Equal(t, ErrorPercentileInvalid, err, "Percentile aggregator percentile must be > 0 and <= 1")
}

// Failure test.
func TestPercentileAggrNegPercentile(t *testing.T) {
	pa := NewPercentileAggregator(-1.0, 100, utils.MINUTES)
	err := pa.Validate()

	assert.Equal(t, ErrorPercentileInvalid, err, "Percentile aggregator percentile must be > 0 and <= 1")
}

// Failure test.
func TestPercentileAggrHigherThanOnePercentile(t *testing.T) {
	pa := NewPercentileAggregator(1.5, 100, utils.MINUTES)
	err := pa.Validate()

	assert.Equal(t, ErrorPercentileInvalid, err, "Percentile aggregator percentile must be > 0 and <= 1")
}

// Failure test.
func TestPercentileAggrZeroValue(t *testing.T) {
	pa := NewPercentileAggregator(0.5, 0, utils.MINUTES)
	err := pa.Validate()

	assert.Equal(t, ErrorSamplingAggrValueInvalid, err, "Percentile aggregator time value must be >= 0")
}

// Failure test.
func TestPercentileAggrNegValue(t *testing.T) {
	pa := NewPercentileAggregator(0.5, -1, utils.MINUTES)
	err := pa.Validate()

	assert.Equal(t, ErrorSamplingAggrValueInvalid, err, "Percentile aggregator time value must be >= 0")
}
