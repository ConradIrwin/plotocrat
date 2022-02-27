package db

import (
	"context"

	"github.com/ConradIrwin/plotocrat/env"
	// sqlx + pgx requires this
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// Open opens a database connection and sets up the database if it's not been done already
func Open(ctx context.Context) *sqlx.DB {
    db, err := sqlx.Connect("pgx", env.Config.ConnectionString)
    if err != nil {
        panic(err)
    }
    if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS charts (
							id TEXT PRIMARY KEY,
							x_axis TEXT NOT NULL,
							data TEXT NOT NULL
						)`); err != nil {
        panic(err)
    }
    return db
}