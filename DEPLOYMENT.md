# AccountStack Deployment Guide

Complete guide for running AccountStack locally and deploying to production.

## Table of Contents
- [Quick Start (Local)](#quick-start-local)
- [Running Individual Services](#running-individual-services)
- [Feature Flag Configuration](#feature-flag-configuration)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Production Deployment](#production-deployment)

---

## Quick Start (Local)

### Prerequisites
- Docker & Docker Compose installed
- (Optional) CloudBees FM API key for real-time flag updates

### Start Everything

```bash
# Clone the repository (if not already done)
git clone https://github.com/CB-AccountStack/AccountStack.git
cd AccountStack

# (Optional) Set CloudBees FM API key
export CLOUDBEES_FM_API_KEY=your-api-key-here

# Start all services
docker compose up --build
```

**Services will be available at:**
- UI: http://localhost:3000
- Accounts API: http://localhost:8001
- Transactions API: http://localhost:8002
- Insights API: http://localhost:8003

### Stop Everything

```bash
# Stop services
docker compose down

# Stop and remove volumes
docker compose down -v
```

---

## Running Individual Services

### React UI (Development Mode)

```bash
cd apps/web

# Install dependencies (first time only)
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Accounts API

```bash
cd apps/api-accounts

# Install dependencies (first time only)
go mod download

# Run directly with Go
export DATA_PATH=../../data/seed
go run cmd/server/main.go

# Or build and run
go build -o bin/accounts-api cmd/server/main.go
DATA_PATH=../../data/seed ./bin/accounts-api

# Or use Makefile
make run-dev
```

### Transactions API

```bash
cd apps/api-transactions

# Run with Makefile
make run-dev

# Or build and run
make build
DATA_PATH=../../data/seed ./bin/transactions-api
```

### Insights API

```bash
cd apps/api-insights

# Run with Makefile
make run-dev

# Or build and run
make build
DATA_PATH=../../data/seed ./bin/insights-api
```

---

## Feature Flag Configuration

### Option 1: CloudBees Feature Management (Recommended)

1. **Get API Key**:
   - Log in to CloudBees platform
   - Navigate to Feature Management
   - Create or copy your API key

2. **Set Environment Variable**:
   ```bash
   export CLOUDBEES_FM_API_KEY=your-api-key-here
   ```

3. **Start Services**:
   ```bash
   docker compose up
   ```

4. **Configure Flags in Dashboard**:
   - UI flags are in `accountstack` namespace
   - API flags are in `api.accounts`, `api.transactions`, `api.insights` namespaces
   - Changes propagate in real-time (no reload needed!)

### Option 2: Environment Variables (Local Development)

Override individual flags:

```bash
# UI Flags
export VITE_FLAG_UI_DASHBOARD_CARDS_V2=false
export VITE_FLAG_UI_INSIGHTS_V2=true
export VITE_FLAG_UI_ALERTS_BANNER=false

# API Flags
export FEATURE_MASK_AMOUNTS=true
export FEATURE_ADVANCED_FILTERS=true
export FEATURE_INSIGHTS_V2=true
export FEATURE_ALERTS_ENABLED=false
```

### Option 3: Hardcoded Defaults (Offline Mode)

If no CloudBees FM API key is set and no environment variables are provided, the application uses hardcoded defaults:

**UI Defaults**:
- `ui.dashboardCardsV2`: true
- `ui.insightsV2`: false
- `ui.alertsBanner`: true
- `ui.transactionsFilters`: true
- `kill.ui.insights`: false

**API Defaults**:
- `api.maskAmounts`: false
- `api.advancedFilters`: false
- `api.insightsV2`: false
- `api.alertsEnabled`: true

---

## Testing

### All Tests

```bash
# From repository root
make test
```

### Unit Tests

```bash
# React UI
cd apps/web
npm run test:unit

# Go APIs
cd apps/api-accounts && go test -v ./...
cd apps/api-transactions && go test -v ./...
cd apps/api-insights && go test -v ./...
```

### Integration Tests

```bash
# Start services first
docker compose up -d

# Wait for services to be ready
sleep 5

# Run integration tests
cd tests/integration
go test -v ./...
```

### E2E Tests (Playwright)

```bash
cd apps/web

# Install Playwright (first time only)
npm install
npx playwright install --with-deps

# Run E2E tests
npm run test:e2e
```

### Health Checks

```bash
# From repository root
make health

# Manual checks
curl http://localhost:8001/healthz  # Accounts
curl http://localhost:8002/healthz  # Transactions
curl http://localhost:8003/healthz  # Insights
curl http://localhost:3000          # UI
```

---

## Troubleshooting

### Services Won't Start

**Issue**: `Error: bind: address already in use`

**Solution**: Another process is using the ports
```bash
# Find and kill process on port 8001
lsof -ti:8001 | xargs kill -9

# Or use different ports
export PORT=9001
docker compose up
```

### API Returns 404

**Issue**: UI can't reach APIs

**Solution**: Verify API is running and accessible
```bash
# Check API directly
curl http://localhost:8001/healthz

# Check Docker network
docker compose ps

# Restart services
docker compose restart
```

### Feature Flags Not Working

**Issue**: Flags not changing behavior

**Solutions**:
1. **Check CloudBees FM connection**:
   ```bash
   # Verify API key is set
   echo $CLOUDBEES_FM_API_KEY

   # Check logs for connection status
   docker compose logs api-accounts | grep -i feature
   ```

2. **Verify flag names**:
   - UI flags: `ui.dashboardCardsV2` (camelCase, dot notation)
   - API flags: `api.maskAmounts` (camelCase, dot notation)

3. **Clear browser cache** (for UI flags):
   - Hard refresh: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)

### Data Not Loading

**Issue**: "No accounts found" or similar

**Solution**: Verify DATA_PATH environment variable
```bash
# Check current value
echo $DATA_PATH

# Set correctly
export DATA_PATH=../../data/seed

# Or use absolute path
export DATA_PATH=/Users/brown/git_orgs/CB-AccountStack/AccountStack/data/seed
```

### Docker Build Fails

**Issue**: Build errors or out of disk space

**Solutions**:
```bash
# Clean up Docker
docker system prune -a

# Rebuild from scratch
docker compose build --no-cache

# Check disk space
df -h
```

### React UI Shows Blank Page

**Solutions**:
1. **Check browser console** for JavaScript errors
2. **Verify API connectivity**:
   ```bash
   # APIs must be running
   curl http://localhost:8001/healthz
   ```
3. **Rebuild UI**:
   ```bash
   cd apps/web
   rm -rf node_modules dist
   npm install
   npm run build
   ```

---

## Production Deployment

### Prerequisites

- CloudBees Unify platform access
- Docker registry access (ghcr.io or private registry)
- Kubernetes cluster (optional)
- CloudBees FM API key

### Build Docker Images

```bash
# Build all images
docker compose build

# Tag for registry
docker tag accountstack-web:latest ghcr.io/cb-accountstack/accountstack-web:v0.2.0
docker tag accountstack-api-accounts:latest ghcr.io/cb-accountstack/accountstack-api-accounts:v0.2.0
docker tag accountstack-api-transactions:latest ghcr.io/cb-accountstack/accountstack-api-transactions:v0.2.0
docker tag accountstack-api-insights:latest ghcr.io/cb-accountstack/accountstack-api-insights:v0.2.0

# Push to registry
docker push ghcr.io/cb-accountstack/accountstack-web:v0.2.0
docker push ghcr.io/cb-accountstack/accountstack-api-accounts:v0.2.0
docker push ghcr.io/cb-accountstack/accountstack-api-transactions:v0.2.0
docker push ghcr.io/cb-accountstack/accountstack-api-insights:v0.2.0
```

### CloudBees Unify Setup

1. **Add Repository to Unify**:
   - Go to CloudBees Unify dashboard
   - Add GitHub repository: `CB-AccountStack/AccountStack`
   - Configure branch protection and triggers

2. **Set Environment Variables**:
   ```yaml
   CLOUDBEES_FM_API_KEY: <your-api-key>
   CLOUDBEES_FM_ENVIRONMENT: production
   DOCKER_REGISTRY: ghcr.io
   IMAGE_PREFIX: cb-accountstack/accountstack
   ```

3. **Trigger Workflow**:
   - Push to `main` branch triggers build
   - Workflow defined in `.cloudbees/workflows/build-and-test.yaml`
   - Parallel builds for all components
   - Automated testing and security scanning

### Kubernetes Deployment (Optional)

Create deployment manifests:

```yaml
# k8s/web-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: accountstack-web
spec:
  replicas: 2
  selector:
    matchLabels:
      app: accountstack-web
  template:
    metadata:
      labels:
        app: accountstack-web
    spec:
      containers:
      - name: web
        image: ghcr.io/cb-accountstack/accountstack-web:v0.2.0
        ports:
        - containerPort: 3000
        env:
        - name: VITE_CLOUDBEES_FM_API_KEY
          valueFrom:
            secretKeyRef:
              name: accountstack-secrets
              key: fm-api-key
---
apiVersion: v1
kind: Service
metadata:
  name: accountstack-web
spec:
  selector:
    app: accountstack-web
  ports:
  - port: 80
    targetPort: 3000
  type: LoadBalancer
```

Apply with:
```bash
kubectl apply -f k8s/
```

### Environment-Specific Configuration

**Development**:
```bash
CLOUDBEES_FM_ENVIRONMENT=dev
LOG_LEVEL=debug
```

**Staging**:
```bash
CLOUDBEES_FM_ENVIRONMENT=staging
LOG_LEVEL=info
```

**Production**:
```bash
CLOUDBEES_FM_ENVIRONMENT=prod
LOG_LEVEL=warn
```

### Monitoring & Health Checks

**Kubernetes Health Probes**:
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8001
  initialDelaySeconds: 10
  periodSeconds: 30

readinessProbe:
  httpGet:
    path: /healthz
    port: 8001
  initialDelaySeconds: 5
  periodSeconds: 10
```

**Logging**:
- All services use structured JSON logging
- Logs are sent to stdout/stderr
- Compatible with CloudBees Analytics

**Metrics** (Future):
- Prometheus endpoints (planned)
- CloudBees Observability integration (planned)

---

## Quick Reference

### Useful Commands

```bash
# Start everything
make up

# Stop everything
make down

# Run tests
make test

# Check service health
make health

# View logs
make logs

# Restart services
make restart

# Clean up
make clean
```

### Port Reference

| Service | Port | Endpoint |
|---------|------|----------|
| Web UI | 3000 | http://localhost:3000 |
| Accounts API | 8001 | http://localhost:8001 |
| Transactions API | 8002 | http://localhost:8002 |
| Insights API | 8003 | http://localhost:8003 |

### Important Files

- `docker-compose.yaml` - Service orchestration
- `Makefile` - Development commands
- `.env.example` - Environment variable template
- `CHANGELOG.md` - Version history
- `config/README.md` - Feature flags reference

### Getting Help

- **Issues**: https://github.com/CB-AccountStack/AccountStack/issues
- **Documentation**: See README files in each component
- **CloudBees Support**: Contact CloudBees SE team

---

**Last Updated**: 2025-12-13
**Version**: 0.2.0
