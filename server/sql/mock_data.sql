USE TradingBot;
GO

-- Insert initial stock data
INSERT INTO Stocks (Symbol, BidPrice, BidVolume, AskPrice, AskVolume)
VALUES 
    ('AAPL', 169.85, 500, 170.15, 300),    -- Apple
    ('MSFT', 378.92, 250, 379.45, 400),    -- Microsoft
    ('GOOGL', 147.75, 150, 148.05, 200),   -- Alphabet
    ('AMZN', 178.25, 300, 178.65, 250),    -- Amazon
    ('NVDA', 875.35, 100, 876.20, 150),    -- NVIDIA
    ('META', 505.75, 200, 506.25, 180),    -- Meta
    ('TSLA', 175.85, 400, 176.25, 350),    -- Tesla
    ('JPM', 182.45, 250, 182.85, 200),     -- JPMorgan Chase
    ('V', 275.65, 150, 276.15, 175),       -- Visa
    ('WMT', 59.85, 300, 60.15, 280);       -- Walmart
GO

-- -- Insert transaction types if not already present
-- INSERT INTO TransactionTypes (TypeId, TypeName)
-- VALUES 
--     (1, 'BUY'),
--     (2, 'SELL');
-- GO

-- -- Insert transaction statuses if not already present
-- INSERT INTO TransactionStatus (StatusId, StatusName)
-- VALUES 
--     (1, 'PENDING'),
--     (2, 'COMPLETED'),
--     (3, 'CANCELLED'),
--     (4, 'FAILED');
-- GO

-- -- Insert mock transactions
-- INSERT INTO Transactions (
--     Symbol,
--     TypeId,
--     StatusId,
--     Quantity,
--     Price,
--     TotalAmount,
--     OrderTime,
--     ExecutionTime,
--     Notes
-- )
-- VALUES 
--     -- Completed Buy of AAPL
--     ('AAPL', 1, 2, 100, 170.15, 17015.00, 
--     DATEADD(MINUTE, -30, GETUTCDATE()), 
--     DATEADD(MINUTE, -29, GETUTCDATE()),
--     'Routine investment in Apple'),

--     -- Completed Sell of MSFT
--     ('MSFT', 2, 2, 50, 378.92, 18946.00,
--     DATEADD(MINUTE, -25, GETUTCDATE()),
--     DATEADD(MINUTE, -24, GETUTCDATE()),
--     'Taking profits on Microsoft position'),

--     -- Pending Buy of NVDA
--     ('NVDA', 1, 1, 25, 876.20, 21905.00,
--     DATEADD(MINUTE, -10, GETUTCDATE()),
--     NULL,
--     'Building position in NVIDIA'),

--     -- Failed Buy of META
--     ('META', 1, 4, 75, 506.25, 37968.75,
--     DATEADD(MINUTE, -15, GETUTCDATE()),
--     NULL,
--     'Failed due to insufficient funds'),

--     -- Cancelled Sell of GOOGL
--     ('GOOGL', 2, 3, 30, 147.75, 4432.50,
--     DATEADD(MINUTE, -20, GETUTCDATE()),
--     NULL,
--     'Cancelled due to price movement');
-- GO