package service

import (
	"context"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/internal/repository"
)

type TransferTxService struct {
	repo repository.Tx
}

func NewTransferService(repo repository.Tx) *TransferTxService {
	return &TransferTxService{
		repo: repo,
	}
}

func (s *TransferTxService) TransferTx(ctx context.Context, arg domain.TransferTxParams) (domain.TransferTxResult, error) {
	return s.repo.TransferTx(ctx, arg)
}
