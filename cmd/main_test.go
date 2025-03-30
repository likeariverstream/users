package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	users "user-service"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundUser(t *testing.T) {
	var testStorage = make(map[string]string)
	handler := users.NewHandler(testStorage)
	r := router(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/aff7323b-7c7c-410d-8079-8fc2ab1e7e84", nil)
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
	}
	b.Message = "not found"
	b.Name = nil
	b.Uuid = "aff7323b-7c7c-410d-8079-8fc2ab1e7e84"
	body, _ := json.Marshal(b)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, string(body), w.Body.String())
}

func TestGetUser(t *testing.T) {
	var testStorage = make(map[string]string)
	testStorage["aff7323b-7c7c-410d-8079-8fc2ab1e7e85"] = "John Doe"
	handler := users.NewHandler(testStorage)
	r := router(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/aff7323b-7c7c-410d-8079-8fc2ab1e7e85", nil)
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
	}
	b.Message = "user exists"
	b.Uuid = "aff7323b-7c7c-410d-8079-8fc2ab1e7e85"
	name := "John Doe"
	b.Name = &name
	body, _ := json.Marshal(b)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(body), w.Body.String())
}

func TestChangeUser(t *testing.T) {
	var testStorage = make(map[string]string)
	testStorage["aff7323b-7c7c-410d-8079-8fc2ab1e7e85"] = "John Doe"
	handler := users.NewHandler(testStorage)
	r := router(handler)
	w := httptest.NewRecorder()
	reqBody := `{"name":"Pavel"}`
	req, _ := http.NewRequest("PUT", "/users/aff7323b-7c7c-410d-8079-8fc2ab1e7e85", strings.NewReader(reqBody))
	r.ServeHTTP(w, req)

	var b struct {
		Message string  `json:"message"`
		Uuid    string  `json:"uuid"`
		Name    *string `json:"name"`
	}
	name := "Pavel"
	b.Message = "user data changed"
	b.Name = &name
	b.Uuid = "aff7323b-7c7c-410d-8079-8fc2ab1e7e85"
	body, _ := json.Marshal(b)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(body), w.Body.String())
}
