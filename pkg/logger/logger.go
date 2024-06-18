package logger

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
	TraceWrap(error) logrus.Fields
}

type Log struct {
	*logrus.Logger
}

func (l *Log) TraceWrap(err error) logrus.Fields {
	st := errors.WithStack(err)
	return logrus.Fields{
		"trace": fmt.Sprintf("%+v", st),
	}
}

func NewLogger(fmter logrus.Formatter, hooks []logrus.Hook) Logger {
	log := &Log{
		logrus.New(),
	}
	log.Formatter = fmter
	for _, h := range hooks {
		log.AddHook(h)
	}

	return log
}
