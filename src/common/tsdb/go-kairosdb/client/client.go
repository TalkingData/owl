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

package client

import (
	"owl/common/tsdb/go-kairosdb/builder"

	"owl/common/tsdb/go-kairosdb/response"
)

type Client interface {
	// Returns a list of all metrics names.
	GetMetricNames() (*response.GetResponse, error)

	// Returns a list of all tag names.
	GetTagNames() (*response.GetResponse, error)

	// Returns a list of all tag values.
	GetTagValues() (*response.GetResponse, error)

	// Queries KairosDB using the query built using builder.
	Query(qb builder.QueryBuilder) (*response.QueryResponse, error)

	// Sends metrics from the builder to the KairosDB server.
	PushMetrics(mb builder.MetricBuilder) (*response.Response, error)

	// Deletes a metric. This is the metric and all its datapoints.
	DeleteMetric(name string) (*response.Response, error)

	// Deletes data in KairosDB using the query built by the builder.
	Delete(builder builder.QueryBuilder) (*response.Response, error)

	// Checks the health of the KairosDB Server.
	HealthCheck() (*response.Response, error)
}
