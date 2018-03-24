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

import "owl/common/tsdb/go-kairosdb/builder/utils"

type percentileAggregator struct {
	*samplingAggregator
	PercentileValue float64
}

func NewPercentileAggregator(percentile float64, value int, unit utils.TimeUnit) *percentileAggregator {
	return &percentileAggregator{
		samplingAggregator: NewSamplingAggregator("percentile", value, unit),
		PercentileValue:    percentile,
	}
}

func (pa *percentileAggregator) Percentile() float64 {
	return pa.PercentileValue
}

func (pa *percentileAggregator) Validate() error {
	if err := pa.samplingAggregator.Validate(); err != nil {
		return err
	}

	if pa.PercentileValue <= 0.0 || pa.PercentileValue > 1.0 {
		return ErrorPercentileInvalid
	}

	return nil
}
