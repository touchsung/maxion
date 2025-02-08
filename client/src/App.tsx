import { useQuery } from "@tanstack/react-query";
import {
  getAllStocks,
  getAllTransactions,
  createTransaction,
  updateTransactionStatus,
} from "./services/TradingBot";
import { Stock, Transaction, TransactionType } from "./types/TradingBot";
import "./App.css";
import { useState } from "react";

const formatPrice = (price: number): string => {
  return price.toFixed(2);
};

const formatVolume = (volume: number): string => {
  return volume.toLocaleString();
};

const formatType = (type: number): string => {
  return type === 1 ? "Buy" : "Sell";
};

const STATUS_OPTIONS = [
  { value: 1, label: "Pending" },
  { value: 2, label: "Completed" },
  { value: 3, label: "Cancelled" },
  { value: 4, label: "Failed" },
];

const handleTransaction = async (stock: Stock, type: TransactionType) => {
  try {
    await createTransaction({
      symbol: stock.Symbol,
      type: type,
      quantity: 100,
      notes: "Automated transaction",
    });
    alert("Waiting for transaction to be processed...");
  } catch (error) {
    console.error("Failed to create transaction:", error);
    alert("Failed to create transaction. Please try again.");
  }
};

const REFRESH_INTERVAL = 15; // seconds

function App() {
  const [showModal, setShowModal] = useState(false);
  const [pendingTransaction, setPendingTransaction] = useState<{
    stock: Stock;
    type: TransactionType;
  } | null>(null);
  const [updatingTransactions, setUpdatingTransactions] = useState<Set<number>>(
    new Set()
  );

  const {
    data: stocks,
    isLoading: stocksLoading,
    error: stocksError,
  } = useQuery<Stock[]>({
    queryKey: ["stocks"],
    queryFn: getAllStocks,
    refetchInterval: 1000,
  });

  const {
    data: transactions,
    isLoading: transactionsLoading,
    error: transactionsError,
  } = useQuery<Transaction[]>({
    queryKey: ["transactions"],
    queryFn: async () => {
      const data = await getAllTransactions();
      setUpdatingTransactions((prev) => {
        const next = new Set(prev);
        for (const id of prev) {
          const transaction = data.find((t) => t.TransactionID === id);
          if (transaction) {
            next.delete(id);
          }
        }
        return next;
      });
      return data;
    },
    refetchInterval: 1000 * REFRESH_INTERVAL,
  });

  const handleTransactionClick = (stock: Stock, type: TransactionType) => {
    setPendingTransaction({ stock, type });
    setShowModal(true);
  };

  const handleConfirmTransaction = async () => {
    if (!pendingTransaction) return;

    await handleTransaction(pendingTransaction.stock, pendingTransaction.type);
    setShowModal(false);
    setPendingTransaction(null);
  };

  const handleStatusChange = async (
    transactionId: number,
    newStatus: number
  ) => {
    try {
      if (
        transactions?.find((t) => t.TransactionID === transactionId)?.Status ===
        newStatus
      ) {
        return;
      }

      setUpdatingTransactions((prev) => new Set(prev).add(transactionId));
      await updateTransactionStatus(transactionId, newStatus);
    } catch (error) {
      console.error("Failed to update transaction status:", error);
      alert("Failed to update status. Please try again.");
    }
  };

  if (stocksLoading || transactionsLoading) {
    return (
      <main className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </main>
    );
  }

  if (stocksError || transactionsError) {
    const error = stocksError || transactionsError;
    return (
      <main className="flex items-center justify-center min-h-screen text-red-600">
        Error loading stocks: {(error as Error).message}
      </main>
    );
  }
  return (
    <main className="h-screen bg-gray-100 py-8 px-4 sm:px-6 lg:px-8 overflow-hidden">
      <div className="max-w-7xl mx-auto h-full flex flex-col">
        <div className="flex-1 grid grid-rows-2 gap-8 min-h-0">
          <div className="flex flex-col min-h-0">
            <header className="mb-6 flex-shrink-0">
              <h1 className="text-3xl font-bold text-gray-900">Stock Market</h1>
            </header>

            <section className="bg-white shadow-sm rounded-lg overflow-hidden text-center flex-1 min-h-0">
              <div className="overflow-auto h-full">
                <table className="min-w-full divide-y divide-gray-200 text-center">
                  <thead className="bg-gray-50 text-center">
                    <tr>
                      <th className="px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider text-center">
                        Symbol
                      </th>
                      <th className="px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider text-center">
                        Bid Price
                      </th>
                      <th className="px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider text-center">
                        Bid Volume
                      </th>
                      <th className="px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider text-center">
                        Ask Price
                      </th>
                      <th className="px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider text-center">
                        Ask Volume
                      </th>

                      <th className="px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider text-center">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {stocks?.map((stock) => {
                      return (
                        <tr key={stock.StockID} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                            {stock.Symbol}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {formatPrice(stock.BidPrice)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {formatVolume(stock.BidVolume)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {formatPrice(stock.AskPrice)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {formatVolume(stock.AskVolume)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                            <button
                              onClick={() =>
                                handleTransactionClick(
                                  stock,
                                  TransactionType.Buy
                                )
                              }
                              className="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded"
                            >
                              Buy
                            </button>
                            <button
                              onClick={() =>
                                handleTransactionClick(
                                  stock,
                                  TransactionType.Sell
                                )
                              }
                              className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                            >
                              Sell
                            </button>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            </section>
          </div>

          <div className="flex flex-col min-h-0">
            <header className="mb-6 flex-shrink-0 flex justify-between items-center">
              <h2 className="text-2xl font-bold text-gray-900">
                Recent Transactions
              </h2>
            </header>

            <section className="bg-white shadow-sm rounded-lg overflow-hidden flex-1 min-h-0">
              <div className="overflow-auto h-full">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Symbol
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Type
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Quantity
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Price
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Total Amount
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {transactions?.map((transaction) => (
                      <tr
                        key={transaction.TransactionID}
                        className="hover:bg-gray-50"
                      >
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                          {transaction.Symbol}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {formatType(transaction.Type)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <select
                            value={transaction.Status}
                            onChange={(e) =>
                              handleStatusChange(
                                transaction.TransactionID,
                                Number(e.target.value)
                              )
                            }
                            disabled={updatingTransactions.has(
                              transaction.TransactionID
                            )}
                            className="w-full px-3 py-1.5 text-sm bg-white border border-gray-300 rounded-md shadow-sm 
                              focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                              cursor-pointer hover:bg-gray-50 transition-colors duration-200
                              appearance-none bg-no-repeat
                              disabled:bg-gray-100 disabled:cursor-not-allowed
                              [background-image:url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2212%22%20height%3D%2212%22%20viewBox%3D%220%200%2012%2012%22%3E%3Cpath%20fill%3D%22%23374151%22%20d%3D%22M6%208.825L1.175%204%202.238%202.938l3.762%203.762%203.762-3.762L10.825%204%206%208.825z%22%2F%3E%3C%2Fsvg%3E')]
                              bg-[length:1em_1em] bg-[right_0.5rem_center] pr-8"
                          >
                            {STATUS_OPTIONS.map((option) => (
                              <option key={option.value} value={option.value}>
                                {option.label}
                              </option>
                            ))}
                          </select>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {formatVolume(transaction.Quantity)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {formatPrice(transaction.Price)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {formatPrice(transaction.TotalAmount)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </section>
          </div>
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
          <div className="bg-white p-6 rounded-lg shadow-xl">
            <h3 className="text-lg font-semibold mb-4">Confirm Transaction</h3>
            <p className="mb-4">
              Are you sure you want to{" "}
              {pendingTransaction?.type === TransactionType.Buy
                ? "buy"
                : "sell"}{" "}
              100 shares of {pendingTransaction?.stock.Symbol}?
            </p>
            <div className="flex justify-end space-x-3">
              <button
                onClick={() => setShowModal(false)}
                className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                onClick={handleConfirmTransaction}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                Confirm
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}

export default App;
