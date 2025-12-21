import { useQuery } from '@tanstack/react-query';
import { Wallet, TrendingUp, TrendingDown, DollarSign } from 'lucide-react';
import { api } from '../services/api';
import AccountCard from '../components/AccountCard';
import AlertBanner from '../components/AlertBanner';
import type { Account } from '../types';

export default function Dashboard() {
  // Fetch accounts data
  const {
    data: accounts,
    isLoading,
    isError,
    error,
  } = useQuery<Account[]>({
    queryKey: ['accounts'],
    queryFn: () => api.getAccounts(),
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  // Calculate summary statistics
  const summary = accounts?.reduce(
    (acc, account) => {
      if (account.status === 'active') {
        acc.totalBalance += account.balance;
        acc.accountCount += 1;

        if (account.accountType === 'checking' || account.accountType === 'savings') {
          acc.liquidAssets += account.balance;
        }
      }
      return acc;
    },
    { totalBalance: 0, accountCount: 0, liquidAssets: 0 }
  );

  // Get currency from user's accounts (all accounts have the same currency)
  const userCurrency = accounts?.[0]?.currency || 'USD';

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: userCurrency,
    }).format(amount);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading your accounts...</p>
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
          title="Error Loading Accounts"
          message={error instanceof Error ? error.message : 'Failed to load account data. Please try again.'}
          dismissible={false}
        />
        <div className="card p-8 text-center">
          <p className="text-gray-600 mb-4">Unable to load your accounts</p>
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
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600 mt-1">Welcome back! Here's your financial overview.</p>
      </div>

      {/* Alert Banner */}
      <AlertBanner
        type="info"
        title="New Feature Available"
        message="Check out the enhanced dashboard cards with improved visuals and insights!"
      />

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Total Balance</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.totalBalance || 0)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Across all accounts</p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <DollarSign className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Active Accounts</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">{summary?.accountCount || 0}</p>
              <p className="text-xs text-gray-500 mt-2">Ready to use</p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <Wallet className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Liquid Assets</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">
                {formatCurrency(summary?.liquidAssets || 0)}
              </p>
              <p className="text-xs text-gray-500 mt-2">Available funds</p>
            </div>
            <div className="bg-blue-100 rounded-full p-3">
              <TrendingUp className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Accounts Section */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900">Your Accounts</h2>
          <button className="btn-primary">
            Add Account
          </button>
        </div>

        {accounts && accounts.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {accounts.map((account) => (
              <AccountCard key={account.id} account={account} />
            ))}
          </div>
        ) : (
          <div className="card p-12 text-center">
            <div className="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <Wallet className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No Accounts Yet</h3>
            <p className="text-gray-600 mb-6">Get started by adding your first account</p>
            <button className="btn-primary">
              Add Your First Account
            </button>
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <button className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors">
            <div className="bg-brand-100 rounded-lg p-2">
              <TrendingUp className="w-5 h-5 text-brand-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">View Transactions</p>
              <p className="text-xs text-gray-500">See recent activity</p>
            </div>
          </button>

          <button className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors">
            <div className="bg-green-100 rounded-lg p-2">
              <DollarSign className="w-5 h-5 text-green-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">Transfer Money</p>
              <p className="text-xs text-gray-500">Between accounts</p>
            </div>
          </button>

          <button className="flex items-center space-x-3 p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors">
            <div className="bg-purple-100 rounded-lg p-2">
              <TrendingDown className="w-5 h-5 text-purple-600" />
            </div>
            <div className="text-left">
              <p className="text-sm font-semibold text-gray-900">Pay Bills</p>
              <p className="text-xs text-gray-500">Quick payments</p>
            </div>
          </button>
        </div>
      </div>
    </div>
  );
}
