package db

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-pg/pg/v10"
)

var _db *pg.DB = nil

var ErrNoRows = pg.ErrNoRows

func PG() *pg.DB {
	return _db
}

func Begin() (*pg.Tx, error) {
	return _db.Begin()
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		return err
	}
	/* 打印SQL执行事件 */
	fmt.Printf("\033[34m%s\n\033[0m", query)
	return nil
}

/* 初始化PostgreSQL */
func InitPG(sURL string, debug bool) error {
	opts, err := pg.ParseURL(sURL)
	if err != nil {
		return err
	}
	opts.PoolSize = runtime.NumCPU() * 10
	_db = pg.Connect(opts)
	/* DEBUG */
	if debug {
		_db.AddQueryHook(dbLogger{})
	}
	return nil
}
