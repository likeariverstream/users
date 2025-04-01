package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/handlers"
)

func TestNotFoundUser(t *testing.T) {
	t.Parallel()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()

	userUuid := uuid.New().String()

	mock.ExpectQuery("SELECT uuid, name, email FROM users WHERE uuid = $1").WithArgs(userUuid).WillReturnError(sql.ErrNoRows)

	handler := handlers.NewHandler(db)
	r := router(handler)
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/users/%s", userUuid)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
		Email   *string `json:"email"`
	}
	b.Message = "user not found: sql: no rows in result set"
	b.Name = nil
	b.Uuid = userUuid
	b.Email = nil
	body, _ := json.Marshal(b)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, string(body), w.Body.String())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	userUuid := uuid.New().String()
	url := fmt.Sprintf("/users/%s", userUuid)
	columns := []string{
		"uuid",
		"name",
		"email",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(userUuid, "John Doe", "john.doe@example.com")
	mock.ExpectQuery("SELECT uuid, name, email FROM users WHERE uuid = $1").WithArgs(userUuid).WillReturnRows(rows)

	handler := handlers.NewHandler(db)
	r := router(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
		Email   *string `json:"email"`
	}

	b.Message = "user exists"
	b.Uuid = userUuid
	name := "John Doe"
	email := "john.doe@example.com"
	b.Name = &name
	b.Email = &email
	body, _ := json.Marshal(b)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(body), w.Body.String())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

type ChReqBody struct {
	Name string `json:"name"`
}

func TestChangeUser(t *testing.T) {
	t.Parallel()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	userUuid := uuid.New().String()
	url := fmt.Sprintf("/users/%s", userUuid)
	columns := []string{"uuid", "name", "email"}

	rows := sqlmock.NewRows(columns).AddRow(userUuid, "Jane Smith", "john.doe@example.com")
	mock.ExpectQuery("UPDATE users SET name = $1, updated_at = $2 WHERE uuid = $3 RETURNING uuid, name, email").
		WithArgs("Jane Smith", sqlmock.AnyArg(), userUuid).
		WillReturnRows(rows)
	handler := handlers.NewHandler(db)
	r := router(handler)
	w := httptest.NewRecorder()

	reqBody := ChReqBody{
		Name: "Jane Smith",
	}

	byteBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(byteBody))
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
		Email   *string `json:"email"`
	}
	name := "Jane Smith"
	email := "john.doe@example.com"
	b.Message = "user data changed"
	b.Email = &email
	b.Name = &name
	b.Uuid = userUuid
	body, _ := json.Marshal(b)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(body), w.Body.String())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFailChangeUser(t *testing.T) {
	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	userUuid := uuid.New().String()
	url := fmt.Sprintf("/users/%s", userUuid)
	handler := handlers.NewHandler(db)
	r := router(handler)
	w := httptest.NewRecorder()

	reqBody := ChReqBody{
		Name: "",
	}

	byteBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(byteBody))
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
		Email   *string `json:"email"`
	}

	b.Message = "name field is required"
	b.Email = nil
	b.Name = nil
	b.Uuid = userUuid
	body, _ := json.Marshal(b)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, string(body), w.Body.String())
}

type CrReqBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	userUuid := uuid.New().String()
	url := "/users"
	columns := []string{"uuid", "name", "email"}

	rows := sqlmock.NewRows(columns).AddRow(userUuid, "Jane Smith", "john.doe@example.com")
	mock.ExpectQuery("INSERT INTO users (uuid, name, email, created_at) VALUES($1, $2, $3, $4) RETURNING uuid, name, email").
		WithArgs(sqlmock.AnyArg(), "Jane Smith", "john.doe@example.com", sqlmock.AnyArg()).
		WillReturnRows(rows)
	handler := handlers.NewHandler(db)
	r := router(handler)
	w := httptest.NewRecorder()

	reqBody := CrReqBody{
		Name:  "Jane Smith",
		Email: "john.doe@example.com",
	}

	byteBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(byteBody))
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
		Email   *string `json:"email"`
	}
	name := "Jane Smith"
	email := "john.doe@example.com"
	b.Message = "user created"
	b.Email = &email
	b.Name = &name
	b.Uuid = userUuid
	body, _ := json.Marshal(b)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, string(body), w.Body.String())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
