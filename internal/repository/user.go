package repository

import (
	"database/sql"
	"github.com/go-redis/redis"
)

type UserRepository struct {
	*UserPostgres
	*UserRedis
}

func NewUserRepository(db *sql.DB, rdb *redis.Client) *UserRepository {
	return &UserRepository{
		UserPostgres: NewUserPostgres(db),
		UserRedis:    NewUserRedis(rdb),
	}
}
