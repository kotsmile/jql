package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
}

type logger struct {
	internal *logrus.Entry
}

func NewTest() Logger {
	return New(true)
}

func New(debug bool) Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	l.SetLevel(logrus.InfoLevel)
	l.SetOutput(os.Stdout)

	if debug {
		l.SetLevel(logrus.DebugLevel)
	}

	return &logger{
		internal: logrus.NewEntry(l),
	}
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{
		internal: l.internal.WithField(key, value),
	}
}

func (l *logger) WithFields(fields Fields) Logger {
	return &logger{
		internal: l.internal.WithFields(logrus.Fields(fields)),
	}
}

func (l *logger) WithError(err error) Logger {
	return &logger{
		internal: l.internal.WithError(err),
	}
}

func (l *logger) Debug(args ...interface{}) {
	l.internal.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.internal.Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.internal.Warn(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.internal.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.internal.Fatal(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.internal.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.internal.Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.internal.Warnf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.internal.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.internal.Fatalf(format, args...)
}
