package services

import (
	"context"

	"github.com/touchsung/maxion-server/internal/core/domain"
	"github.com/touchsung/maxion-server/internal/core/ports"
)

type tradingService struct {
	stockRepo       ports.StockRepository
	transactionRepo ports.TransactionRepository
	cacheService    *CacheService
}

func NewTradingService(
	stockRepo ports.StockRepository, 
	transactionRepo ports.TransactionRepository,
	cacheService *CacheService,
) ports.TradingService {
	return &tradingService{
		stockRepo:       stockRepo,
		transactionRepo: transactionRepo,
		cacheService:    cacheService,
	}
}

func (s *tradingService) GetAllStocks() ([]domain.Stock, error) {
	return s.stockRepo.GetAllStocks()
}

func (s *tradingService) GetAllTransactions() ([]domain.Transaction, error) {
	return s.cacheService.GetAllTransactions(context.Background())
}

func (s *tradingService) CreateTransaction(tx *domain.Transaction) error {
	stock, err := s.stockRepo.GetStockBySymbol(tx.Symbol)
	if err != nil {
		return err
	}

	if tx.Type == domain.Buy {
		tx.Price = stock.AskPrice
	} else {
		tx.Price = stock.BidPrice
	}

	tx.TotalAmount = float64(tx.Quantity) * tx.Price

	return s.cacheService.CacheTransaction(tx)
}

func (s *tradingService) UpdateTransactionStatus(id int64, status domain.TransactionStatus) error {
	return s.cacheService.CacheTransactionUpdate(id, status)
} 