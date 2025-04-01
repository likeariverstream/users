package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type StDb struct {
	db *sql.DB
}

type User struct {
	Uuid  string `sql:"uuid"`
	Name  string `sql:"name"`
	Email string `sql:"email"`
}

func NewStorage(db *sql.DB) *StDb {
	return &StDb{db: db}
}

func (st *StDb) AddUser(name string, email string) (*User, error) {
	var user User
	newUuid := uuid.New().String()
	row := st.db.QueryRow("INSERT INTO users (uuid, name, email, created_at) VALUES($1, $2, $3, $4) RETURNING uuid, name, email",
		newUuid, name, email, time.Now())

	if err := row.Scan(&user.Uuid, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(err.Error())
		}
		return nil, err
	}
	return &user, nil

}

func (st *StDb) GetUser(uuid string) (*User, error) {
	var user User

	row := st.db.QueryRow("SELECT uuid, name, email FROM users WHERE uuid = $1", uuid)

	if err := row.Scan(&user.Uuid, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", sql.ErrNoRows)
		}
		return nil, err
	}
	return &user, nil
}

func (st *StDb) ChangeUser(uuid string, name string) (*User, error) {
	var user User

	row := st.db.QueryRow("UPDATE users SET name = $1, updated_at = $2 WHERE uuid = $3 RETURNING uuid, name, email", name, time.Now(), uuid)

	if err := row.Scan(&user.Uuid, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", sql.ErrNoRows)
		}
		return nil, err
	}

	return &user, nil
}
