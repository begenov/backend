package v1

import (
	"database/sql"
	"net/http"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pkg/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/create", h.createUser)
		users.POST("/login", h.loginUser)
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
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

	res := newUserResponse(user)

	ctx.JSON(http.StatusOK, res)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) loginUser(ctx *gin.Context) {
	var inp loginUserRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	res, err := h.service.User.GetUserByUsername(ctx, inp.Username, inp.Password)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			newResponse(ctx, http.StatusNotFound, "Invalid db:"+err.Error())
			return
		case e.ErrInvalidToken:
			newResponse(ctx, http.StatusUnauthorized, "Invalid db:"+err.Error())
			return
		case e.ErrPassword:
			newResponse(ctx, http.StatusBadRequest, "Invalid db:"+err.Error())
			return
		default:
			newResponse(ctx, http.StatusInternalServerError, "Invalid db:"+err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, res)
}

func newUserResponse(user domain.User) domain.UserResponse {
	return domain.UserResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}
