package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/yanpavel/wallet_api/internal/env"
	"github.com/yanpavel/wallet_api/internal/store"
)

type LoginPayload struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload LoginPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	//validate the payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// set context timeout
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*8)
	defer cancel()

	// check if the user exists
	u, err := app.store.UsersStore.GetUser(ctx, payload.Login)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestError(w, r, errors.New("Login/Password invalid"))
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	// check passwords from db and payload
	if !comparePasswords(u.Password, []byte(payload.Password)) {
		app.badRequestError(w, r, errors.New("Login/Password invalid"))
		return
	}

	// create jwt token
	secret := []byte(env.GetString("JWTSecret", JwtSecret))
	token, err := CreateJWT(secret, u.Id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send token in response
	if err := app.jsonResponse(w, http.StatusOK, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload LoginPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	//validate the payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	// check if the user exists
	_, err := app.store.UsersStore.GetUser(ctx, payload.Login)
	if err == nil {
		app.badRequestError(w, r, errors.New("user already exists"))
		return
	}

	hashedPassword, err := hashPassword(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// if it doesnt we create the new user
	id, err := app.store.UsersStore.CreateUser(ctx, payload.Login, hashedPassword)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	uuid, err := app.store.WalletStore.CreateWallet(ctx, id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	responseBody := make(map[string]any)
	responseBody["userId"] = id
	responseBody["walletId"] = uuid
	// send user id in response
	if err := app.jsonResponse(w, http.StatusCreated, responseBody); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
