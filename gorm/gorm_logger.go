// Package gorm provides logger that implements gorm's logger.Interface
package gorm

import (
	"context"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger ezutil.Logger
	level  logger.LogLevel
}

func NewGormLogger(l ezutil.Logger) *GormLogger {
	return &GormLogger{logger: l, level: logger.Silent}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{logger: g.logger, level: level}
}

func (g *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	g.logger.Infof(msg, args...)
}

func (g *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	g.logger.Warnf(msg, args...)
}

func (g *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	g.logger.Errorf(msg, args...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && g.level >= logger.Error:
		g.logger.Errorf("[%.3fms | %d rows] %s — %v", float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
	case elapsed > 200*time.Millisecond && g.level >= logger.Warn:
		g.logger.Warnf("[%.3fms | %d rows] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case g.level >= logger.Info:
		g.logger.Debugf("[%.3fms | %d rows] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
