package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.]+$`)

type TimeSeriesData struct {
	Metric    string            `json:"metric"`    //sys.cpu.idle
	DataType  string            `json:"data_type"` //COUNTER,GAUGE,DERIVE
	Value     float64           `json:"value"`     //99.00
	Timestamp int64             `json:"timestamp"` //unix timestamp
	Cycle     int               `json:"cycle,omitempty"`
	Tags      map[string]string `json:"tags"` //{"product":"app01", "group":"dev02"}
}

func (m *TimeSeriesData) Validate() bool {
	if !reg.MatchString(m.Metric) || m.Metric == "" {
		return false
	}
	switch strings.ToLower(m.DataType) {
	case "gauge", "counter", "derive":
	default:
		return false
	}
	return true
}

func (tsd *TimeSeriesData) Encode() []byte {
	data, _ := json.Marshal(tsd)
	return data
}

func (tsd *TimeSeriesData) Decode(data []byte) error {
	return json.Unmarshal(data, &tsd)
}

func (tsd TimeSeriesData) String() string {
	return fmt.Sprintf("{metric:%s, data_type:%s, value:%.2f, time:%d, cycle:%d, tags:%s}",
		tsd.Metric,
		tsd.DataType,
		tsd.Value,
		tsd.Timestamp,
		tsd.Cycle,
		tsd.Tags2String(),
	)
}

func (tsd *TimeSeriesData) Tags2String() string {
	if len(tsd.Tags) == 0 {
		return ""
	}
	taglen := len(tsd.Tags)
	keys := make([]string, taglen)
	i := 0
	for k := range tsd.Tags {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	ret := ""
	for _, k := range keys {
		ret += fmt.Sprintf("%s=%s,", strings.TrimSpace(k), strings.TrimSpace(tsd.Tags[k]))
	}
	return strings.Trim(ret, ",")
}

func (tsd *TimeSeriesData) PK() string {
	return fmt.Sprintf("%s.%s", tsd.Metric, tsd.Tags2String())
}

func (tsd *TimeSeriesData) GetMetric() string {
	metric := tsd.Metric
	if len(tsd.Tags2String()) > 0 {
		metric = fmt.Sprintf("%s/%s", metric, tsd.Tags2String())
	}
	return metric
}

func (tsd *TimeSeriesData) AddTags(tags map[string]string) {
	if tsd.Tags == nil {
		tsd.Tags = tags
		return
	}
	for k, v := range tags {
		tsd.Tags[k] = v
	}
}

//tag1=v1,tag2=v2,tag3=v3
//{"tag1":v1,"tag2":v2,"tag3":v3}
func ParseTags(name string) map[string]string {
	res := make(map[string]string)
	kv := strings.Split(name, ",")
	for _, v := range kv {
		tmp := strings.Split(v, "=")
		if len(tmp) != 2 {
			continue
		}
		res[tmp[0]] = tmp[1]
	}
	return res
}
