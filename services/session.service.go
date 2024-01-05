package services

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type SessionService struct {
	db squirrel.StatementBuilderType
}

func NewSessionService(db *sql.DB) *SessionService {
	return &SessionService{
		db: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}
