package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	WalletStore interface {
		Deposit(context.Context, string, float64) (*float64, error)
		Withdraw(context.Context, string, float64) (*float64, error)
		GetWallet(context.Context, string) (*WalletData, error)
		CreateWallet(context.Context, *int64) (*string, error)
	}
	UsersStore interface {
		GetUser(context.Context, string) (*User, error)
		CreateUser(context.Context, string, string) (*int64, error)
		GetUserByID(context.Context, int64) (*User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		WalletStore: &WalletStore{db: db},
		UsersStore:  &UsersStore{db: db},
	}
}

func NewMockStorage() Storage {
	return Storage{
		WalletStore: &WalletMockStore{},
		UsersStore:  &UserMockStore{},
	}
}

func NewMockStorage2() Storage {
	return Storage{
		WalletStore: &WalletMockStore{},
		UsersStore:  &UserMockStore2{},
	}
}

type WalletMockStore struct {
	balances float64
}

func (wm *WalletMockStore) Deposit(ctx context.Context, uuid string, amount float64) (*float64, error) {
	wm.balances += amount
	return &wm.balances, nil
}

func (wm *WalletMockStore) Withdraw(context.Context, string, float64) (*float64, error) {
	return nil, nil
}

func (wm *WalletMockStore) GetWallet(context.Context, string) (*WalletData, error) {
	return &WalletData{
		UserId: 1,
		Id:     "3990064f-ec54-48a2-ad34-1e5e82c57b4b",
	}, nil
}

func (wm *WalletMockStore) CreateWallet(context.Context, *int64) (*string, error) {
	return nil, nil
}

type UserMockStore struct {
}

func (um *UserMockStore) GetUser(context.Context, string) (*User, error) {
	return &User{
		Id:       1,
		Username: "test",
		Password: "$2a$10$JUDX2Qc9Go9pF6/xkPcT6.gMYURvFatumZUm6qWIhDlm8U9F8foVu",
	}, nil
}
func (um *UserMockStore) CreateUser(context.Context, string, string) (*int64, error) {
	return nil, nil
}
func (um *UserMockStore) GetUserByID(context.Context, int64) (*User, error) {
	return &User{
		Id:       1,
		Username: "test",
		Password: "$2a$10$JUDX2Qc9Go9pF6/xkPcT6.gMYURvFatumZUm6qWIhDlm8U9F8foVu"}, nil
}

type UserMockStore2 struct {
}

func (um *UserMockStore2) GetUser(context.Context, string) (*User, error) {
	return nil, errors.New("error")
}
func (um *UserMockStore2) CreateUser(context.Context, string, string) (*int64, error) {
	id := int64(1)
	id2 := &id
	return id2, nil
}
func (um *UserMockStore2) GetUserByID(context.Context, int64) (*User, error) {
	return nil, errors.New("error")
}
