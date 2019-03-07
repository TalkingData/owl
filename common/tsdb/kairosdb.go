package tsdb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"owl/common/tsdb/go-kairosdb/builder"

	"owl/common/tsdb/go-kairosdb/builder/utils"
	"owl/common/tsdb/go-kairosdb/client"
)

type KairosDbClient struct {
	rawClient client.Client
}

func NewKairosDbClient(addr string) (*KairosDbClient, error) {
	if !strings.Contains(addr, "http://") {
		addr = fmt.Sprintf("http://%s", addr)
	}
	cli := client.NewHttpClient(addr)
	_, err := cli.HealthCheck()
	return &KairosDbClient{cli}, err
}

func (k *KairosDbClient) newQueryBuilder(start, end, rawTags, aggregator, metric string, isRelative bool) builder.QueryBuilder {
	tags := make(map[string][]string)
	tagks := []string{}
	if rawTags != "" {
		tagsPairs := strings.Split(rawTags, ",")
		for _, tagPair := range tagsPairs {
			tagKV := strings.Split(tagPair, "=")
			if _, ok := tags[tagKV[0]]; !ok {
				tagks = append(tagks, tagKV[0])
			}
			if tagKV[1] == "*" {
				continue
			}
			tags[tagKV[0]] = append(tags[tagKV[0]], strings.Split(tagKV[1], "|")...)
		}
	}
	qb := builder.NewQueryBuilder()
	if isRelative {
		started, _ := strconv.Atoi(start)
		qb.SetRelativeStart(started, utils.MINUTES)
	} else {
		timeStart, _ := time.ParseInLocation("2006/01/02-15:04:05", start, time.Local)
		timeEnd, _ := time.ParseInLocation("2006/01/02-15:04:05", end, time.Local)
		qb.SetAbsoluteStart(timeStart.Local())
		qb.SetAbsoluteEnd(timeEnd.Local())
	}
	qb.AddMetric(metric).AddTags(tags).AddGrouper(tagks)
	return qb
}

func (k *KairosDbClient) Query(start, end, rawTags, aggregator, metric string, isRelative bool) (results []Result, err error) {
	results = make([]Result, 0)
	qb := k.newQueryBuilder(start, end, rawTags, aggregator, metric, isRelative)
	queryResp, err := k.rawClient.Query(qb)
	if err != nil {
		return
	}
	errs := queryResp.GetErrors()
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "|"))
	}
	if len(queryResp.QueriesArr) > 0 {
		for _, r := range queryResp.QueriesArr[0].ResultsArr {
			result := new(Result)
			result.Tags = make(map[string]string)
			result.Dps = make(map[string]float64)
			result.Metric = r.Name
			for k, v := range r.Tags {
				result.Tags[k] = v[0]
			}
			for _, d := range r.DataPoints {
				value, _ := d.Float64Value()
				result.Dps[fmt.Sprintf("%d", d.Timestamp()/1000)] = value
			}
			results = append(results, *result)
		}
	}
	return
}
