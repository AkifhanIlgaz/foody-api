package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/AkifhanIlgaz/foody-api/cfg"
	"github.com/AkifhanIlgaz/foody-api/database"
	"github.com/AkifhanIlgaz/foody-api/models"
	"github.com/AkifhanIlgaz/foody-api/utils"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection = "users"

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(client *mongo.Client, config *cfg.Config) *UserService {

	return &UserService{
		collection: client.Database(config.MongoDbName).Collection(usersCollection),
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
		Suffix("RETURNING id").
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
