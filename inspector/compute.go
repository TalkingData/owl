package main

import (
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
	NODATA_METHOD = "nodata"
)

func maxMethod(host_id string, cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{make([]*types.TriggerResult, 0), false}

	params := NewQueryParams(host_id, fmt.Sprintf("%dm-ago", cycle), "", trigger.Tags, "sum", trigger.Metric)
	results, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
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
			lg.Error(err.Error())
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func minMethod(host_id string, cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{make([]*types.TriggerResult, 0), false}

	params := NewQueryParams(host_id, fmt.Sprintf("%dm-ago", cycle), "", trigger.Tags, "sum", trigger.Metric)
	results, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
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
			lg.Error(err.Error())
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func ratioMethod(host_id string, cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		lg.Error(err.Error())
		return nil, err
	}

	trigger_result_set := &types.TriggerResultSet{make([]*types.TriggerResult, 0), false}
	params := NewQueryParams(host_id, fmt.Sprintf("%dm-ago", cycle), "", trigger.Tags, "sum", trigger.Metric)
	results, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	start_time, end_time := unit(cycle, trigger.Number)
	params = NewQueryParams(host_id, start_time, end_time, trigger.Tags, "sum", trigger.Metric)
	results_ago, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	for index, result := range results {
		if len(result.Dps) == 0 {
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
			lg.Error(err.Error())
			return trigger_result_set, err
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func topMethod(host_id string, cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		lg.Error(err.Error())
		return nil, err
	}

	trigger_result_set := &types.TriggerResultSet{make([]*types.TriggerResult, 0), false}
	params := NewQueryParams(host_id, fmt.Sprintf("%dm-ago", cycle), "", trigger.Tags, "sum", trigger.Metric)
	results, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
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
			count += 1
			if count > trigger.Number {
				break
			}
			sum += value
			parameters := make(map[string]interface{}, 8)
			parameters["current_threshold"] = value
			parameters["threshold"] = trigger.Threshold
			expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
			one_result, err := compute(parameters, expression)
			if err != nil {
				lg.Error(err.Error())
				continue
			}

			if !one_result {
				trigger_result = false
				current_threshold = value
			}
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
			current_threshold = sum / float64(trigger.Number)
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func bottomMethod(host_id string, cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		lg.Error(err.Error())
		return nil, err
	}

	trigger_result_set := &types.TriggerResultSet{make([]*types.TriggerResult, 0), false}
	params := NewQueryParams(host_id, fmt.Sprintf("%dm-ago", cycle), "", trigger.Tags, "sum", trigger.Metric)
	results, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
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
			count += 1
			if count > trigger.Number {
				break
			}
			sum += value
			parameters := make(map[string]interface{}, 8)
			parameters["current_threshold"] = value
			parameters["threshold"] = trigger.Threshold
			expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
			one_result, err := compute(parameters, expression)
			if err != nil {
				lg.Error(err.Error())
				continue
			}

			if !one_result {
				trigger_result = false
				current_threshold = value
			}
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
			current_threshold = sum / float64(trigger.Number)
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func lastMethod(host_id string, cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	if trigger.Number == 0 {
		err := errors.New("number can not be 0")
		lg.Error(err.Error())
		return nil, err
	}

	trigger_result_set := &types.TriggerResultSet{make([]*types.TriggerResult, 0), false}
	params := NewQueryParams(host_id, fmt.Sprintf("%dm-ago", cycle), "", trigger.Tags, "sum", trigger.Metric)
	results, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	for _, result := range results {
		if len(result.Dps) == 0 {
			continue
		}

		value_times := make([]string, 0)
		for value_time, _ := range result.Dps {
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
			count += 1
			if count > trigger.Number {
				break
			}
			sum += value
			parameters := make(map[string]interface{}, 8)
			parameters["current_threshold"] = value
			parameters["threshold"] = trigger.Threshold
			expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
			one_result, err := compute(parameters, expression)
			if err != nil {
				lg.Error(err.Error())
				continue
			}

			if !one_result {
				trigger_result = false
				current_threshold = value
			}
		}

		if !trigger_result_set.Triggered && trigger_result {
			trigger_result_set.Triggered = trigger_result
			current_threshold = sum / float64(trigger.Number)
		}

		trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, result.Tags, result.AggregateTags, current_threshold, trigger_result))
	}

	return trigger_result_set, nil
}

func nodataMethod(host_id string, cycle int, trigger *types.Trigger) (*types.TriggerResultSet, error) {
	trigger_result_set := &types.TriggerResultSet{make([]*types.TriggerResult, 0), false}

	params := NewQueryParams(host_id, fmt.Sprintf("%dm-ago", cycle), "", trigger.Tags, "sum", trigger.Metric)
	results, err := tsdbClient.Query(params)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	var current_threshold float64 = 1

	if len(results) != 0 && len(results[0].Dps) != 0 {
		current_threshold = 0
	}

	parameters := make(map[string]interface{}, 8)
	parameters["current_threshold"] = current_threshold
	parameters["threshold"] = trigger.Threshold
	expression := fmt.Sprintf("current_threshold %s threshold", trigger.Symbol)
	trigger_result, err := compute(parameters, expression)
	if err != nil {
		return nil, err
	}

	trigger_result_set.Triggered = trigger_result

	trigger_result_set.TriggerResults = append(trigger_result_set.TriggerResults, types.NewTriggerResult(trigger.Index, make(map[string]string), make([]string, 0), current_threshold, trigger_result))

	return trigger_result_set, nil
}

func compute(params map[string]interface{}, express string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(express)
	if err != nil {
		lg.Error("Parse the expression error %s", err.Error())
		return false, err
	}

	result, err := expression.Evaluate(params)
	if err != nil {
		lg.Error("Evaluate the expression error %s", err.Error())
		return false, err
	}
	value, ok := result.(bool)
	if !ok {
		lg.Error("result value is not bool type")
	}

	return value, nil
}

func unit(cycle, number int) (start_time, end_time string) {
	start := time.Now().Add(-time.Duration(time.Duration(number) * 24 * time.Hour)).Add(-time.Duration(time.Duration(cycle) * time.Minute))
	end := time.Now().Add(-time.Duration(time.Duration(number) * 24 * time.Hour))

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
