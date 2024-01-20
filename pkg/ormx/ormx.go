package ormx

import (
	"context"
	"february/common"
	"february/pkg/logx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"runtime/debug"
	"time"
)

func Init() error {
	c := common.GlobalConfig.DB
	db, err := gorm.Open(mysql.Open(c.Dsn), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to connect database")
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "faild to open db")
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	common.Ormx = db
	return nil
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
	logx.InfoF(ctx, str, zap.Any("args", args))
}

func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	logx.WarnF(ctx, str, zap.Any("args", args))
}

func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	logx.ErrorF(ctx, str, zap.Any("args", args))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.ZapLogger.Core().Enabled(zap.ErrorLevel):
		sql, rows := fc()
		debug.PrintStack()
		logx.ErrorF(ctx, "SQL::",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
		)
	case elapsed > 100*time.Millisecond && l.ZapLogger.Core().Enabled(zap.WarnLevel):
		sql, rows := fc()
		logx.WarnF(ctx, "SQL-SLOW::",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
		)
	case l.ZapLogger.Core().Enabled(zap.DebugLevel):
		sql, rows := fc()
		logx.DebugF(ctx, "SQL::",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
		)
	}
}