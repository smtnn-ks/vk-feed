package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxConnection struct {
	Client *pgxpool.Pool
}

func (conn PgxConnection) CreateUser(name, password string) (id int, err error) {
	query := "INSERT INTO usrs (name, pass) VALUES ($1, $2) RETURNING id"
	err = conn.Client.QueryRow(context.Background(), query, name, password).Scan(&id)
	return
}
