import { useState } from 'react';
import { ArrowUpRight, ArrowDownLeft, Filter, Search } from 'lucide-react';
import { format } from 'date-fns';
import type { Transaction } from '../types';
import useRoxFlag from '../hooks/useRoxFlag';

interface TransactionListProps {
  transactions: Transaction[];
  currency?: string;
}

export default function TransactionList({ transactions, currency = 'USD' }: TransactionListProps) {
  const transactionsFilters = useRoxFlag('transactionsFilters');
  const [searchTerm, setSearchTerm] = useState('');
  const [typeFilter, setTypeFilter] = useState<string>('all');
  const [categoryFilter, setCategoryFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');

  // Get unique categories and statuses for filters
  const categories = Array.from(new Set(transactions.map(t => t.category)));
  const statuses = Array.from(new Set(transactions.map(t => t.status)));

  // Filter transactions
  const filteredTransactions = transactions.filter(transaction => {
    const matchesSearch =
      transaction.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      transaction.merchant?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      transaction.category.toLowerCase().includes(searchTerm.toLowerCase());

    const matchesType = typeFilter === 'all' || transaction.type === typeFilter;
    const matchesCategory = categoryFilter === 'all' || transaction.category === categoryFilter;
    const matchesStatus = statusFilter === 'all' || transaction.status === statusFilter;

    return matchesSearch && matchesType && matchesCategory && matchesStatus;
  });

  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(amount);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'badge-success';
      case 'pending':
        return 'badge-warning';
      case 'failed':
        return 'badge-error';
      default:
        return 'badge-info';
    }
  };

  return (
    <div className="space-y-4">
      {/* Filters - Only shown if feature flag is enabled */}
      {transactionsFilters && (
        <div className="card p-4 space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-semibold text-gray-900 flex items-center">
              <Filter className="w-4 h-4 mr-2" />
              Filters
            </h3>
            <button
              onClick={() => {
                setSearchTerm('');
                setTypeFilter('all');
                setCategoryFilter('all');
                setStatusFilter('all');
              }}
              className="text-sm text-brand-600 hover:text-brand-700 font-medium"
            >
              Clear All
            </button>
          </div>

          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search transactions..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="input pl-10"
            />
          </div>

          {/* Filter dropdowns */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label className="block text-xs font-medium text-gray-700 mb-1">Type</label>
              <select
                value={typeFilter}
                onChange={(e) => setTypeFilter(e.target.value)}
                className="input"
              >
                <option value="all">All Types</option>
                <option value="debit">Debit</option>
                <option value="credit">Credit</option>
              </select>
            </div>

            <div>
              <label className="block text-xs font-medium text-gray-700 mb-1">Category</label>
              <select
                value={categoryFilter}
                onChange={(e) => setCategoryFilter(e.target.value)}
                className="input"
              >
                <option value="all">All Categories</option>
                {categories.map(cat => (
                  <option key={cat} value={cat}>{cat}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-xs font-medium text-gray-700 mb-1">Status</label>
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="input"
              >
                <option value="all">All Statuses</option>
                {statuses.map(status => (
                  <option key={status} value={status} className="capitalize">{status}</option>
                ))}
              </select>
            </div>
          </div>

          {/* Results count */}
          <div className="text-sm text-gray-600">
            Showing {filteredTransactions.length} of {transactions.length} transactions
          </div>
        </div>
      )}

      {/* Transaction List */}
      <div className="card divide-y divide-gray-200">
        {filteredTransactions.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-gray-500">No transactions found</p>
            {transactionsFilters && (searchTerm || typeFilter !== 'all' || categoryFilter !== 'all' || statusFilter !== 'all') && (
              <button
                onClick={() => {
                  setSearchTerm('');
                  setTypeFilter('all');
                  setCategoryFilter('all');
                  setStatusFilter('all');
                }}
                className="mt-4 text-brand-600 hover:text-brand-700 text-sm font-medium"
              >
                Clear filters
              </button>
            )}
          </div>
        ) : (
          filteredTransactions.map((transaction) => (
            <div
              key={transaction.id}
              className="p-4 hover:bg-gray-50 transition-colors cursor-pointer"
            >
              <div className="flex items-start justify-between">
                <div className="flex items-start space-x-3 flex-1 min-w-0">
                  {/* Icon */}
                  <div
                    className={`flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center ${
                      transaction.type === 'credit'
                        ? 'bg-green-100 text-green-600'
                        : 'bg-red-100 text-red-600'
                    }`}
                  >
                    {transaction.type === 'credit' ? (
                      <ArrowDownLeft className="w-5 h-5" />
                    ) : (
                      <ArrowUpRight className="w-5 h-5" />
                    )}
                  </div>

                  {/* Transaction details */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-4">
                      <div className="min-w-0 flex-1">
                        <p className="text-sm font-medium text-gray-900 truncate">
                          {transaction.description}
                        </p>
                        {transaction.merchant && (
                          <p className="text-sm text-gray-500 truncate mt-0.5">
                            {transaction.merchant}
                          </p>
                        )}
                        <div className="flex items-center gap-2 mt-2">
                          <span className="text-xs text-gray-500">
                            {format(new Date(transaction.date), 'MMM d, yyyy')}
                          </span>
                          <span className="text-xs text-gray-300">•</span>
                          <span className="text-xs text-gray-500 capitalize">{transaction.category}</span>
                          <span className={`badge ${getStatusColor(transaction.status)}`}>
                            {transaction.status}
                          </span>
                        </div>
                      </div>

                      {/* Amount */}
                      <div className="text-right flex-shrink-0">
                        <p
                          className={`text-sm font-semibold ${
                            transaction.type === 'credit'
                              ? 'text-green-600'
                              : 'text-red-600'
                          }`}
                        >
                          {transaction.type === 'credit' ? '+' : '-'}
                          {formatAmount(Math.abs(transaction.amount))}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
