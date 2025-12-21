# Feature Flags Reference

This document lists all feature flags used in AccountStack and shows how they're declared in code.

## How Feature Flags Work

1. **Flags are declared in application code** with default values
2. **Application registers flags** with CloudBees FM SDK on startup
3. **Flags automatically appear** in CloudBees FM dashboard
4. **Dashboard is used to override** defaults in real-time (no reload needed)

## UI Feature Flags

Declared in: `apps/web/src/features/flags.ts`

```typescript
import Rox from 'rox-browser';

export const Flags = {
  ui: {
    // Dashboard Features
    dashboardCardsV2: new Rox.Flag(true),
    // Description: Enable redesigned account cards with enhanced visuals
    // Default: true (new cards shown by default)
    // Used in: Dashboard.tsx

    // Insights Features
    insightsV2: new Rox.Flag(false),
    // Description: Enable new insights engine with ML predictions
    // Default: false (legacy insights by default)
    // Used in: InsightsPage.tsx

    // Alert Features
    alertsBanner: new Rox.Flag(true),
    // Description: Show alert banner at top of dashboard
    // Default: true (alerts visible by default)
    // Used in: DashboardLayout.tsx

    // Transaction Features
    transactionsFilters: new Rox.Flag(true),
    // Description: Enable advanced transaction filtering
    // Default: true (filters enabled by default)
    // Used in: TransactionsPage.tsx

    // Kill Switches
    killInsights: new Rox.Flag(false),
    // Description: Emergency kill switch for insights feature
    // Default: false (insights enabled by default)
    // Used in: InsightsPage.tsx (emergency disable)
  }
};

// Initialize and register with CloudBees FM
export async function initializeFeatureFlags() {
  Rox.register('accountstack', Flags);
  await Rox.setup(import.meta.env.VITE_CLOUDBEES_FM_API_KEY);
}
```

### Usage Example (React)

```typescript
// In a React component
import { Flags } from '@/features/flags';

function Dashboard() {
  const showNewCards = Flags.ui.dashboardCardsV2.isEnabled();

  return (
    <div>
      {showNewCards ? <AccountCardsV2 /> : <AccountCardsV1 />}
    </div>
  );
  // Component automatically re-renders when flag changes in FM!
}
```

## API Feature Flags

Declared in each API service's `internal/features/flags.go`

### Accounts API

```go
// apps/api-accounts/internal/features/flags.go
package features

import (
    "os"
    rox "github.com/rollout/rox-go/v5/core"
)

var (
    // MaskAmounts - Mask sensitive dollar amounts in API responses
    // Default: false (amounts shown in clear)
    // Used in: All account/transaction endpoints
    MaskAmounts = rox.NewRoxFlag(false)
)

func Initialize() {
    rox.Register("api.accounts", &MaskAmounts)
    rox.Setup(os.Getenv("CLOUDBEES_FM_API_KEY"), rox.RoxOptions{})
}
```

### Transactions API

```go
// apps/api-transactions/internal/features/flags.go
package features

import (
    "os"
    rox "github.com/rollout/rox-go/v5/core"
)

var (
    // AdvancedFilters - Enable complex transaction filtering
    // Default: false (basic filtering only)
    // Used in: Transaction list endpoint
    AdvancedFilters = rox.NewRoxFlag(false)
)

func Initialize() {
    rox.Register("api.transactions", &AdvancedFilters)
    rox.Setup(os.Getenv("CLOUDBEES_FM_API_KEY"), rox.RoxOptions{})
}
```

### Insights API

```go
// apps/api-insights/internal/features/flags.go
package features

import (
    "os"
    rox "github.com/rollout/rox-go/v5/core"
)

var (
    // InsightsV2Enabled - Use new insights calculation engine
    // Default: false (legacy algorithm)
    // Used in: Insights calculation logic
    InsightsV2Enabled = rox.NewRoxFlag(false)

    // AlertsEnabled - Enable alert generation
    // Default: true (alerts generated)
    // Used in: Alert generation logic
    AlertsEnabled = rox.NewRoxFlag(true)
)

func Initialize() {
    rox.Register("api.insights", &InsightsV2Enabled)
    rox.Register("api.insights", &AlertsEnabled)
    rox.Setup(os.Getenv("CLOUDBEES_FM_API_KEY"), rox.RoxOptions{})
}
```

### Usage Example (Go)

```go
// In an API handler
func (h *InsightsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
    if features.InsightsV2Enabled.IsEnabled() {
        return h.getInsightsV2(w, r)
    }
    return h.getInsightsV1(w, r)
    // Behavior changes instantly when flag toggled in FM dashboard!
}
```

## Deployment Flags

Used in CloudBees Unify workflows to control component deployment.

Declared in: `.cloudbees/workflows/flags.yaml` (workflow configuration)

```yaml
# These control which components get deployed
deploy:
  ui: true                    # Deploy web UI
  api-accounts: true          # Deploy accounts API
  api-transactions: true      # Deploy transactions API
  api-insights: true          # Deploy insights API
```

## Adding New Flags

To add a new feature flag:

1. **Declare in code** with a sensible default:
   ```typescript
   // UI flag
   myNewFeature: new Rox.Flag(false)
   ```
   ```go
   // API flag
   MyNewFeature = rox.NewRoxFlag(false)
   ```

2. **Register with SDK** (add to existing registration):
   ```typescript
   Rox.register('accountstack', Flags);
   ```
   ```go
   rox.Register("api", &MyNewFeature)
   ```

3. **Use in code**:
   ```typescript
   if (Flags.myNewFeature.isEnabled()) { /* ... */ }
   ```
   ```go
   if MyNewFeature.IsEnabled() { /* ... */ }
   ```

4. **Start app** - flag automatically appears in CloudBees FM dashboard

5. **Document here** - add to this README with description and default

## Flag Naming Conventions

- Use dot notation for namespacing: `ui.dashboard.cardsV2`
- Use camelCase for multi-word names
- Prefix kill switches with `kill.`: `kill.ui.insights`
- Boolean flags: `enabled`, `shown`, `v2`
- Be descriptive: `transactionsFilters` not `filters`

## Environment Variables

Configure CloudBees FM connection:

```bash
# Required
CLOUDBEES_FM_API_KEY=your-api-key-here

# Optional
CLOUDBEES_FM_ENVIRONMENT=dev  # Options: dev, staging, prod
```

## Offline/Demo Mode

When CloudBees FM is unavailable (no internet, wrong API key, etc.):
- Application uses the default values declared in code
- No errors or degraded experience
- Perfect for offline demos

Example:
```typescript
new Rox.Flag(true)  // ← This "true" is used when offline
```

## Testing with Flags

Write tests for both flag states:

```typescript
describe('Dashboard', () => {
  it('shows new cards when flag enabled', () => {
    Flags.ui.dashboardCardsV2.setValue(true);
    // Test with new cards
  });

  it('shows old cards when flag disabled', () => {
    Flags.ui.dashboardCardsV2.setValue(false);
    // Test with old cards
  });
});
```

## Flag List Summary

| Flag | Type | Default | Component | Description |
|------|------|---------|-----------|-------------|
| `ui.dashboardCardsV2` | Boolean | `true` | Web UI | New account card design |
| `ui.insightsV2` | Boolean | `false` | Web UI | New insights engine |
| `ui.alertsBanner` | Boolean | `true` | Web UI | Alert banner visibility |
| `ui.transactionsFilters` | Boolean | `true` | Web UI | Advanced filtering |
| `kill.ui.insights` | Boolean | `false` | Web UI | Emergency kill switch |
| `api.maskAmounts` | Boolean | `false` | Accounts API | Mask dollar amounts |
| `api.insights.v2` | Boolean | `false` | Insights API | New calculation engine |
| `api.alerts.enabled` | Boolean | `true` | Insights API | Alert generation |
| `api.transactions.advancedFilters` | Boolean | `false` | Transactions API | Complex filters |

## Resources

- [CloudBees Feature Management Docs](https://docs.cloudbees.com/docs/cloudbees-feature-management/)
- [Rox SDK - JavaScript](https://github.com/rollout/rox-js)
- [Rox SDK - Go](https://github.com/rollout/rox-go)
- [Architecture Documentation](../docs/ARCHITECTURE.md)
