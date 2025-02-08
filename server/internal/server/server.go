package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
	"github.com/touchsung/maxion-server/internal/handlers"
	"github.com/touchsung/maxion-server/internal/repositories"
	"github.com/touchsung/maxion-server/internal/core/services"
	"github.com/go-redis/redis/v8"
	"os"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"context"
)

type Server struct {
	app          *fiber.App
	db           *gorm.DB
	redis        *redis.Client
	handlers     *handlers.TradingHandlers
	cacheService *services.CacheService
	stockUpdater *services.StockUpdater
}

func NewServer(db *gorm.DB) *Server {
	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	// Initialize repositories
	tradingRepo := repositories.NewTradingRepository(db)

	// Initialize cache service
	cacheService := services.NewCacheService(redisClient, tradingRepo)

	// Initialize services
	tradingService := services.NewTradingService(tradingRepo, tradingRepo, cacheService)

	// Initialize handlers
	handlers := handlers.NewTradingHandlers(tradingService)

	// Initialize stock updater
	stockUpdater := services.NewStockUpdater(tradingRepo)

	return &Server{
		app:          fiber.New(),
		db:           db,
		redis:        redisClient,
		handlers:     handlers,
		cacheService: cacheService,
		stockUpdater: stockUpdater,
	}
}

func (s *Server) setupRoutes() {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	s.app.Use(logger.New())
	
	// Trading routes
	s.app.Get("/stocks", s.handlers.GetAllStocks)
	s.app.Get("/transactions", s.handlers.GetAllTransactions)
	s.app.Post("/transactions", s.handlers.CreateTransaction)
	s.app.Put("/transactions/:id/status", s.handlers.UpdateTransactionStatus)
}

func (s *Server) Start(addr string) error {
	// Start the stock updater
	ctx := context.Background()
	s.stockUpdater.Start(ctx)
	
	s.setupRoutes()
	return s.app.Listen(addr)
} 