package types

import "time"

type AdFeed struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"iamgeUrl"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	AuthorId  int       `json:"authorId"`
	IsYours   bool      `json:"isYours"`
}
