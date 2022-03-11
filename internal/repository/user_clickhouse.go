package repository

import (
	"database/sql"
	"fmt"
	"grpc/internal/entity"
	"log"
)

type UserClickHouse struct {
	clickDB *sql.DB
}

func NewUserClickHouse(clickDB *sql.DB) *UserClickHouse {
	return &UserClickHouse{
		clickDB: clickDB,
	}
}

func (c *UserClickHouse) LogUser(user *entity.User) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (userId, name, email) VALUES ($1, $2, $3)`,
		usersTable,
	)
	_, err := c.clickDB.Exec(query, user.Id, user.Email, user.Name)
	log.Println("UserClickHouse err: ", err)
	return err
}
