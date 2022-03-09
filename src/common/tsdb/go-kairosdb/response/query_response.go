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

package response

import "owl/common/tsdb/go-kairosdb/builder"

type GroupResult struct {
	Name string `json:"name,omitempty"`
}

type Results struct {
	Name       string              `json:"name,omitempty"`
	DataPoints []builder.DataPoint `json:"values,omitempty"`
	Tags       map[string][]string `json:"tags,omitempty"`
	Group      []GroupResult       `json:"group_by,omitempty"`
}

type Queries struct {
	SampleSize int64     `json:"sample_size,omitempty"`
	ResultsArr []Results `json:"results,omitempty"`
}

type QueryResponse struct {
	*Response
	QueriesArr []Queries `json:"queries",omitempty`
}

func NewQueryResponse(code int) *QueryResponse {
	qr := &QueryResponse{
		Response: &Response{},
	}

	qr.SetStatusCode(code)
	return qr
}
