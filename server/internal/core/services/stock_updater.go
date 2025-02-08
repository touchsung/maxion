package services

import (
	"context"
	"math/rand"
	"time"
	"github.com/touchsung/maxion-server/internal/core/ports"
	"github.com/touchsung/maxion-server/internal/core/domain"
)

type StockUpdater struct {
	stockRepo ports.StockRepository
	done      chan bool
}

func NewStockUpdater(stockRepo ports.StockRepository) *StockUpdater {
	return &StockUpdater{
		stockRepo: stockRepo,
		done:      make(chan bool),
	}
}

func (su *StockUpdater) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				su.updateStockPrices()
			case <-ctx.Done():
				ticker.Stop()
				close(su.done)
				return
			}
		}
	}()
}

func (su *StockUpdater) Stop() {
	<-su.done
}

func (su *StockUpdater) updateStockPrices() {
	stocks, err := su.stockRepo.GetAllStocks()
	if err != nil {
		return
	}

	for _, stock := range stocks {
		updatedStock := domain.Stock{
			StockID:     stock.StockID,
			Symbol:      stock.Symbol,
			BidPrice:    stock.BidPrice * (1 + (rand.Float64() - 0.5) * 0.01),
			AskPrice:    stock.AskPrice * (1 + (rand.Float64() - 0.5) * 0.01),
			BidVolume:   int(float64(stock.BidVolume) * (1 + (rand.Float64() - 0.5) * 0.2)),
			AskVolume:   int(float64(stock.AskVolume) * (1 + (rand.Float64() - 0.5) * 0.2)),
			LastUpdated: time.Now(),
		}
		
		su.stockRepo.UpdateStock(&updatedStock)
	}
} 