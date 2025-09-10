package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yanpavel/wallet_api/internal/store"
)

func TestUserLoginHandler(t *testing.T) {
	store := store.NewMockStorage()
	cfg := config{}
	app := application{
		store:  store,
		config: cfg,
	}

	t.Run("status OK test", func(t *testing.T) {
		payload := LoginPayload{
			Login:    "test",
			Password: "1234",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/login", app.loginHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d",
				http.StatusOK, rr.Code)
		}
	})

	t.Run("bad request test", func(t *testing.T) {
		payload := LoginPayload{
			Login: "test",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/login", app.loginHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
	})
}

func TestUserRegisterHandler(t *testing.T) {
	store := store.NewMockStorage2()
	cfg := config{}
	app := application{
		store:  store,
		config: cfg,
	}

	t.Run("status OK test", func(t *testing.T) {
		payload := LoginPayload{
			Login:    "test",
			Password: "1234",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/register", app.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d",
				http.StatusCreated, rr.Code)
		}
	})

	t.Run("bad request test", func(t *testing.T) {
		payload := LoginPayload{
			Login: "test",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/register", app.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
	})
}
