# Maxion Trading Server

A high-performance trading server built with Go, featuring real-time stock data management and transaction processing with caching capabilities.

## Features

- Real-time stock data management with temporal support
- Transaction processing with status tracking
- Redis-based caching system for improved performance
- SQL Server database with temporal tables
- RESTful API endpoints for trading operations
- Docker containerization for easy deployment

## Tech Stack

- **Go** with Fiber web framework
- **SQL Server** for persistent storage
- **Redis** for caching
- **Docker** for containerization
- **GORM** for database operations

## Prerequisites

- Docker and Docker Compose
- Go 1.22 or higher (for local development)

## Start the application using Docker Compose

```bash
docker-compose up -d
```

The server will be available at `http://localhost:3000`

## API Endpoints

### Stocks

- `GET /stocks` - Get all available stocks

### Transactions

- `GET /transactions` - Get all transactions
- `POST /transactions` - Create a new transaction
- `PUT /transactions/:id/status` - Update transaction status

## Database Schema

The database includes the following main tables:

- `Stocks` - Stock market data with temporal support
- `Transactions` - Trading transaction records
- `TransactionTypes` - Transaction type enumerations (BUY/SELL)
- `TransactionStatus` - Transaction status enumerations

## Architecture

The project follows a clean architecture pattern with the following components:

- **Handlers** - HTTP request handlers
- **Services** - Business logic implementation
- **Repositories** - Data access layer
- **Domain** - Core business entities
- **Ports** - Interface definitions

## Caching Strategy

The application implements a write-through caching strategy using Redis:

- Transaction creations and updates are first cached
- Background process syncs cached data to the database
- Cache duration: 30 seconds
- Sync interval: 10 seconds
