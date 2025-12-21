# AccountStack Helm Chart

This Helm chart deploys the AccountStack application to Kubernetes.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- kubectl configured to access your cluster

## Installation

### 1. Configure values

Create a `custom-values.yaml` file with your configuration:

```yaml
# IMPORTANT: Change these values in production!
auth:
  jwtSecret: "your-super-secret-jwt-key-minimum-32-characters-long"
  username: "your-admin-email@example.com"
  password: "your-secure-password"

cloudbees:
  apiKey: "your-cloudbees-api-key"
  environment: "production"

web:
  ingress:
    enabled: true
    hosts:
      - host: accountstack.yourdomain.com
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: accountstack-tls
        hosts:
          - accountstack.yourdomain.com
```

### 2. Install the chart

```bash
helm install accountstack ./helm/accountstack -f custom-values.yaml
```

### 3. Upgrade the deployment

```bash
helm upgrade accountstack ./helm/accountstack -f custom-values.yaml
```

### 4. Uninstall

```bash
helm uninstall accountstack
```

## Configuration

The following table lists the configurable parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `auth.jwtSecret` | JWT secret for token signing | `dev-secret-key...` |
| `auth.username` | Admin username | `demo@accountstack.com` |
| `auth.password` | Admin password | `demo123` |
| `cloudbees.apiKey` | CloudBees Feature Management API key | `dev-mode` |
| `cloudbees.environment` | CloudBees environment | `production` |
| `web.replicaCount` | Number of web replicas | `2` |
| `apiAccounts.replicaCount` | Number of accounts API replicas | `2` |
| `apiTransactions.replicaCount` | Number of transactions API replicas | `2` |
| `apiInsights.replicaCount` | Number of insights API replicas | `2` |

## Security Notes

- **IMPORTANT**: The default credentials are for demo purposes only
- **ALWAYS** change `auth.jwtSecret`, `auth.username`, and `auth.password` in production
- Use Kubernetes Secrets or a secret management solution (e.g., HashiCorp Vault, AWS Secrets Manager)
- Enable TLS/HTTPS for all ingress endpoints
- Consider using cert-manager for automatic TLS certificate management

## Monitoring

The deployments include liveness and readiness probes on `/healthz` endpoints.

## Support

For issues and questions, please refer to the main project documentation.
