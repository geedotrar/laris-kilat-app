package repository

import (
	"context"
	"errors"
	"time"

	"github.com/geedotrar/lariskilat/checkout-engine/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) NaiveDecrementStock(ctx context.Context, itemID, userID string) error {
	var currentStock int
	err := r.db.QueryRow(ctx, "SELECT stock FROM flash_items WHERE id = $1", itemID).Scan(&currentStock)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrItemNotFound
	}
	if err != nil {
		return err
	}
	if currentStock <= 0 {
		return domain.ErrItemSoldOut
	}

	time.Sleep(15 * time.Millisecond)

	_, err = r.db.Exec(ctx, "UPDATE flash_items SET stock = stock - 1, updated_at = NOW() WHERE id = $1", itemID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, "INSERT INTO stock_reservations (item_id, user_id) VALUES ($1, $2)", itemID, userID)
	return err
}

func (r *PostgresRepository) ResetStock(ctx context.Context, itemID string, initialStock int) error {
	_, err := r.db.Exec(ctx, "UPDATE flash_items SET stock = $1, updated_at = NOW() WHERE id = $2", initialStock, itemID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, "TRUNCATE TABLE stock_reservations RESTART IDENTITY")
	return err
}
