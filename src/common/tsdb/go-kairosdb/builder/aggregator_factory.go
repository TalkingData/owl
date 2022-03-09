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
	"owl/common/tsdb/go-kairosdb/builder/aggregator"
	"owl/common/tsdb/go-kairosdb/builder/utils"
)

type TrimType string

const (
	TRIM_FIRST TrimType = "first"
	TRIM_LAST  TrimType = "last"
	TRIM_BOTH  TrimType = "both"
)

// Creates an aggregator that returns the minimum values for each time period as specified.
// For example, "5 minutes" would returns the minimum value for each 5 minute period.
//
// @param value value for time period.
// @param unit unit of time
// @return min aggregator
func CreateMinAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("min", value, unit)
}

// Creates an aggregator that returns the maximum values for each time period as specified.
// For example, "5 minutes" would returns the maximum value for each 5 minute period.
//
// @param value value for time period.
// @param unit unit of time
// @return max aggregator
func CreateMaxAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("max", value, unit)
}

// Creates an aggregator that returns the average values for each time period as specified.
// For example, "5 minutes" would returns the average value for each 5 minute period.
//
// @param value value for time period.
// @param unit unit of time
// @return average aggregator
func CreateAverageAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("avg", value, unit)
}

// Creates an aggregator that returns the standard deviation values for each time period
// as specified. For example, "5 minutes" would returns the standard deviation value for
// each 5 minute period.
//
// @param value value for time period.
// @param unit unit of time
// @return standard deviation aggregator
func CreateStandardDeviationAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("dev", value, unit)
}

// Creates an aggregator that returns the sum of all values over each time period as
// specified. For example, "5 minutes" would returns the sum value for each 5 minute
// period.
//
// @param value value for time period.
// @param unit unit of time
// @return sum aggregator
func CreateSumAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("sum", value, unit)
}

// Creates an aggregator that returns the count of all values for each time period as
// specified. For example, "5 minutes" would returns the count of data points for each
// 5 minute period.
//
// @param value value for time period.
// @param unit unit of time
// @return count aggregator
func CreateCountAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("count", value, unit)
}

// Creates an aggregator that returns the last data point for the time range.
//
// @param value value for time period.
// @param unit unit of time
// @return last aggregator
func CreateLastAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("last", value, unit)
}

// Creates an aggregator that returns the first data point for the time range.
//
// @param value value for time period.
// @param unit unit of time
// @return first aggregator
func CreateFirstAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("first", value, unit)
}

// Creates an aggregator that marks gaps in data according to sampling rate with a null
// data point.
//
// @param value value for time period.
// @param unit unit of time
// @return gap marking aggregator
func CreateDataGapsMarkingAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("gaps", value, unit)
}

// Creates an aggregator that returns a best fit line through the datapoints using the
// least squares algorithm..
//
// @param value value for time period.
// @param unit unit of time
// @return least squares aggregator
func CreateLeastSquaresAggregator(value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewSamplingAggregator("least_squares", value, unit)
}

// Creates an aggregator that returns the percentile value for a given percentage of all
// values over each time period as specified. For example, "0.5" and "5 minutes" would
// returns the median of data points for each 5 minute period.
//
// @param value percentage
// @param unit unit of time
// @return percentile aggregator
func CreatePercentileAggregator(percentile float64, value int, unit utils.TimeUnit) Aggregator {
	return aggregator.NewPercentileAggregator(percentile, value, unit)
}

// Creates an aggregator that computes the difference between successive data points.
//
// @return diff aggregator
func CreateDiffAggregator() Aggregator {
	return aggregator.NewBasicAggregator("diff")
}

// Creates an aggregator that computes the sampling rate of change for the data points.
//
// @return sampler aggregator
func CreateSamplerAggregator() Aggregator {
	return aggregator.NewBasicAggregator("sampler")
}

// Creates an aggregator that returns the rate of change between each pair of data points
//
// @param unit unit of time
// @return rate aggregator
func CreateRateAggregator(unit utils.TimeUnit) Aggregator {
	return aggregator.NewRateAggregator(unit)
}

// Creates an aggregator that divides each value by the divisor.
//
// @param divisor divisor.
// @return div aggregator
func CreateDivAggregator(divisor float64) Aggregator {
	m := make(map[string]interface{})
	m["name"] = "div"
	m["divisor"] = divisor
	return aggregator.NewCustomAggregator(m)
}

// Creates an aggregator that scales each data point by a factor.
//
// @param factor factor to scale by
// @return sampler aggregator
func CreateScaleAggregator(factor float64) Aggregator {
	m := make(map[string]interface{})
	m["name"] = "scale"
	m["factor"] = factor
	return aggregator.NewCustomAggregator(m)
}

// Creates an aggregator that saves the results of the query to a new metric.
//
// @param newMetricName metric to save results to
// @return save as aggregator
func CreateSaveAsAggregator(newMetricName string) Aggregator {
	m := make(map[string]interface{})
	m["name"] = "save_as"
	m["metric_name"] = newMetricName
	return aggregator.NewCustomAggregator(m)
}

// Creates an aggregator that trim of the first, last, or both data points returned by
// the query.
//
// @param trim what to trim
// @return trim aggregator
func CreateTrimAggregator(trim TrimType) Aggregator {
	m := make(map[string]interface{})
	m["name"] = "trim"
	m["trim"] = trim
	return aggregator.NewCustomAggregator(m)
}
