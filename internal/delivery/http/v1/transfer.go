package v1

import (
	"fmt"
	"net/http"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pkg/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTransferTxRoutes(api *gin.RouterGroup) {
	transfers := api.Group("/transfers")
	{
		transfers.POST("/create", h.createTransfer)
	}
}

type transferRequest struct {
	FromAccountID int    `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int    `json:"to_account_id" binding:"required,min=1"`
	Amount        int    `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (h *Handler) createTransfer(ctx *gin.Context) {
	var inp transferRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorect input"+err.Error())
		return
	}

	if !h.validAccount(ctx, inp.FromAccountID, inp.Currency) {
		return
	}

	if !h.validAccount(ctx, inp.ToAccountID, inp.Currency) {
		return
	}

	arg := domain.TransferTxParams{
		FromAccountID: inp.FromAccountID,
		ToAccountID:   inp.ToAccountID,
		Amount:        inp.Amount,
	}

	result, err := h.service.TransferTx.TransferTx(ctx, arg)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Incorect db:"+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) validAccount(ctx *gin.Context, accountID int, currency string) bool {

	account, err := h.service.Account.GetAccountByID(ctx, accountID)
	if err != nil {
		if err == e.ErrRecordNotFound {
			newResponse(ctx, http.StatusNotFound, err.Error())
			return false
		}
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return false
	}

	return true

}
