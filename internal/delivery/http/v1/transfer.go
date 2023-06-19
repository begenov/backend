package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/begenov/backend/internal/domain"
	"github.com/begenov/backend/pkg/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTransferTxRoutes(api *gin.RouterGroup) {
	transfers := api.Group("/transfers", h.userIdentity)
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

	account, ok := h.validAccount(ctx, inp.FromAccountID, inp.Currency)
	if !ok {
		return
	}

	username := ctx.MustGet(userCtx).(string)

	if account.Owner != username {
		log.Println(account.Owner, username)
		newResponse(ctx, http.StatusUnauthorized, "from account doent't belong to the authenticated user")
		return
	}

	_, ok = h.validAccount(ctx, inp.ToAccountID, inp.Currency)
	if !ok {
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

func (h *Handler) validAccount(ctx *gin.Context, accountID int, currency string) (domain.Account, bool) {

	account, err := h.service.Account.GetAccountByID(ctx, accountID)
	if err != nil {
		if err == e.ErrRecordNotFound {
			newResponse(ctx, http.StatusNotFound, err.Error())
			return account, false
		}
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		newResponse(ctx, http.StatusBadRequest, err.Error())
		return account, false
	}

	return account, true

}
