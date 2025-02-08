package services

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
    "github.com/touchsung/maxion-server/internal/core/domain"
    "github.com/touchsung/maxion-server/internal/core/ports"
)

const (
    CACHE_DURATION = 30 * time.Second
    SYNC_INTERVAL = 15 * time.Second
    PENDING_CREATE_PREFIX = "pending_create_tx:"
    PENDING_UPDATE_PREFIX = "pending_update_tx:"
    ALL_TRANSACTIONS_KEY = "all_transactions"
)

type CacheService struct {
    redis *redis.Client
    db    ports.TransactionRepository
}

type syncResult struct {
    key string
    err error
}

func generateCreateKey(txID int64) string {
    now := time.Now().UTC().Format(time.RFC3339)
    return fmt.Sprintf("%s%d_%s", PENDING_CREATE_PREFIX, txID, now)
}

func generateUpdateKey(txID int64) string {
    now := time.Now().UTC().Format(time.RFC3339)
    return fmt.Sprintf("%s%d_%s", PENDING_UPDATE_PREFIX, txID, now)
}

func NewCacheService(redisClient *redis.Client, db ports.TransactionRepository) *CacheService {
    service := &CacheService{
        redis: redisClient,
        db:    db,
    }
    
    go service.startBackgroundSync()
    
    return service
}

func (s *CacheService) CacheTransaction(tx *domain.Transaction) error {
    txJSON, err := json.Marshal(tx)
    if err != nil {
        return fmt.Errorf("failed to marshal transaction: %w", err)
    }
    
    key := generateCreateKey(tx.TransactionID)
    if err := s.setCache(context.Background(), key, txJSON); err != nil {
        return err
    }

    s.redis.Del(context.Background(), ALL_TRANSACTIONS_KEY)
    
    return nil
}

func (s *CacheService) CacheTransactionUpdate(id int64, status domain.TransactionStatus) error {
    update := struct {
        ID     int64                   `json:"id"`
        Status domain.TransactionStatus `json:"status"`
    }{
        ID:     id,
        Status: status,
    }
    
    updateJSON, err := json.Marshal(update)
    if err != nil {
        return fmt.Errorf("failed to marshal update: %w", err)
    }
    
    key := generateUpdateKey(id)
    if err := s.setCache(context.Background(), key, updateJSON); err != nil {
        return err
    }

    s.redis.Del(context.Background(), ALL_TRANSACTIONS_KEY)
    
    return nil
}

func (s *CacheService) handleCreateSync(ctx context.Context, key string, txJSON string) syncResult {
    var tx domain.Transaction
    if err := json.Unmarshal([]byte(txJSON), &tx); err != nil {
        return syncResult{key, fmt.Errorf("failed to unmarshal transaction: %w", err)}
    }

    if err := s.db.CreateTransaction(&tx); err != nil {
        return syncResult{key, fmt.Errorf("failed to create transaction: %w", err)}
    }

    return syncResult{key, nil}
}

func (s *CacheService) handleUpdateSync(ctx context.Context, key string, updateJSON string) syncResult {
    var update struct {
        ID     int64                   `json:"id"`
        Status domain.TransactionStatus `json:"status"`
    }
    
    if err := json.Unmarshal([]byte(updateJSON), &update); err != nil {
        return syncResult{key, fmt.Errorf("failed to unmarshal update: %w", err)}
    }

    if err := s.db.UpdateTransactionStatus(update.ID, update.Status); err != nil {
        return syncResult{key, fmt.Errorf("failed to update transaction: %w", err)}
    }

    return syncResult{key, nil}
}

func (s *CacheService) syncCreates(ctx context.Context) {
    keys, err := s.redis.Keys(ctx, PENDING_CREATE_PREFIX+"*").Result()
    if err != nil {
        return
    }

    for _, key := range keys {
        txJSON, err := s.redis.Get(ctx, key).Result()
        if err != nil {
            continue
        }

        result := s.handleCreateSync(ctx, key, txJSON)
        if result.err == nil {
            s.redis.Del(ctx, result.key)
        }
    }
}

func (s *CacheService) syncUpdates(ctx context.Context) {
    keys, err := s.redis.Keys(ctx, PENDING_UPDATE_PREFIX+"*").Result()
    if err != nil {
        return
    }

    for _, key := range keys {
        updateJSON, err := s.redis.Get(ctx, key).Result()
        if err != nil {
            continue
        }

        result := s.handleUpdateSync(ctx, key, updateJSON)
        if result.err == nil {
            s.redis.Del(ctx, result.key)
        }
    }
}

func (s *CacheService) setCache(ctx context.Context, key string, value []byte) error {
    return s.redis.Set(ctx, key, value, CACHE_DURATION).Err()
}

func (s *CacheService) startBackgroundSync() {
    ticker := time.NewTicker(SYNC_INTERVAL)
    ctx := context.Background()

    for range ticker.C {
        s.syncCreates(ctx)
        s.syncUpdates(ctx)
    }
}

func (s *CacheService) GetAllTransactions(ctx context.Context) ([]domain.Transaction, error) {
    txJSON, err := s.redis.Get(ctx, ALL_TRANSACTIONS_KEY).Result()
    if err == nil {
        var transactions []domain.Transaction
        if err := json.Unmarshal([]byte(txJSON), &transactions); err != nil {
            return nil, fmt.Errorf("failed to unmarshal cached transactions: %w", err)
        }
        return transactions, nil
    }

    transactions, err := s.db.GetAllTransactions()
    if err != nil {
        return nil, err
    }

    txJSONBytes, err := json.Marshal(transactions)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal transactions: %w", err)
    }
    
    if err := s.setCache(ctx, ALL_TRANSACTIONS_KEY, txJSONBytes); err != nil {
        fmt.Printf("failed to cache transactions: %v\n", err)
    }

    return transactions, nil
} 