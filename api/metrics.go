package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type MetricIndex struct {
	ID       int    `json:"id"`
	Metric   string `json:"metric"`
	TagIndex []*TagIndex
}

type TagIndex struct {
	ID   int    `json:"id"`
	Tagk string `json:"tagk"`
	Tagv string `json:"tagv"`
}

func suggestMetric(c *gin.Context) {
	q := c.Query("q")
	metricIndex := []MetricIndex{}
	mydb := mydb.Table("metric_index")
	if len(q) > 0 {
		q = fmt.Sprintf("%%%s%%", q)
		mydb = mydb.Where("metric like ?", q)
	}
	mydb.Find(&metricIndex)
	metrics := make([]string, len(metricIndex))
	for i, m := range metricIndex {
		metrics[i] = m.Metric
	}
	c.JSON(http.StatusOK, metrics)
}

func suggestTagk(c *gin.Context) {
	tags := []TagIndex{}
	mydb.Table("tag_index").Select("distinct tagk").Find(&tags)
	result := make([]string, len(tags))
	for i, t := range tags {
		result[i] = t.Tagk
	}
	c.JSON(http.StatusOK, result)
}

func suggestTagv(c *gin.Context) {
	tagk := c.Query("tagk")
	tags := []TagIndex{}
	mydb := mydb.Table("tag_index")
	if tagk != "" {
		mydb = mydb.Where("tagk = ?", tagk)
	}
	mydb.Find(&tags)
	result := make([]string, len(tags))
	for i, t := range tags {
		result[i] = t.Tagv
	}
	c.JSON(http.StatusOK, result)
}

func BuildMetricAndTagIndex(c *gin.Context) {
	buildMetricAndTagIndex()
	c.String(http.StatusOK, "%s", "ok")
}
func buildMetricAndTagIndex() {
	metrics := []Metric{}
	mydb.Table("metric").Select("distinct name").Find(&metrics)
	metricMap := make(map[string]struct{})
	tagsMap := make(map[string]map[string]struct{})
	for _, m := range metrics {
		var (
			metric string
			tagk   string
			tagv   string
		)
		idx := strings.IndexRune(m.Name, 47) // 47 == /
		if idx == -1 {                       // not found
			metric = m.Name
		} else {
			metric = m.Name[:idx]
		}
		if _, ok := metricMap[metric]; !ok {
			metricMap[metric] = struct{}{}
		}
		if idx != -1 {
			fields := strings.Split(m.Name[idx+1:], ",")
			for _, f := range fields {
				kv := strings.Split(f, "=")
				if len(kv) != 2 {
					continue
				}
				tagk = kv[0]
				tagv = kv[1]
				if _, ok := tagsMap[tagk]; !ok {
					tagsMap[tagk] = make(map[string]struct{})
				}
				if _, ok := tagsMap[tagk][tagv]; !ok {
					tagsMap[tagk][tagv] = struct{}{}
				}
			}
		}
	}
	for metric, _ := range metricMap {
		metricIndex := &MetricIndex{
			Metric: metric,
		}
		mydb.Table("metric_index").Where("metric = ?", metric).FirstOrCreate(&metricIndex)
	}
	for tagk, v := range tagsMap {
		for tagv, _ := range v {
			tagIndex := &TagIndex{
				Tagk: tagk,
				Tagv: tagv,
			}
			mydb.Table("tag_index").Where("tagk = ? and tagv = ?",
				tagIndex.Tagk, tagIndex.Tagv).FirstOrCreate(&tagIndex)

		}
	}
}

func autoBuildMetricAndTagIndex() {
	go func() {
		for {
			buildMetricAndTagIndex()
			time.Sleep(time.Minute * time.Duration(GlobalConfig.AUTO_BUILD_INTERVAL))
		}
	}()
}
