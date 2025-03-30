package users

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

type Handler struct {
	storage map[string]string
}

type User struct {
	Name string `json:"name,required"`
}

func NewHandler(storage map[string]string) *Handler {
	return &Handler{storage}
}

func (h *Handler) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userUuid := c.Param("uuid")
		name, ok := h.storage[userUuid]
		if ok {
			c.JSON(http.StatusOK, gin.H{"message": "user exists", "uuid": userUuid, "name": name})
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found", "uuid": userUuid})
		}
	}
}

func (h *Handler) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userUuid := uuid.New().String()
		h.storage[userUuid] = user.Name
		name, ok := h.storage[userUuid]

		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "recording error", "uuid": userUuid})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "user created", "uuid": userUuid, "name": name})
	}
}

func (h *Handler) ChangeUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user User
		userUuid := c.Param("uuid")
		_, ok := h.storage[userUuid]

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found", "uuid": userUuid})
			return
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if user.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "name field is required", "uuid": userUuid})
			return
		}
		mu := &sync.Mutex{}
		mu.Lock()
		h.storage[userUuid] = user.Name
		mu.Unlock()
		name, ok := h.storage[userUuid]

		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "recording error", "uuid": userUuid})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user data changed", "uuid": userUuid, "name": name})
	}
}
