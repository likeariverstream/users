package users

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

type Handler struct {
	Storage map[string]string
}

type UserReq struct {
	Name string `json:"name,required"`
}
type UserResp struct {
	Message string  `json:"message,required"`
	Uuid    string  `json:"uuid,required"`
	Name    *string `json:"name,required"`
}

func NewHandler(storage map[string]string) *Handler {
	return &Handler{storage}
}

func (h *Handler) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userUuid := c.Param("uuid")
		name, ok := h.Storage[userUuid]

		if ok {
			r := &UserResp{
				Message: "user exists",
				Uuid:    userUuid,
				Name:    &name,
			}
			c.JSON(http.StatusOK, r)
			return
		} else {
			r := &UserResp{
				Message: "not found",
				Uuid:    userUuid,
				Name:    nil,
			}
			c.JSON(http.StatusNotFound, r)
		}
	}
}

func (h *Handler) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user UserReq
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userUuid := uuid.New().String()
		h.Storage[userUuid] = user.Name
		name, ok := h.Storage[userUuid]

		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "recording error", "uuid": userUuid})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "user created", "uuid": userUuid, "name": name})
	}
}

func (h *Handler) ChangeUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user UserReq
		userUuid := c.Param("uuid")
		_, ok := h.Storage[userUuid]

		if !ok {
			r := &UserResp{
				Message: "user not found",
				Uuid:    userUuid,
				Name:    nil,
			}
			c.JSON(http.StatusNotFound, r)
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
		h.Storage[userUuid] = user.Name
		mu.Unlock()
		name, ok := h.Storage[userUuid]

		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "recording error", "uuid": userUuid})
			return
		}
		r := &UserResp{
			Message: "user data changed",
			Uuid:    userUuid,
			Name:    &name,
		}
		c.JSON(http.StatusOK, r)
	}
}
