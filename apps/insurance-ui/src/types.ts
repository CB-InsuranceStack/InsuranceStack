// Type definitions for AccountStack application

export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

export interface Account {
  id: string;
  userId: string;
  accountNumber: string;
  accountType: 'checking' | 'savings' | 'credit' | 'investment';
  balance: number;
  currency: string;
  nickname?: string;
  status: 'active' | 'inactive' | 'frozen';
  createdAt: string;
  updatedAt: string;
}

export interface Transaction {
  id: string;
  accountId: string;
  amount: number;
  type: 'debit' | 'credit';
  category: string;
  description: string;
  merchant?: string;
  status: 'pending' | 'completed' | 'failed';
  date: string;
}

export interface Insight {
  id: string;
  userId: string;
  type: 'spending' | 'saving' | 'alert' | 'recommendation';
  severity: 'info' | 'warning' | 'critical';
  title: string;
  description: string;
  category?: string;
  amount?: number;
  actionable: boolean;
  actionText?: string;
  dismissed: boolean;
  createdAt: string;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  timestamp: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  hasMore: boolean;
}
