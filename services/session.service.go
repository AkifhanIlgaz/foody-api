package services

import (
	"context"
	"fmt"

	"github.com/AkifhanIlgaz/foody-api/cfg"
	"github.com/AkifhanIlgaz/foody-api/models"
	"github.com/AkifhanIlgaz/foody-api/utils"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionService struct {
	ctx                context.Context
	sessionsCollection *mongo.Collection
	usersCollection    *mongo.Collection
}

const (
	bytesPerToken      int    = 32
	sessionsCollection string = "sessions"
)

func NewSessionService(ctx context.Context, client *mongo.Client, config *cfg.Config) *SessionService {
	return &SessionService{
		ctx:                ctx,
		sessionsCollection: client.Database(config.MongoDbName).Collection(sessionsCollection),
		usersCollection:    client.Database(config.MongoDbName).Collection(usersCollection),
	}
}

func (service *SessionService) Create(uid primitive.ObjectID) (*models.Session, error) {
	token := randstr.String(bytesPerToken)

	session := models.Session{
		Uid:       uid,
		Token:     token,
		TokenHash: utils.HashToken(token),
	}

	_, err := service.sessionsCollection.InsertOne(service.ctx, session)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	return &session, nil
}

func (service *SessionService) User(token string) (*models.User, error) {
	tokenHash := utils.HashToken(token)
	sessionFilter := bson.M{"tokenHash": tokenHash}

	var session models.Session
	err := service.sessionsCollection.FindOne(service.ctx, sessionFilter).Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("get user by session: %w", err)
	}

	userFilter := bson.M{"_id": session.Uid}
	var user models.User
	err = service.usersCollection.FindOne(service.ctx, userFilter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("get user by session: %w", err)
	}

	return &user, nil
}

func (service *SessionService) Delete(token string) error {
	tokenHash := utils.HashToken(token)

	filter := bson.M{"tokenHash": tokenHash}

	res, err := service.sessionsCollection.DeleteOne(service.ctx, filter)
	if err != nil {
		return fmt.Errorf("delete token: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("token cannot found: %v", token)
	}

	return nil
}
