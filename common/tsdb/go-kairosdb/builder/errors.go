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

import "errors"

var (
	// Metric Errors.
	ErrorMetricNameInvalid = errors.New("Metric name empty")
	ErrorTagNameInvalid    = errors.New("Tag name empty")
	ErrorTagValueInvalid   = errors.New("Tag value empty")
	ErrorTTLInvalid        = errors.New("TTL value invalid")

	// Data Point Errors.
	ErrorDataPointInt64   = errors.New("Not an int64 data value")
	ErrorDataPointFloat64 = errors.New("Not a float64 data value")

	// Query Metric Errors.
	ErrorQMetricNameInvalid     = errors.New("Query Metric name empty")
	ErrorQMetricTagNameInvalid  = errors.New("Query Metric Tag name empty")
	ErrorQMetricTagValueInvalid = errors.New("Query Metric Tag value empty")
	ErrorQMetricLimitInvalid    = errors.New("Query Metric Limit must be >= 0")

	// Query Builder Errors.
	ErrorAbsRelativeStartSet      = errors.New("Both absolute and relative start times cannot be set")
	ErrorRelativeStartTimeInvalid = errors.New("Relative start time duration must be > 0")
	ErrorAbsRelativeEndSet        = errors.New("Both absolute and relative end times cannot be set")
	ErrorRelativeEndTimeInvalid   = errors.New("Relative end time duration must be > 0")
	ErrorStartTimeNotSpecified    = errors.New("Start time not specified")
)
