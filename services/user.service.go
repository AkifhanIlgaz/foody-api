package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/AkifhanIlgaz/foody-api/database"
	"github.com/AkifhanIlgaz/foody-api/models"
	"github.com/AkifhanIlgaz/foody-api/utils"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type UserService struct {
	db squirrel.StatementBuilderType
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}

func (service *UserService) Create(email, password string) (*models.User, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user := models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err = service.db.Insert(database.TableUsers).
		Columns(database.ColumnEmail, database.ColumnPasswordHash).
		Values(user.Email, user.PasswordHash).
		Suffix("RETURNING \"id\"").
		QueryRow().
		Scan(&user.Id)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (service *UserService) Authenticate(email, password string) (*models.User, error) {
	user := models.User{
		Email: email,
	}

	err := service.db.Select(database.ColumnId, database.ColumnPasswordHash).
		From(database.TableUsers).
		Where(squirrel.Eq{database.ColumnEmail: email}).
		QueryRow().
		Scan(&user.Id, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = utils.VerifyPassword(user.PasswordHash, password)
	if err != nil {
		return nil, ErrWrongPassword
	}

	return &user, nil
}
