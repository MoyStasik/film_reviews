package models

type User struct {
	Id       int64
	Email    string
	Name     string
	PassHash []byte
}
