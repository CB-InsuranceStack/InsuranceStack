# AccountStack Architecture

This document provides detailed technical architecture, design decisions, and customization guidance for AccountStack.

## Table of Contents

- [System Architecture](#system-architecture)
- [Component Design](#component-design)
- [Feature Management System](#feature-management-system)
- [Testing Strategy](#testing-strategy)
- [Local Development](#local-development)
- [CI/CD Pipeline](#cicd-pipeline)
- [Customization Guide](#customization-guide)
- [Technology Stack](#technology-stack)

---

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        AccountStack                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────┐         ┌──────────────────────────────┐  │
│  │   Web UI    │────────▶│      API Gateway (Nginx)      │  │
│  │  (React)    │         └──────────────────────────────┘  │
│  └─────────────┘                      │                      │
│                                        │                      │
│                    ┌───────────────────┼───────────────┐     │
│                    │                   │               │     │
│                    ▼                   ▼               ▼     │
│         ┌──────────────────┐  ┌──────────────┐  ┌─────────┐│
│         │  Accounts API    │  │Transactions  │  │Insights ││
│         │   (Go - 8001)    │  │API (Go-8002) │  │API(8003)││
│         └──────────────────┘  └──────────────┘  └─────────┘│
│                    │                   │               │     │
│                    └───────────────────┴───────────────┘     │
│                                        │                      │
│                           ┌────────────▼──────────────┐      │
│                           │   Shared Data Layer       │      │
│                           │   (Seeded JSON/SQLite)    │      │
│                           └───────────────────────────┘      │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │     CloudBees Feature Management (Real-time)          │ │
│  │     (Fallback to hardcoded defaults when offline)     │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

**Web UI (apps/web)**:
- React-based single-page application
- Account overview, transactions, insights views
- Feature-flag driven UI components
- Responsive design for desktop and mobile

**Accounts API (apps/api-accounts)**:
- Account listing and details
- Balance information
- Account metadata

**Transactions API (apps/api-transactions)**:
- Transaction history
- Transaction filtering and search
- Transaction categorization

**Insights API (apps/api-insights)**:
- Spending analytics
- Account insights and recommendations
- Alert generation

### Data Flow

1. User accesses Web UI
2. UI initializes CloudBees FM SDK (connects to real-time flag updates)
3. UI conditionally renders components based on current flag values
4. User interactions trigger API calls to respective services
5. APIs check feature flags for behavior modifications (via CloudBees FM SDK)
6. APIs return data from seeded data layer
7. UI displays formatted data to user
8. **Flag changes propagate in real-time** (UI/API updates without reload)

---

## Component Design

### apps/web (React UI)

**Structure**:
```
apps/web/
├── src/
│   ├── components/
│   │   ├── AccountCard/      # Feature: ui.dashboard.cardsV2
│   │   ├── TransactionList/  # Feature: ui.transactions.filters
│   │   ├── InsightsPanel/    # Feature: ui.insights.v2
│   │   └── AlertBanner/      # Feature: ui.alerts.banner
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── Transactions.tsx
│   │   ├── Insights.tsx
│   │   └── Settings.tsx
│   ├── services/
│   │   ├── api.ts
│   │   └── featureFlags.ts   # FM client
│   ├── hooks/
│   │   └── useFeatureFlag.ts
│   └── App.tsx
├── tests/
│   ├── unit/
│   └── e2e/
├── package.json
├── vite.config.ts
└── Dockerfile
```

**Testing**:
- ~70 unit tests (components, hooks, utilities)
- ~20 e2e tests (user journeys, flag variations)

### apps/api-accounts (Go)

**Structure**:
```
apps/api-accounts/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── accounts.go      # GET /accounts, /accounts/{id}
│   │   ├── health.go        # GET /healthz
│   │   └── me.go            # GET /me
│   ├── services/
│   │   └── accounts.go
│   ├── repository/
│   │   └── accounts.go
│   └── features/
│       └── flags.go         # FM client
├── tests/
│   ├── unit/
│   └── integration/
├── go.mod
└── Dockerfile
```

**Testing**:
- ~50 unit tests (handlers, services, repository)
- ~30 integration tests (endpoints, feature flags)

### apps/api-transactions (Go)

Similar structure to api-accounts, focused on:
- Transaction listing with pagination
- Transaction filtering (date range, amount, category)
- Transaction search
- Feature flags: `api.transactions.advancedFilters`

**Testing**:
- ~50 unit tests
- ~25 integration tests

### apps/api-insights (Go)

Similar structure to api-accounts, focused on:
- Spending analytics
- Insight generation
- Alert recommendations
- Feature flags: `api.insights.v2`, `api.alerts.enabled`

**Testing**:
- ~40 unit tests
- ~25 integration tests

---

## Feature Management System

### CloudBees FM Integration

**All feature flags use CloudBees Feature Management** with real-time updates.

**Key Characteristics**:
- **Real-time updates**: Flag changes propagate instantly without application reload
- **CloudBees FM SDK**: Integrated in all components (UI and APIs)
- **Fallback to defaults**: When FM is unavailable (offline demos), hardcoded defaults are used
- **No local FM system**: Flags are defined in CloudBees FM, not local config files

### How It Works

**Flag Resolution**:
1. **CloudBees FM** (primary): Flags fetched from CloudBees Feature Management
2. **Hardcoded defaults** (fallback): Used when FM is unavailable (offline, connection issues)

**Real-time Updates**:
- UI components subscribe to flag changes via CloudBees FM SDK
- API services receive flag updates through SDK event streams
- **No page reload or service restart required**
- Updates propagate in milliseconds

### Flag Declaration Flow

**Flags are defined in application code, not in CloudBees FM dashboard.**

1. Developer declares flag in code with default value
2. Application starts and registers flag with CloudBees FM SDK
3. Flag automatically appears in CloudBees FM dashboard
4. Dashboard is used to override defaults in real-time

See `config/README.md` for a complete list of flags with code examples.

**Example flags** (declared in code):
```yaml
# Application Branding (Fork-Friendly)
app.name: "AccountStack"                    # Default: "AccountStack"
app.brandColor: "#0066cc"                   # Default: "#0066cc"
app.logo: "default"                         # Default: "default"
org.name: "CB-AccountStack"                 # Default: "CB-AccountStack"

# Component Deployment
deploy.ui: true                             # Default: true
deploy.api.accounts: true                   # Default: true
deploy.api.transactions: true               # Default: true
deploy.api.insights: true                   # Default: true

# UI Feature Flags
ui.dashboard.cardsV2: true                  # Default: true
ui.insights.v2: false                       # Default: false
ui.alerts.banner: true                      # Default: true
ui.transactions.filters: true               # Default: true
kill.ui.insights: false                     # Default: false (kill switch)

# API Feature Flags
api.maskAmounts: false                      # Default: false
api.insights.v2: false                      # Default: false
api.alerts.enabled: true                    # Default: true
api.transactions.advancedFilters: false     # Default: false
```

### SDK Integration

**JavaScript/TypeScript** (React UI):
```typescript
// apps/web/src/features/flags.ts
import Rox from 'rox-browser';

// DECLARE flags here with default values
// These will auto-register with CloudBees FM on startup
export const Flags = {
  ui: {
    dashboardCardsV2: new Rox.Flag(true),        // Default: true
    insightsV2: new Rox.Flag(false),             // Default: false
    alertsBanner: new Rox.Flag(true),            // Default: true
    transactionsFilters: new Rox.Flag(true),     // Default: true
    killInsights: new Rox.Flag(false),           // Default: false (kill switch)
  }
};

// Initialize CloudBees FM SDK
export async function initializeFeatureFlags() {
  // Register all flags
  Rox.register('accountstack', Flags);

  // Connect to CloudBees FM (flags auto-appear in dashboard)
  await Rox.setup(import.meta.env.VITE_CLOUDBEES_FM_API_KEY);

  // Flags are now live and update in real-time!
}

// Usage in components
function Dashboard() {
  // Access declared flag directly
  const showCardsV2 = Flags.ui.dashboardCardsV2.isEnabled();

  return showCardsV2 ? <NewCards /> : <OldCards />;
  // Component automatically re-renders when flag changes in FM!
}
```

**Go** (API Services):
```go
// apps/api-accounts/internal/features/flags.go
package features

import (
    "os"
    rox "github.com/rollout/rox-go/v5/core"
)

// DECLARE flags here with default values
// These will auto-register with CloudBees FM on startup
var (
    MaskAmounts       = rox.NewRoxFlag(false)  // Default: false
    InsightsV2Enabled = rox.NewRoxFlag(false)  // Default: false
    AlertsEnabled     = rox.NewRoxFlag(true)   // Default: true
    AdvancedFilters   = rox.NewRoxFlag(false)  // Default: false
)

func Initialize() {
    // Register flags with CloudBees FM SDK
    rox.Register("api", &MaskAmounts)
    rox.Register("api", &InsightsV2Enabled)
    rox.Register("api", &AlertsEnabled)
    rox.Register("api", &AdvancedFilters)

    // Connect to CloudBees FM (flags auto-appear in dashboard)
    options := rox.RoxOptions{}
    <-rox.Setup(os.Getenv("CLOUDBEES_FM_API_KEY"), options)

    // Flags are now live and update in real-time!
}

// Usage in handlers
func (h *InsightsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
    // Access declared flag directly
    if InsightsV2Enabled.IsEnabled() {
        // Use new insights engine
        return h.getInsightsV2(w, r)
    }
    // Use legacy insights
    return h.getInsightsV1(w, r)
    // Behavior changes instantly when flag is toggled in FM dashboard!
}
```

### Demo Resilience (Offline Mode)

**When CloudBees FM is unavailable**:
- SDK automatically uses the defaults declared in code
- Application continues to function normally
- No errors or degraded experience
- Perfect for demos without internet connectivity

**Fallback Behavior**:
The defaults you declare in code are used automatically:

```typescript
// UI: Uses default from flag declaration
new Rox.Flag(true)  // ← This "true" is used when offline

// API: Uses default from flag declaration
rox.NewRoxFlag(false)  // ← This "false" is used when offline
```

**No configuration file needed** - defaults are in the code itself.

### Environment Configuration

**CloudBees FM Connection**:
```bash
# Required for CloudBees FM
CLOUDBEES_FM_API_KEY=your-api-key-here
CLOUDBEES_FM_ENVIRONMENT=dev  # Options: dev, staging, prod

# Optional configuration
CLOUDBEES_FM_TIMEOUT=5000     # Connection timeout (ms)
CLOUDBEES_FM_CACHE_TTL=300    # Cache TTL in seconds
```

**Docker Compose** (local development):
```yaml
services:
  web:
    environment:
      - CLOUDBEES_FM_API_KEY=${CLOUDBEES_FM_API_KEY}
      - CLOUDBEES_FM_ENVIRONMENT=dev

  api-accounts:
    environment:
      - CLOUDBEES_FM_API_KEY=${CLOUDBEES_FM_API_KEY}
      - CLOUDBEES_FM_ENVIRONMENT=dev
```

### Real-time Update Flow

1. **Flag changed in CloudBees FM dashboard**
2. **FM service pushes update** to connected SDK clients (SSE/WebSocket)
3. **SDK updates flag value** in memory
4. **UI components re-render** automatically (React hooks detect change)
5. **API handlers use new value** on next request (no restart needed)
6. **User sees changes immediately** - no page reload required

**Demo Impact**:
- Toggle flag in FM dashboard during demo
- UI updates instantly in front of audience
- Shows power of feature management in real-time
- No deployments, no restarts, no waiting

### Flag Types and Naming

**Boolean Flags** (most common):
```typescript
ui.dashboard.cardsV2: boolean        // Enable/disable features
api.insights.v2: boolean             // API behavior switches
kill.ui.insights: boolean            // Emergency kill switches
```

**String Flags** (configuration):
```typescript
app.name: string                     // Application branding
app.brandColor: string               // Theme colors
```

**Number Flags** (thresholds):
```typescript
api.performance.rateLimitPerMin: number  // Rate limits
api.alerts.threshold: number             // Alert thresholds
```

### Best Practices

1. **Declare flags in code**: All flags with defaults live in code, not config files
2. **Sensible defaults**: Choose defaults that work for most users/environments
3. **Test both states**: Write tests for flag enabled AND disabled
4. **Use kill switches**: Prefix emergency flags with `kill.`
5. **Document in code**: Add comments explaining what each flag controls
6. **Real-time ready**: Design features to handle mid-session flag changes
7. **Single source**: Flag declaration in code → auto-appears in FM dashboard

---

## Testing Strategy

### Test Volume Goals

Total: **300-440 tests** to demonstrate SmartTests effectiveness

| Test Type     | Count   | Purpose                          |
|---------------|---------|----------------------------------|
| Unit          | 200-300 | Component/function level testing |
| Integration   | 80-100  | API and data flow testing        |
| End-to-End    | 30-40   | User journey testing             |

### Test Organization

**Unit Tests**:
```
tests/unit/
├── web/
│   ├── components/      # 50-70 tests
│   ├── hooks/           # 20-30 tests
│   └── utils/           # 30-40 tests
├── api-accounts/        # 50 tests
├── api-transactions/    # 50 tests
└── api-insights/        # 40 tests
```

**Integration Tests**:
```
tests/integration/
├── api-accounts/        # 30 tests
├── api-transactions/    # 25 tests
├── api-insights/        # 25 tests
└── feature-flags/       # 20 tests
```

**E2E Tests**:
```
tests/e2e/
├── critical-flows/      # 15-20 tests
├── feature-variations/  # 10-15 tests
└── cross-component/     # 5-10 tests
```

### Test Characteristics

All tests must be:
- **Deterministic**: Same input = same output
- **Fast**: Unit <100ms, Integration <2s, E2E <10s
- **Isolated**: No shared state between tests
- **Meaningful**: Tests real functionality, not trivial
- **Mapped**: Clear traceability to source files for impact analysis

### SmartTests Demonstration

AccountStack's test suite is designed to showcase:
1. **Test Subsetting**: Run only tests affected by code changes
2. **Impact Analysis**: Show which tests cover which components
3. **Time Savings**: Demonstrate reduced test execution time
4. **Confidence**: Prove comprehensive coverage despite subsetting

**Example Scenarios**:
- Change to `AccountCard.tsx` runs ~15 tests instead of 350
- Change to `api-accounts/handlers` runs ~80 tests instead of 350
- Feature flag change runs all flag-dependent tests (~120)

---

## Local Development

### Prerequisites

- Docker & Docker Compose
- Node.js 18+ (for local UI development)
- Go 1.21+ (for local API development)
- Make

### Running Locally

**Full Stack**:
```bash
docker compose up --build
```

**Individual Components**:
```bash
# UI only
cd apps/web && npm install && npm run dev

# Accounts API only
cd apps/api-accounts && go run cmd/server/main.go

# All APIs
docker compose up api-accounts api-transactions api-insights
```

### Environment Configuration

**Local Overrides** (`.env.local`):
```bash
# Override branding
APP_NAME="MyBank Portal"
APP_BRAND_COLOR="#ff6600"

# Disable specific components
DEPLOY_API_INSIGHTS=false

# Enable specific features
UI_INSIGHTS_V2=true
API_MASK_AMOUNTS=true
```

### Hot Reloading

- **UI**: Vite dev server with HMR
- **APIs**: Air (Go live reload) in development mode
- **Docker Compose**: Volume mounts for source code

### Database/Data

**Seeded Data** (`data/seed/`):
- Accounts: `accounts.json`
- Transactions: `transactions.json`
- Insights: `insights.json`
- Users: `users.json`

Data is loaded on startup, persisted in Docker volumes.

---

## CI/CD Pipeline

### CloudBees Unify Workflows

**On Push** (`.cloudbees/workflows/main.yaml`):
```yaml
- Build all components
- Run full test suite (with SmartTests optimization)
- Security scans (SAST, SCA, container scanning)
- Collect evidence and artifacts
- Deploy to dev environment
```

**On Pull Request**:
```yaml
- Build changed components
- Run affected tests only (SmartTests)
- Code review checks
- Preview deployment
```

**Promotion Gates** (dev → staging → prod):
```yaml
- All tests pass
- Security scans clean
- Manual approval (staging → prod)
- Evidence collection for audit
```

### Workflow Configuration

**Component Detection**:
Workflows automatically detect which components changed:
```yaml
changes:
  - apps/web/**         → Build & deploy UI
  - apps/api-accounts/** → Build & deploy Accounts API
  - config/features.yaml → Redeploy all components
```

**Parallel Builds**:
Independent components build in parallel for speed.

**Selective Deployment**:
Only changed components are deployed.

---

## Customization Guide

### Forking AccountStack

AccountStack is designed to be forked and customized for specific customer scenarios or demonstrations.

#### Step 1: Fork Repository

```bash
# Fork on GitHub to your org
# Clone locally
git clone git@github.com:YourOrg/AccountStack.git
cd AccountStack
```

#### Step 2: Update Branding

Edit `config/features.yaml`:
```yaml
app.name: "YourBank Portal"
app.brandColor: "#your-color"
app.logo: "yourbank"
org.name: "YourOrg"
```

#### Step 3: Auto-Detection

Application automatically detects:
- Org name from git remote
- Repo name from git remote
- Environment from context

No hardcoded references to `CB-AccountStack`.

#### Step 4: Customize Features

Enable/disable features as needed:
```yaml
ui.insights.v2: true           # Enable new insights
api.maskAmounts: true          # Enable amount masking
deploy.api.insights: false     # Disable insights service
```

#### Step 5: Update Documentation

Update `README.md` and `docs/SETUP.md` with customer-specific details.

### Adding New Components

To add a new API service:

1. Create directory: `apps/api-newservice/`
2. Follow existing API structure
3. Add to `docker-compose.yaml`
4. Add deployment flag: `deploy.api.newservice`
5. Add tests: `tests/unit/api-newservice/`
6. Update workflows: `.cloudbees/workflows/`

### Adding New Feature Flags

1. Add to `config/features.yaml`
2. Document in `docs/ARCHITECTURE.md`
3. Implement flag checks in code
4. Add tests for flag variations
5. Update UI if user-facing

---

## Technology Stack

### Frontend

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router** - Client-side routing
- **TanStack Query** - Data fetching and caching
- **Tailwind CSS** - Styling
- **Jest** - Unit testing
- **Playwright** - E2E testing

### Backend

- **Go 1.21+** - API services
- **Gorilla Mux** - HTTP routing
- **Go standard library** - Most functionality
- **SQLite** (or JSON) - Local data storage
- **Go testing** - Unit and integration tests

### Infrastructure

- **Docker** - Containerization
- **Docker Compose** - Local orchestration
- **Kubernetes** - Production deployment (optional)
- **Nginx** - Reverse proxy and static file serving

### CI/CD

- **CloudBees Unify** - Workflow orchestration
- **CloudBees SmartTests** - Test optimization
- **CloudBees Feature Management** - Feature flags (optional)
- **CloudBees Security** - Scanning and evidence

### Monitoring (Future)

- **CloudBees Observability** - Metrics and tracing
- **CloudBees Analytics** - Deployment insights

---

## Design Principles

### 1. Local-First
Everything must work without internet connectivity. Cloud services enhance, but are not required.

### 2. Demo-Ready
Application must be impressive out of the box. Realistic data, polished UI, smooth interactions.

### 3. Fork-Friendly
Zero hardcoded assumptions. Easy to rebrand and customize for any customer.

### 4. Test-Rich
High test volume with clear impact mapping for SmartTests demonstrations.

### 5. Production-Like
Architecture mirrors real-world regulated enterprise applications.

### 6. Simple but Not Trivial
Complex enough to be interesting, simple enough to understand quickly.

---

## Security Considerations

### Authentication

- Mock JWT authentication for demo purposes
- No real user credentials
- Session management examples

### Authorization

- Role-based access control examples
- Feature flag-based access control

### Data Protection

- Feature flag for sensitive data masking (`api.maskAmounts`)
- No real PII or financial data
- HTTPS in production deployments

### Secrets Management

- Environment variables for sensitive config
- No secrets in code or config files
- CloudBees Secrets integration for workflows

---

## Performance

### Local Development
- Full stack starts in <30 seconds
- UI hot reload <200ms
- API reload <2 seconds

### Test Execution
- Full test suite: ~5-7 minutes
- With SmartTests: ~1-3 minutes (typical change)
- Unit tests only: ~30 seconds

### Build Times
- UI production build: ~30 seconds
- API builds (Go): ~20 seconds each
- Docker images: ~2-3 minutes

---

## Future Enhancements

Potential additions (not in initial scope):

- Real database (PostgreSQL)
- GraphQL API option
- Mobile app (React Native)
- Advanced analytics
- Multi-tenancy
- Internationalization (i18n)
- Accessibility (a11y) enhancements
- Performance monitoring
- Distributed tracing

---

## Support and Contribution

### Getting Help

- GitHub Issues for bugs and questions
- Internal CloudBees Slack for SE collaboration
- Documentation updates via pull requests

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Update documentation
6. Submit a pull request

### Code Standards

- **Go**: Follow `gofmt` and `golint`
- **TypeScript**: ESLint + Prettier
- **Tests**: Required for all new features
- **Documentation**: Update for all changes

---

## Changelog

See [CHANGELOG.md](../CHANGELOG.md) for version history and migration guides.
