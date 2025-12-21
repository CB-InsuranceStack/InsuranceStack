import { AlertCircle, TrendingUp, TrendingDown, Lightbulb, X, CheckCircle } from 'lucide-react';
import { format } from 'date-fns';
import type { Insight } from '../types';
import useRoxFlag from '../hooks/useRoxFlag';

interface InsightsPanelProps {
  insights: Insight[];
  onDismiss?: (insightId: string) => void;
  onAction?: (insightId: string) => void;
}

// V1 - Simple list design
function InsightsPanelV1({ insights, onDismiss, onAction }: InsightsPanelProps) {
  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'border-red-500 bg-red-50';
      case 'warning':
        return 'border-yellow-500 bg-yellow-50';
      case 'info':
        return 'border-blue-500 bg-blue-50';
      default:
        return 'border-gray-500 bg-gray-50';
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case 'critical':
        return <AlertCircle className="w-5 h-5 text-red-600" />;
      case 'warning':
        return <AlertCircle className="w-5 h-5 text-yellow-600" />;
      case 'info':
        return <Lightbulb className="w-5 h-5 text-blue-600" />;
      default:
        return <Lightbulb className="w-5 h-5 text-gray-600" />;
    }
  };

  const activeInsights = insights.filter(i => !i.dismissed);

  if (activeInsights.length === 0) {
    return (
      <div className="card p-8 text-center">
        <CheckCircle className="w-12 h-12 text-green-500 mx-auto mb-3" />
        <h3 className="text-lg font-semibold text-gray-900">All Caught Up!</h3>
        <p className="text-gray-500 mt-1">No new insights at the moment.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {activeInsights.map((insight) => (
        <div
          key={insight.id}
          className={`border-l-4 rounded-lg p-4 ${getSeverityColor(insight.severity)}`}
        >
          <div className="flex items-start justify-between">
            <div className="flex items-start space-x-3 flex-1">
              <div className="flex-shrink-0 mt-0.5">
                {getSeverityIcon(insight.severity)}
              </div>
              <div className="flex-1 min-w-0">
                <h4 className="text-sm font-semibold text-gray-900">{insight.title}</h4>
                <p className="text-sm text-gray-600 mt-1">{insight.description}</p>
                {insight.amount && (
                  <p className="text-sm font-medium text-gray-900 mt-2">
                    Amount: ${insight.amount.toFixed(2)}
                  </p>
                )}
                <p className="text-xs text-gray-500 mt-2">
                  {format(new Date(insight.createdAt), 'MMM d, yyyy h:mm a')}
                </p>
                {insight.actionable && insight.actionText && (
                  <button
                    onClick={() => onAction?.(insight.id)}
                    className="mt-3 btn-primary text-xs"
                  >
                    {insight.actionText}
                  </button>
                )}
              </div>
            </div>
            {onDismiss && (
              <button
                onClick={() => onDismiss(insight.id)}
                className="text-gray-400 hover:text-gray-600 flex-shrink-0 ml-2"
                aria-label="Dismiss insight"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}

// V2 - Enhanced card-based design with better visuals
function InsightsPanelV2({ insights, onDismiss, onAction }: InsightsPanelProps) {
  const getSeverityStyles = (severity: string) => {
    switch (severity) {
      case 'critical':
        return {
          bg: 'bg-gradient-to-br from-red-500 to-red-600',
          icon: AlertCircle,
          badge: 'bg-red-100 text-red-800',
        };
      case 'warning':
        return {
          bg: 'bg-gradient-to-br from-yellow-500 to-yellow-600',
          icon: AlertCircle,
          badge: 'bg-yellow-100 text-yellow-800',
        };
      case 'info':
        return {
          bg: 'bg-gradient-to-br from-blue-500 to-blue-600',
          icon: Lightbulb,
          badge: 'bg-blue-100 text-blue-800',
        };
      default:
        return {
          bg: 'bg-gradient-to-br from-gray-500 to-gray-600',
          icon: Lightbulb,
          badge: 'bg-gray-100 text-gray-800',
        };
    }
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'spending':
        return TrendingDown;
      case 'saving':
        return TrendingUp;
      default:
        return Lightbulb;
    }
  };

  const activeInsights = insights.filter(i => !i.dismissed);

  if (activeInsights.length === 0) {
    return (
      <div className="card p-12 text-center">
        <div className="bg-green-100 rounded-full w-20 h-20 flex items-center justify-center mx-auto mb-4">
          <CheckCircle className="w-10 h-10 text-green-600" />
        </div>
        <h3 className="text-xl font-bold text-gray-900">All Caught Up!</h3>
        <p className="text-gray-500 mt-2">No new insights at the moment. Keep up the great work!</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {activeInsights.map((insight) => {
        const styles = getSeverityStyles(insight.severity);
        const Icon = styles.icon;
        const TypeIcon = getTypeIcon(insight.type);

        return (
          <div key={insight.id} className="card overflow-hidden hover:shadow-lg transition-shadow">
            {/* Header with gradient */}
            <div className={`${styles.bg} p-4 text-white relative`}>
              <div className="flex items-start justify-between">
                <div className="flex items-center space-x-3">
                  <div className="bg-white/20 backdrop-blur-sm rounded-lg p-2">
                    <Icon className="w-5 h-5" />
                  </div>
                  <div>
                    <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${styles.badge}`}>
                      {insight.severity.toUpperCase()}
                    </span>
                  </div>
                </div>
                {onDismiss && (
                  <button
                    onClick={() => onDismiss(insight.id)}
                    className="text-white hover:bg-white/20 p-1 rounded-md transition-colors"
                    aria-label="Dismiss insight"
                  >
                    <X className="w-5 h-5" />
                  </button>
                )}
              </div>
            </div>

            {/* Content */}
            <div className="p-4 space-y-3">
              <div>
                <div className="flex items-center space-x-2 mb-2">
                  <TypeIcon className="w-4 h-4 text-gray-400" />
                  <span className="text-xs text-gray-500 uppercase tracking-wider">{insight.type}</span>
                </div>
                <h4 className="text-base font-semibold text-gray-900">{insight.title}</h4>
                <p className="text-sm text-gray-600 mt-2 leading-relaxed">{insight.description}</p>
              </div>

              {insight.amount && (
                <div className="bg-gray-50 rounded-lg p-3">
                  <p className="text-xs text-gray-500">Amount</p>
                  <p className="text-lg font-bold text-gray-900 mt-0.5">
                    ${insight.amount.toFixed(2)}
                  </p>
                </div>
              )}

              {insight.category && (
                <div className="flex items-center text-sm text-gray-600">
                  <span className="font-medium">Category:</span>
                  <span className="ml-2 capitalize">{insight.category}</span>
                </div>
              )}

              <div className="flex items-center justify-between pt-2 border-t border-gray-100">
                <span className="text-xs text-gray-500">
                  {format(new Date(insight.createdAt), 'MMM d, yyyy')}
                </span>
                {insight.actionable && insight.actionText && (
                  <button
                    onClick={() => onAction?.(insight.id)}
                    className="btn-primary text-xs py-1.5 px-3"
                  >
                    {insight.actionText}
                  </button>
                )}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

// Main component that switches between versions based on feature flag
export default function InsightsPanel({ insights, onDismiss, onAction }: InsightsPanelProps) {
  const insightsV2 = useRoxFlag('insightsV2');

  return insightsV2 ? (
    <InsightsPanelV2 insights={insights} onDismiss={onDismiss} onAction={onAction} />
  ) : (
    <InsightsPanelV1 insights={insights} onDismiss={onDismiss} onAction={onAction} />
  );
}
