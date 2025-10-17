package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/sharifahmad2061/trip-grpc-go/internal/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	_ "github.com/lib/pq"
)

func Initialize(ctx context.Context) (*sql.DB, error) {
	conf := config.Load()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DbName,
		conf.Database.SslMode,
	)
	db, err := otelsql.Open("postgres", connStr, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL),
		otelsql.WithSQLCommenter(true))
	if err != nil {
		return nil, err
	}
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
