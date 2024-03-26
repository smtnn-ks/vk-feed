package db

import (
	"context"
	"fmt"
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

func (conn PgxConnection) GetAds(userId int, params types.GetAdParams) (res []types.AdFeed, err error) {
	query := fmt.Sprintf(
		"SELECT * FROM ads WHERE price >= $1 AND price <= $2 ORDER BY %s %s OFFSET 10*$3 LIMIT 10",
		params.SortBy,
		params.OrderBy,
	)
	rows, err := conn.Client.Query(context.Background(), query, params.MinPrice, params.MaxPrice, params.Page)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ad types.AdFeed
		rows.Scan(&ad.Id, &ad.Title, &ad.Content, &ad.ImageUrl, &ad.Price, &ad.AuthorId, &ad.CreatedAt)
		if ad.AuthorId == userId {
			ad.IsYours = true
		}
		res = append(res, ad)
	}
	return
}
