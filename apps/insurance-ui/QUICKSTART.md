# AccountStack Web - Quick Start Guide

Get the AccountStack web application up and running in minutes!

## Prerequisites

- Node.js 18+ installed
- npm or yarn package manager

## Installation & Setup

### 1. Install Dependencies

```bash
cd apps/web
npm install
```

### 2. Configure Environment

```bash
# Copy the example environment file
cp .env.example .env
```

Edit `.env` and configure your settings:

```env
# API Configuration (defaults to /api, proxied by Vite)
VITE_API_BASE_URL=/api

# CloudBees Feature Management API Key (optional for local dev)
# Get your key from: https://app.cloudbees.io/
VITE_ROX_API_KEY=

# Environment
VITE_ENV=development
```

**Note:** The app will work without a CloudBees API key using default flag values. To connect to CloudBees FM for dynamic flag control, sign up at https://app.cloudbees.io/ and add your API key.

### 3. Start Development Server

```bash
npm run dev
```

The application will be available at: **http://localhost:3000**

## What You'll See

### Dashboard (/)
- Overview of all accounts
- Summary cards showing total balance, active accounts, and liquid assets
- Account cards with V1 or V2 styling (controlled by feature flag)
- Quick action buttons

### Transactions (/transactions)
- Complete list of all transactions
- Statistics showing total transactions, income, and expenses
- Advanced filtering (when feature flag enabled)
- Export to CSV functionality

### Insights (/insights)
- Financial insights and recommendations
- Alert management
- Actionable recommendations
- V1 or V2 panel layout (controlled by feature flag)

## Feature Flags

The app is pre-configured with these feature flags:

| Feature | Default | What it does |
|---------|---------|-------------|
| Dashboard Cards V2 | ✅ ON | Shows enhanced cards with gradients and animations |
| Insights V2 | ❌ OFF | Shows new card-based insights layout |
| Alerts Banner | ✅ ON | Displays info banners at top of pages |
| Transactions Filters | ✅ ON | Enables advanced search and filtering |
| Kill Insights | ❌ OFF | Emergency kill switch for insights feature |

### Testing Feature Flags Locally

Without CloudBees FM connection, you can test flags by editing:
`/Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/web/src/features/flags.ts`

Change the default values:
```typescript
public dashboardCardsV2 = new Rox.Flag(false); // Change to false
public insightsV2 = new Rox.Flag(true);        // Change to true
```

### Using CloudBees FM Dashboard

1. Sign up at https://app.cloudbees.io/
2. Create a new application
3. Create feature flags matching the names above
4. Get your API key and add to `.env`
5. Restart the dev server
6. Control flags from the CloudBees dashboard in real-time!

## API Integration

The app expects three backend services:

- **Accounts API** (port 8001): `/api/accounts`
- **Transactions API** (port 8002): `/api/transactions`
- **Insights API** (port 8003): `/api/insights`

Vite's dev server automatically proxies these requests. In Docker Compose, the services are:
- `api-accounts:8001`
- `api-transactions:8002`
- `api-insights:8003`

## Common Commands

```bash
# Development
npm run dev              # Start dev server

# Building
npm run build            # Build for production
npm run preview          # Preview production build

# Testing
npm run test             # Run tests in watch mode
npm run test:unit        # Run tests once
npm run test:coverage    # Generate coverage report
npm run test:e2e         # Run E2E tests

# Code Quality
npm run lint             # Run ESLint
npm run format           # Format code with Prettier
```

## Troubleshooting

### Port Already in Use
If port 3000 is taken, edit `vite.config.ts`:
```typescript
server: {
  port: 3001, // Change to any available port
  // ...
}
```

### API Connection Issues
1. Ensure backend services are running
2. Check proxy configuration in `vite.config.ts`
3. Verify `VITE_API_BASE_URL` in `.env`

### Feature Flags Not Working
1. Check browser console for Rox initialization messages
2. Verify API key is correct in `.env`
3. Ensure you restarted dev server after changing `.env`

### Build Errors
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Vite cache
rm -rf node_modules/.vite
npm run dev
```

## Next Steps

- **Customize Styling**: Edit `tailwind.config.js` and `src/styles/index.css`
- **Add Features**: Create new components in `src/components/`
- **Add Pages**: Create new pages in `src/pages/` and update `App.tsx`
- **Configure Flags**: Add more feature flags in `src/features/flags.ts`
- **Write Tests**: Add tests in `__tests__` directories

## Resources

- [React Documentation](https://react.dev/)
- [Vite Documentation](https://vitejs.dev/)
- [TanStack Query](https://tanstack.com/query/latest)
- [Tailwind CSS](https://tailwindcss.com/)
- [CloudBees FM Docs](https://docs.cloudbees.com/docs/cloudbees-feature-management/)

## Need Help?

- Check the main README.md for detailed documentation
- Review example components in `src/components/`
- Look at test files for usage examples
- Check browser console for errors and warnings
