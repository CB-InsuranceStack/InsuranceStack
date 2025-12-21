import { useQuery } from '@tanstack/react-query';
import { CreditCard, Download, RefreshCw } from 'lucide-react';
import { api } from '../services/api';
import TransactionList from '../components/TransactionList';
import AlertBanner from '../components/AlertBanner';
import useRoxFlag from '../hooks/useRoxFlag';
import type { Transaction } from '../types';

export default function Transactions() {
  const transactionsFilters = useRoxFlag('transactionsFilters');

  // Fetch accounts to get user's currency
  const { data: accounts } = useQuery({
    queryKey: ['accounts'],
    queryFn: () => api.getAccounts(),
  });

  // Fetch transactions data
  const {
    data: transactions,
    isLoading,
    isError,
    error,
    refetch,
    isFetching,
  } = useQuery<Transaction[]>({
    queryKey: ['transactions'],
    queryFn: () => api.getTransactions(),
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  // Get currency from user's accounts (all accounts have the same currency)
  const userCurrency = accounts?.[0]?.currency || 'USD';

  // Calculate statistics
  const stats = transactions?.reduce(
    (acc, transaction) => {
      if (transaction.status === 'completed') {
        if (transaction.type === 'credit') {
          acc.totalCredits += transaction.amount;
        } else {
          acc.totalDebits += Math.abs(transaction.amount);
        }
      }
      acc.totalTransactions += 1;
      return acc;
    },
    { totalCredits: 0, totalDebits: 0, totalTransactions: 0 }
  );

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: userCurrency,
    }).format(amount);
  };

  const handleExport = () => {
    // Export transactions to CSV
    if (!transactions) return;

    const csv = [
      ['Date', 'Description', 'Merchant', 'Category', 'Type', 'Amount', 'Status'].join(','),
      ...transactions.map((t) =>
        [
          t.date,
          `"${t.description}"`,
          `"${t.merchant || ''}"`,
          t.category,
          t.type,
          t.amount,
          t.status,
        ].join(',')
      ),
    ].join('\n');

    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `transactions-${new Date().toISOString().split('T')[0]}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading transactions...</p>
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
          title="Error Loading Transactions"
          message={
            error instanceof Error ? error.message : 'Failed to load transactions. Please try again.'
          }
          dismissible={false}
        />
        <div className="card p-8 text-center">
          <p className="text-gray-600 mb-4">Unable to load transactions</p>
          <button onClick={() => refetch()} className="btn-primary">
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Transactions</h1>
          <p className="text-gray-600 mt-1">
            View and manage all your transactions across accounts
          </p>
        </div>
        <div className="flex items-center space-x-3">
          <button
            onClick={() => refetch()}
            disabled={isFetching}
            className="btn-secondary"
            title="Refresh transactions"
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${isFetching ? 'animate-spin' : ''}`} />
            Refresh
          </button>
          <button onClick={handleExport} className="btn-primary" title="Export to CSV">
            <Download className="w-4 h-4 mr-2" />
            Export
          </button>
        </div>
      </div>

      {/* Feature Flag Alert */}
      {transactionsFilters && (
        <AlertBanner
          type="info"
          title="Advanced Filters Enabled"
          message="Use the filters below to search and organize your transactions by type, category, and status."
        />
      )}

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Transactions</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">{stats?.totalTransactions || 0}</p>
              <p className="text-xs text-gray-500 mt-2">All time</p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <CreditCard className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Income</p>
              <p className="text-2xl font-bold text-green-600 mt-2">
                {formatCurrency(stats?.totalCredits || 0)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Completed credits</p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <CreditCard className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Expenses</p>
              <p className="text-2xl font-bold text-red-600 mt-2">
                {formatCurrency(stats?.totalDebits || 0)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Completed debits</p>
            </div>
            <div className="bg-red-100 rounded-full p-3">
              <CreditCard className="w-6 h-6 text-red-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Transaction List */}
      {transactions && transactions.length > 0 ? (
        <TransactionList transactions={transactions} currency={userCurrency} />
      ) : (
        <div className="card p-12 text-center">
          <div className="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
            <CreditCard className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No Transactions</h3>
          <p className="text-gray-600">Your transactions will appear here once you start using your accounts</p>
        </div>
      )}
    </div>
  );
}
