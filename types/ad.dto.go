package types

type AdDto struct {
	Title    string `json:"title" validate:"min=2,max=255"`
	Content  string `json:"content" validate:"min=2,max=1000"`
	ImageUrl string `json:"imageUrl" validate:"url"`
	Price    int    `json:"price" validate:"min=1,max=1000000"`
}
