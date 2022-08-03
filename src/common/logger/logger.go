package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

type Logger struct {
	lg *logrus.Logger
	fw *lumberjack.Logger

	opts Options
}

// NewLogger 新建Logger
func NewLogger(options ...Option) (*Logger, error) {
	opts := newOptions(options...)

	// 创建日志路径
	if err := os.Mkdir(opts.LogPath, 0755); err != nil {
		fmt.Println("Skipped mkdir error in NewLogger: ", err.Error())
	}

	lg := &Logger{
		lg: logrus.New(),
		fw: &lumberjack.Logger{
			Filename: filepath.Join(opts.LogPath, opts.ServiceName+".log"),

			MaxSize:    opts.LogSize,
			MaxAge:     opts.LogAge,
			MaxBackups: opts.LogBackups,
			Compress:   opts.LogBackupCompress,

			LocalTime: true,
		},

		opts: opts,
	}

	lg.lg.SetOutput(io.MultiWriter(os.Stdout, lg.fw))

	lg.lg.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: opts.TimestampFormat,
	})

	logLvl, err := logrus.ParseLevel(opts.LogLevel)
	if err != nil {
		return nil, err
	}
	lg.lg.SetLevel(logLvl)

	return lg, nil
}

// NewDefaultLogger 新建DefaultLogger
func NewDefaultLogger() *Logger {
	return &Logger{
		lg: logrus.New(),
	}
}

func (l *Logger) Close() {
	if l.fw == nil {
		return
	}
	_ = l.fw.Close()
}
