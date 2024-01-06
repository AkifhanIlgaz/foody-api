package services

import (
	"database/sql"
	"fmt"

	"github.com/AkifhanIlgaz/foody-api/database"
	"github.com/AkifhanIlgaz/foody-api/models"
	"github.com/AkifhanIlgaz/foody-api/utils"
	"github.com/Masterminds/squirrel"
	"github.com/thanhpk/randstr"
)

type SessionService struct {
	db squirrel.StatementBuilderType
}

const BytesPerToken int = 32

func NewSessionService(db *sql.DB) *SessionService {
	return &SessionService{
		db: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}

func (service *SessionService) Create(uid int) (*models.Session, error) {
	token := randstr.String(BytesPerToken)

	session := models.Session{
		Uid:       uid,
		Token:     token,
		TokenHash: utils.HashToken(token),
	}

	err := service.db.Insert(database.TableSessions).
		Columns(database.ColumnUserId, database.ColumnTokenHash).
		Values(session.Uid, session.TokenHash).
		Suffix(`ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id`).
		QueryRow().
		Scan(&session.Id)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	return &session, nil
}

func (service *SessionService) User(token string) (*models.User, error) {

	return nil, nil
}

func (service *SessionService) Delete(token string) error {

	return nil
}
