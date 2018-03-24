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

import "encoding/json"

type customAggregator struct {
	KeyVal map[string]interface{}
}

func NewCustomAggregator(kv map[string]interface{}) *customAggregator {
	return &customAggregator{
		KeyVal: kv,
	}
}

func (ca *customAggregator) Name() string {
	name, ok := ca.KeyVal["name"].(string)
	if !ok {
		return ""
	} else {
		return name
	}
}

func (ca *customAggregator) Validate() error {
	return nil
}

func (ca *customAggregator) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.KeyVal)
}
