package main

import (
	"database/sql"
	"flag"
	"github.com/go-sql-driver/mysql"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web/pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
