package repository

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
)

type UserRepository struct {
	*UserPostgres
	*UserRedis
	*UserClickHouse
}


func NewUserRepository(db *sql.DB, rdb *redis.Client,clickDB *sql.DB) *UserRepository {
	return &UserRepository{
		UserPostgres: NewUserPostgres(db),
		UserRedis:    NewUserRedis(rdb),
		UserClickHouse: NewUserClickHouse(clickDB),
	}
}
