package env

import (
	"context"
	"os"

	"github.com/jmoiron/sqlx"
)

// Config is global to the app, and initialized by Init()
var Config struct {
	Env string
	Listen string
	ConnectionString string
	ServerOrigin string
}

// DB is a shared database handle for the app to use
var DB *sqlx.DB

var ctx context.Context
var cancel func() 

// Init sets up the environment and returns the root context
func Init () (context.Context, func()) {
	ctx, cancel = context.WithCancel(context.Background())

	if os.Getenv("APP_ENV") == "production" {
		Config.Env = "production"
		Config.Listen = ":" + os.Getenv("PORT")
		Config.ConnectionString = os.Getenv("DATABASE_URL")
		Config.ServerOrigin = "https://plotocrat.com"
	} else {
		Config.Env = "development"
		Config.Listen = "localhost:" + os.Getenv("PORT")
		Config.ConnectionString = "postgres://plotocrat:plotocrat@localhost:5432/plotocrat"
		Config.ServerOrigin = "http://localhost:3000"
	}
	return ctx, cancel
}