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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQBMetricNameEmpty(t *testing.T) {
	qb := NewQueryBuilder()
	qb.SetRelativeStart(1, "MONTHS").AddMetric("")

	j, err := qb.Build()
	assert.Equal(t, ErrorQMetricNameInvalid, err, "Query Metric name cannot be empty")
	assert.Nil(t, j, "No output expected")
}

func TestQBAbsRelativeStartSet(t *testing.T) {
	qb := NewQueryBuilder()
	qb.SetAbsoluteStart(time.Now()).SetRelativeStart(2, "MONTHS")

	j, err := qb.Build()
	assert.Equal(t, ErrorAbsRelativeStartSet, err, "Both absolute & relative start times cannot be set")
	assert.Nil(t, j, "No output expected")
}

func TestQBStartTimeNotSet(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddMetric("qm1")

	j, err := qb.Build()
	assert.Equal(t, ErrorStartTimeNotSpecified, err, "Start time must be set")
	assert.Nil(t, j, "No output expected")
}

func TestQBRelativeStartDurationZero(t *testing.T) {
	qb := NewQueryBuilder()
	qb.SetRelativeStart(0, "MONTHS")

	j, err := qb.Build()
	assert.Equal(t, ErrorRelativeStartTimeInvalid, err, "Relative start durartion cannot be zero")
	assert.Nil(t, j, "No output expected")
}

func TestQBRelativeStartDurationNeg(t *testing.T) {
	qb := NewQueryBuilder()
	qb.SetRelativeStart(-1, "MONTHS")

	j, err := qb.Build()
	assert.Equal(t, ErrorRelativeStartTimeInvalid, err, "Relative start durartion cannot be negative")
	assert.Nil(t, j, "No output expected")
}

func TestQBAbsRelativeEndSet(t *testing.T) {
	qb := NewQueryBuilder()
	qb.SetAbsoluteEnd(time.Now()).SetRelativeEnd(2, "MONTHS")

	j, err := qb.Build()
	assert.Equal(t, ErrorAbsRelativeEndSet, err, "Both absolute & relative end times cannot be set")
	assert.Nil(t, j, "No output expected")
}

func TestQBRelativeEndDurationZero(t *testing.T) {
	qb := NewQueryBuilder()
	qb.SetRelativeEnd(0, "MONTHS")

	j, err := qb.Build()
	assert.Equal(t, ErrorRelativeEndTimeInvalid, err, "Relative end durartion cannot be zero")
	assert.Nil(t, j, "No output expected")
}

func TestQBRelativeEndDurationNeg(t *testing.T) {
	qb := NewQueryBuilder()
	qb.SetRelativeEnd(-1, "MONTHS")

	j, err := qb.Build()
	assert.Equal(t, ErrorRelativeEndTimeInvalid, err, "Relative end durartion cannot be negative")
	assert.Nil(t, j, "No output expected")
}
