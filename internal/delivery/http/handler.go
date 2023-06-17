package http

import (
	"github.com/begenov/backend/internal/config"
	v1 "github.com/begenov/backend/internal/delivery/http/v1"
	"github.com/begenov/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	h.init(router)

	return router
}

func (h *Handler) init(router *gin.Engine) {
	handlerv1 := v1.NewHandler(h.service)
	api := router.Group("/api")
	{
		handlerv1.Init(api)
	}
}
