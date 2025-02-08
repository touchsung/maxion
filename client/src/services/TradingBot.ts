import axios from 'axios';
import { Stock, Transaction, TransactionRequest, TransactionStatus } from '../types/TradingBot';

const api = axios.create();

api.interceptors.request.use((config) => {
  config.baseURL = import.meta.env.VITE_API_URL || "http://localhost:3000";
  return config;
});

export const getAllStocks = async (): Promise<Stock[]> => {
    const response = await api.get<Stock[]>('/stocks');
    return response.data;
};

export const getAllTransactions = async (): Promise<Transaction[]> => {
    const response = await api.get<Transaction[]>('/transactions');
    return response.data;
};

export const createTransaction = async (
    transaction: TransactionRequest
): Promise<Transaction> => {
    const response = await api.post<Transaction>('/transactions', transaction);
    return response.data;
};

export const updateTransactionStatus = async (
    transactionId: number, 
    status: TransactionStatus
): Promise<Transaction> => {
    const response = await api.put<Transaction>(
        `/transactions/${transactionId}/status`, 
        { status }
    );
    return response.data;
};
