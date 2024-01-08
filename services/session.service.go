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
	tokenHash := utils.HashToken(token)
	var user models.User

	err := service.db.Select(
		columnWithDot(database.TableUsers, database.ColumnId),
		columnWithDot(database.TableUsers, database.ColumnEmail),
		columnWithDot(database.TableUsers, database.ColumnPasswordHash)).
		From(database.TableSessions).
		Join("users ON users.id = sessions.user5_id").
		Where(squirrel.Eq{columnWithDot(database.TableSessions, database.ColumnTokenHash): tokenHash}).
		QueryRow().
		Scan(&user.Id, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("get user by session: %w", err)
	}

	return &user, nil
}

func (service *SessionService) Delete(token string) error {
	tokenHash := utils.HashToken(token)

	_, err := service.db.Delete(database.TableSessions).
		Where(squirrel.Eq{database.ColumnTokenHash: tokenHash}).
		Exec()
	if err != nil {
		return fmt.Errorf("delete token: %w", err)
	}

	return nil
}
