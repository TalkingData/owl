package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

func (l *Logger) Debug(args ...interface{}) {
	if l.lg.Level >= logrus.DebugLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		entry.Data["file"] = l.getFileInfo()
		entry.Debug(args...)
	}
}

func (l *Logger) DebugWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.DebugLevel {
		entry := l.lg.WithFields(f.RusFields(l))
		entry.Data["file"] = l.getFileInfo()
		entry.Debug(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.lg.Level >= logrus.InfoLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		entry.Data["file"] = l.getFileInfo()
		entry.Info(args...)
	}
}

func (l *Logger) InfoWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.InfoLevel {
		entry := l.lg.WithFields(f.RusFields(l))
		entry.Data["file"] = l.getFileInfo()
		entry.Info(args...)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	if l.lg.Level >= logrus.WarnLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		entry.Data["file"] = l.getFileInfo()
		entry.Warn(args...)
	}
}

func (l *Logger) WarnWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.WarnLevel {
		entry := l.lg.WithFields(f.RusFields(l))
		entry.Data["file"] = l.getFileInfo()
		entry.Warn(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.lg.Level >= logrus.ErrorLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		entry.Data["file"] = l.getFileInfo()
		entry.Error(args...)
	}
}

func (l *Logger) ErrorWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.ErrorLevel {
		entry := l.lg.WithFields(f.RusFields(l))
		entry.Data["file"] = l.getFileInfo()
		entry.Error(args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.lg.Level >= logrus.FatalLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		entry.Data["file"] = l.getFileInfo()
		entry.Fatal(args...)
	}
}

func (l *Logger) FatalWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.FatalLevel {
		entry := l.lg.WithFields(f.RusFields(l))
		entry.Data["file"] = l.getFileInfo()
		entry.Fatal(args...)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	entry := l.lg.WithFields(logrus.Fields{})
	entry.Data["file"] = l.getFileInfo()
	entry.Panic(args...)
}

func (l *Logger) PanicWithFields(f Fields, args ...interface{}) {
	entry := l.lg.WithFields(f.RusFields(l))
	entry.Data["file"] = l.getFileInfo()
	entry.Panic(args...)
}

func (l *Logger) getFileInfo() string {
	_, file, line, ok := runtime.Caller(l.opts.SkipFiles)
	if !ok {
		return "<Unknown>:-1"
	}

	srcPathIdx := strings.LastIndex(file, l.opts.SourceCodePath)
	if srcPathIdx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return fmt.Sprintf(".%s:%d", file[srcPathIdx:], line)
}
