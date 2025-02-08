# Maxion Trading Platform

A full-stack trading platform featuring real-time stock data management and transaction processing, built with Go and React.

## Features

- Real-time stock data management with temporal support
- Transaction processing with status tracking
- Redis-based caching system for improved performance
- SQL Server database with temporal tables
- RESTful API endpoints for trading operations
- Dockerized deployment

## Tech Stack

### Backend

- Go with Fiber web framework
- SQL Server for persistent storage
- Redis for caching
- GORM for database operations
- Docker for containerization

### Frontend

- React with TypeScript
- Vite for build tooling
- TailwindCSS for styling
- React Query for data fetching

## Prerequisites

- Docker and Docker Compose
- Go 1.22 or higher (for local development)
- Node.js 20+ (for frontend development)

## Getting Started

1. Clone the repository
2. Start the application using Docker Compose:

```bash
docker-compose up -d
```

The application will be available at:

- Frontend: http://localhost:80
- Backend API: http://localhost:3000

## API Endpoints

### Stocks

- `GET /stocks` - Get all available stocks

### Transactions

- `GET /transactions` - Get all transactions
- `POST /transactions` - Create a new transaction
- `PUT /transactions/:id/status` - Update transaction status

## Project Structure

### Backend (`/server`)

- `/cmd` - Application entry points
- `/internal` - Internal packages
  - `/core` - Domain models and interfaces
  - `/handlers` - HTTP request handlers
  - `/repositories` - Data access layer
  - `/services` - Business logic
  - `/config` - Configuration

### Frontend (`/client`)

- Standard Vite + React + TypeScript structure
- TailwindCSS for styling

## Database Schema

The database includes the following main tables:

- `Stocks` - Stock market data with temporal support
- `Transactions` - Trading transaction records
- `TransactionTypes` - Transaction type enumerations (BUY/SELL)
- `TransactionStatus` - Transaction status enumerations

## Caching Strategy

The application implements a write-through caching strategy using Redis:

- Transaction creations and updates are first cached
- Background process syncs cached data to the database
- Cache duration: 30 seconds
- Sync interval: 10 seconds
