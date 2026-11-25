package service

import (
	"context"

	"github.com/geedotrar/lariskilat/checkout-engine/internal/domain"
)

type StockStorage interface {
	NaiveDecrementStock(ctx context.Context, itemID, userID string) error
	ResetStock(ctx context.Context, itemID string, initialStock int) error
}

type AtomicStockCache interface {
	AtomicReserve(ctx context.Context, itemID, userID string) (int, error)
	PrewarmStock(ctx context.Context, itemID string, stock int) error
}

type CheckoutService struct {
	pgRepo    StockStorage
	redisRepo AtomicStockCache
}

func NewCheckoutService(pgRepo StockStorage, redisRepo AtomicStockCache) *CheckoutService {
	return &CheckoutService{pgRepo: pgRepo, redisRepo: redisRepo}
}

// naive without locking
func (s *CheckoutService) ProcessNaiveCheckout(ctx context.Context, req domain.CheckoutRequest) (*domain.CheckoutResult, error) {
	if err := s.pgRepo.NaiveDecrementStock(ctx, req.ItemID, req.UserID); err != nil {
		return nil, err
	}
	return &domain.CheckoutResult{
		Status: "RESERVED_NAIVE",
		UserID: req.UserID,
		ItemID: req.ItemID,
	}, nil
}

func (s *CheckoutService) ProcessAtomicCheckout(ctx context.Context, req domain.CheckoutRequest) (*domain.CheckoutResult, error) {
	remaining, err := s.redisRepo.AtomicReserve(ctx, req.ItemID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &domain.CheckoutResult{
		Status:         "RESERVED_ATOMIC",
		UserID:         req.UserID,
		ItemID:         req.ItemID,
		RemainingStock: remaining,
	}, nil
}

func (s *CheckoutService) ResetSimulation(ctx context.Context, itemID string, stock int) error {
	if err := s.pgRepo.ResetStock(ctx, itemID, stock); err != nil {
		return err
	}
	return s.redisRepo.PrewarmStock(ctx, itemID, stock)
}
