package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yanpavel/wallet_api/internal/env"
	"github.com/yanpavel/wallet_api/internal/store"
)

type contextKey string

const UserKey contextKey = "userID"

var JwtSecret string = "nosecrets"

func CreateJWT(secret []byte, userID int64) (string, error) {
	expiration := time.Second * time.Duration(env.GetInt("JWTExpirationTime", 900))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   strconv.Itoa(int(userID)),
		"expireAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetTokenFromRequest(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func GetUserIdFromContext(ctx context.Context) int64 {
	userId, ok := ctx.Value(UserKey).(int64)
	if !ok {
		return -1
	}

	return userId
}

func WithAuthJWT(handlerFunc http.HandlerFunc, store store.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			unauthorized(w, r, err)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			unauthorized(w, r, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		str := claims["userID"].(string)

		userID, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			unauthorized(w, r, err)
			log.Printf("failed to convert userID to omt")
			return
		}

		u, err := store.UsersStore.GetUserByID(r.Context(), userID)
		if err != nil {
			unauthorized(w, r, err)
			log.Printf("failed to get user by id: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, u.Id)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.GetString("JWTSecret", JwtSecret)), nil
	})
}
