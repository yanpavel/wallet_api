package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yanpavel/wallet_api/internal/store"
)

func TestGetWallets(t *testing.T) {
	store := store.NewMockStorage()
	cfg := config{}
	app := application{
		store:  store,
		config: cfg,
	}

	t.Run("status 200 test", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/wallet/3990064f-ec54-48a2-ad34-1e5e82c57b4b/", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), UserKey, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/wallet/{id}/", app.getBalanceHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d",
				http.StatusOK, rr.Code)
		}
	})

	t.Run("status 403 test", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/wallet/3990064f-ec54-48a2-ad34-1e5e82c57b4b/", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), UserKey, int64(2))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/wallet/{id}/", app.getBalanceHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status code %d, got %d",
				http.StatusForbidden, rr.Code)
		}
	})
}

func TestWalletOperations(t *testing.T) {
	store := store.NewMockStorage()
	cfg := config{}
	app := application{
		store:  store,
		config: cfg,
	}

	t.Run("status 200 test", func(t *testing.T) {
		payload := WalletOperationPayload{
			WalletId:      "8cfa863b-44b4-4fe1-82cc-ebc6223811a1",
			OperationType: 0,
			Amount:        100,
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), UserKey, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/wallet", app.changeBalanceHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d",
				http.StatusOK, rr.Code)
		}
	})

	t.Run("status 400 test", func(t *testing.T) {
		payload := WalletOperationPayload{
			WalletId:      "8cfa863b-44b4-4fe1-82cc-ebc6223811a1",
			OperationType: 1,
			Amount:        10000000,
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), UserKey, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/wallet", app.changeBalanceHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("status 403 test", func(t *testing.T) {
		payload := WalletOperationPayload{
			WalletId:      "8cfa863b-44b4-4fe1-82cc-ebc6223811a1",
			OperationType: 1,
			Amount:        10000000,
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), UserKey, int64(2))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/wallet", app.changeBalanceHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status code %d, got %d",
				http.StatusForbidden, rr.Code)
		}
	})
}
