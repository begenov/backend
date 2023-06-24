package gapi

import (
	"context"
	"log"
	"time"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pb"
	"github.com/begenov/backend/pkg/e"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.ResponseAccount, error) {
	username := getUsernameFromContext(ctx)
	log.Println(username)

	arg := domain.CreateAccountParams{
		Owner:    username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := h.service.Account.CreateAccount(ctx, arg)
	if err != nil {

		if e.ErrorCode(err) == e.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "method CreateUser Already Exists", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)

	}

	createdTime := time.Now().Unix()
	response := &pb.ResponseAccount{
		ID:       int32(account.ID),
		Owner:    account.Owner,
		Balance:  int32(account.Balance),
		Currency: account.Currency,
		CreatedAt: &timestamp.Timestamp{
			Seconds: createdTime,
			Nanos:   0,
		},
	}

	return response, nil
}

func getUsernameFromContext(ctx context.Context) string {
	if username, ok := ctx.Value("user").(string); ok {
		return username
	}
	return ""
}
