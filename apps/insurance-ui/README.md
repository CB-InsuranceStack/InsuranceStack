# AccountStack Web Application

Modern React-based web application for AccountStack, featuring CloudBees Feature Management integration for dynamic feature control.

## Features

- **Dashboard**: Overview of all accounts with real-time balance information
- **Transactions**: Comprehensive transaction list with advanced filtering
- **Insights**: AI-powered financial insights and recommendations
- **Feature Flags**: Dynamic feature control using CloudBees Feature Management (Rox)
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Real-time Updates**: Automatic data refresh with TanStack Query
- **Type Safety**: Full TypeScript support

## Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router** - Client-side routing
- **TanStack Query** - Data fetching and caching
- **Axios** - HTTP client
- **Tailwind CSS** - Styling
- **CloudBees FM (Rox)** - Feature flag management
- **Vitest** - Unit testing
- **Playwright** - E2E testing

## Feature Flags

The application uses **CloudBees Feature Management** with **fully reactive, real-time updates**. Flag changes appear instantly in the UI without page refresh.

### Available Flags

| Flag | Default | Description |
|------|---------|-------------|
| `dashboardCardsV2` | `true` | Enhanced dashboard card design with gradients and better visuals |
| `insightsV2` | `false` | New insights panel with improved layout and card-based design |
| `alertsBanner` | `true` | Top banner for displaying important alerts and announcements |
| `transactionsFilters` | `true` | Advanced filtering options for transactions (search, type, category, status) |
| `killInsights` | `false` | Kill switch to disable the entire insights feature |

### Reactive Pattern

Flags use a **snapshot + listener pattern** for instant updates:

- **Snapshot**: Current state of all flags (evaluated once)
- **Listeners**: Components subscribe to flag changes
- **Updates**: When FM fetches new config, snapshot rebuilds and notifies all listeners
- **Re-renders**: React components using `useRoxFlag()` automatically re-render

This enables **zero-latency updates** perfect for live demos.

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
# Install dependencies
npm install
```

### Environment Configuration

**For standalone development** (without Docker), create a `.env` file:

```bash
cp .env.example .env
```

Configure the following variables:

```env
# CloudBees Feature Management API Key (optional)
# Get your key from: https://app.cloudbees.io/
VITE_ROX_API_KEY=your_api_key_here
```

**For Docker Compose**, create `.env` in the project root:

```bash
# In /Users/you/AccountStack/.env (project root)
CLOUDBEES_FM_API_KEY=your_api_key_here
```

**Note**:
- Without an FM key, the app works perfectly with default flag values
- The `.env` file is gitignored for security
- In production, FM key is injected via Helm at deployment time

### Development

Start the development server:

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Building for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

### Preview Production Build

```bash
npm run preview
```

## Testing

### Unit Tests

```bash
# Run tests
npm run test

# Run tests in watch mode
npm run test:unit

# Generate coverage report
npm run test:coverage
```

### E2E Tests

```bash
npm run test:e2e
```

## Code Quality

### Linting

```bash
npm run lint
```

### Formatting

```bash
npm run format
```

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Layout.tsx      # Main layout with header, nav, footer
│   ├── AccountCard.tsx # Account card with V1/V2 variants
│   ├── TransactionList.tsx # Transaction list with filtering
│   ├── InsightsPanel.tsx   # Insights panel with V1/V2 variants
│   └── AlertBanner.tsx     # Alert banner component
├── pages/              # Page components
│   ├── Dashboard.tsx   # Dashboard page
│   ├── Transactions.tsx # Transactions page
│   └── Insights.tsx    # Insights page
├── hooks/              # Custom React hooks
│   └── useRoxFlag.ts   # Reactive feature flag hook
├── features/           # Feature-specific code
│   └── flags.ts        # CloudBees FM integration & snapshot pattern
├── services/           # API and external services
│   └── api.ts          # Axios API client
├── styles/             # Global styles
│   └── index.css       # Tailwind CSS imports and custom styles
├── test/               # Test configuration
│   └── setup.ts        # Vitest setup
├── types.ts            # TypeScript type definitions
├── App.tsx             # Main app component with routing
├── main.tsx            # Application entry point
└── vite-env.d.ts       # Vite environment types
```

## API Integration

The application connects to the following API endpoints:

- `GET /api/accounts/me` - Get current user information
- `GET /api/accounts` - List all accounts
- `GET /api/transactions` - List all transactions
- `GET /api/insights` - Get financial insights

### API Proxy Configuration

The Vite dev server is configured to proxy API requests:

- `/api/accounts` → `http://api-accounts:8001`
- `/api/transactions` → `http://api-transactions:8002`
- `/api/insights` → `http://api-insights:8003`

## CloudBees Feature Management Integration

### Setup

1. Sign up for CloudBees Feature Management at https://app.cloudbees.io/
2. Create a new application
3. Get your API key
4. Add the API key to your `.env` file (or project root for Docker)

### Using Feature Flags in Components (Reactive)

**✅ Recommended: Use `useRoxFlag()` hook for reactive updates**

```tsx
import useRoxFlag from '@/hooks/useRoxFlag';

function MyComponent() {
  // Component automatically re-renders when flag changes in FM dashboard
  const dashboardCardsV2 = useRoxFlag('dashboardCardsV2');
  const transactionsFilters = useRoxFlag('transactionsFilters');

  return (
    <div>
      {dashboardCardsV2 ? <EnhancedCard /> : <BasicCard />}
      {transactionsFilters && <Filters />}
    </div>
  );
}
```

**⚠️ Legacy: Static helper functions (not reactive)**

```tsx
import {
  isDashboardCardsV2Enabled,
  isTransactionsFiltersEnabled
} from '@/features/flags';

// These work but DON'T trigger re-renders on flag changes
if (isDashboardCardsV2Enabled()) {
  // Use V2 implementation
}
```

### How Reactive Flags Work

1. **Component mounts** → `useRoxFlag()` reads current snapshot
2. **Component subscribes** → Listens for flag changes
3. **FM config updates** → `configurationFetchedHandler` fires
4. **Snapshot rebuilds** → All flags re-evaluated
5. **Listeners notified** → Components with `useRoxFlag()` re-render
6. **UI updates instantly** → No polling, no page refresh 🎉

### Testing Flag Changes Live

1. Start the application: `npm run dev` or `docker compose up`
2. Open browser: http://localhost:3000
3. Open CloudBees FM dashboard
4. Toggle `transactionsFilters` flag
5. Watch filters appear/disappear in UI instantly

## Styling

The application uses Tailwind CSS with a custom brand color scheme:

- Primary Brand Color: `#0066cc`
- Color palette: `brand-{50-900}`

Custom component classes are available:
- `.card` - Base card styling
- `.btn-primary` - Primary button
- `.btn-secondary` - Secondary button
- `.badge-*` - Status badges
- `.input` - Form inputs

## Contributing

1. Follow the existing code style
2. Run linting and formatting before committing
3. Write tests for new features
4. Update documentation as needed

## License

Proprietary - All rights reserved
