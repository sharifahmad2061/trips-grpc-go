package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sharifahmad2061/trip-grpc-go/internal/config"

	_ "github.com/lib/pq"
)

func initialize(ctx context.Context) (*sql.DB, error) {
	conf := config.Load()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DbName,
		conf.Database.SslMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
