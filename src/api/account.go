package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/beatzoid/simple-bank/db/sqlc"
	"github.com/beatzoid/simple-bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

const (
	mismatchingAccountOwnerError = "account doesn't belong to the authenticated user"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,isValidCurrency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	// ShouldBindJSON checks that the request body is a valid JSON and
	// that it has the correct fields as defined by the createAccountRequest struct.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	createdAccount, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createdAccount)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	// ShouldBindUri checks that the URI is valid and
	// that it has the correct fields as defined by the getAccountRequest struct.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify the authenticated user owns the account their trying to get
	account, statusCode, err := verifyAccountOwner(ctx, server.store, req.ID)

	if err != nil {
		ctx.JSON(statusCode, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest

	// ShouldBindQuery checks that the query values are valid and
	// that they have the correct fields as defined by the listAccountsRequest struct.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required,min=0"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest

	// ShouldBindJSON checks that the request body is a valid JSON and
	// that it has the correct fields as defined by the updateAccountRequest struct.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify the authenticated user owns the account their trying to update
	_, statusCode, err := verifyAccountOwner(ctx, server.store, req.ID)

	if err != nil {
		ctx.JSON(statusCode, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	updatedAccount, err := server.store.UpdateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedAccount)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest

	// ShouldBindUri checks that the URI is valid and
	// that it has the correct fields as defined by the deleteAccountRequest struct.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify the authenticated user owns the account their trying to delete
	_, statusCode, err := verifyAccountOwner(ctx, server.store, req.ID)

	if err != nil {
		ctx.JSON(statusCode, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "successfully deleted account"})
}

func verifyAccountOwner(ctx *gin.Context, store db.Store, accountID int64) (db.Account, int, error) {
	// Get authorization data from the context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Get the account
	account, err := store.GetAccount(ctx, accountID)

	if err != nil {

		if err == sql.ErrNoRows {
			return db.Account{}, http.StatusNotFound, errors.New("account doesn't exist")
		}

		// Other internal error
		return db.Account{}, http.StatusInternalServerError, err
	}

	if account.Owner != authPayload.Username {
		// print("verify second if - account doesn't belong to the authenticated user", "\n")
		return db.Account{}, http.StatusUnauthorized, errors.New(mismatchingAccountOwnerError)
	}

	return account, http.StatusOK, nil
}
