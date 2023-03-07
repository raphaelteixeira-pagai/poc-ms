package services

import (
	"context"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/entities"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/repository"
)

type IWalletService interface {
	Withdraw(ctx context.Context, amount int64) error // saque
	Deposit(ctx context.Context, amount int64) error  // deposito
	Create(ctx context.Context, wallet entities.Wallet) error
	Delete(ctx context.Context, owner string) error
}

type wallet struct {
	repo repository.IWalletRepository
}

func (w *wallet) Withdraw(ctx context.Context, amount int64) error {
	//TODO implement me
	panic("implement me")
}

func (w *wallet) Deposit(ctx context.Context, amount int64) error {
	//TODO implement me
	panic("implement me")
}

func (w *wallet) Create(ctx context.Context, wallet entities.Wallet) error {
	//TODO implement me
	panic("implement me")
}

func (w *wallet) Delete(ctx context.Context, owner string) error {
	//TODO implement me
	panic("implement me")
}

func NewWalletService(repo repository.IWalletRepository) IWalletService {
	return &wallet{repo: repo}
}
