package db

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type LibsqlDB struct {
	url   string
	token string

	db *sql.DB
}

func NewLibsqlDB(ctx context.Context, parentLog *slog.Logger, url, token string) (*LibsqlDB, error) {
	log := parentLog.WithGroup("db")

	log.InfoContext(ctx, "connecting to libsql database...")
	db, err := sql.Open("libsql", url+"?authToken="+token)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.InfoContext(ctx, "connected to libsql database")

	ldb := &LibsqlDB{
		url:   url,
		token: token,
		db:    db,
	}

	log.InfoContext(ctx, "migrating libsql database...")
	if err := ldb.Migrate(ctx); err != nil {
		return nil, err
	}
	log.InfoContext(ctx, "migrated libsql database")

	return ldb, nil
}

func (db *LibsqlDB) DB() *sql.DB {
	return db.db
}

func (db *LibsqlDB) Close(_ context.Context) error {
	return db.db.Close()
}

func (db *LibsqlDB) Migrate(ctx context.Context) error {
	goose.SetBaseFS(migrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db.db, "migrations"); err != nil {
		return err
	}

	return nil
}
