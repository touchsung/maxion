package repositories

import (
	"gorm.io/gorm"
	"github.com/touchsung/maxion-server/internal/core/domain"
)

type tradingRepository struct {
	db *gorm.DB
}

func NewTradingRepository(db *gorm.DB) *tradingRepository {
	return &tradingRepository{db: db}
}

func (r *tradingRepository) GetAllStocks() ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.Find(&stocks).Error
	return stocks, err
}

func (r *tradingRepository) GetStockBySymbol(symbol string) (*domain.Stock, error) {
	var stock domain.Stock
	err := r.db.Where("Symbol = ?", symbol).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *tradingRepository) GetAllTransactions() ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Preload("Stock").Find(&transactions).Error
	return transactions, err
}

func (r *tradingRepository) CreateTransaction(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *tradingRepository) UpdateTransactionStatus(id int64, status domain.TransactionStatus) error {
	return r.db.Model(&domain.Transaction{}).
		Where("TransactionId = ?", id).
		Update("StatusId", status).Error
}

func (r *tradingRepository) UpdateStock(stock *domain.Stock) error {
	return r.db.Model(&domain.Stock{}).
		Where("Symbol = ?", stock.Symbol).
		Updates(map[string]interface{}{
			"BidPrice":  stock.BidPrice,
			"BidVolume": stock.BidVolume,
			"AskPrice":  stock.AskPrice,
			"AskVolume": stock.AskVolume,
		}).Error
} 