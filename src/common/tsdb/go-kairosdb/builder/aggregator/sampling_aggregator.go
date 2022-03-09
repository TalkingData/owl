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

type sampling struct {
	Value int            `json:"value,omitempty"`
	Unit  utils.TimeUnit `json:"unit,omitempty"`
}

type samplingAggregator struct {
	*basicAggregator
	AlignStartTimeBool bool     `json:"align_start_time,omitempty"`
	AlignSamplingBool  bool     `json:"align_sampling,omitempty"`
	StartTimeValue     int64    `json:"start_time,omitempty"`
	Sample             sampling `json:"sampling,omitempty"`
}

func NewSamplingAggregator(name string, value int, unit utils.TimeUnit) *samplingAggregator {
	return &samplingAggregator{
		basicAggregator: NewBasicAggregator(name),
		Sample: sampling{
			Value: value,
			Unit:  unit,
		},
	}
}

// Alignment based on the sampling size. For example if your sample size is either
// milliseconds, seconds, minutes or hours then the start of the range will always
// be at the top of the hour.  The effect of setting this to true is that your data
// will take the same shape when graphed as you refresh the data.
//
// Only one alignment type can be used.
func (sa *samplingAggregator) SetSamplingAlignment() *samplingAggregator {
	sa.AlignSamplingBool = true
	return sa
}

// Alignment based on the aggregation range rather than the value of the first
// data point within that range.
//
// Only one alignment type can be used.
func (sa *samplingAggregator) SetStartTimeAlignmentOnly() *samplingAggregator {
	sa.AlignStartTimeBool = true
	return sa
}

// Alignment that starts based on the specified time. For example, if startTime
// is set to noon today,then alignment starts at noon today.
//
// Only one alignment type can be used.
func (sa *samplingAggregator) SetStartTimeAlignment(startTime int64) *samplingAggregator {
	sa.AlignStartTimeBool = true
	sa.StartTimeValue = startTime
	return sa
}

func (sa *samplingAggregator) AlignSampling() bool {
	return sa.AlignSamplingBool
}

func (sa *samplingAggregator) AlignStartTime() bool {
	return sa.AlignStartTimeBool
}

func (sa *samplingAggregator) StartTime() int64 {
	return sa.StartTimeValue
}

func (sa *samplingAggregator) Value() int {
	return sa.Sample.Value
}

func (sa *samplingAggregator) Unit() utils.TimeUnit {
	return sa.Sample.Unit
}

func (sa *samplingAggregator) Validate() error {
	err := sa.basicAggregator.Validate()
	if err != nil {
		return err
	}

	if sa.Sample.Value <= 0 {
		return ErrorSamplingAggrValueInvalid
	}

	if sa.StartTimeValue < 0 {
		return ErrorSamplingAggrStartTimeInvalid
	}

	return nil
}
