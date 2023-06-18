package v1

import (
	"net/http"
	"time"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pkg/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/create", h.createUser)
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username" binding:"required,alphanum"`
	FullName          string    `json:"full_name" binding:"required"`
	Email             string    `json:"email" binding:"required,email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (h *Handler) createUser(ctx *gin.Context) {
	var inp createUserRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Invalid input"+err.Error())
		return
	}
	arg := domain.CreateUserParams{
		Username:       inp.Username,
		HashedPassword: inp.Password,
		FullName:       inp.FullName,
		Email:          inp.Email,
	}

	user, err := h.service.User.CreateUser(ctx, arg)
	if err != nil {
		if e.ErrorCode(err) == e.UniqueViolation {
			newResponse(ctx, http.StatusForbidden, "Internal db:"+err.Error())
			return
		}

		newResponse(ctx, http.StatusInternalServerError, "Invalid db:"+err.Error())
		return
	}

	res := createUserResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}

	ctx.JSON(http.StatusOK, res)
}
