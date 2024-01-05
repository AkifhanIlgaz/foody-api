package models

type Session struct {
	Id        int
	Uid       int
	Token     string
	TokenHash string
}
