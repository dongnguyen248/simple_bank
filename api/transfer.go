package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	if !server.validAccount(ctx, arg.FromAccountID, req.Currency) || !server.validAccount(ctx, arg.ToAccountID, req.Currency) {
		return
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAcount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d has currency %s, not %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return false
	}
	return true
}
