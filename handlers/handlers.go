package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/db"
)

type Storage interface {
	AddUser(name string, email string) (*db.User, error)
	GetUser(uuid string) (*db.User, error)
	ChangeUser(uuid string, name string) (*db.User, error)
}
type Handler struct {
	Storage Storage
}

type UserReq struct {
	Name  string `json:"name,required" binding:"required"`
	Email string `json:"email,required" binding:"required,email"`
}
type UserResp struct {
	Message string  `json:"message"`
	Uuid    string  `json:"uuid"`
	Name    *string `json:"name"`
	Email   *string `json:"email"`
}
type Param struct {
	uuid string `binding:"uuid"`
}

func NewHandler(storage *sql.DB) *Handler {
	return &Handler{Storage: db.NewStorage(storage)}
}

// GetUser godoc
// @Summary Get user
// @Description Get user
// @Tags Users
// @Accept json
// @Produce json
// @Param uuid path string true "User uuid"
// @Success 200 {object} UserResp "Get successfully"
// @Failure 400 {object} UserResp
// @Failure 404 {object} UserResp
// @Router /users/{uuid} [get]
func (h *Handler) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userUuid := c.Param("uuid")
		if err := c.ShouldBindUri(&Param{
			uuid: userUuid,
		}); err != nil {
			r := UserResp{
				Message: err.Error(),
			}
			c.JSON(http.StatusBadRequest, r)
		}
		user, err := h.Storage.GetUser(userUuid)

		if err == nil {
			r := &UserResp{
				Message: "user exists",
				Uuid:    user.Uuid,
				Name:    &user.Name,
				Email:   &user.Email,
			}
			c.JSON(http.StatusOK, r)
			return
		} else {
			r := &UserResp{
				Message: err.Error(),
				Uuid:    userUuid,
				Name:    nil,
				Email:   nil,
			}
			c.JSON(http.StatusNotFound, r)
			return
		}
	}
}

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags Users
// @Accept json
// @Produce json
// @Param data body UserReq true "User data"
// @Success 201 {object} UserResp "Create successfully"
// @Failure 400 {object} UserResp "Bad request"
// @Failure 422 {object} UserResp "Unprocessable"
// @Router /users [post]
func (h *Handler) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user UserReq
		r := &UserResp{
			Message: "",
			Uuid:    "",
			Name:    nil,
			Email:   nil,
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			r.Message = err.Error()
			c.JSON(http.StatusBadRequest, r)
			return
		}
		res, err := h.Storage.AddUser(user.Name, user.Email)

		if err != nil {
			r.Message = err.Error()
			c.JSON(http.StatusUnprocessableEntity, r)
			return
		}
		r.Message = "user created"
		r.Uuid = res.Uuid
		r.Name = &res.Name
		r.Email = &res.Email

		c.JSON(http.StatusCreated, r)
		return
	}
}

// ChangeUser godoc
// @Summary Change user
// @Description Change user
// @Tags Users
// @Accept json
// @Produce json
// @Param uuid path string true "User uuid"
// @Param data body UserReq true "User data"
// @Success 200 {object} UserResp "Change successfully"
// @Failure 400 {object} UserResp "Bad request"
// @Failure 422 {object} UserResp "Unprocessable"
// @Router /users/{uuid} [put]
func (h *Handler) ChangeUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user UserReq
		userUuid := c.Param("uuid")
		if err := c.ShouldBindUri(&Param{
			uuid: userUuid,
		}); err != nil {
			r := UserResp{
				Message: err.Error(),
			}
			c.JSON(http.StatusBadRequest, r)
		}
		r := &UserResp{
			Message: "",
			Uuid:    userUuid,
			Name:    nil,
			Email:   nil,
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			r.Message = err.Error()
			r.Name = &user.Name
			c.JSON(http.StatusBadRequest, r)
			return
		}

		if user.Name == "" {
			r.Message = "name field is required"
			c.JSON(http.StatusBadRequest, r)
			return
		}

		res, err := h.Storage.ChangeUser(userUuid, user.Name)

		if err != nil {
			r.Message = err.Error()
			r.Name = &user.Name
			c.JSON(http.StatusUnprocessableEntity, r)
			return
		}
		r = &UserResp{
			Message: "user data changed",
			Uuid:    res.Uuid,
			Name:    &res.Name,
			Email:   &res.Email,
		}
		c.JSON(http.StatusOK, r)
		return
	}
}
