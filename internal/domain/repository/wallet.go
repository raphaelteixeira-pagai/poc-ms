package repository

import (
	"context"
	"errors"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/entities"
	"github.com/raphaelteixeira-pagai/poc-ms/pkg/database"
)

type IWalletRepository interface {
	Get(ctx context.Context, owner string) (entities.Wallet, error)
	Update(ctx context.Context, wallet entities.Wallet) error
	Create(ctx context.Context, wallet entities.Wallet) error
	Delete(ctx context.Context, owner string) error
}

var (
	ErrInternal = errors.New("fatal error")
)

type wallet struct {
	pool database.DBPool
}

func (w *wallet) Get(ctx context.Context, owner string) (entities.Wallet, error) {
	sess, err := w.pool.Acquire()
	if err != nil {
		return entities.Wallet{}, ErrInternal
	}
	defer w.pool.Release(sess)

	var res entities.Wallet
	errc := sess.Select("balance").
		From("wallet").
		Where("owner = ?", owner).
		LoadOneContext(ctx, &res)
	return res, errc
}

func (w *wallet) Update(ctx context.Context, wallet entities.Wallet) error {
	sess, err := w.pool.Acquire()
	if err != nil {
		return ErrInternal
	}
	defer w.pool.Release(sess)

	_, err = sess.Update("wallet").
		Set("balance", wallet.Balance).
		Where("id = ?", wallet.ID).
		ExecContext(ctx)
	return err
}

func (w *wallet) Create(ctx context.Context, wallet entities.Wallet) error {
	sess, err := w.pool.Acquire()
	if err != nil {
		return ErrInternal
	}
	defer w.pool.Release(sess)

	_, err = sess.InsertInto("wallet").
		Pair("balance", wallet.Balance).
		Pair("owner", wallet.Owner).
		ExecContext(ctx)
	return err
}

func (w *wallet) Delete(ctx context.Context, owner string) error {
	sess, err := w.pool.Acquire()
	if err != nil {
		return ErrInternal
	}
	defer w.pool.Release(sess)

	_, err = sess.DeleteFrom("wallet").
		Where("owner = ?", owner).
		ExecContext(ctx)
	return err
}

func NewWalletRepository(pool database.DBPool) IWalletRepository {
	return &wallet{pool: pool}
}
