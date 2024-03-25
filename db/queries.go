package db

import (
	"context"
	"vk-feed/types"

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

func (conn PgxConnection) GetUserByName(name string) (id int, password string, err error) {
	query := "SELECT id, pass FROM usrs WHERE name = $1"
	err = conn.Client.QueryRow(context.Background(), query, name).Scan(&id, &password)
	return
}

func (conn PgxConnection) CreateAd(dto types.AdDto, userId int) (id int, err error) {
	query := "INSERT INTO ads (title, content, image_url, price, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = conn.Client.QueryRow(context.Background(), query, dto.Title, dto.Content, dto.ImageUrl, dto.Price, userId).Scan(&id)
	return
}
