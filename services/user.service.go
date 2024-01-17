package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/AkifhanIlgaz/foody-api/cfg"
	"github.com/AkifhanIlgaz/foody-api/models"
	"github.com/AkifhanIlgaz/foody-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection = "users"

type UserService struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewUserService(ctx context.Context, client *mongo.Client, config *cfg.Config) *UserService {
	return &UserService{
		ctx:        ctx,
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

	_, err = service.collection.InsertOne(service.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (service *UserService) Authenticate(email, password string) (*models.User, error) {
	user := models.User{
		Email: email,
	}

	filter := bson.M{"email": email}

	err := service.collection.FindOne(service.ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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
