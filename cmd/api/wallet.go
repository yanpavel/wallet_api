package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yanpavel/wallet_api/internal/store"
)

type WalletOperationPayload struct {
	WalletId      string  `json:"walletId" validate:"required"`
	OperationType int     `json:"operationType"`
	Amount        float64 `json:"amount" validate:"required"`
}

func (app *application) changeBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// parsing request body
	var payload WalletOperationPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// validating body fields
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// set context timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// check if wallet exists
	wallet, err := app.store.WalletStore.GetWallet(ctx, payload.WalletId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	// check if user has permission
	if wallet.UserId != GetUserIdFromContext(ctx) {
		log.Printf("%v and %v", wallet.UserId, GetUserIdFromContext(ctx))
		app.permissionDenied(w, r)
		return
	}

	var balance *float64

	// do operation via db
	switch payload.OperationType {
	case 0:
		balance, err = app.store.WalletStore.Deposit(ctx, payload.WalletId, payload.Amount)
	case 1:
		// checking if funds are enough
		if wallet.Balance < payload.Amount {
			app.badRequestError(w, r, errors.New("insufficient funds"))
			return
		}
		balance, err = app.store.WalletStore.Withdraw(ctx, payload.WalletId, payload.Amount)
	default:
		app.badRequestError(w, r, errors.New("invalid operation type"))
		return
	}

	// if operation is failed in db - internal error
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send response
	if err := app.jsonResponse(w, http.StatusOK, balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// parse param header
	param := chi.URLParam(r, "id")
	// checking if format is valid
	if _, err := uuid.Parse(param); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// set context timeout
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*8)
	defer cancel()

	// get wallet
	wallet, err := app.store.WalletStore.GetWallet(ctx, param)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	// check if user has permission
	if wallet.UserId != GetUserIdFromContext(ctx) {
		log.Printf("%v and %v", wallet.UserId, GetUserIdFromContext(ctx))
		app.permissionDenied(w, r)
	}

	// send response
	if err := app.jsonResponse(w, http.StatusOK, wallet.Balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
