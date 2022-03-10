package repository

import (
	"database/sql"
	"fmt"
	"grpc/internal/entity"
	"log"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db}
}

func (r *UserPostgres) Create(user *entity.User) (uint64, error) {
	var id uint64
	query := fmt.Sprintf(
		"INSERT INTO %s (name, email, password) values ($1,$2,$3) RETURNING id",
		usersTable,
	)
	row := r.db.QueryRow(query, user.Name, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserPostgres) GetAll() ([]entity.User, error) {
	var users []entity.User
	query := fmt.Sprintf(
		"SELECT * FROM %s",
		usersTable,
	)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.Id, &user.Email, &user.Name, &user.Password)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (r *UserPostgres) Delete(userId uint64) error {
	query := fmt.Sprintf(
		`DELETE FROM %s u WHERE u.id = $1`,
		usersTable,
	)
	_, err := r.db.Exec(query, userId)
	return err
}
