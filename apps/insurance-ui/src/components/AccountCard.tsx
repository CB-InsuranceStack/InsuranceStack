import { Wallet, TrendingUp, TrendingDown, MoreVertical } from 'lucide-react';
import type { Account } from '../types';
import useRoxFlag from '../hooks/useRoxFlag';

interface AccountCardProps {
  account: Account;
}

// V1 - Simple card design
function AccountCardV1({ account }: AccountCardProps) {
  const formatBalance = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: account.currency,
    }).format(amount);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'inactive':
        return 'bg-gray-100 text-gray-800';
      case 'frozen':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="card card-hover p-6">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="font-semibold text-lg text-gray-900">
            {account.nickname || account.accountType.charAt(0).toUpperCase() + account.accountType.slice(1)}
          </h3>
          <p className="text-sm text-gray-500 mt-1">••••{account.accountNumber.slice(-4)}</p>
        </div>
        <span className={`badge ${getStatusColor(account.status)}`}>
          {account.status}
        </span>
      </div>

      <div className="mt-4">
        <p className="text-sm text-gray-500">Balance</p>
        <p className="text-2xl font-bold text-gray-900 mt-1">{formatBalance(account.balance)}</p>
      </div>

      <div className="mt-4 pt-4 border-t border-gray-100 flex items-center justify-between">
        <span className="text-xs text-gray-500 capitalize">{account.accountType} Account</span>
        <button className="text-brand-500 hover:text-brand-600 text-sm font-medium">
          View Details
        </button>
      </div>
    </div>
  );
}

// V2 - Enhanced card design with better visuals and icons
function AccountCardV2({ account }: AccountCardProps) {
  const formatBalance = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: account.currency,
    }).format(amount);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'text-green-600';
      case 'inactive':
        return 'text-gray-600';
      case 'frozen':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const getAccountIcon = () => {
    switch (account.accountType) {
      case 'checking':
        return <Wallet className="w-6 h-6" />;
      case 'savings':
        return <TrendingUp className="w-6 h-6" />;
      case 'credit':
        return <TrendingDown className="w-6 h-6" />;
      case 'investment':
        return <TrendingUp className="w-6 h-6" />;
      default:
        return <Wallet className="w-6 h-6" />;
    }
  };

  const getGradient = () => {
    switch (account.accountType) {
      case 'checking':
        return 'from-blue-500 to-blue-600';
      case 'savings':
        return 'from-green-500 to-green-600';
      case 'credit':
        return 'from-purple-500 to-purple-600';
      case 'investment':
        return 'from-orange-500 to-orange-600';
      default:
        return 'from-brand-500 to-brand-600';
    }
  };

  return (
    <div className="card card-hover overflow-hidden group">
      {/* Header with gradient */}
      <div className={`bg-gradient-to-br ${getGradient()} p-6 text-white relative`}>
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <div className="bg-white/20 backdrop-blur-sm rounded-lg p-2">
              {getAccountIcon()}
            </div>
            <div>
              <h3 className="font-semibold text-lg">
                {account.nickname || account.accountType.charAt(0).toUpperCase() + account.accountType.slice(1)}
              </h3>
              <p className="text-sm opacity-90 mt-0.5">••••{account.accountNumber.slice(-4)}</p>
            </div>
          </div>
          <button className="text-white hover:bg-white/20 p-1 rounded-md transition-colors">
            <MoreVertical className="w-5 h-5" />
          </button>
        </div>

        <div className="mt-6">
          <p className="text-sm opacity-90">Available Balance</p>
          <p className="text-3xl font-bold mt-1">{formatBalance(account.balance)}</p>
        </div>

        {/* Decorative pattern */}
        <div className="absolute top-0 right-0 opacity-10">
          <svg width="200" height="200" viewBox="0 0 200 200">
            <circle cx="150" cy="50" r="60" fill="currentColor" />
            <circle cx="180" cy="120" r="40" fill="currentColor" />
          </svg>
        </div>
      </div>

      {/* Footer */}
      <div className="p-4 bg-white">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className={`w-2 h-2 rounded-full ${getStatusColor(account.status)} animate-pulse`}></div>
            <span className="text-sm text-gray-600 capitalize">{account.status}</span>
          </div>
          <span className="text-xs text-gray-500 uppercase tracking-wide">{account.accountType}</span>
        </div>
      </div>
    </div>
  );
}

// Main component that switches between versions based on feature flag
export default function AccountCard({ account }: AccountCardProps) {
  const dashboardCardsV2 = useRoxFlag('dashboardCardsV2');

  return dashboardCardsV2 ? <AccountCardV2 account={account} /> : <AccountCardV1 account={account} />;
}
