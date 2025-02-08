package domain

import (
	"time"
)

type Stock struct {
	StockID     int64     `gorm:"column:StockId;primaryKey;autoIncrement"`
	Symbol      string    `gorm:"column:Symbol;uniqueIndex"`
	BidPrice    float64   `gorm:"column:BidPrice"`
	BidVolume   int       `gorm:"column:BidVolume"`
	AskPrice    float64   `gorm:"column:AskPrice"`
	AskVolume   int       `gorm:"column:AskVolume"`
	LastUpdated time.Time `gorm:"column:LastUpdated"`
}

type TransactionType int

const (
	Buy  TransactionType = 1
	Sell TransactionType = 2
)

type TransactionStatus int

const (
	Pending   TransactionStatus = 1
	Completed TransactionStatus = 2
	Cancelled TransactionStatus = 3
	Failed    TransactionStatus = 4
)

type Transaction struct {
	TransactionID  int64             `gorm:"column:TransactionId;primaryKey;autoIncrement"`
	Symbol        string            `gorm:"column:Symbol"`
	Type          TransactionType   `gorm:"column:TypeId"`
	Status        TransactionStatus `gorm:"column:StatusId"`
	Quantity      int              `gorm:"column:Quantity"`
	Price         float64          `gorm:"column:Price"`
	TotalAmount   float64          `gorm:"column:TotalAmount"`
	OrderTime     time.Time        `gorm:"column:OrderTime"`
	ExecutionTime *time.Time       `gorm:"column:ExecutionTime"`
	Notes         *string          `gorm:"column:Notes"`
	Stock         Stock            `gorm:"foreignKey:Symbol;references:Symbol"`
}

func (Stock) TableName() string {
	return "Stocks"
}

func (Transaction) TableName() string {
	return "Transactions"
}

