package store

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type WalletStore struct {
	db *sql.DB
	mx sync.RWMutex
}

type WalletData struct {
	Id      string  `json:"id"`
	Balance float64 `json:"balance"`
	UserId  int64   `json:"userId"`
}

func (w *WalletStore) Deposit(ctx context.Context, Id string, amount float64) (*float64, error) {
	w.mx.Lock()
	defer w.mx.Unlock()

	query := `
		INSERT INTO wallet_operations(wallet_id, operation_type, amount, balance)
		VALUES ($1, 0, $2,
		COALESCE((SELECT balance FROM wallet_operations WHERE wallet_id=$1 ORDER BY id DESC LIMIT 1), 0) + $2);
	`
	var data WalletData
	_, err := w.db.ExecContext(ctx, query, Id, amount)
	if err != nil {
		return nil, err
	}

	query2 := `
		UPDATE wallets SET balance=
		(SELECT balance FROM wallet_operations WHERE wallet_id=$1 ORDER BY id DESC LIMIT 1)
		WHERE id=$1 RETURNING balance;
	`

	if err := w.db.QueryRowContext(ctx, query2, Id).Scan(&data.Balance); err != nil {
		return nil, err
	}

	return &data.Balance, err
}

func (w *WalletStore) Withdraw(ctx context.Context, Id string, amount float64) (*float64, error) {
	w.mx.Lock()
	defer w.mx.Unlock()

	query := `
		INSERT INTO wallet_operations(wallet_id, operation_type, amount, balance)
		VALUES ($1, 0, $2,
		COALESCE((SELECT balance FROM wallet_operations WHERE wallet_id=$1 ORDER BY id DESC LIMIT 1), 0) - $2);
	`
	var data WalletData
	_, err := w.db.ExecContext(ctx, query, Id, amount)
	if err != nil {
		return nil, err
	}

	query2 := `
		UPDATE wallets SET balance =
		(SELECT balance FROM wallet_operations WHERE wallet_id=$1 ORDER BY id DESC LIMIT 1)
		WHERE id=$1 RETURNING balance;
	`

	if err := w.db.QueryRowContext(ctx, query2, Id).Scan(&data.Balance); err != nil {
		return nil, err
	}

	return &data.Balance, err
}

func (w *WalletStore) GetWallet(ctx context.Context, Id string) (*WalletData, error) {
	w.mx.RLock()
	defer w.mx.RUnlock()

	query := `
	SELECT id, balance, user_id FROM wallets WHERE id=$1 FOR UPDATE;
	`
	var data WalletData
	err := w.db.QueryRowContext(ctx, query, Id).Scan(&data.Id, &data.Balance, &data.UserId)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &data, nil
}

func (w *WalletStore) CreateWallet(ctx context.Context, id *int64) (*string, error) {
	w.mx.RLock()
	defer w.mx.RUnlock()

	walletId := uuid.New()
	query := `
		INSERT INTO wallets (id, balance, user_id)
		VALUES ($1, $2, $3) RETURNING id;
	`
	var wallet WalletData
	err := w.db.QueryRowContext(ctx, query, walletId, 0, id).Scan(&wallet.Id)
	if err != nil {
		return nil, err
	}
	return &wallet.Id, nil
}
