-- Create the database
CREATE DATABASE TradingBot;
GO

USE TradingBot;
GO

-- Create the Stocks table with temporal support
CREATE TABLE Stocks (
    StockId INT IDENTITY(1,1) PRIMARY KEY,
    Symbol VARCHAR(10) NOT NULL,
    BidPrice DECIMAL(18,4) NOT NULL,
    BidVolume INT NOT NULL,
    AskPrice DECIMAL(18,4) NOT NULL,
    AskVolume INT NOT NULL,
    LastUpdated DATETIME2(7) GENERATED ALWAYS AS ROW START,
    ValidTo DATETIME2(7) GENERATED ALWAYS AS ROW END,
    PERIOD FOR SYSTEM_TIME (LastUpdated, ValidTo)
)
WITH (
    SYSTEM_VERSIONING = ON (
        HISTORY_TABLE = dbo.StocksHistory
    )
);
GO

-- Create an index on Symbol for faster lookups
CREATE UNIQUE INDEX IX_Stocks_Symbol ON Stocks(Symbol);
GO

-- Create a view for the latest stock changes
CREATE VIEW vw_StockChanges
AS
SELECT 
    s.Symbol,
    s.BidPrice,
    s.BidVolume,
    s.AskPrice,
    s.AskVolume,
    s.LastUpdated
FROM Stocks s;
GO

-- Create trigger for real-time notifications
CREATE TRIGGER tr_Stocks_Changes
ON Stocks
AFTER UPDATE, INSERT
AS
BEGIN
    SET NOCOUNT ON;
    INSERT INTO StockChangeLog (
        Symbol,
        BidPrice,
        BidVolume,
        AskPrice,
        AskVolume,
        ChangeType
    )
    SELECT 
        i.Symbol,
        i.BidPrice,
        i.BidVolume,
        i.AskPrice,
        i.AskVolume,
        CASE 
            WHEN d.Symbol IS NULL THEN 'INSERT'
            ELSE 'UPDATE'
        END
    FROM inserted i
    LEFT JOIN deleted d ON i.StockId = d.StockId
END;
GO

-- Create log table for changes
CREATE TABLE StockChangeLog (
    LogId INT IDENTITY(1,1) PRIMARY KEY,
    Symbol VARCHAR(10) NOT NULL,
    BidPrice DECIMAL(18,4) NOT NULL,
    BidVolume INT NOT NULL,
    AskPrice DECIMAL(18,4) NOT NULL,
    AskVolume INT NOT NULL,
    ChangeType VARCHAR(10) NOT NULL,
    LoggedAt DATETIME2(7) DEFAULT GETUTCDATE()
);
GO

-- Create enum table for transaction types
CREATE TABLE TransactionTypes (
    TypeId INT PRIMARY KEY,
    TypeName VARCHAR(20) NOT NULL UNIQUE
);

-- Insert transaction types
INSERT INTO TransactionTypes (TypeId, TypeName) 
VALUES 
    (1, 'BUY'),
    (2, 'SELL');
GO

-- Create enum table for transaction status
CREATE TABLE TransactionStatus (
    StatusId INT PRIMARY KEY,
    StatusName VARCHAR(20) NOT NULL UNIQUE
);

-- Insert status types
INSERT INTO TransactionStatus (StatusId, StatusName)
VALUES 
    (1, 'PENDING'),
    (2, 'COMPLETED'),
    (3, 'CANCELLED'),
    (4, 'FAILED');
GO

-- Create the main transactions table
CREATE TABLE Transactions (
    TransactionId BIGINT IDENTITY(1,1) PRIMARY KEY,
    Symbol VARCHAR(10) NOT NULL,
    TypeId INT NOT NULL,
    StatusId INT NOT NULL,
    Quantity INT NOT NULL,
    Price DECIMAL(18,4) NOT NULL,
    TotalAmount DECIMAL(18,4) NOT NULL,
    OrderTime DATETIME2(7) DEFAULT GETUTCDATE(),
    ExecutionTime DATETIME2(7),
    Notes NVARCHAR(500),
    CONSTRAINT FK_Transactions_Stock FOREIGN KEY (Symbol) REFERENCES Stocks(Symbol),
    CONSTRAINT FK_Transactions_Type FOREIGN KEY (TypeId) REFERENCES TransactionTypes(TypeId),
    CONSTRAINT FK_Transactions_Status FOREIGN KEY (StatusId) REFERENCES TransactionStatus(StatusId),
    CONSTRAINT CHK_Transactions_Quantity CHECK (Quantity > 0),
    CONSTRAINT CHK_Transactions_Price CHECK (Price > 0)
);
GO

-- Create indexes for better performance
CREATE INDEX IX_Transactions_Symbol ON Transactions(Symbol);
CREATE INDEX IX_Transactions_OrderTime ON Transactions(OrderTime);
CREATE INDEX IX_Transactions_Status ON Transactions(StatusId);
GO

-- Create view for transaction details
CREATE VIEW vw_TransactionDetails
AS
SELECT 
    t.TransactionId,
    t.Symbol,
    tt.TypeName as TransactionType,
    ts.StatusName as Status,
    t.Quantity,
    t.Price,
    t.TotalAmount,
    t.OrderTime,
    t.ExecutionTime,
    t.Notes
FROM Transactions t
JOIN TransactionTypes tt ON t.TypeId = tt.TypeId
JOIN TransactionStatus ts ON t.StatusId = ts.StatusId;
GO

-- Create trigger to update ExecutionTime when status changes to COMPLETED
CREATE TRIGGER tr_Transactions_Completion
ON Transactions
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;
    
    IF EXISTS (
        SELECT 1 
        FROM inserted i 
        JOIN deleted d ON i.TransactionId = d.TransactionId
        WHERE i.StatusId = 2 -- COMPLETED
        AND d.StatusId != 2
    )
    BEGIN
        UPDATE Transactions
        SET ExecutionTime = GETUTCDATE()
        FROM Transactions t
        JOIN inserted i ON t.TransactionId = i.TransactionId
        WHERE i.StatusId = 2;
    END
END;
GO