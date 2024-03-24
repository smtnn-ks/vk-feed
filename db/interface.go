package db

type DBConnection interface {
	CreateUser(name, password string) (int, error)
	GetUserByName(name string) (int, string, error)
}
