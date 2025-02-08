package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/touchsung/maxion-server/internal/core/domain"
	"github.com/touchsung/maxion-server/internal/core/ports"
)

type TradingHandlers struct {
	tradingService ports.TradingService
}

func NewTradingHandlers(tradingService ports.TradingService) *TradingHandlers {
	return &TradingHandlers{
		tradingService: tradingService,
	}
}

func (h *TradingHandlers) GetAllStocks(c *fiber.Ctx) error {
	stocks, err := h.tradingService.GetAllStocks()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch stocks"})
	}
	return c.JSON(stocks)
}

func (h *TradingHandlers) GetAllTransactions(c *fiber.Ctx) error {
	transactions, err := h.tradingService.GetAllTransactions()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch transactions"})
	}
	return c.JSON(transactions)
}

func (h *TradingHandlers) CreateTransaction(c *fiber.Ctx) error {
	var req struct {
		Symbol   string  `json:"symbol"`
		Type     int     `json:"type"`
		Quantity int     `json:"quantity"`
		Notes    *string `json:"notes,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tx := &domain.Transaction{
		Symbol:   req.Symbol,
		Type:     domain.TransactionType(req.Type),
		Quantity: req.Quantity,
		Notes:    req.Notes,
		Status:   domain.Pending,
	}

	if err := h.tradingService.CreateTransaction(tx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	return c.Status(201).JSON(tx)
}

func (h *TradingHandlers) UpdateTransactionStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	var req struct {
		Status int `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err = h.tradingService.UpdateTransactionStatus(int64(id), domain.TransactionStatus(req.Status))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update transaction"})
	}

	return c.SendStatus(200)
} 