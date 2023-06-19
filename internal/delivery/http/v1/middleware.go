package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "Authorization"
	userCtx                = "userId"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	username, err := h.parseAuthHeader(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ctx.Set(userCtx, username)
	ctx.Next()
}

func (h *Handler) parseAuthHeader(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader(authorizationHeaderKey)
	if header == "" {
		return "", errors.New("empty auth header")
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		fmt.Println(header)
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.token.Parse(headerParts[1])
}
