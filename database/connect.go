package database

import (
	"context"
	"fmt"

	"github.com/AkifhanIlgaz/foody-api/cfg"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Databases struct {
	Mongo *mongo.Client
	Redis *redis.Client
}

func Connect(config *cfg.Config) (*Databases, error) {
	mongo, err := connectToMongo(context.TODO(), config)
	if err != nil {
		return nil, fmt.Errorf("connect to databases: %w", err)
	}

	redis, err := connectToRedis(config)
	if err != nil {
		return nil, fmt.Errorf("connect to databases: %w", err)
	}

	return &Databases{
		Mongo: mongo,
		Redis: redis,
	}, nil
}

func connectToMongo(ctx context.Context, config *cfg.Config) (*mongo.Client, error) {
	mongoConn := options.Client().ApplyURI(config.MongoUri)

	client, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		return nil, fmt.Errorf(" connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	fmt.Println("MongoDB successfully connected...")
	return client, nil
}

func connectToRedis(config *cfg.Config) (*redis.Client, error) {
	ctx := context.TODO()

	client := redis.NewClient(&redis.Options{
		Addr: config.RedisUrl,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	err := client.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	fmt.Println("Redis client connected successfully!")
	return client, nil
}
