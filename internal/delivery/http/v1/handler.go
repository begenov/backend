package v1

import (
	"github.com/begenov/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	v1 := api.Group("/v1")

	{
		h.initAccountsRoutes(v1)
		h.initTransferTxRoutes(v1)
		h.initUsersRoutes(v1)
	}
}
