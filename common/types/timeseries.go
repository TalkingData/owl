package types

import (
	"encoding/json"
	"fmt"
	"owl/common/utils"
	"sort"
	"strings"
)

type TimeSeriesData struct {
	Metric    string            `json:"metric"`    //sys.cpu.idle
	DataType  string            `json:"data_type"` //COUNTER,GAUGE,DERIVE
	Value     float64           `json:"value"`     //99.00
	Timestamp int64             `json:"timestamp"` //unix timestamp
	Cycle     int               `json:"cycle"`
	Tags      map[string]string `json:"tags"` //{"product":"app01", "group":"dev02"}
}

func (this *TimeSeriesData) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *TimeSeriesData) Decode(data []byte) error {
	return json.Unmarshal(data, &this)
}

func (this TimeSeriesData) String() string {
	return fmt.Sprintf("{metric:%s, data_type:%s, value:%.2f, time:%d, cycle:%d, tags:%s}",
		this.Metric,
		this.DataType,
		this.Value,
		this.Timestamp,
		this.Cycle,
		this.Tags2String(),
	)
}

func (this *TimeSeriesData) Tags2String() string {
	if len(this.Tags) == 0 {
		return ""
	}
	taglen := len(this.Tags)
	keys := make([]string, taglen)
	i := 0
	for k := range this.Tags {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	ret := ""
	for _, k := range keys {
		ret += fmt.Sprintf("%s=%s,", strings.TrimSpace(k), strings.TrimSpace(this.Tags[k]))
	}
	return strings.Trim(ret, ",")
}

func (this *TimeSeriesData) PK() string {
	return utils.Md5(fmt.Sprintf("%s.%s", this.Metric, this.Tags2String()))
}

func (this *TimeSeriesData) GetMetric() string {
	metric := this.Metric
	if len(this.Tags2String()) > 0 {
		metric = fmt.Sprintf("%s/%s", metric, this.Tags2String())
	}
	return metric
}

type MetricConfig struct {
	HostID     string         `json:"host_id"`
	SeriesData TimeSeriesData `json:"time_series_data"`
}

func (this *MetricConfig) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *MetricConfig) Decode(data []byte) error {
	return json.Unmarshal(data, &this)
}

//tag1=v1,tag2=v2,tag3=v3
//{"tag1":v1,"tag2":v2,"tag3":v3}
func ParseTags(name string) map[string]string {
	res := make(map[string]string)
	kv := strings.Split(name, ",")
	for _, v := range kv {
		tmp := strings.Split(v, "=")
		res[tmp[0]] = tmp[1]
	}
	return res
}
