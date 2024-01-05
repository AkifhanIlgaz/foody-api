package services

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type UserService struct {
	db squirrel.StatementBuilderType
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}
