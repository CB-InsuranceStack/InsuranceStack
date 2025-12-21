// API service using Axios for AccountStack
import axios, { AxiosInstance, AxiosError } from 'axios';
import type {
  User,
  Account,
  Transaction,
  Insight,
} from '../types';

// Create axios instance with base configuration
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth tokens, logging, etc.
apiClient.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response) {
      // Server responded with error status
      console.error('[API Error]', {
        status: error.response.status,
        data: error.response.data,
        url: error.config?.url,
      });
    } else if (error.request) {
      // Request made but no response
      console.error('[Network Error]', error.message);
    } else {
      // Error in request setup
      console.error('[Request Error]', error.message);
    }
    return Promise.reject(error);
  }
);

// API Service class
class ApiService {
  // User endpoints
  async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<User>('accounts/me');
    return response.data;
  }

  // Account endpoints
  async getAccounts(): Promise<Account[]> {
    const response = await apiClient.get<Account[]>('accounts');
    return response.data;
  }

  async getAccount(accountId: string): Promise<Account> {
    const response = await apiClient.get<Account>(`accounts/${accountId}`);
    return response.data;
  }

  async createAccount(accountData: Partial<Account>): Promise<Account> {
    const response = await apiClient.post<Account>('accounts', accountData);
    return response.data;
  }

  async updateAccount(accountId: string, accountData: Partial<Account>): Promise<Account> {
    const response = await apiClient.put<Account>(
      `accounts/${accountId}`,
      accountData
    );
    return response.data;
  }

  async deleteAccount(accountId: string): Promise<void> {
    await apiClient.delete(`accounts/${accountId}`);
  }

  // Transaction endpoints
  async getTransactions(params?: {
    accountId?: string;
    type?: string;
    category?: string;
    startDate?: string;
    endDate?: string;
    page?: number;
    pageSize?: number;
  }): Promise<Transaction[]> {
    const response = await apiClient.get<Transaction[]>('transactions', {
      params,
    });
    return response.data;
  }

  async getTransaction(transactionId: string): Promise<Transaction> {
    const response = await apiClient.get<Transaction>(
      `transactions/${transactionId}`
    );
    return response.data;
  }

  async createTransaction(transactionData: Partial<Transaction>): Promise<Transaction> {
    const response = await apiClient.post<Transaction>(
      'transactions',
      transactionData
    );
    return response.data;
  }

  // Insights endpoints
  async getInsights(params?: {
    type?: string;
    severity?: string;
    dismissed?: boolean;
  }): Promise<Insight[]> {
    const response = await apiClient.get<Insight[]>('insights', {
      params,
    });
    return response.data;
  }

  async getInsight(insightId: string): Promise<Insight> {
    const response = await apiClient.get<Insight>(`insights/${insightId}`);
    return response.data;
  }

  async dismissInsight(insightId: string): Promise<void> {
    await apiClient.patch(`insights/${insightId}/dismiss`);
  }

  async takeAction(insightId: string): Promise<void> {
    await apiClient.post(`insights/${insightId}/action`);
  }
}

// Export singleton instance
export const api = new ApiService();
export default apiClient;
