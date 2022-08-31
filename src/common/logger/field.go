package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type Fields logrus.Fields

// String 将logger.Fields转换为字符串形式
func (f Fields) String(l *Logger) string {
	var data []string
	for k, v := range f {
		data = append(data, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(data, l.opts.LogDataSeparator)
}

// Trans2LogrusFields 将logger.Fields转换为logrus.Fields
func (f Fields) Trans2LogrusFields(l *Logger) logrus.Fields {
	res := logrus.Fields{}

	// 从data中提取code字段，code字段不会出现在最终输入的data中
	if f["code"] != nil {
		res["code"] = f["code"]
		delete(f, "code")
	}

	// 从data中提取error字段，error字段不会出现在最终输入的data中
	if f["error"] != nil {
		switch errType := f["error"].(type) {
		// 如果是error类型，设置最终error字段为err.Error()
		case error:
			res["error"] = errType.Error()
		// 如果是其他类型，设置最终error字段为error（不变）
		default:
			res["error"] = f["error"]
		}
		delete(f, "error")
	}

	// 从data中提取trace_id字段，trace_id字段不会出现在最终输入的data中
	if f["trace_id"] != nil {
		res["trace_id"] = f["trace_id"]
		delete(f, "trace_id")
	}

	// 将data转为字符串形式
	res["data"] = f.String(l)

	return res
}
