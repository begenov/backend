package v1

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/begenov/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) initAccountsRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/accounts", h.userIdentity)
	{
		accounts.POST("/create", h.createAccount)
		accounts.GET("/:id", h.getAccountByID)
		accounts.GET("", h.listAccount)
	}
}

type createAccountRequest struct {
	Owner    string
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (h *Handler) createAccount(ctx *gin.Context) {

	var inp createAccountRequest
	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input"+err.Error())
		return
	}

	username := ctx.MustGet(userCtx).(string)
	log.Println(username)
	arg := domain.CreateAccountParams{
		Owner:    username,
		Currency: inp.Currency,
		Balance:  0,
	}

	account, err := h.service.Account.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				newResponse(ctx, http.StatusForbidden, "Incorrect db:"+err.Error())
				return
			}
		}
		newResponse(ctx, http.StatusInternalServerError, "Incorrect db:"+err.Error())
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func (h *Handler) getAccountByID(ctx *gin.Context) {

	var inp getAccountRequest
	if err := ctx.BindUri(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input:"+err.Error())
		return
	}

	account, err := h.service.Account.GetAccountByID(ctx, inp.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			newResponse(ctx, http.StatusNotFound, "Incorrect db:"+err.Error())
			return
		}
		newResponse(ctx, http.StatusInternalServerError, "Incorrect db:"+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=5,max=10"`
}

func (h *Handler) listAccount(ctx *gin.Context) {
	var inp listAccountRequest
	if err := ctx.BindQuery(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input:"+err.Error())
		return
	}

	arg := domain.ListAccountsParams{
		Limit:  inp.PageSize,
		Offset: (inp.PageID - 1) * inp.PageSize,
	}
	accounts, err := h.service.Account.ListAccounts(ctx, arg)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Incorrect db:"+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
