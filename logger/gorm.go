//nolint:govet
package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/ExecutionLab/gopkg/model"
)

type GormLogger struct {
	logger                    *logrus.Logger
	level                     gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	OnlyLogErrorQuery         bool
}

func NewGorm(onlyLogErrorQuery bool) *GormLogger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})

	return &GormLogger{
		logger:            l,
		OnlyLogErrorQuery: onlyLogErrorQuery,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.level = level

	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.logger.WithContext(ctx).Infof(s, args...)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logger.WithContext(ctx).Warnf(s, args)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.logger.WithContext(ctx).Errorf(s, args)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent || (l.OnlyLogErrorQuery && err == nil) {
		return
	}

	logger := l.logger.WithContext(ctx)

	ddTraceID := ctx.Value("DataDogTraceId")
	if ddTraceID != nil {
		logger = logger.WithField("trace_id", ddTraceID)
	}

	ddSpanID := ctx.Value("DatadogSpanId")
	if ddSpanID != nil {
		logger = logger.WithField("span_id", ddSpanID)
	}

	traceID := ctx.Value(model.KeyContextTraceID)
	if traceID != nil {
		logger = logger.WithField("trace_id", traceID)
	}

	spanID := ctx.Value(model.KeyContextSpanID)
	if spanID != nil {
		logger = logger.WithField("span_id", spanID)
	}

	elapsed := time.Since(begin)

	switch {
	case err != nil && l.level >= gormlogger.Error &&
		(!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.WithFields(logrus.Fields{
				"line":     utils.FileWithLineNum(),
				"sql":      sql,
				"duration": float64(elapsed.Nanoseconds()) / 1e6,
			}).Error(err)
		} else {
			logger.WithFields(logrus.Fields{
				"line":     utils.FileWithLineNum(),
				"sql":      sql,
				"rows":     rows,
				"duration": float64(elapsed.Nanoseconds()) / 1e6,
			}).Error(err)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.level >= gormlogger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)

		sql, rows := fc()
		if rows == -1 {
			logger.WithFields(logrus.Fields{
				"line":     utils.FileWithLineNum(),
				"sql":      sql,
				"slowLog":  slowLog,
				"duration": float64(elapsed.Nanoseconds()) / 1e6,
			}).Warn(err)
		} else {
			logger.WithFields(logrus.Fields{
				"line":     utils.FileWithLineNum(),
				"sql":      sql,
				"slowLog":  slowLog,
				"rows":     rows,
				"duration": float64(elapsed.Nanoseconds()) / 1e6,
			}).Warn(err)
		}
	case l.level == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.WithFields(logrus.Fields{
				"line":     utils.FileWithLineNum(),
				"sql":      sql,
				"duration": float64(elapsed.Nanoseconds()) / 1e6,
			}).Info(err)
		} else {
			logger.WithFields(logrus.Fields{
				"line":     utils.FileWithLineNum(),
				"sql":      sql,
				"rows":     rows,
				"duration": float64(elapsed.Nanoseconds()) / 1e6,
			}).Info(err)
		}
	}
}
