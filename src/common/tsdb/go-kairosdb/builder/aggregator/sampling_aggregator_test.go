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

// Failure test.
func TestSamplingAggrEmptyName(t *testing.T) {
	sa := NewSamplingAggregator("", 1, utils.HOURS).SetStartTimeAlignment(100)
	err := sa.Validate()
	assert.Equal(t, ErrorAggrNameInvalid, err, "Invalid aggregator name error expected")
}

// Failure test.
func TestSamplingAggrNegValue(t *testing.T) {
	sa := NewSamplingAggregator("test", -1, utils.HOURS).SetStartTimeAlignment(100)
	err := sa.Validate()
	assert.Equal(t, ErrorSamplingAggrValueInvalid, err, "Sampling aggregator value cannot be -ve")
}

// Failure test.
func TestSamplingAggrZeroValue(t *testing.T) {
	sa := NewSamplingAggregator("test", 0, utils.HOURS).SetStartTimeAlignment(100)
	err := sa.Validate()
	assert.Equal(t, ErrorSamplingAggrValueInvalid, err, "Sampling aggregator value cannot be 0")
}

// Failure test.
func TestSamplingAggrNegStartTime(t *testing.T) {
	sa := NewSamplingAggregator("test", 1, utils.HOURS).SetStartTimeAlignment(-100)
	err := sa.Validate()
	assert.Equal(t, ErrorSamplingAggrStartTimeInvalid, err, "Sampling aggregator start time cannot be -ve")
}

// Success test.
func TestSamplingAggrSamplingAlign(t *testing.T) {
	sa := NewSamplingAggregator("test", 1, utils.HOURS).SetSamplingAlignment()
	assert.EqualValues(t, "test", sa.Name(), "Sampling aggregator name must be 'test'")
	assert.True(t, sa.AlignSampling(), "Sampling Alignment must be true")
	assert.False(t, sa.AlignStartTime(), "Sampling Align Start Time must be false")
}

// Success test.
func TestSamplingAggrStartTimeAlignOnly(t *testing.T) {
	sa := NewSamplingAggregator("test", 1, utils.HOURS).SetStartTimeAlignmentOnly()
	assert.True(t, sa.AlignStartTime(), "Sampling Align Start Time must be true")
	assert.False(t, sa.AlignSampling(), "Sampling Alignment must be false")
}

// Success test.
func TestSamplingAggrStartTimeAlign(t *testing.T) {
	sa := NewSamplingAggregator("test", 1, utils.HOURS).SetStartTimeAlignment(int64(123))
	assert.True(t, sa.AlignStartTime(), "Sampling Align Start Time must be true")
	assert.False(t, sa.AlignSampling(), "Sampling Alignment must be false")
	assert.Equal(t, int64(123), sa.StartTime(), "Sampling Alignment Start Time value is not the same")
}
