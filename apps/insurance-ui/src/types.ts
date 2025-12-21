// Type definitions for InsuranceStack application

export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

export interface Customer {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  dateOfBirth?: string;
  address?: string;
  status: 'active' | 'inactive' | 'suspended';
  createdAt: string;
  updatedAt: string;
}

export interface Policy {
  id: string;
  customerId: string;
  policyNumber: string;
  policyType: 'auto' | 'home' | 'life' | 'health';
  status: 'active' | 'expired' | 'cancelled' | 'pending';
  premium: number;
  coverage: number;
  deductible?: number;
  startDate: string;
  endDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface Claim {
  id: string;
  policyId: string;
  claimNumber: string;
  claimType: string;
  status: 'submitted' | 'under_review' | 'approved' | 'denied' | 'paid';
  amount: number;
  description: string;
  dateOfLoss: string;
  submittedDate: string;
  resolvedDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Payment {
  id: string;
  policyId: string;
  amount: number;
  paymentType: 'premium' | 'claim' | 'refund';
  status: 'pending' | 'completed' | 'failed' | 'refunded';
  paymentMethod?: string;
  paymentDate: string;
  createdAt: string;
}

export interface Quote {
  id: string;
  customerId?: string;
  policyType: 'auto' | 'home' | 'life' | 'health';
  premium: number;
  coverage: number;
  deductible?: number;
  status: 'draft' | 'provided' | 'accepted' | 'expired';
  expiryDate?: string;
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
