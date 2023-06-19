package http

import (
	"github.com/begenov/backend/internal/config"
	v1 "github.com/begenov/backend/internal/delivery/http/v1"
	"github.com/begenov/backend/internal/service"
	"github.com/begenov/backend/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
	token   auth.TokenManager
}

func NewHandler(service *service.Service, token auth.TokenManager) *Handler {
	return &Handler{
		service: service,
		token:   token,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	h.init(router)

	return router
}

func (h *Handler) init(router *gin.Engine) {
	handlerv1 := v1.NewHandler(h.service, h.token)
	api := router.Group("/api")
	{
		handlerv1.Init(api)
	}
}
