package main

import (
	"context"
	"log"
	"os"

	"github.com/geedotrar/lariskilat/checkout-engine/internal/handler"
	"github.com/geedotrar/lariskilat/checkout-engine/internal/repository"
	"github.com/geedotrar/lariskilat/checkout-engine/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func main() {
	_ = godotenv.Load("../../.env", ".env")

	ctx := context.Background()

	port := getEnv("PORT", "8081")
	pgDSN := getEnv("POSTGRES_DSN", "postgres://laris_admin:laris_secret@localhost:5432/lariskilat_db?sslmode=disable")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")

	pgConfig, err := pgxpool.ParseConfig(pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	pgConfig.MaxConns = 25

	pgPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		PoolSize: 50,
	})
	defer rdb.Close()

	pgRepo := repository.NewPostgresRepository(pgPool)
	redisRepo := repository.NewRedisRepository(rdb)
	checkoutSvc := service.NewCheckoutService(pgRepo, redisRepo)
	checkoutHandler := handler.NewCheckoutHandler(checkoutSvc)

	// Pre-warm stok redis
	_ = checkoutSvc.ResetSimulation(ctx, "item-sneaker-01", 50)

	app := fiber.New()
	v1 := app.Group("/api/v1")
	v1.Post("/checkout/naive", checkoutHandler.HandleNaiveCheckout)
	v1.Post("/checkout/atomic", checkoutHandler.HandleAtomicCheckout)
	v1.Post("/admin/reset", checkoutHandler.HandleReset)

	log.Printf("LarisKilat Checkout Engine running on :%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
