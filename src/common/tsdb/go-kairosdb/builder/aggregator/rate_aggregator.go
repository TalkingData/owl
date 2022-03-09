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

type rateAggregator struct {
	*basicAggregator
	UnitVal utils.TimeUnit `json:"unit,omitempty"`
}

func NewRateAggregator(unit utils.TimeUnit) *rateAggregator {
	return &rateAggregator{
		basicAggregator: NewBasicAggregator("rate"),
		UnitVal:         unit,
	}
}

func (ra *rateAggregator) Unit() utils.TimeUnit {
	return ra.UnitVal
}

func (ra *rateAggregator) Validate() error {
	err := ra.basicAggregator.Validate()
	if err != nil {
		return err
	}

	return nil
}
