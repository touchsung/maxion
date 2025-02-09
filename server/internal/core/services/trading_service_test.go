package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/touchsung/maxion-server/internal/core/domain"
)

// MockStockRepository mocks the StockRepository interface
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) GetAllStocks() ([]domain.Stock, error) {
	args := m.Called()
	return args.Get(0).([]domain.Stock), args.Error(1)
}

func (m *MockStockRepository) GetStockBySymbol(symbol string) (*domain.Stock, error) {
	args := m.Called(symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Stock), args.Error(1)
}

func (m *MockStockRepository) UpdateStock(stock *domain.Stock) error {
	args := m.Called(stock)
	return args.Error(0)
}

func TestTradingService_GetAllStocks(t *testing.T) {
	mockStockRepo := new(MockStockRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	cacheService := NewCacheService(redisClient, mockTransactionRepo)
	tradingService := NewTradingService(mockStockRepo, mockTransactionRepo, cacheService)

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedStocks := []domain.Stock{
		{
			StockID:     1,
			Symbol:      "AAPL",
			BidPrice:    150.00,
			BidVolume:   1000,
			AskPrice:    150.50,
			AskVolume:   800,
			LastUpdated: fixedTime,
		},
		{
			StockID:     2,
			Symbol:      "GOOGL",
			BidPrice:    2800.00,
			BidVolume:   500,
			AskPrice:    2801.00,
			AskVolume:   300,
			LastUpdated: fixedTime,
		},
	}

	mockStockRepo.On("GetAllStocks").Return(expectedStocks, nil)

	stocks, err := tradingService.GetAllStocks()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(expectedStocks), len(stocks))
	for i, stock := range stocks {
		assert.Equal(t, expectedStocks[i].StockID, stock.StockID)
		assert.Equal(t, expectedStocks[i].Symbol, stock.Symbol)
		assert.Equal(t, expectedStocks[i].BidPrice, stock.BidPrice)
		assert.Equal(t, expectedStocks[i].AskPrice, stock.AskPrice)
		assert.Equal(t, expectedStocks[i].BidVolume, stock.BidVolume)
		assert.Equal(t, expectedStocks[i].AskVolume, stock.AskVolume)
	}
}

func TestTradingService_GetAllTransactions(t *testing.T) {
	mockStockRepo := new(MockStockRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	cacheService := NewCacheService(redisClient, mockTransactionRepo)
	tradingService := NewTradingService(mockStockRepo, mockTransactionRepo, cacheService)

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedTransactions := []domain.Transaction{
		{
			TransactionID: 1,
			Symbol:        "AAPL",
			Type:          domain.Buy,
			Status:        domain.Completed,
			Quantity:      100,
			Price:         150.50,
			TotalAmount:   15050.00,
			OrderTime:     fixedTime,
		},
	}

	mockTransactionRepo.On("GetAllTransactions").Return(expectedTransactions, nil)

	transactions, err := tradingService.GetAllTransactions()

	assert.NoError(t, err)
	assert.Equal(t, len(expectedTransactions), len(transactions))
	for i, tx := range transactions {
		assert.Equal(t, expectedTransactions[i].TransactionID, tx.TransactionID)
		assert.Equal(t, expectedTransactions[i].Symbol, tx.Symbol)
		assert.Equal(t, expectedTransactions[i].Type, tx.Type)
		assert.Equal(t, expectedTransactions[i].Status, tx.Status)
		assert.Equal(t, expectedTransactions[i].Quantity, tx.Quantity)
		assert.Equal(t, expectedTransactions[i].Price, tx.Price)
		assert.Equal(t, expectedTransactions[i].TotalAmount, tx.TotalAmount)
	}
}

func TestTradingService_CreateTransaction(t *testing.T) {
	mockStockRepo := new(MockStockRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	cacheService := NewCacheService(redisClient, mockTransactionRepo)
	tradingService := NewTradingService(mockStockRepo, mockTransactionRepo, cacheService)

	stock := &domain.Stock{
		StockID:   1,
		Symbol:    "AAPL",
		BidPrice:  150.00,
		AskPrice:  150.50,
		BidVolume: 1000,
		AskVolume: 800,
	}

	testCases := []struct {
		name          string
		transaction   *domain.Transaction
		expectedPrice float64
	}{
		{
			name: "Buy Transaction",
			transaction: &domain.Transaction{
				Symbol:   "AAPL",
				Type:     domain.Buy,
				Quantity: 100,
			},
			expectedPrice: 150.50, // AskPrice
		},
		{
			name: "Sell Transaction",
			transaction: &domain.Transaction{
				Symbol:   "AAPL",
				Type:     domain.Sell,
				Quantity: 100,
			},
			expectedPrice: 150.00, // BidPrice
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStockRepo.On("GetStockBySymbol", tc.transaction.Symbol).Return(stock, nil)

			err := tradingService.CreateTransaction(tc.transaction)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPrice, tc.transaction.Price)
			assert.Equal(t, tc.expectedPrice*float64(tc.transaction.Quantity), tc.transaction.TotalAmount)
		})
	}
}

func TestTradingService_UpdateTransactionStatus(t *testing.T) {
	mockStockRepo := new(MockStockRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	cacheService := NewCacheService(redisClient, mockTransactionRepo)
	tradingService := NewTradingService(mockStockRepo, mockTransactionRepo, cacheService)

	transactionID := int64(1)
	newStatus := domain.Completed

	err := tradingService.UpdateTransactionStatus(transactionID, newStatus)

	assert.NoError(t, err)

	keys, err := redisClient.Keys(context.Background(), PENDING_UPDATE_PREFIX+"*").Result()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keys))
}
