# Feature Flags Guide - CloudBees FM Integration

This guide explains how feature flags are implemented in AccountStack and how to use them effectively.

## Overview

AccountStack uses **CloudBees Feature Management (powered by Rox SDK)** to provide dynamic feature control without redeployment. This allows you to:

- Roll out features gradually
- A/B test different UI variations
- Instantly disable problematic features
- Target specific user segments
- Reduce deployment risk

## Implemented Feature Flags

### 1. Dashboard Cards V2 (`ui.dashboardCardsV2`)

**Default:** `true` (enabled)

**What it does:**
- Toggles between two versions of account cards on the dashboard
- V1: Simple card design with basic information
- V2: Enhanced design with gradients, icons, animations, and better visual hierarchy

**Visual Differences:**

**V1 (Simple):**
```
┌──────────────────────────┐
│ Checking Account   ●Active│
│ ••••1234                 │
│                          │
│ Balance                  │
│ $5,234.56               │
│ ─────────────────────── │
│ Checking • View Details  │
└──────────────────────────┘
```

**V2 (Enhanced):**
```
┌──────────────────────────┐
│ ╔════════════════════╗  │
│ ║ [icon] Checking  ⋮  ║  │ ← Gradient header
│ ║ ••••1234            ║  │
│ ║                     ║  │
│ ║ Available Balance   ║  │
│ ║ $5,234.56          ║  │
│ ╚════════════════════╝  │
│ ● Active    CHECKING    │
└──────────────────────────┘
```

**Files affected:**
- `/Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/web/src/components/AccountCard.tsx`

**Code example:**
```typescript
import { useFeatureFlags } from '@/features/flags';

function MyComponent() {
  const { dashboardCardsV2 } = useFeatureFlags();

  return dashboardCardsV2 ? <AccountCardV2 /> : <AccountCardV1 />;
}
```

---

### 2. Insights V2 (`ui.insightsV2`)

**Default:** `false` (disabled)

**What it does:**
- Toggles between list-based and card-based insights layouts
- V1: Vertical list with color-coded borders
- V2: Grid-based card layout with enhanced visuals and better organization

**Visual Differences:**

**V1 (List):**
```
┌───────────────────────────────────┐
│ ⚠ High Spending Alert            │
│   You spent 20% more than usual   │
│   [Take Action]                   │
└───────────────────────────────────┘
┌───────────────────────────────────┐
│ 💡 Saving Opportunity             │
│   You could save $150 this month  │
│   [View Details]                  │
└───────────────────────────────────┘
```

**V2 (Cards):**
```
┌─────────────────┐  ┌─────────────────┐
│ ╔═════════════╗ │  │ ╔═════════════╗ │
│ ║ ⚠ WARNING   ║ │  │ ║ 💡 INFO     ║ │
│ ╚═════════════╝ │  │ ╚═════════════╝ │
│                 │  │                 │
│ High Spending   │  │ Saving Opp.     │
│ Alert           │  │                 │
│                 │  │ You could save  │
│ You spent 20%...│  │ $150 this month │
│                 │  │                 │
│ [Take Action]   │  │ [View Details]  │
└─────────────────┘  └─────────────────┘
```

**Files affected:**
- `/Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/web/src/components/InsightsPanel.tsx`

**Code example:**
```typescript
const { insightsV2 } = useFeatureFlags();

return insightsV2 ?
  <InsightsPanelV2 insights={insights} /> :
  <InsightsPanelV1 insights={insights} />;
```

---

### 3. Alerts Banner (`ui.alertsBanner`)

**Default:** `true` (enabled)

**What it does:**
- Shows/hides informational banners at the top of pages
- Used for announcements, warnings, and important messages
- Dismissible by users

**Visual Example:**
```
┌─────────────────────────────────────────────────┐
│ ℹ️ New Feature Available                      ✕ │
│   Check out the enhanced dashboard cards!       │
└─────────────────────────────────────────────────┘

[Page content below...]
```

**When disabled:** Banners don't render at all, saving vertical space.

**Files affected:**
- `/Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/web/src/components/AlertBanner.tsx`

**Code example:**
```typescript
<AlertBanner
  type="info"
  title="New Feature Available"
  message="Check out the enhanced dashboard cards!"
/>
// Banner only renders if alertsBanner flag is true
```

---

### 4. Transactions Filters (`ui.transactionsFilters`)

**Default:** `true` (enabled)

**What it does:**
- Enables advanced filtering interface for transactions
- Adds search, type filter, category filter, and status filter
- Shows result count

**Visual Differences:**

**Without Filters (Disabled):**
```
Transactions Page
─────────────────────
[Transaction 1]
[Transaction 2]
[Transaction 3]
...
```

**With Filters (Enabled):**
```
Transactions Page
─────────────────────
┌─────────────────────────────────────┐
│ 🔍 Filters                 Clear All│
│ ─────────────────────────────────── │
│ [Search transactions...]            │
│ [Type ▼] [Category ▼] [Status ▼]  │
│ Showing 15 of 42 transactions       │
└─────────────────────────────────────┘

[Filtered Results:]
[Transaction 1]
[Transaction 2]
...
```

**Files affected:**
- `/Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/web/src/components/TransactionList.tsx`

**Code example:**
```typescript
const { transactionsFilters } = useFeatureFlags();

return (
  <div>
    {transactionsFilters && (
      <FilterPanel
        onSearch={handleSearch}
        onFilter={handleFilter}
      />
    )}
    <TransactionItems items={filteredItems} />
  </div>
);
```

---

### 5. Kill Insights (`kill.ui.insights`)

**Default:** `false` (not killed)

**What it does:**
- Emergency kill switch for the entire insights feature
- When enabled, completely disables insights page and API calls
- Shows maintenance message to users

**Visual Example:**

**Normal (Kill Switch Off):**
```
Insights Page
─────────────
[Statistics Cards]
[Insights Panel]
[Action Buttons]
```

**Killed (Kill Switch On):**
```
Insights Page
─────────────
┌───────────────────────────────────┐
│ ⚠️ Feature Temporarily Unavailable│
│                                   │
│ The insights feature is currently │
│ disabled for maintenance. Please  │
│ check back later.                 │
└───────────────────────────────────┘
```

**Files affected:**
- `/Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/web/src/pages/Insights.tsx`

**Code example:**
```typescript
const { killInsights } = useFeatureFlags();

if (killInsights) {
  return <MaintenanceMessage />;
}

// Normal insights rendering
return <InsightsPage />;
```

---

## How to Use Feature Flags

### In React Components (Recommended)

Use the `useFeatureFlags()` hook:

```typescript
import { useFeatureFlags } from '@/features/flags';

function MyComponent() {
  const flags = useFeatureFlags();

  return (
    <div>
      {flags.alertsBanner && <Alert />}
      {flags.transactionsFilters && <Filters />}
    </div>
  );
}
```

### Programmatically

Use helper functions for conditional logic:

```typescript
import {
  isDashboardCardsV2Enabled,
  isInsightsV2Enabled,
  isAlertsBannerEnabled,
  isTransactionsFiltersEnabled,
  isInsightsKilled,
} from '@/features/flags';

// In any TypeScript file
if (isDashboardCardsV2Enabled()) {
  // Use V2 implementation
}
```

### Direct Access

Access flag values directly:

```typescript
import { flags } from '@/features/flags';

const isV2 = flags.dashboardCardsV2.isEnabled();
```

## Configuration

### Local Development (No API Key)

The app works with default values defined in the code:

```typescript
// src/features/flags.ts
export class FeatureFlags {
  public dashboardCardsV2 = new Rox.Flag(true);   // Default: ON
  public insightsV2 = new Rox.Flag(false);        // Default: OFF
  public alertsBanner = new Rox.Flag(true);       // Default: ON
  public transactionsFilters = new Rox.Flag(true);// Default: ON
  public killInsights = new Rox.Flag(false);      // Default: OFF
}
```

To test different values locally, modify these defaults.

### With CloudBees FM (Production)

1. **Get API Key**
   - Sign up at https://app.cloudbees.io/
   - Create a new application
   - Navigate to Settings → API Keys
   - Copy your API key

2. **Configure Environment**
   ```bash
   # .env file
   VITE_ROX_API_KEY=your_api_key_here
   ```

3. **Create Flags in CloudBees Dashboard**
   - Go to Feature Flags section
   - Create flags with these exact names:
     - `accountstack.dashboardCardsV2`
     - `accountstack.insightsV2`
     - `accountstack.alertsBanner`
     - `accountstack.transactionsFilters`
     - `accountstack.killInsights`

4. **Control Flags Remotely**
   - Toggle flags on/off in real-time
   - Create targeting rules
   - Schedule flag changes
   - View analytics

## Testing Flag Combinations

### Recommended Test Scenarios

1. **All Flags ON**
   - Dashboard Cards V2: ✅
   - Insights V2: ✅
   - Alerts Banner: ✅
   - Transactions Filters: ✅
   - Kill Insights: ❌

2. **All Flags OFF**
   - Dashboard Cards V2: ❌
   - Insights V2: ❌
   - Alerts Banner: ❌
   - Transactions Filters: ❌
   - Kill Insights: ❌

3. **Emergency Mode** (Kill Switch Active)
   - Kill Insights: ✅
   - Others: Any state

4. **Gradual Rollout**
   - Start: Dashboard Cards V2 only
   - Next: Add Transactions Filters
   - Then: Enable Insights V2
   - Finally: All features

## Best Practices

### 1. Feature Flag Naming
- Use descriptive names
- Prefix with component/feature area
- Use kebab-case: `ui.dashboard-cards-v2`

### 2. Default Values
- Default to safer option (usually `false` for new features)
- Consider backward compatibility
- Document why each default was chosen

### 3. Kill Switches
- Always have kill switches for critical features
- Monitor after enabling flags
- Have rollback plan

### 4. Cleanup
- Remove flags after feature is stable
- Don't accumulate technical debt
- Archive old flags in CloudBees

### 5. Testing
- Test both flag states
- Test flag transitions
- Test with various user segments

## Monitoring

### Check Flag Status

In browser console:
```javascript
// See all flag values
console.log({
  dashboardCardsV2: window.Rox?.flags?.dashboardCardsV2?.isEnabled(),
  insightsV2: window.Rox?.flags?.insightsV2?.isEnabled(),
  alertsBanner: window.Rox?.flags?.alertsBanner?.isEnabled(),
  transactionsFilters: window.Rox?.flags?.transactionsFilters?.isEnabled(),
  killInsights: window.Rox?.flags?.killInsights?.isEnabled(),
});
```

### CloudBees Dashboard

View in CloudBees:
- Flag usage statistics
- User targeting effectiveness
- Error rates by flag state
- A/B test results

## Troubleshooting

### Flags Not Updating

1. Check API key is correct
2. Verify flag names match exactly
3. Check browser console for Rox errors
4. Clear browser cache
5. Verify network requests to CloudBees

### Default Values Being Used

- API key not set → Uses defaults
- Network error → Falls back to defaults
- Invalid flag name → Uses default
- Timeout → Uses defaults

### Flag Changes Not Visible

- Hard refresh browser (Cmd+Shift+R / Ctrl+Shift+R)
- Check if flag is cached
- Verify flag is published in CloudBees
- Check targeting rules

## Resources

- [CloudBees FM Documentation](https://docs.cloudbees.com/docs/cloudbees-feature-management/)
- [Rox SDK Reference](https://docs.cloudbees.com/docs/cloudbees-feature-management/latest/rox-sdks)
- [Feature Flag Best Practices](https://docs.cloudbees.com/docs/cloudbees-feature-management/latest/best-practices)

## Summary

AccountStack implements 5 feature flags that control:
- ✅ UI variations (Cards V2, Insights V2)
- ✅ Feature visibility (Alerts Banner, Filters)
- ✅ Emergency controls (Kill Insights)

All flags work locally without CloudBees and can be controlled remotely with CloudBees FM for dynamic feature management in production.
