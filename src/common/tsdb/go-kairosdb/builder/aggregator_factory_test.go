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
	"fmt"
	"testing"

	"owl/common/tsdb/go-kairosdb/builder/utils"

	"github.com/stretchr/testify/assert"
)

type sampling_aggr_desc struct {
	aggrFunc func(int, utils.TimeUnit) Aggregator
	name     string
}

// Constructing an array of all sampling aggregators so that we can
// run a basic success test on all of them in a loop.
var samplingAggrDescArray = []sampling_aggr_desc{
	{CreateMinAggregator, "min"},
	{CreateMaxAggregator, "max"},
	{CreateAverageAggregator, "avg"},
	{CreateStandardDeviationAggregator, "dev"},
	{CreateSumAggregator, "sum"},
	{CreateCountAggregator, "count"},
	{CreateLastAggregator, "last"},
	{CreateFirstAggregator, "first"},
	{CreateDataGapsMarkingAggregator, "gaps"},
	{CreateLeastSquaresAggregator, "least_squares"},
}

// Success test.
func TestSamplingAggregatorFuncs(t *testing.T) {
	for _, adesc := range samplingAggrDescArray {
		a := adesc.aggrFunc(1, utils.MINUTES)
		err := a.Validate()

		assert.Nil(t, err, "No error expected")
		assert.EqualValues(t, adesc.name, a.Name(), fmt.Sprintf("Aggregator name must be set to '%s'", adesc.name))
	}
}
