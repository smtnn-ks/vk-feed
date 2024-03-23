package db

type DBConnection interface {
	CreateUser(name, password string) (id int, err error)
}
