package repository

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"grpc/internal/entity"
)

type Repository struct {
	User
}

type User interface {
	Create(user *entity.User) (uint64, error)
	GetAll() ([]entity.User, error)
	Delete(userId uint64) error
	CacheUsers([]entity.User)
	GetCachedUsers() ([]entity.User, error)
	HasCachedUsers() bool
	FlushCachedUsers()
	LogUser(user *entity.User) error
}

func NewRepository(db *sql.DB, rdb *redis.Client, clickDB *sql.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db, rdb, clickDB),
	}
}
