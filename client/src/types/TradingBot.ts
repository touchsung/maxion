// Stock represents the current market data for a stock
export interface Stock {
    StockID: number;
    Symbol: string;
    BidPrice: number;
    BidVolume: number;
    AskPrice: number;
    AskVolume: number;
    LastUpdated: string; // ISO date string
}

// TransactionType represents the type of transaction (BUY/SELL)
export enum TransactionType {
    Buy = 1,
    Sell = 2
}

// TransactionStatus represents the status of a transaction
export enum TransactionStatus {
    Pending = 1,
    Completed = 2,
    Cancelled = 3,
    Failed = 4
}

// Transaction represents a trading transaction
export interface Transaction {
    TransactionID: number;
    Symbol: string;
    Type: TransactionType;
    Status: TransactionStatus;
    Quantity: number;
    Price: number;
    TotalAmount: number;
    OrderTime: string; // ISO date string
    ExecutionTime?: string; // Optional ISO date string
    Notes?: string;
    Stock: Stock;
}

export interface TransactionRequest {
    symbol: string;
    type: TransactionType;
    quantity: number;
    notes: string;
}
