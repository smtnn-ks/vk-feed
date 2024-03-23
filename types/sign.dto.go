package types

type SignDto struct {
	Name     string `json:"name" validate:"min=8,max=16"`
	Password string `json:"password" validate:"min=8,max=16"`
}
