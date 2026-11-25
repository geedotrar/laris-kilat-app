package repository

import (
	"context"
	"fmt"

	"github.com/geedotrar/lariskilat/checkout-engine/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

var reserveLuaScript = redis.NewScript(`
	if redis.call("SISMEMBER", KEYS[2], ARGV[1]) == 1 then
		return -1
	end

	local stock = tonumber(redis.call("GET", KEYS[1]))
	if not stock or stock <= 0 then
		return 0
	end

	local remaining = redis.call("DECR", KEYS[1])
	redis.call("SADD", KEYS[2], ARGV[1])
	return remaining + 1
`)

func (r *RedisRepository) AtomicReserve(ctx context.Context, itemID, userID string) (int, error) {
	stockKey := fmt.Sprintf("lariskilat:stock:%s", itemID)
	buyersKey := fmt.Sprintf("lariskilat:buyers:%s", itemID)

	res, err := reserveLuaScript.Run(ctx, r.client, []string{stockKey, buyersKey}, userID).Int()
	if err != nil {
		return 0, err
	}

	switch res {
	case -1:
		return 0, domain.ErrAlreadyPurchased
	case 0:
		return 0, domain.ErrItemSoldOut
	default:
		return res - 1, nil
	}
}

func (r *RedisRepository) PrewarmStock(ctx context.Context, itemID string, stock int) error {
	stockKey := fmt.Sprintf("lariskilat:stock:%s", itemID)
	buyersKey := fmt.Sprintf("lariskilat:buyers:%s", itemID)

	pipe := r.client.Pipeline()
	pipe.Set(ctx, stockKey, stock, 0)
	pipe.Del(ctx, buyersKey)
	_, err := pipe.Exec(ctx)
	return err
}
