package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/touchsung/maxion-server/internal/core/domain"
)

// MockTradingService implements ports.TradingService for testing
type MockTradingService struct {
	mock.Mock
}

func (m *MockTradingService) GetAllStocks() ([]domain.Stock, error) {
	args := m.Called()
	return args.Get(0).([]domain.Stock), args.Error(1)
}

func (m *MockTradingService) GetAllTransactions() ([]domain.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTradingService) CreateTransaction(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTradingService) UpdateTransactionStatus(id int64, status domain.TransactionStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func setupTest() (*fiber.App, *MockTradingService) {
	app := fiber.New()
	mockService := new(MockTradingService)
	handlers := NewTradingHandlers(mockService)

	app.Get("/stocks", handlers.GetAllStocks)
	app.Get("/transactions", handlers.GetAllTransactions)
	app.Post("/transactions", handlers.CreateTransaction)
	app.Put("/transactions/:id/status", handlers.UpdateTransactionStatus)

	return app, mockService
}

func TestGetAllStocks(t *testing.T) {
	// Setup
	app, mockService := setupTest()
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test data
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

	mockService.On("GetAllStocks").Return(expectedStocks, nil)

	// Execute
	req := httptest.NewRequest("GET", "/stocks", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result []domain.Stock
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedStocks), len(result))
	assert.Equal(t, expectedStocks[0].Symbol, result[0].Symbol)
	assert.Equal(t, expectedStocks[0].BidPrice, result[0].BidPrice)
}

func TestGetAllTransactions(t *testing.T) {
	// Setup
	app, mockService := setupTest()
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test data
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
		},
	}

	mockService.On("GetAllTransactions").Return(expectedTxs, nil)

	// Execute
	req := httptest.NewRequest("GET", "/transactions", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result []domain.Transaction
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedTxs), len(result))
	assert.Equal(t, expectedTxs[0].Symbol, result[0].Symbol)
	assert.Equal(t, expectedTxs[0].Type, result[0].Type)
}

func TestCreateTransaction(t *testing.T) {
	// Setup
	app, mockService := setupTest()

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid Transaction",
			requestBody: map[string]interface{}{
				"symbol":   "AAPL",
				"type":     1,
				"quantity": 100,
				"notes":    "Test transaction",
			},
			expectedStatus: 201,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectation if needed
			if tc.expectedStatus == 201 {
				mockService.On("CreateTransaction", mock.AnythingOfType("*domain.Transaction")).Return(nil).Once()
			}

			// Create request body
			jsonBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Execute
			resp, err := app.Test(req)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

func TestUpdateTransactionStatus(t *testing.T) {
	// Setup
	app, mockService := setupTest()

	// Test cases
	testCases := []struct {
		name           string
		transactionID  string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:          "Valid Update",
			transactionID: "1",
			requestBody: map[string]interface{}{
				"status": 2,
			},
			expectedStatus: 200,
		},
		{
			name:          "Invalid Transaction ID",
			transactionID: "invalid",
			requestBody: map[string]interface{}{
				"status": 2,
			},
			expectedStatus: 400,
		},
		{
			name:          "Invalid Status",
			transactionID: "1",
			requestBody: map[string]interface{}{
				"status": "invalid",
			},
			expectedStatus: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectation if needed
			if tc.expectedStatus == 200 {
				mockService.On("UpdateTransactionStatus", int64(1), domain.TransactionStatus(2)).Return(nil).Once()
			}

			// Create request body
			jsonBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("PUT", "/transactions/"+tc.transactionID+"/status", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Execute
			resp, err := app.Test(req)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
