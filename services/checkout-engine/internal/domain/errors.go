package domain

import "errors"

var (
	ErrItemNotFound = errors.New("flash sale item not found")
	ErrItemSoldOut = errors.New("item is sold out")
	ErrAlreadyPurchased = errors.New("user has already purchased this flash sale item")
)

