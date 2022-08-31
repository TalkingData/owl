package dto

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"owl/common/utils"
	proxyProto "owl/proxy/proto"
	"regexp"
	"sort"
	"strings"
)

const (
	TsDataTypeGauge   = "GAUGE"
	TsDataTypeCounter = "COUNTER"
	TsDataTypeDerive  = "DERIVE"
)

var (
	RegMetric   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]+$`)
	RegTagKey   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]+$`)
	RegTagValue = regexp.MustCompile(`[a-zA-Z0-9_.-/]+$`)

	AllowedTsDataType = map[string]string{
		TsDataTypeGauge:   "",
		TsDataTypeCounter: "",
		TsDataTypeDerive:  "",
	}
)

type TsData struct {
	Metric    string            `json:"metric" binding:"required"`
	DataType  string            `json:"data_type" binding:"required"`
	Value     float64           `json:"value" binding:"required"`
	Timestamp int64             `json:"timestamp"`
	Cycle     int32             `json:"cycle,omitempty" binding:"required"`
	Tags      map[string]string `json:"tags"`
}

func NewTsData(metric, dataType string, val float64, ts int64, cycle int32, tags map[string]string) *TsData {
	return &TsData{
		Metric:    metric,
		DataType:  strings.ToUpper(dataType),
		Value:     val,
		Timestamp: utils.AlignTimestamp(ts, cycle),
		Cycle:     cycle,
		Tags:      tags,
	}
}

func (tsData *TsData) Encode() []byte {
	data, _ := json.Marshal(tsData)
	return data
}

func (tsData *TsData) Validate() error {
	tsData.arrange()

	if !RegMetric.MatchString(tsData.Metric) || tsData.Metric == "" {
		return fmt.Errorf("invalid metric %s, must complie %s", tsData.Metric, RegMetric)
	}

	if _, ok := AllowedTsDataType[tsData.DataType]; !ok {
		return fmt.Errorf("invalid data type %s, only allowed [GAUGE, COUNTER, DERIVE]", tsData.DataType)
	}

	for tKey, tVal := range tsData.Tags {
		if !RegTagKey.MatchString(tKey) {
			return fmt.Errorf("invalid tag key %s, must complie %s", tKey, RegTagKey)
		}
		if !RegTagValue.MatchString(tVal) {
			return fmt.Errorf("invalid tag value %s, must complie %s", tVal, RegTagValue)
		}
	}

	return nil
}

// arrange 整理数据
func (tsData *TsData) arrange() {
	tsData.DataType = strings.ToUpper(tsData.DataType)
	tsData.Timestamp = utils.AlignTimestamp(tsData.Timestamp, tsData.Cycle)
	// 使tsData.Value只保留两位小数
	tsData.Value, _ = decimal.NewFromFloat(tsData.Value).Round(2).Float64()
}

// MergeTags 依照map[string]string合并tag，Key相同的会被覆盖
func (tsData *TsData) MergeTags(tags map[string]string) {
	if tsData.Tags == nil {
		tsData.Tags = tags
		return
	}

	for k, v := range tags {
		tsData.Tags[k] = v
	}
}

// PutTag 增加tag，Key相同的会被覆盖
func (tsData *TsData) PutTag(key, val string) {
	tsData.Tags[key] = val
}

func (tsData *TsData) GetPk() string {
	return fmt.Sprintf("%s.%s", tsData.Metric, tsData.Tags2Str())
}

func (tsData *TsData) Tags2Str() (res string) {
	if len(tsData.Tags) == 0 {
		return
	}

	keyArr := []string{}
	for k := range tsData.Tags {
		tagStr := fmt.Sprintf("%s=%s", strings.TrimSpace(k), strings.TrimSpace(tsData.Tags[k]))
		keyArr = append(keyArr, tagStr)
	}

	sort.Strings(keyArr)
	return strings.Join(keyArr, ",")
}

func (tsData *TsData) Trans2ProxyTsData() *proxyProto.TsData {
	return &proxyProto.TsData{
		Metric:    tsData.Metric,
		DataType:  tsData.DataType,
		Value:     tsData.Value,
		Timestamp: tsData.Timestamp,
		Cycle:     tsData.Cycle,
		Tags:      tsData.Tags,
	}
}
