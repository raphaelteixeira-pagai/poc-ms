package services

import (
	"context"
	"errors"
	"github.com/gocraft/dbr/v2"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/entities"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/repository"
)

type IWalletService interface {
	Fetch(ctx context.Context, owner string) (entities.Wallet, error)
	Withdraw(ctx context.Context, amount int64, owner string) error // saque
	Deposit(ctx context.Context, amount int64, owner string) error  // deposito
	Create(ctx context.Context, wallet entities.Wallet) error
	Delete(ctx context.Context, owner string) error
}

type wallet struct {
	repo repository.IWalletRepository
}

func (w *wallet) Fetch(ctx context.Context, owner string) (entities.Wallet, error) {
	walt, err := w.repo.Get(ctx, owner)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		return entities.Wallet{}, err
	}
	return walt, nil
}

func (w *wallet) Withdraw(ctx context.Context, amount int64, owner string) error {
	if amount <= 0 {
		return errors.New("amount should be greater than zero")
	}

	if owner == "" {
		return errors.New("owner should not be empty")
	}
	walt, err := w.repo.Get(ctx, owner)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		return err
	}

	if walt.Balance-amount < 0 {
		return errors.New("insufficient funds")
	}

	walt.Balance -= amount

	return w.repo.Update(ctx, walt)
}

func (w *wallet) Deposit(ctx context.Context, amount int64, owner string) error {
	if amount <= 0 {
		return errors.New("amount should be greater than zero")
	}

	if owner == "" {
		return errors.New("owner should not be empty")
	}

	walt, err := w.repo.Get(ctx, owner)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		return err
	}

	walt.Balance += amount

	return w.repo.Update(ctx, walt)
}

func (w *wallet) Create(ctx context.Context, wallet entities.Wallet) error {
	if wallet.Owner == "" {
		return errors.New("owner should not be empty")
	}

	walt, err := w.repo.Get(ctx, wallet.Owner)
	if err != nil && !errors.Is(err, dbr.ErrNotFound) {
		return err
	}

	if walt.ID != 0 {
		return errors.New("wallet already exists")
	}

	if err := w.repo.Create(ctx, wallet); err != nil {
		return errors.New("internal error cannot create new wallet")
	}

	return nil
}

func (w *wallet) Delete(ctx context.Context, owner string) error {
	if owner == "" {
		return errors.New("owner should not be empty")
	}

	_, err := w.repo.Get(ctx, owner)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return errors.New("wallet does not exists")
		}
		return err
	}

	if err := w.repo.Delete(ctx, owner); err != nil {
		return errors.New("internal error cannot create new wallet")
	}

	return nil
}

func NewWalletService(repo repository.IWalletRepository) IWalletService {
	return &wallet{repo: repo}
}
