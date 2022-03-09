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

package utils

import "time"

type RelativeTime struct {
	RTvalue int      `json:"value,omitempty"`
	RTunit  TimeUnit `json:"unit,omitempty"`
}

func NewRelativeTime(value int, unit TimeUnit) *RelativeTime {
	return &RelativeTime{
		RTvalue: value,
		RTunit:  unit,
	}
}

func (rt *RelativeTime) Value() int {
	return rt.RTvalue
}

func (rt *RelativeTime) Unit() TimeUnit {
	return rt.RTunit
}

func (rt *RelativeTime) RelativeTimeTo(t time.Time) time.Time {
	var newTime time.Time

	switch rt.RTunit {
	case YEARS:
		newTime = t.AddDate(-rt.RTvalue, 0, 0)
	case MONTHS:
		newTime = t.AddDate(0, -rt.RTvalue, 0)
	case WEEKS:
		days := rt.RTvalue * 7
		newTime = t.AddDate(0, 0, -days)
	case DAYS:
		newTime = t.AddDate(0, 0, -rt.RTvalue)
	case HOURS:
		newTime = t.Add(-(time.Duration(rt.RTvalue) * time.Hour))
	case MINUTES:
		newTime = t.Add(-(time.Duration(rt.RTvalue) * time.Minute))
	case SECONDS:
		newTime = t.Add(-(time.Duration(rt.RTvalue) * time.Second))
	}

	return newTime
}
