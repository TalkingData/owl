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
		l.fillingFields(entry.Data)
		entry.Debug(args...)
	}
}

func (l *Logger) DebugWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.DebugLevel {
		entry := l.lg.WithFields(f.Trans2LogrusFields(l))
		l.fillingFields(entry.Data)
		entry.Debug(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.lg.Level >= logrus.InfoLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		l.fillingFields(entry.Data)
		entry.Info(args...)
	}
}

func (l *Logger) InfoWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.InfoLevel {
		entry := l.lg.WithFields(f.Trans2LogrusFields(l))
		l.fillingFields(entry.Data)
		entry.Info(args...)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	if l.lg.Level >= logrus.WarnLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		l.fillingFields(entry.Data)
		entry.Warn(args...)
	}
}

func (l *Logger) WarnWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.WarnLevel {
		entry := l.lg.WithFields(f.Trans2LogrusFields(l))
		l.fillingFields(entry.Data)
		entry.Warn(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.lg.Level >= logrus.ErrorLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		l.fillingFields(entry.Data)
		entry.Error(args...)
	}
}

func (l *Logger) ErrorWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.ErrorLevel {
		entry := l.lg.WithFields(f.Trans2LogrusFields(l))
		l.fillingFields(entry.Data)
		entry.Error(args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.lg.Level >= logrus.FatalLevel {
		entry := l.lg.WithFields(logrus.Fields{})
		l.fillingFields(entry.Data)
		entry.Fatal(args...)
	}
}

func (l *Logger) FatalWithFields(f Fields, args ...interface{}) {
	if l.lg.Level >= logrus.FatalLevel {
		entry := l.lg.WithFields(f.Trans2LogrusFields(l))
		l.fillingFields(entry.Data)
		entry.Fatal(args...)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	entry := l.lg.WithFields(logrus.Fields{})
	l.fillingFields(entry.Data)
	entry.Panic(args...)
}

func (l *Logger) PanicWithFields(f Fields, args ...interface{}) {
	entry := l.lg.WithFields(f.Trans2LogrusFields(l))
	l.fillingFields(entry.Data)
	entry.Panic(args...)
}

func (l *Logger) fillingFields(f logrus.Fields) {
	pc, file, line, ok := runtime.Caller(l.opts.SkipFiles)
	if !ok {
		f["file"] = "<Unknown>:-1"
		f["func"] = "<Unknown>"
		return
	}

	srcPathIdx := strings.LastIndex(file, l.opts.SourceCodePath)
	if srcPathIdx < 0 {
		f["file"] = fmt.Sprintf("%s:%d", file, line)
	} else {
		f["file"] = fmt.Sprintf(".%s:%d", file[srcPathIdx:], line)
	}

	f["func"] = runtime.FuncForPC(pc).Name()
}
