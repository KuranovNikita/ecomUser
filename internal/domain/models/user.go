package models

type User struct {
	ID       int64
	Email    string
	Login    string
	PassHash []byte
}
