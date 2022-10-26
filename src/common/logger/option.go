package logger

import (
	"owl/common/global"
	"strings"
)

const (
	defaultFilename          = "undefined"
	defaultFilenameSeparator = "-"

	defaultLogPath  = "./logs"
	defaultLogLevel = "debug"

	defaultLogSize           = 500
	defaultLogAge            = 3
	defaultLogBackups        = 3
	defaultLogBackupCompress = true

	defaultTimestampFormat  = global.TimestampFormat
	defaultLogDataSeparator = "; "
	defaultSourceCodePath   = "/src/"

	defaultSkipFiles = 2
)

type Option func(o *Options)

// Options struct
type Options struct {
	Filename string

	LogPath, LogLevel           string
	LogSize, LogAge, LogBackups int
	LogBackupCompress           bool

	TimestampFormat  string
	LogDataSeparator string
	SourceCodePath   string

	SkipFiles int
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Filename: defaultFilename,

		LogPath:  defaultLogPath,
		LogLevel: defaultLogLevel,

		LogSize:           defaultLogSize,
		LogAge:            defaultLogAge,
		LogBackups:        defaultLogBackups,
		LogBackupCompress: defaultLogBackupCompress,

		TimestampFormat:  defaultTimestampFormat,
		LogDataSeparator: defaultLogDataSeparator,
		SourceCodePath:   defaultSourceCodePath,

		SkipFiles: defaultSkipFiles,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Filename 设置Filename，影响日志文件名，建议设置为服务全名
func Filename(in string, suffix ...string) Option {
	return func(o *Options) {
		fields := []string{in}
		for _, s := range suffix {
			fields = append(fields, s)
		}
		o.Filename = strings.Join(fields, defaultFilenameSeparator)
	}
}

// LogPath 设置LogPath
func LogPath(in string) Option {
	return func(o *Options) {
		o.LogPath = in
	}
}

// LogLevel 设置LogLevel，默认为debug
func LogLevel(in string) Option {
	return func(o *Options) {
		o.LogLevel = in
	}
}

// LogSize 设置LogSize，日志大小到达LogSize(MB)就开始backup
func LogSize(in int) Option {
	return func(o *Options) {
		o.LogSize = in
	}
}

// LogAge 设置LogAge，旧日志保存的最大天数
func LogAge(in int) Option {
	return func(o *Options) {
		o.LogAge = in
	}
}

// LogBackups 设置LogBackups，旧日志保存的最大数量
func LogBackups(in int) Option {
	return func(o *Options) {
		o.LogBackups = in
	}
}

// LogBackupCompress 设置LogBackupCompress，是否对backup的日志进行压缩
func LogBackupCompress(in bool) Option {
	return func(o *Options) {
		o.LogBackupCompress = in
	}
}

// TimestampFormat 设置时间戳格式，影响日志time字段，不建议修改
func TimestampFormat(in string) Option {
	return func(o *Options) {
		o.TimestampFormat = in
	}
}

// LogDataSeparator 设置数据字段分割符，影响日志data字段，不建议修改
func LogDataSeparator(in string) Option {
	return func(o *Options) {
		o.LogDataSeparator = in
	}
}

// SourceCodePath 设置项目起始的相对路径，影响日志path字段，不建议修改
func SourceCodePath(in string) Option {
	return func(o *Options) {
		o.SourceCodePath = in
	}
}

// SkipFiles 设置SkipFiles，影响日志path字段，非常不建议修改
func SkipFiles(in int) Option {
	return func(o *Options) {
		o.SkipFiles = in
	}
}
