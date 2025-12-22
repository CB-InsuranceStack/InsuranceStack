// CloudBees Feature Management (Rox) integration
import Rox, { type FetcherResults, type RoxSetupOptions } from 'rox-browser';

// Define feature flags with default values
export class FeatureFlags {
  // UI Dashboard Cards V2 - Enhanced card design with better visuals
  public dashboardCardsV2 = new Rox.Flag(true);

  // UI Insights V2 - New insights panel with improved analytics
  public insightsV2 = new Rox.Flag(false);

  // UI Alerts Banner - Top banner for important alerts
  public alertsBanner = new Rox.Flag(true);

  // UI Transactions Filters - Advanced filtering for transactions
  public transactionsFilters = new Rox.Flag(true);

  // Kill switch for insights feature
  public killInsights = new Rox.Flag(false);
}

// Create feature flags instance
export const flags = new FeatureFlags();

// Configuration for CloudBees FM
interface RoxConfig {
  apiKey?: string;
  devModeSecret?: string;
}

// Initialize Rox with the feature flags
export async function initializeFeatureFlags(config: RoxConfig = {}): Promise<void> {
  // Register the feature flags container
  Rox.register('insurancestack', flags);

  // Setup Rox with configuration
  const roxConfig: RoxSetupOptions = {
    // Note: debugLevel removed for security - prevents API key from being logged to console
    configurationFetchedHandler: (fetcherResults: FetcherResults) => {
      console.log('[FeatureFlags] Configuration fetched:', {
        hasChanges: fetcherResults.hasChanges,
        source: fetcherResults.fetcherStatus,
      });
      // Update snapshot when configuration changes (reactive pattern)
      setFlagsSnapshot('fetched');
    },
  };

  try {
    // Try to fetch FM key from runtime config file (deployed via Helm)
    // Falls back to build-time env var for local development
    let apiKey = config.apiKey || import.meta.env.VITE_ROX_API_KEY || '';

    // In production (deployed via Helm), fetch from runtime config
    // Use relative path to respect Vite's base path
    if (!apiKey) {
      try {
        const response = await fetch('config/fm.json');
        if (response.ok) {
          const fmConfig = await response.json();
          apiKey = fmConfig.envKey || '';
          if (apiKey) {
            console.log('[FeatureFlags] Loaded FM key from runtime config');
          }
        }
      } catch (fetchError) {
        console.log('[FeatureFlags] No runtime config found, using defaults');
      }
    }

    if (apiKey) {
      await Rox.setup(apiKey, roxConfig);
      console.log('[FeatureFlags] CloudBees FM initialized successfully');
    } else {
      console.warn(
        '[FeatureFlags] No API key provided, using default flag values. ' +
        'Set VITE_ROX_API_KEY environment variable to connect to CloudBees FM.'
      );
      // In dev mode without API key, we can still use the default values
      await Rox.setup('', roxConfig);
    }

    // Initialize snapshot after setup
    setFlagsSnapshot('initialized');
  } catch (error) {
    console.error('[FeatureFlags] Failed to initialize CloudBees FM:', error);
    // Continue with default values if setup fails
    // Initialize snapshot even if setup fails
    setFlagsSnapshot('error');
  }
}

// Helper functions to check flag values
export function isDashboardCardsV2Enabled(): boolean {
  return flags.dashboardCardsV2.isEnabled();
}

export function isInsightsV2Enabled(): boolean {
  return flags.insightsV2.isEnabled() && !flags.killInsights.isEnabled();
}

export function isAlertsBannerEnabled(): boolean {
  return flags.alertsBanner.isEnabled();
}

export function isTransactionsFiltersEnabled(): boolean {
  return flags.transactionsFilters.isEnabled();
}

export function isInsightsKilled(): boolean {
  return flags.killInsights.isEnabled();
}

// Reactive feature flags pattern (inspired by squid-ui)
// Snapshot of current flag values
let _snapshot: Record<string, boolean> = {};

// Listeners for flag changes
const listeners = new Set<(reason: string, snapshot: Record<string, boolean>) => void>();

// Build snapshot by evaluating all flags once
function buildSnapshot(): Record<string, boolean> {
  return {
    dashboardCardsV2: flags.dashboardCardsV2.isEnabled(),
    insightsV2: flags.insightsV2.isEnabled() && !flags.killInsights.isEnabled(),
    alertsBanner: flags.alertsBanner.isEnabled(),
    transactionsFilters: flags.transactionsFilters.isEnabled(),
    killInsights: flags.killInsights.isEnabled(),
  };
}

// Get current snapshot
export function getFlagsSnapshot(): Record<string, boolean> {
  return _snapshot;
}

// Update snapshot and notify listeners
export function setFlagsSnapshot(reason: string): void {
  _snapshot = buildSnapshot();
  console.log('[FeatureFlags] Snapshot updated:', reason, _snapshot);
  listeners.forEach((listener) => {
    try {
      listener(reason, _snapshot);
    } catch (error) {
      console.error('[FeatureFlags] Listener error:', error);
    }
  });
}

// Subscribe to flag changes
export function subscribeFlags(
  callback: (reason: string, snapshot: Record<string, boolean>) => void
): () => void {
  listeners.add(callback);
  return () => {
    listeners.delete(callback);
  };
}

// Hook for React components to use feature flags
export function useFeatureFlags() {
  return {
    dashboardCardsV2: isDashboardCardsV2Enabled(),
    insightsV2: isInsightsV2Enabled(),
    alertsBanner: isAlertsBannerEnabled(),
    transactionsFilters: isTransactionsFiltersEnabled(),
    killInsights: isInsightsKilled(),
  };
}
