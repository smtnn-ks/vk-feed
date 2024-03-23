package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Init(url string) (PgxConnection, error) {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return PgxConnection{}, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return PgxConnection{}, err
	}
	return PgxConnection{Client: pool}, nil
}
