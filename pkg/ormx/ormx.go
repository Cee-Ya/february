package ormx

import (
	"ai-report/common"
	"ai-report/config/log"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitOrmx() {
	c := common.GlobalConfig.DB
	db, err := gorm.Open(mysql.Open(c.Dsn), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	if c.Debug {
		db = db.Debug()
	}
	common.Ormx = db
}

type GormLogger struct {
	ZapLogger *zap.Logger
}

func NewGormLogger() *GormLogger {
	return &GormLogger{ZapLogger: common.Logger}
}

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	log.InfoF(ctx, str, zap.Any("args", args))
}

func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	log.WarnF(ctx, str, zap.Any("args", args))
}

func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	log.ErrorF(ctx, str, zap.Any("args", args))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.ZapLogger.Core().Enabled(zap.ErrorLevel):
		sql, rows := fc()
		log.ErrorF(ctx, "SQL::",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
		)
	case elapsed > 100*time.Millisecond && l.ZapLogger.Core().Enabled(zap.WarnLevel):
		sql, rows := fc()
		log.WarnF(ctx, "SQL-SLOW::",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
		)
	case l.ZapLogger.Core().Enabled(zap.DebugLevel):
		sql, rows := fc()
		log.DebugF(ctx, "SQL::",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
		)
	}
}
