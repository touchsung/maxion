package services

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/touchsung/maxion-server/internal/core/domain"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) UpdateTransactionStatus(id int64, status domain.TransactionStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetAllTransactions() ([]domain.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func setupRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}

func cleanupRedis(client *redis.Client, t *testing.T) {
	err := client.FlushDB(context.Background()).Err()
	if err != nil {
		t.Errorf("Failed to flush Redis DB: %v", err)
	}
	client.Close()
}

func TestCacheService_CacheTransaction(t *testing.T) {
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	mockDB := new(MockTransactionRepository)
	cacheService := NewCacheService(redisClient, mockDB)

	tx := &domain.Transaction{
		TransactionID: 1,
		Symbol:        "AAPL",
		Type:          domain.Buy,
		Status:        domain.Pending,
		Quantity:      100,
		Price:         150.50,
		TotalAmount:   15050.00,
		OrderTime:     time.Now(),
	}

	err := cacheService.CacheTransaction(tx)
	assert.NoError(t, err)

	keys, err := redisClient.Keys(context.Background(), PENDING_CREATE_PREFIX+"*").Result()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keys))
}

func TestCacheService_CacheTransactionUpdate(t *testing.T) {
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	mockDB := new(MockTransactionRepository)
	cacheService := NewCacheService(redisClient, mockDB)

	err := cacheService.CacheTransactionUpdate(1, domain.Completed)
	assert.NoError(t, err)

	keys, err := redisClient.Keys(context.Background(), PENDING_UPDATE_PREFIX+"*").Result()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keys))
}

func TestCacheService_GetAllTransactions(t *testing.T) {
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	mockDB := new(MockTransactionRepository)
	cacheService := NewCacheService(redisClient, mockDB)

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedTxs := []domain.Transaction{
		{
			TransactionID: 1,
			Symbol:        "AAPL",
			Type:          domain.Buy,
			Status:        domain.Completed,
			Quantity:      100,
			Price:         150.50,
			TotalAmount:   15050.00,
			OrderTime:     fixedTime,
			Stock: domain.Stock{
				Symbol:      "AAPL",
				LastUpdated: fixedTime,
			},
		},
		{
			TransactionID: 2,
			Symbol:        "GOOGL",
			Type:          domain.Sell,
			Status:        domain.Pending,
			Quantity:      50,
			Price:         2800.75,
			TotalAmount:   140037.50,
			OrderTime:     fixedTime,
			Stock: domain.Stock{
				Symbol:      "GOOGL",
				LastUpdated: fixedTime,
			},
		},
	}

	mockDB.On("GetAllTransactions").Return(expectedTxs, nil)

	ctx := context.Background()
	txs, err := cacheService.GetAllTransactions(ctx)
	assert.NoError(t, err)

	for i := range expectedTxs {
		assert.Equal(t, expectedTxs[i].TransactionID, txs[i].TransactionID)
		assert.Equal(t, expectedTxs[i].Symbol, txs[i].Symbol)
		assert.Equal(t, expectedTxs[i].Type, txs[i].Type)
		assert.Equal(t, expectedTxs[i].Status, txs[i].Status)
		assert.Equal(t, expectedTxs[i].Quantity, txs[i].Quantity)
		assert.Equal(t, expectedTxs[i].Price, txs[i].Price)
		assert.Equal(t, expectedTxs[i].TotalAmount, txs[i].TotalAmount)
	}

	cachedTxs, err := cacheService.GetAllTransactions(ctx)
	assert.NoError(t, err)

	for i := range expectedTxs {
		assert.Equal(t, expectedTxs[i].TransactionID, cachedTxs[i].TransactionID)
		assert.Equal(t, expectedTxs[i].Symbol, cachedTxs[i].Symbol)
		assert.Equal(t, expectedTxs[i].Type, cachedTxs[i].Type)
		assert.Equal(t, expectedTxs[i].Status, cachedTxs[i].Status)
		assert.Equal(t, expectedTxs[i].Quantity, cachedTxs[i].Quantity)
		assert.Equal(t, expectedTxs[i].Price, cachedTxs[i].Price)
		assert.Equal(t, expectedTxs[i].TotalAmount, cachedTxs[i].TotalAmount)
	}

	mockDB.AssertNumberOfCalls(t, "GetAllTransactions", 1)
}

func TestCacheService_BackgroundSync(t *testing.T) {
	redisClient := setupRedis(t)
	defer cleanupRedis(redisClient, t)

	mockDB := new(MockTransactionRepository)
	cacheService := NewCacheService(redisClient, mockDB)

	tx := &domain.Transaction{
		TransactionID: 1,
		Symbol:        "AAPL",
		Type:          domain.Buy,
		Status:        domain.Pending,
		Quantity:      100,
		Price:         150.50,
		TotalAmount:   15050.00,
		OrderTime:     time.Now(),
	}

	mockDB.On("CreateTransaction", mock.AnythingOfType("*domain.Transaction")).Return(nil)
	mockDB.On("UpdateTransactionStatus", int64(1), domain.Completed).Return(nil)

	err := cacheService.CacheTransaction(tx)
	assert.NoError(t, err)

	err = cacheService.CacheTransactionUpdate(1, domain.Completed)
	assert.NoError(t, err)

	time.Sleep(SYNC_INTERVAL + time.Second)

	createKeys, err := redisClient.Keys(context.Background(), PENDING_CREATE_PREFIX+"*").Result()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(createKeys))

	updateKeys, err := redisClient.Keys(context.Background(), PENDING_UPDATE_PREFIX+"*").Result()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(updateKeys))

	mockDB.AssertExpectations(t)
}
