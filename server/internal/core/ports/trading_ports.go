package ports

import (
	"github.com/touchsung/maxion-server/internal/core/domain"
)

type StockRepository interface {
	GetAllStocks() ([]domain.Stock, error)
	GetStockBySymbol(symbol string) (*domain.Stock, error)
	UpdateStock(stock *domain.Stock) error
}

type TransactionRepository interface {
	GetAllTransactions() ([]domain.Transaction, error)
	CreateTransaction(tx *domain.Transaction) error
	UpdateTransactionStatus(id int64, status domain.TransactionStatus) error
}

type TradingService interface {
	GetAllStocks() ([]domain.Stock, error)
	GetAllTransactions() ([]domain.Transaction, error)
	CreateTransaction(tx *domain.Transaction) error
	UpdateTransactionStatus(id int64, status domain.TransactionStatus) error
} 