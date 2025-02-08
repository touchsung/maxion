package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	CacheDuration time.Duration
	RateLimit     rate.Limit // requests per second
	BurstSize     int       // maximum burst size
}

type App struct {
	redis      *redis.Client
	config     Config
	limiter    *rate.Limiter
	router     *mux.Router
}

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

func newApp() *App {
	config := Config{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		CacheDuration: 5 * time.Minute,
		RateLimit:     100, 
		BurstSize:     150,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	return &App{
		redis:   rdb,
		config:  config,
		limiter: rate.NewLimiter(config.RateLimit, config.BurstSize),
		router:  mux.NewRouter(),
	}
}

func (a *App) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) getCachedData(ctx context.Context, key string) (string, error) {
	return a.redis.Get(ctx, key).Result()
}

func (a *App) setCachedData(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return a.redis.Set(ctx, key, data, a.config.CacheDuration).Err()
}

func (a *App) exampleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cacheKey := "example_key"

	if cachedData, err := a.getCachedData(ctx, cacheKey); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	}

	response := Response{
		Data:    map[string]interface{}{"message": "Hello, World!"},
		Success: true,
	}

	if err := a.setCachedData(ctx, cacheKey, response); err != nil {
		log.Printf("Error caching data: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *App) setupRoutes() {
	a.router.Use(a.rateLimitMiddleware)

	a.router.HandleFunc("/api/example", a.exampleHandler).Methods("GET")
}

func main() {
	app := newApp()
	app.setupRoutes()

	ctx := context.Background()
	if err := app.redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", app.router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
