package gorm_logrus

import (
	"context"
	"errors"
	"time"

	logrusLogger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type WrapLogger struct {
	logrusLogger          *logrusLogger.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewLogger(logger *logrusLogger.Logger) *WrapLogger {
	return &WrapLogger{
		logrusLogger:          logger,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *WrapLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.logrusLogger.SetLevel(logrusLogger.Level(level))
	return l
}

func (l *WrapLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.logrusLogger.WithContext(ctx).Infof(s, args)
}

func (l *WrapLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logrusLogger.WithContext(ctx).Warnf(s, args)
}

func (l *WrapLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.logrusLogger.WithContext(ctx).Errorf(s, args)
}

func (l *WrapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrusLogger.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrusLogger.ErrorKey] = err
		l.logrusLogger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logrusLogger.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	if l.Debug {
		l.logrusLogger.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}
