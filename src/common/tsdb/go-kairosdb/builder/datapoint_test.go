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

func TestTimeStampNegValue(t *testing.T) {
	dp := NewDataPoint(-100, 3)

	assert.Equal(t, int64(-100), dp.Timestamp(), "Got incorrect timestamp")
}

func TestTimeStampZeroValue(t *testing.T) {
	dp := NewDataPoint(0, 3)

	assert.Equal(t, int64(0), dp.Timestamp(), "Got incorrect timestamp")
}

func TestDataPointInt64Value(t *testing.T) {
	dp := NewDataPoint(12345678, int64(1024))

	assert.Equal(t, int64(12345678), dp.Timestamp(), "Got incorrect timestamp")

	val, err := dp.Int64Value()
	assert.Nil(t, err, "Didn't expect an error")
	assert.Equal(t, int64(1024), val, "Got different value")

	_, err = dp.Float64Value()
	assert.Equal(t, ErrorDataPointFloat64, err, "Expecting an error")
}

func TestDataPointFloat64Value(t *testing.T) {
	dp := NewDataPoint(12345678, float64(1024.00))

	assert.Equal(t, int64(12345678), dp.Timestamp(), "Got incorrect timestamp")

	val, err := dp.Float64Value()
	assert.Nil(t, err, "Didn't expect an error")
	assert.Equal(t, float64(1024.00), val, "Got different value")

	_, err = dp.Int64Value()
	assert.Equal(t, ErrorDataPointInt64, err, "Expecting an error")
}

func TestDataPointInvalidValue(t *testing.T) {
	dp := NewDataPoint(12345678, "abc")

	_, err := dp.Int64Value()
	assert.Equal(t, ErrorDataPointInt64, err, "Expecting an error")

	_, err = dp.Float64Value()
	assert.Equal(t, ErrorDataPointFloat64, err, "Expecting an error")
}
