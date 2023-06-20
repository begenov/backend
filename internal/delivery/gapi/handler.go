package gapi

import (
	"github.com/begenov/backend/internal/service"
	"github.com/begenov/backend/pb"
	"github.com/begenov/backend/pkg/auth"
)

type Handler struct {
	pb.UnimplementedSimpleBankServer
	service *service.Service
	token   auth.TokenManager
}

func NewHandler(service *service.Service, token auth.TokenManager) *Handler {
	return &Handler{
		service: service,
		token:   token,
	}
}
