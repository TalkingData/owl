package main

import (
	"fmt"
	"log"
	"net/http"
	"owl/common/tsdb"
	"owl/common/types"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var tsdbClient tsdb.TsdbClient

func initTSDB() error {
	var err error
	switch config.TimeSeriesStorage {
	case "opentsdb":
		tsdbClient, err = tsdb.NewOpenTsdbClient(config.OpentsdbAddr, time.Duration(config.OpenttsdbReadTimeout)*time.Second)
	case "kairosdb":
		tsdbClient, err = tsdb.NewKairosDbClient(config.KairosdbAddr)
	default:
		err = fmt.Errorf("%s timeseries storage not support", config.TimeSeriesStorage)
	}
	return err
}

func queryTimeSeriesData(c *gin.Context) {
	response := gin.H{}
	defer c.JSON(http.StatusOK, response)
	metric := c.Query("metric")
	tags := c.Query("tags")
	tagMap := types.ParseTags(tags)
	if groupNames, exist := tagMap["host_group"]; exist {
		productIDStr, ok := c.GetQuery("product_id")
		if !ok {
			response["code"] = http.StatusNotFound
			response["message"] = "product id must provide"
			return
		}
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			response["code"] = http.StatusInternalServerError
			response["message"] = "invalid product id "
			return
		}
		delete(tagMap, "host_group")
		var hostSet []string
		for _, groupName := range strings.Split(groupNames, "|") {
			hostSet = append(hostSet, getHostnameTagsFromProductGroup(productID, groupName)...)
		}
		if len(hostSet) == 0 {
			response["code"] = http.StatusBadRequest
			response["message"] = "all group has no host"
			return
		}
		hosts := strings.Join(hostSet, "|")
		// 如果存在 tag host， merge
		if host, ok := tagMap["host"]; ok {
			hosts = hosts + "|" + host
		}
		tagMap["host"] = hosts
		tags = Tags2String(tagMap)
	}

	start := c.DefaultQuery("start", time.Now().Add(-time.Hour).Format("2006/01/02-15:04:05"))
	end := c.DefaultQuery("end", time.Now().Format("2006/01/02-15:04:05"))
	log.Println("query time series data, metric:", metric, "tags:", tags, "start_time:", start, "end_time:", end)
	result, err := tsdbClient.Query(start, end, tags, "sum", metric, false)
	if err != nil {
		response["message"] = err.Error()
		response["code"] = http.StatusInternalServerError
		return
	}
	response["data"] = result
	response["code"] = http.StatusOK
}

func getHostnameTagsFromProductGroup(productID int, groupName string) []string {
	hostnameSet := []string{}
	hostGroup := mydb.getProductHostGroupByName(productID, groupName)
	if hostGroup.ID == 0 {
		return hostnameSet
	}
	_, hosts := mydb.getProductHostGroupHosts(productID, hostGroup.ID, false, "", "", 0, 0)
	for _, host := range hosts {
		hostnameSet = append(hostnameSet, host.Hostname)
	}
	return hostnameSet
}

func Tags2String(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}
	taglen := len(tags)
	keys := make([]string, taglen)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	ret := ""
	for _, k := range keys {
		ret += fmt.Sprintf("%s=%s,", strings.TrimSpace(k), strings.TrimSpace(tags[k]))
	}
	return strings.Trim(ret, ",")
}
