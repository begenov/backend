package gapi

import (
	"context"
	"database/sql"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pb"
	"github.com/begenov/backend/pkg/e"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	arg := domain.CreateUserParams{
		Username:       req.Username,
		HashedPassword: req.Password,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := h.service.User.CreateUser(ctx, arg)
	if err != nil {

		if e.ErrorCode(err) == e.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "method CreateUser Already Exists", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)

	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func (h *Handler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	res, err := h.service.User.GetUserByUsername(ctx, req.Username, req.Password)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, status.Errorf(codes.NotFound, "failed to create user: %v", err)
		case e.ErrInvalidToken:
			return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
		case e.ErrPassword:
			return nil, status.Errorf(codes.NotFound, "failed to create user: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
		}
	}

	rq := &pb.LoginUserResponse{
		User: &pb.User{
			Username:          res.User.Username,
			FullName:          res.User.FullName,
			PasswordChangedAt: timestamppb.New(res.User.PasswordChangedAt),
			CreatedAt:         timestamppb.New(res.User.CreatedAt),
		},
		AccessToken: res.AccessToken,
	}

	return rq, nil
}

func convertUser(user domain.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
