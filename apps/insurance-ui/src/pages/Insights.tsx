import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Lightbulb, AlertCircle, TrendingUp, RefreshCw } from 'lucide-react';
import { api } from '../services/api';
import InsightsPanel from '../components/InsightsPanel';
import AlertBanner from '../components/AlertBanner';
import useRoxFlag from '../hooks/useRoxFlag';
import type { Insight } from '../types';

export default function Insights() {
  const queryClient = useQueryClient();
  const insightsV2 = useRoxFlag('insightsV2');
  const killInsights = useRoxFlag('killInsights');

  // Fetch insights data
  const {
    data: insights,
    isLoading,
    isError,
    error,
    refetch,
    isFetching,
  } = useQuery<Insight[]>({
    queryKey: ['insights'],
    queryFn: () => api.getInsights(),
    refetchInterval: 60000, // Refetch every minute
    enabled: !killInsights, // Don't fetch if kill switch is enabled
  });

  // Mutation for dismissing insights
  const dismissMutation = useMutation({
    mutationFn: (insightId: string) => api.dismissInsight(insightId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['insights'] });
    },
  });

  // Mutation for taking action on insights
  const actionMutation = useMutation({
    mutationFn: (insightId: string) => api.takeAction(insightId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['insights'] });
    },
  });

  const handleDismiss = (insightId: string) => {
    dismissMutation.mutate(insightId);
  };

  const handleAction = (insightId: string) => {
    actionMutation.mutate(insightId);
  };

  // Calculate statistics
  const stats = insights?.reduce(
    (acc, insight) => {
      if (!insight.dismissed) {
        acc.active += 1;
        if (insight.severity === 'critical') acc.critical += 1;
        if (insight.actionable) acc.actionable += 1;
      }
      return acc;
    },
    { active: 0, critical: 0, actionable: 0 }
  );

  // If kill switch is enabled, show disabled message
  if (killInsights) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Insights</h1>
          <p className="text-gray-600 mt-1">Smart financial insights and recommendations</p>
        </div>

        <AlertBanner
          type="warning"
          title="Insights Feature Temporarily Unavailable"
          message="The insights feature is currently disabled for maintenance. Please check back later."
          dismissible={false}
        />

        <div className="card p-12 text-center">
          <div className="bg-yellow-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
            <AlertCircle className="w-8 h-8 text-yellow-600" />
          </div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Feature Disabled</h3>
          <p className="text-gray-600">
            The insights feature is temporarily unavailable. We'll be back soon with improved analytics!
          </p>
        </div>
      </div>
    );
  }

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="spinner border-brand-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading insights...</p>
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
          title="Error Loading Insights"
          message={
            error instanceof Error ? error.message : 'Failed to load insights. Please try again.'
          }
          dismissible={false}
        />
        <div className="card p-8 text-center">
          <p className="text-gray-600 mb-4">Unable to load insights</p>
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
          <h1 className="text-3xl font-bold text-gray-900">Insights</h1>
          <p className="text-gray-600 mt-1">
            Smart financial insights and personalized recommendations
          </p>
        </div>
        <button
          onClick={() => refetch()}
          disabled={isFetching}
          className="btn-primary"
          title="Refresh insights"
        >
          <RefreshCw className={`w-4 h-4 mr-2 ${isFetching ? 'animate-spin' : ''}`} />
          Refresh
        </button>
      </div>

      {/* Feature Flag Alert */}
      {insightsV2 && (
        <AlertBanner
          type="info"
          title="Enhanced Insights Available"
          message="You're using the new insights panel with improved visuals and better organization!"
        />
      )}

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Active Insights</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">{stats?.active || 0}</p>
              <p className="text-xs text-gray-500 mt-2">Waiting for your review</p>
            </div>
            <div className="bg-brand-100 rounded-full p-3">
              <Lightbulb className="w-6 h-6 text-brand-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Critical Alerts</p>
              <p className="text-2xl font-bold text-red-600 mt-2">{stats?.critical || 0}</p>
              <p className="text-xs text-gray-500 mt-2">Require attention</p>
            </div>
            <div className="bg-red-100 rounded-full p-3">
              <AlertCircle className="w-6 h-6 text-red-600" />
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 font-medium">Actionable Items</p>
              <p className="text-2xl font-bold text-green-600 mt-2">{stats?.actionable || 0}</p>
              <p className="text-xs text-gray-500 mt-2">Ready to act on</p>
            </div>
            <div className="bg-green-100 rounded-full p-3">
              <TrendingUp className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Insights Panel */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Your Insights</h2>
        {insights && insights.length > 0 ? (
          <InsightsPanel
            insights={insights}
            onDismiss={handleDismiss}
            onAction={handleAction}
          />
        ) : (
          <div className="card p-12 text-center">
            <div className="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
              <Lightbulb className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No Insights Available</h3>
            <p className="text-gray-600">
              Check back soon for personalized insights based on your financial activity
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
