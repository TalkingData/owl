package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"owl/common/types"

	"github.com/Knetic/govaluate"
)

const (
	MAX_METHOD    = "max"
	MIN_METHOD    = "min"
	RATIO_METHOD  = "ratio"
	TOP_METHOD    = "top"
	BOTTOM_METHOD = "bottom"
	LAST_METHOD   = "last"
	DIFF_METHOD   = "diff"
	NODATA_METHOD = "nodata"
	AVG_METHOD    = "avg"
)

func maxMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}
		values := make([]float64, 0)
		for _, value := range result.Dps {
			values = append(values, value)
		}

		sort.Float64s(values)
		parameters := make(map[string]interface{}, 8)
		current_threshold := values[len(values)-1]
		parameters["current_threshold"] = current_threshold
		parameters["threshold"] = trigger.Threshold
		expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
		trigger_result, err := compute(parameters, expression)
		if err != nil {
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func minMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}
		values := make([]float64, 0)
		for _, value := range result.Dps {
			values = append(values, value)
		}

		sort.Float64s(values)
		parameters := make(map[string]interface{}, 8)
		current_threshold := values[0]
		parameters["current_threshold"] = current_threshold
		parameters["threshold"] = trigger.Threshold
		expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
		trigger_result, err := compute(parameters, expression)
		if err != nil {
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func ratioMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		return nil, err
	}
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}

	start_time, end_time := unit(cycle, trigger.Number)
	results_ago, err := tsdbClient.Query(start_time, end_time, trigger.Tags, "sum", trigger.Metric, false)
	if err != nil {
		return nil, err
	}

	for index, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}
		data := avg(result.Dps)
		if index > len(results_ago)-1 {
			break
		}
		data_ago := avg(results_ago[index].Dps)

		current_threshold := ((data - data_ago) / data_ago) * 100

		parameters := make(map[string]interface{}, 8)
		parameters["current_threshold"] = current_threshold
		parameters["threshold"] = trigger.Threshold
		expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
		trigger_result, err := compute(parameters, expression)
		if err != nil {
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func topMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		return nil, err
	}
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}
		values := make([]float64, 0)
		for _, value := range result.Dps {
			values = append(values, value)
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(values)))
		trigger_result := true
		var count int
		var sum float64
		var current_threshold float64
		for _, value := range values {
			if count >= trigger.Number {
				break
			}
			count += 1
			sum += value
			parameters := make(map[string]interface{}, 8)
			parameters["current_threshold"] = value
			parameters["threshold"] = trigger.Threshold
			expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
			one_result, err := compute(parameters, expression)
			if err != nil {
				return trigger_result_set, err
			}
			current_threshold = value
			if !one_result {
				trigger_result = false
				break
			}
		}
		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
			current_threshold = sum / float64(count)
		}
		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func bottomMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		return nil, err
	}
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}
		values := make([]float64, 0)
		for _, value := range result.Dps {
			values = append(values, value)
		}
		sort.Float64s(values)
		trigger_result := true
		var count int
		var sum float64
		var current_threshold float64
		for _, value := range values {
			if count >= trigger.Number {
				break
			}
			count += 1
			sum += value
			parameters := make(map[string]interface{}, 8)
			parameters["current_threshold"] = value
			parameters["threshold"] = trigger.Threshold
			expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
			one_result, err := compute(parameters, expression)
			if err != nil {
				return trigger_result_set, err
			}
			current_threshold = value
			if !one_result {
				trigger_result = false
				break
			}
		}
		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
			current_threshold = sum / float64(count)
		}
		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func lastMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		return nil, err
	}
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}
		value_times := make([]string, 0)
		for value_time := range result.Dps {
			value_times = append(value_times, value_time)
		}
		sort.Strings(value_times)
		values := make([]float64, len(value_times))
		for index, value_time := range value_times {
			values[len(value_times)-(index+1)] = result.Dps[value_time]
		}
		trigger_result := true
		var count int
		var sum float64
		var current_threshold float64
		for _, value := range values {
			if count >= trigger.Number {
				break
			}
			count += 1
			sum += value
			parameters := make(map[string]interface{}, 8)
			parameters["current_threshold"] = value
			parameters["threshold"] = trigger.Threshold
			expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
			one_result, err := compute(parameters, expression)
			if err != nil {
				return trigger_result_set, err
			}
			current_threshold = value
			if !one_result {
				trigger_result = false
				break
			}
		}
		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
			current_threshold = sum / float64(count)
		}
		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func diffMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}

		values := make([]float64, 0)
		for _, value := range result.Dps {
			values = append(values, value)
		}

		var current_threshold float64 = 0

		if len(values) > 1 {
			value_tmp := values[0]
			for _, value := range values[1:] {
				if value != value_tmp {
					current_threshold = 1
					break
				}
			}
		}

		parameters := make(map[string]interface{}, 8)
		parameters["current_threshold"] = current_threshold
		parameters["threshold"] = trigger.Threshold
		expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
		trigger_result, err := compute(parameters, expression)
		if err != nil {
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func nodataMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}
	lg.Debug("query trigger result, trigger:%s, result:%s, error:%v",
		covert2JSONString(trigger), covert2JSONString(results), err)
	var current_threshold float64 = 1

	if len(results) == 0 {
		tags := map[string]string{}
		if len(trigger.Tags) != 0 {
			tags = types.ParseTags(trigger.Tags)
		}
		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults,
			types.NewTriggerResult(trigger.Index, tags, make([]string, 0), current_threshold, true))
		trigger_result_set.Triggered = true
		return trigger_result_set, nil
	}

	for _, result := range results {
		if len(result.Dps) != 0 {
			current_threshold = 0
		}

		parameters := make(map[string]interface{}, 8)
		parameters["current_threshold"] = current_threshold
		parameters["threshold"] = trigger.Threshold
		expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
		triggered, err := compute(parameters, expression)
		if err != nil {
			return nil, err
		}
		if !trigger_result_set.Triggered && triggered {
			trigger_result_set.Triggered = triggered
		}
		if len(result.Tags) == 0 {
			result.Tags = types.ParseTags(trigger.Tags)
		}
		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults,
			types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, triggered))
	}

	return trigger_result_set, nil
}

func avgMethod(cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{TriggerResults: make([]*types.TriggerResult, 0), Triggered: false}

	results, err := tsdbClient.Query(fmt.Sprintf("%d", cycle), "", trigger.Tags, "sum", trigger.Metric, true)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
			lg.Warn("metric:%s, tags:%s %d no data", trigger.Metric, trigger.Tags, cycle)
			continue
		}

		parameters := make(map[string]interface{}, 8)
		current_threshold := avg(result.Dps)
		parameters["current_threshold"] = current_threshold
		parameters["threshold"] = trigger.Threshold
		expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
		trigger_result, err := compute(parameters, expression)
		if err != nil {
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func compute(params map[string]interface{}, express string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(express)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(params)
	if err != nil {
		return false, err
	}
	value, ok := result.(bool)
	if !ok {
		return false, errors.New("result value is not bool type")
	}

	return value, nil
}

func unit(cycle, number int) (start_time, end_time string) {
	start := time.Now().Add(-time.Duration(number) * time.Minute).Add(-time.Duration(cycle) * time.Minute)
	end := time.Now().Add(-time.Duration(number) * time.Minute)

	start_time = start.Format("2006/01/02 15:04:05")
	end_time = end.Format("2006/01/02 15:04:05")

	return start_time, end_time
}

func avg(values map[string]float64) float64 {
	var sum float64
	count := 0
	for _, value := range values {
		sum += value
		count += 1
	}

	return sum / float64(count)
}

func covert2JSONString(s interface{}) string {
	data, _ := json.Marshal(s)
	return string(data)
}
