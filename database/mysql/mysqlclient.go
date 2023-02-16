package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	databaseconfig "github.com/forbole/juno/v4/database/config"
	"github.com/forbole/juno/v4/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg *databaseconfig.Config) (*gorm.DB, error) {
	if cfg.Secrets != nil {
		secret, err := databaseconfig.GetString(cfg.Secrets)
		if err != nil {
			log.Errorf("invalid secrets %+v err:%v", cfg.Secrets, err)
			return nil, err
		}
		cfg.DSN = secret
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN),
		&gorm.Config{
			Logger:                                   &loggerAdaptor{slowThreshold: time.Duration(cfg.SlowThreshold)},
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)
	if err != nil {
		log.Errorw("failed to open database", "err", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorw("failed to get database", "err", err)
		return nil, err
	}

	if cfg.MaxOpenConnections <= 0 {
		cfg.MaxOpenConnections = 256
	}
	if cfg.MaxIdleConnections <= 0 {
		cfg.MaxIdleConnections = cfg.MaxOpenConnections
	}
	if cfg.ConnMaxIdleTime <= databaseconfig.Duration(time.Minute) {
		cfg.ConnMaxIdleTime = databaseconfig.Duration(5 * time.Minute)
	}
	if cfg.ConnMaxLifetime <= databaseconfig.Duration(time.Minute) {
		cfg.ConnMaxLifetime = databaseconfig.Duration(time.Hour)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime))
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime))

	return db, nil
}

type loggerAdaptor struct {
	slowThreshold time.Duration
}

func (la *loggerAdaptor) LogMode(logger.LogLevel) logger.Interface {
	return &loggerAdaptor{slowThreshold: la.slowThreshold}
}

func (*loggerAdaptor) Info(ctx context.Context, fmt string, args ...interface{}) {
	log.With("module", "gorm").AddCallerSkip(1).CtxInfof(ctx, fmt, args...)
}

func (*loggerAdaptor) Warn(ctx context.Context, fmt string, args ...interface{}) {
	log.With("module", "gorm").AddCallerSkip(1).CtxWarnf(ctx, fmt, args...)
}

func (*loggerAdaptor) Error(ctx context.Context, fmt string, args ...interface{}) {
	log.With("module", "gorm").AddCallerSkip(1).CtxErrorf(ctx, fmt, args...)
}
func (la *loggerAdaptor) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil:
		// ignore not found
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		strSql, rows := fc()
		log.CtxErrorw(ctx, "error sql", "err", err, "elapsed", elapsed, "sql", strSql, "rows", rows)
	case elapsed > la.slowThreshold && la.slowThreshold != 0:
		strSql, rows := fc()
		log.CtxWarnw(ctx, "slow sql", "elapsed", elapsed, "sql", strSql, "rows", rows)
	}
}
