package types

type Ad struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageUrl string `json:"imageUrl"`
	Price    int    `json:"price"`
}
