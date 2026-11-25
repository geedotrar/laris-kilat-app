package domain

type CheckoutRequest struct {
	UserID string `json:"user_id"`
	ItemID string `json:"item_id"`
}

type CheckoutResult struct {
	Status         string `json:"status"`
	UserID         string `json:"user_id"`
	ItemID         string `json:"item_id"`
	RemainingStock int    `json:"remaining_stock,omitempty"`
}