package handler

import (
	"errors"

	"github.com/geedotrar/lariskilat/checkout-engine/internal/domain"
	"github.com/geedotrar/lariskilat/checkout-engine/internal/service"
	"github.com/gofiber/fiber/v3"
)

type CheckoutHandler struct {
	checkoutService *service.CheckoutService
}

func NewCheckoutHandler(checkoutService *service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{checkoutService: checkoutService}
}

func (h *CheckoutHandler) HandleNaiveCheckout(c fiber.Ctx) error {
	var req domain.CheckoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	res, err := h.checkoutService.ProcessNaiveCheckout(c.Context(), req)
	if err != nil {
		return mapDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *CheckoutHandler) HandleAtomicCheckout(c fiber.Ctx) error {
	var req domain.CheckoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	res, err := h.checkoutService.ProcessAtomicCheckout(c.Context(), req)
	if err != nil {
		return mapDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *CheckoutHandler) HandleReset(c fiber.Ctx) error {
	if err := h.checkoutService.ResetSimulation(c.Context(), "item-sneaker-01", 50); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Postgres and Redis stock reset to 50"})
}

func mapDomainError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrItemNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, domain.ErrItemSoldOut), errors.Is(err, domain.ErrAlreadyPurchased):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
