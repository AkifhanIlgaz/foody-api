package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AkifhanIlgaz/foody-api/cfg"
	"github.com/go-redis/redis/v8"

	_ "github.com/lib/pq"
)

type Databases struct {
	Postgres *sql.DB
	Redis    *redis.Client
}

func Connect(config *cfg.Config) (*Databases, error) {
	postgres, err := connectToPostgres(config)
	if err != nil {
		return nil, fmt.Errorf("connect to databases: %w", err)
	}

	redis, err := connectToRedis(config)
	if err != nil {
		return nil, fmt.Errorf("connect to databases: %w", err)
	}

	return &Databases{
		Postgres: postgres,
		Redis:    redis,
	}, nil
}

func connectToPostgres(config *cfg.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%v user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPassword, config.PostgresDBName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("connect to Postgres: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to ping postgres: %w", err)
	}

	fmt.Println("Postgres connected successfully!")
	return db, nil
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
