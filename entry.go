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

type WrapEntry struct {
	logrusEntry           *logrusLogger.Entry
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewEntry(entry *logrusLogger.Entry) *WrapEntry {
	return &WrapEntry{
		logrusEntry:           entry,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *WrapEntry) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.logrusEntry.Logger.SetLevel(logrusLogger.Level(level))
	return l
}

func (l *WrapEntry) Info(ctx context.Context, s string, args ...interface{}) {
	l.logrusEntry.WithContext(ctx).Infof(s, args)
}

func (l *WrapEntry) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logrusEntry.WithContext(ctx).Warnf(s, args)
}

func (l *WrapEntry) Error(ctx context.Context, s string, args ...interface{}) {
	l.logrusEntry.WithContext(ctx).Errorf(s, args)
}

func (l *WrapEntry) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrusLogger.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrusLogger.ErrorKey] = err
		l.logrusEntry.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logrusEntry.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	if l.Debug {
		l.logrusEntry.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}
