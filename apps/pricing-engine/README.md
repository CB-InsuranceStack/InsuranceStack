# Insights API Service

A production-ready Go API service for managing financial insights and alerts with CloudBees Feature Management integration.

## Features

- RESTful API for insights and alerts management
- Real-time feature flag system using CloudBees Feature Management (Rox SDK)
- Feature flag: `api.insightsV2` - dynamically switch between insight calculation algorithms
- Feature flag: `api.alertsEnabled` - enable/disable alerts endpoint at runtime
- Proper error handling and structured logging
- CORS support
- Graceful shutdown
- Health check endpoint
- JSON structured logging

## Project Structure

```
api-insights/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── health.go           # Health check handler
│   │   ├── insights.go         # Insights endpoints
│   │   └── alerts.go           # Alerts endpoints
│   ├── services/                # Business logic
│   │   ├── insights_service.go # Insights business logic
│   │   └── alerts_service.go   # Alerts business logic
│   ├── repository/              # Data access layer
│   │   └── repository.go       # Repository implementation
│   ├── features/                # Feature flags
│   │   └── flags.go            # CloudBees FM/Rox integration
│   ├── models/                  # Data models
│   │   ├── insight.go          # Insight model
│   │   └── alert.go            # Alert model
│   └── middleware/              # HTTP middleware
│       ├── logging.go          # Request logging
│       ├── cors.go             # CORS configuration
│       └── auth.go             # Authentication
├── go.mod                       # Go module definition
├── Dockerfile                   # Docker configuration
├── Makefile                     # Build automation
└── README.md                    # This file
```

## API Endpoints

### Health Check

**GET /healthz**

Returns the health status of the service.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-12-13T10:30:00Z",
  "service": "api-insights"
}
```

### List Insights

**GET /insights**

Returns all insights for the authenticated user. When `insightsV2=true`, titles include "(V2)" suffix.

**Headers:**
- `X-User-ID` (optional): User ID, defaults to `user-001` if not provided

**Response (when insightsV2 = false):**
```json
[
  {
    "id": "insight-001",
    "userId": "user-001",
    "type": "spending_alert",
    "category": "food_dining",
    "title": "Higher than usual dining spending",
    "description": "You've spent $162.25 on dining this week, which is 35% higher than your average.",
    "severity": "medium",
    "createdAt": "2024-12-13T08:00:00Z",
    "actionable": true,
    "recommendation": "Consider meal planning to reduce dining out costs."
  }
]
```

**Response (when insightsV2 = true):**
```json
[
  {
    "id": "insight-001",
    "userId": "user-001",
    "type": "spending_alert",
    "category": "food_dining",
    "title": "Higher than usual dining spending (V2)",
    "description": "You've spent $162.25 on dining this week, which is 35% higher than your average.",
    "severity": "medium",
    "createdAt": "2024-12-13T08:00:00Z",
    "actionable": true,
    "recommendation": "Consider meal planning to reduce dining out costs."
  }
]
```

### Get Insight by ID

**GET /insights/{id}**

Returns a specific insight by ID. The insight must belong to the authenticated user.

**Headers:**
- `X-User-ID` (optional): User ID, defaults to `user-001` if not provided

**Parameters:**
- `id` (path): Insight ID

**Response:**
```json
{
  "id": "insight-001",
  "userId": "user-001",
  "type": "spending_alert",
  "category": "food_dining",
  "title": "Higher than usual dining spending",
  "description": "You've spent $162.25 on dining this week, which is 35% higher than your average.",
  "severity": "medium",
  "createdAt": "2024-12-13T08:00:00Z",
  "actionable": true,
  "recommendation": "Consider meal planning to reduce dining out costs."
}
```

**Error Responses:**
- `404 Not Found` - Insight does not exist
- `403 Forbidden` - Insight does not belong to the user

### List Alerts

**GET /alerts**

Returns all alerts for the authenticated user. Returns `503 Service Unavailable` if `alertsEnabled=false`.

**Headers:**
- `X-User-ID` (optional): User ID, defaults to `user-001` if not provided

**Response (when alertsEnabled = true):**
```json
[
  {
    "id": "alert-001",
    "userId": "user-001",
    "type": "spending_alert",
    "title": "Higher than usual dining spending",
    "message": "You've spent $162.25 on dining this week, which is 35% higher than your average.",
    "priority": "high",
    "createdAt": "2024-12-13T08:00:00Z",
    "read": false
  }
]
```

**Response (when alertsEnabled = false):**
```json
{
  "error": "Service Unavailable",
  "message": "Alerts feature is currently disabled"
}
```

**Status Code:** `503 Service Unavailable` when feature is disabled

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8003` |
| `DATA_PATH` | Path to seed data directory | `../../data/seed` |
| `CLOUDBEES_FM_API_KEY` | CloudBees Feature Management API key | `dev-mode` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `FEATURE_INSIGHTS_V2` | Enable V2 algorithm in dev mode (true/false) | `false` |
| `FEATURE_ALERTS_ENABLED` | Enable alerts in dev mode (true/false) | `true` |

## Feature Flags

### api.insightsV2

**Default:** `false`

When enabled, the V2 insights calculation algorithm is used. In this demo implementation, V2 insights have "(V2)" appended to their titles to demonstrate the feature flag is active.

**Use Cases:**
- A/B testing different insight algorithms
- Rolling out new calculation logic gradually
- Instant rollback if issues are detected

### api.alertsEnabled

**Default:** `true`

Controls whether the alerts endpoint is available. When disabled, `GET /alerts` returns `503 Service Unavailable`.

**Use Cases:**
- Disable alerts during system maintenance
- Control feature access by customer tier
- Emergency feature kill switch

### CloudBees Integration

The service uses the CloudBees Rox SDK (`github.com/rollout/rox-go/v5/core`) for real-time feature flag management. Flags can be toggled instantly without redeploying the service.

**Development Mode:** When no CloudBees API key is provided, the service falls back to environment variables (`FEATURE_INSIGHTS_V2` and `FEATURE_ALERTS_ENABLED`).

**Production Mode:** Provide `CLOUDBEES_FM_API_KEY` to use CloudBees Feature Management for centralized control and real-time updates.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Access to seed data files at `/data/seed/insights.json`
- CloudBees Feature Management account (optional, for production)

### Installation

1. Install dependencies:

```bash
cd /Users/brown/git_orgs/CB-AccountStack/AccountStack/apps/api-insights
go mod download
```

2. Set up environment variables:

```bash
export PORT=8003
export DATA_PATH=/Users/brown/git_orgs/CB-AccountStack/AccountStack/data/seed
export CLOUDBEES_FM_API_KEY=your-api-key-here
export LOG_LEVEL=info
```

3. Run the server:

```bash
go run cmd/server/main.go
```

Or build and run:

```bash
go build -o bin/insights-api cmd/server/main.go
./bin/insights-api
```

### Using Make

```bash
# Build the application
make build

# Run in development mode
make run-dev

# Run tests
make test

# Format code
make fmt

# Build Docker image
make docker-build

# Run in Docker
make docker-run

# Show all available commands
make help
```

## Testing the API

### Health Check
```bash
curl http://localhost:8003/healthz
```

### Get All Insights (default algorithm)
```bash
curl http://localhost:8003/insights
```

### Get All Insights (V2 algorithm)
```bash
# Start server with V2 enabled
export FEATURE_INSIGHTS_V2=true
go run cmd/server/main.go

# In another terminal
curl http://localhost:8003/insights
# Titles will have "(V2)" suffix
```

### Get Specific Insight
```bash
curl http://localhost:8003/insights/insight-001
```

### Get All Alerts
```bash
curl http://localhost:8003/alerts
```

### Test Alerts Disabled
```bash
# Start server with alerts disabled
export FEATURE_ALERTS_ENABLED=false
go run cmd/server/main.go

# In another terminal
curl http://localhost:8003/alerts
# Returns 503 Service Unavailable
```

### Test with Different User
```bash
curl -H "X-User-ID: user-002" http://localhost:8003/insights
```

## Development

### Project Layout

The service follows Go best practices with a clean architecture:

1. **Handlers Layer** (`internal/handlers/`): HTTP request/response handling
2. **Services Layer** (`internal/services/`): Business logic and feature flag application
3. **Repository Layer** (`internal/repository/`): Data access abstraction
4. **Models Layer** (`internal/models/`): Domain models and data structures
5. **Features Layer** (`internal/features/`): Feature flag management

### Middleware

- **Logging**: Logs all HTTP requests with method, path, status, and duration
- **CORS**: Handles cross-origin resource sharing
- **Auth**: Extracts and validates user authentication (X-User-ID header)

### Feature Flag Architecture

Feature flags are initialized on startup and can be updated in real-time via CloudBees. The service checks flag status on each request, allowing instant behavior changes without downtime.

```go
// Check if V2 algorithm should be used
if flags.IsInsightsV2Enabled() {
    insights = applyV2Algorithm(insights)
}

// Check if alerts are enabled
if !flags.IsAlertsEnabled() {
    return 503 // Service Unavailable
}
```

## Docker

### Build Image
```bash
docker build -t accountstack/api-insights:latest .
```

### Run Container
```bash
docker run -p 8003:8003 \
  -e PORT=8003 \
  -e LOG_LEVEL=info \
  -e CLOUDBEES_FM_API_KEY=your-key \
  -v /path/to/data/seed:/data/seed \
  -e DATA_PATH=/data/seed \
  accountstack/api-insights:latest
```

## Production Considerations

1. **Authentication**: The current implementation uses a simple `X-User-ID` header for demo purposes. In production, implement proper JWT token validation or OAuth2.

2. **CORS**: The CORS middleware currently allows all origins (`*`). In production, specify exact allowed origins.

3. **Database**: Data is loaded from JSON files. In production, integrate with a proper database (PostgreSQL, MySQL, etc.).

4. **Monitoring**: Add metrics collection (Prometheus), distributed tracing (OpenTelemetry), and error tracking (Sentry).

5. **Rate Limiting**: Implement rate limiting to prevent abuse.

6. **TLS**: Enable HTTPS with proper certificates.

7. **Secrets Management**: Use a secrets manager (AWS Secrets Manager, HashiCorp Vault) for sensitive configuration.

8. **Feature Flag Management**: Use CloudBees Feature Management dashboard to control flags across environments.

## Data Model

### Insight
- `id`: Unique identifier
- `userId`: Owner user ID
- `type`: Insight type (spending_alert, savings_opportunity, etc.)
- `category`: Category (food_dining, utilities, etc.)
- `title`: Insight title
- `description`: Detailed description
- `severity`: Severity level (info, low, medium, high)
- `createdAt`: Creation timestamp
- `actionable`: Whether user can take action
- `recommendation`: Optional recommendation text

### Alert
- `id`: Unique identifier
- `userId`: Owner user ID
- `type`: Alert type
- `title`: Alert title
- `message`: Alert message
- `priority`: Priority level (medium, high, critical)
- `createdAt`: Creation timestamp
- `read`: Read status
- `actionUrl`: Optional action URL

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests manually
go test -v -race ./...
```

## License

Copyright 2024 CB-AccountStack

## Quick Start

The fastest way to get started:

```bash
# From the api-insights directory
export DATA_PATH=../../data/seed
go run cmd/server/main.go
```

Then test the API:

```bash
# Health check
curl http://localhost:8003/healthz

# Get insights (default algorithm)
curl http://localhost:8003/insights

# Get insights (V2 algorithm)
export FEATURE_INSIGHTS_V2=true
go run cmd/server/main.go &
curl http://localhost:8003/insights
# Titles will have "(V2)" suffix

# Get alerts
curl http://localhost:8003/alerts

# Test alerts disabled
export FEATURE_ALERTS_ENABLED=false
go run cmd/server/main.go &
curl http://localhost:8003/alerts
# Returns 503 Service Unavailable
```

## CloudBees Feature Management Setup

1. Sign up for CloudBees Feature Management
2. Create a new application
3. Get your API key
4. Set the environment variable:
   ```bash
   export CLOUDBEES_FM_API_KEY=your-actual-api-key
   ```
5. Run the service - it will automatically register flags with CloudBees
6. Toggle flags in the CloudBees dashboard - changes apply instantly!

## Support

For issues or questions, please contact the development team.
