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
	"encoding/json"
	"time"

	"owl/common/tsdb/go-kairosdb/builder/utils"
)

type QueryBuilder interface {
	// The beginning time in the time range.
	SetAbsoluteStart(date time.Time) QueryBuilder

	// The beginning time of the time range relative to now.
	SetRelativeStart(duration int, unit utils.TimeUnit) QueryBuilder

	// The ending value of the time range. Must be later in time than the
	// start time. An end time is not required and default to now.
	SetAbsoluteEnd(date time.Time) QueryBuilder

	// The ending time of the time range relative to now.
	SetRelativeEnd(duration int, unit utils.TimeUnit) QueryBuilder

	// How long to cache this exact query. The default is to never cache.
	SetCacheTime(cacheTimeMs int) QueryBuilder

	// The metric to query for.
	AddMetric(name string) QueryMetric

	// Returns the absolute range start time.
	AbsoluteStart() time.Time

	// Returns the relative range start time.
	RelativeStart() *utils.RelativeTime

	// Returns the absolute range end time.
	AbsoluteEnd() time.Time

	// Returns the relative range end time.
	RelativeEnd() *utils.RelativeTime

	// Returns the Cache time.
	CacheTime() int

	// Returns array of metrics.
	Metrics() []QueryMetric

	// Encodes the QueryBuilder into JSON.
	Build() ([]byte, error)
}

// Type that implements the QueryBuilder interface.v
type qBuilder struct {
	StartAbs    int64               `json:"start_absolute,omitempty"`
	EndAbs      int64               `json:"end_absolute,omitempty"`
	StartRel    *utils.RelativeTime `json:"start_relative,omitempty"`
	EndRel      *utils.RelativeTime `json:"end_relative,omitempty"`
	CacheTimeMs int                 `json:"cache_time,omitempty"`
	MetricsArr  []QueryMetric       `json:"metrics,omitempty"`
}

func NewQueryBuilder() QueryBuilder {
	return &qBuilder{
		MetricsArr: make([]QueryMetric, 0),
	}
}

func (qb *qBuilder) timeInMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func (qb *qBuilder) SetAbsoluteStart(date time.Time) QueryBuilder {
	qb.StartAbs = qb.timeInMs(date)
	return qb
}

func (qb *qBuilder) SetRelativeStart(duration int, unit utils.TimeUnit) QueryBuilder {
	qb.StartRel = utils.NewRelativeTime(duration, unit)
	return qb
}

func (qb *qBuilder) SetAbsoluteEnd(date time.Time) QueryBuilder {
	qb.EndAbs = qb.timeInMs(date)
	return qb
}

func (qb *qBuilder) SetRelativeEnd(duration int, unit utils.TimeUnit) QueryBuilder {
	qb.EndRel = utils.NewRelativeTime(duration, unit)
	return qb
}

func (qb *qBuilder) SetCacheTime(cacheTimeMs int) QueryBuilder {
	qb.CacheTimeMs = cacheTimeMs
	return qb
}

func (qb *qBuilder) AddMetric(name string) QueryMetric {
	qm := NewQueryMetric(name)
	qb.MetricsArr = append(qb.MetricsArr, qm)
	return qm
}

func (qb *qBuilder) AbsoluteStart() time.Time {
	return time.Unix(0, qb.StartAbs*int64(time.Millisecond))
}

func (qb *qBuilder) RelativeStart() *utils.RelativeTime {
	return qb.StartRel
}

func (qb *qBuilder) AbsoluteEnd() time.Time {
	return time.Unix(0, qb.EndAbs*int64(time.Millisecond))
}

func (qb *qBuilder) RelativeEnd() *utils.RelativeTime {
	return qb.EndRel
}

func (qb *qBuilder) CacheTime() int {
	return qb.CacheTimeMs
}

func (qb *qBuilder) Metrics() []QueryMetric {
	return qb.MetricsArr
}

func (qb *qBuilder) Build() ([]byte, error) {
	if qb.StartAbs != 0 && qb.StartRel != nil {
		return nil, ErrorAbsRelativeStartSet
	}

	if qb.StartRel != nil && qb.StartRel.Value() <= 0 {
		return nil, ErrorRelativeStartTimeInvalid
	}

	if qb.EndAbs != 0 && qb.EndRel != nil {
		return nil, ErrorAbsRelativeEndSet
	}

	if qb.EndRel != nil && qb.EndRel.Value() <= 0 {
		return nil, ErrorRelativeEndTimeInvalid
	}

	if qb.StartAbs == 0 && qb.StartRel == nil {
		return nil, ErrorStartTimeNotSpecified
	}

	for _, qm := range qb.MetricsArr {
		err := qm.Validate()
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(qb)
}
