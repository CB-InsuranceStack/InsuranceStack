import { useQuery } from '@tanstack/react-query';
import { AlertCircle, CheckCircle, Clock, XCircle, DollarSign, Search, Filter } from 'lucide-react';
import { api } from '../services/api';
import AlertBanner from '../components/AlertBanner';
import ClaimList from '../components/ClaimList';
import type { Claim } from '../types';
import { useState } from 'react';

export default function Claims() {
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');

  // Fetch claims data
  const {
    data: claims,
    isLoading,
    isError,
    error,
  } = useQuery<Claim[]>({
    queryKey: ['claims'],
    queryFn: () => api.getClaims(),
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  // Filter claims based on search and status
  const filteredClaims = claims?.filter((claim) => {
    const matchesSearch =
      claim.claimNumber.toLowerCase().includes(searchQuery.toLowerCase()) ||
      claim.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      claim.claimType.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesStatus = !statusFilter || claim.status === statusFilter;

    return matchesSearch && matchesStatus;
  });

  // Calculate summary statistics
  const summary = claims?.reduce(
    (acc, claim) => {
      acc.totalClaims += 1;
      acc.totalAmount += claim.amount;

      if (claim.status === 'submitted' || claim.status === 'under_review') {
        acc.pendingClaims += 1;
      }
      if (claim.status === 'approved' || claim.status === 'paid') {
        acc.approvedAmount += claim.amount;
      }

      return acc;
    },
    { totalClaims: 0, totalAmount: 0, pendingClaims: 0, approvedAmount: 0 }
  );

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading claims...</p>
          </div>
        </div>
      </div>
    );
  }

  // Error state
  if (isError) {
    return (
      <div className="space-y-6">
        <AlertBanner
          type="critical"
          title="Error Loading Claims"
          message={error instanceof Error ? error.message : 'Failed to load claims data. Please try again.'}
          dismissible={false}
        />
        <div className="card p-8 text-center">
          <p className="text-gray-600 mb-4">Unable to load claims</p>
          <button
            onClick={() => window.location.reload()}
            className="btn-primary"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Claims</h1>
          <p className="text-gray-600 mt-1">Track and manage your insurance claims.</p>
        </div>
        <button className="btn-primary">
          File New Claim
        </button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Claims</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary?.totalClaims || 0}
              </p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <AlertCircle className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Pending Review</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {summary?.pendingClaims || 0}
              </p>
            </div>
            <div className="bg-yellow-100 rounded-full p-3">
              <Clock className="w-6 h-6 text-yellow-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Amount</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.totalAmount || 0)}
              </p>
            </div>
            <div className="bg-blue-100 rounded-full p-3">
              <DollarSign className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Approved</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.approvedAmount || 0)}
              </p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <CheckCircle className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="card p-6">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search claims by number, type, or description..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-10 w-full"
            />
          </div>
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2">
              <Filter className="w-5 h-5 text-gray-400" />
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="input min-w-[150px]"
              >
                <option value="">All Status</option>
                <option value="submitted">Submitted</option>
                <option value="under_review">Under Review</option>
                <option value="approved">Approved</option>
                <option value="denied">Denied</option>
                <option value="paid">Paid</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Claims List */}
      <div className="card">
        {filteredClaims && filteredClaims.length > 0 ? (
          <ClaimList claims={filteredClaims} />
        ) : (
          <div className="p-12 text-center">
            <div className="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <AlertCircle className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              {claims?.length === 0 ? 'No Claims Filed' : 'No Claims Found'}
            </h3>
            <p className="text-gray-600 mb-6">
              {claims?.length === 0
                ? 'You have not filed any claims yet'
                : 'Try adjusting your search or filter criteria'}
            </p>
            {claims?.length === 0 && (
              <button className="btn-primary">
                File Your First Claim
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
