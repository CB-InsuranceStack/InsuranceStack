# AccountStack + Mimic Integration

**Meeting Date:** [Tomorrow]
**Attendees:** Logan, Stu Brown
**Subject:** Integrating AccountStack demo scenario into Mimic

---

## What is AccountStack?

AccountStack is a complete financial platform demo showcasing CloudBees Unify capabilities:

- **Architecture:** Monorepo with 4 applications
  - React web UI with Feature Management integration
  - 3 Go microservices: api-accounts, api-transactions, api-insights
- **Single CloudBees Unify component** (simplified for demos)
- **No database dependencies** (uses in-memory data)
- **166+ comprehensive tests** across all APIs (great for Unify test reporting demos)
- **Current Repository:** `CB-AccountStack/AccountStack` (GitHub)

**Demo Value:**
- Multi-service application with feature flags
- Rich test reporting (166 tests)
- Artifact tracking (4 deployed components)
- Helm-based Kubernetes deployment
- CI/CD workflows already configured

---

## Current Mimic Scenario Status

We've created a working Mimic scenario (`accountstack-demo.yaml` - see below) that:

✅ Creates GitHub repository from template
✅ Creates CloudBees Unify component
✅ Creates CloudBees Unify environment
✅ Creates CloudBees Unify application
✅ Creates 3 feature flags (alerts-banner, transaction-insights, dark-mode)

**Current Manual Post-Setup Required:**
1. Add `KUBECONFIG` secret to CloudBees Unify environment
2. Add `FM_KEY` secret to GitHub repository

---

## Infrastructure Requirements

### 1. **Nginx Ingress Controller with DNS01 Challenge**

AccountStack requires:
- **Ingress controller:** nginx
- **TLS/SSL:** cert-manager with LetsEncrypt
- **DNS challenge:** DNS01 provisioner (not HTTP01)
- **Why DNS01:** We use wildcard certificates for flexibility

Current Helm configuration:
```yaml
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: accountstack.se-main-demo.sa-demo.beescloud.com
      paths:
        - path: /{org}/{env}
          pathType: Prefix
  tls:
    - secretName: accountstack-tls
      hosts:
        - accountstack.se-main-demo.sa-demo.beescloud.com
```

### 2. **Suffix-Based Routing per Instance/Environment**

Each deployment uses a unique path prefix:
- **Pattern:** `/{github-org}/{environment-name}`
- **Examples:**
  - `https://domain.com/cb-demos/dev` → cb-demos org, dev environment
  - `https://domain.com/stubrowncloudbees/staging` → stubrowncloudbees org, staging environment

**Why:** Multiple users can demo simultaneously without conflicts

**Deploy workflow calculates this automatically:**
```yaml
PATH_PREFIX="/${ORG}/${NAMESPACE}"
```

### 3. **Feature Management Integration**

AccountStack web UI uses CloudBees Feature Management (ROX SDK):
- **Requirement:** FM_KEY secret in GitHub repository
- **Usage:** Controls feature flags (alerts-banner, transaction-insights, dark-mode)
- **Fallback:** Deploy workflow has `local-mode` fallback if FM_KEY missing

**Current scenario:** User must manually add FM_KEY to GitHub repo after creation

**Question for Logan:** Can Mimic automatically add GitHub repo secrets during scenario run?

### 4. **Kubernetes Deployment (KUBECONFIG)**

Each CloudBees Unify environment needs its own KUBECONFIG secret:
- **Format:** Base64-encoded kubeconfig file content
- **Storage:** CloudBees Unify environment-scoped secret (NOT org-level)
- **Usage:** Deploy workflow uses `guru-actions/kubeconfig@1.16` action

**Current situation:**
- Stu has a dedicated cluster for demo purposes
- Can generate KUBECONFIGs for demo users
- Currently manual: users must add KUBECONFIG to environment after Mimic creates it

**Question for Logan:**
- Can Mimic pre-populate environment secrets during scenario creation?
- Or should we provide a shared KUBECONFIG for all demo instances?

---

## Questions for Logan

### 1. **Repository Location**
Where do you want AccountStack scenario repositories created?
- **Current:** `CB-AccountStack` GitHub org
- **Alternative:** Could fork to `cb-demos` or another org
- **Your preference?**

### 2. **Scenario Pack Distribution**
Should we:
- **Option A:** Add `accountstack-demo.yaml` to official `mimic-scenarios` repo
- **Option B:** Create separate `cb-demos/accountstack-scenarios` repo
- **Option C:** Keep it internal/private

### 3. **Secret Management Automation**
Can we automate these secrets in Mimic scenarios?
- **KUBECONFIG** → CloudBees Unify environment secret (required for K8s deploy)
- **FM_KEY** → GitHub repository secret (required for feature management)

Currently both are manual post-setup steps.

### 4. **Shared Infrastructure**
- Can we provide a shared demo Kubernetes cluster KUBECONFIG?
- Or should each user bring their own cluster?
- Stu has a cluster available for this purpose

### 5. **DNS Configuration**
- What base domain should we use? (current: `accountstack.se-main-demo.sa-demo.beescloud.com`)
- Can Mimic parameterize the base domain?
- Do you have preferred demo infrastructure domains?

---

## Technical Architecture

### Deployment Flow
```
1. Mimic creates GitHub repo from CB-AccountStack/AccountStack
2. Mimic creates CloudBees Unify component + environment + app + flags
3. [MANUAL] User adds KUBECONFIG to environment
4. [MANUAL] User adds FM_KEY to GitHub repo
5. Push to main branch triggers CloudBees workflow
6. Workflow:
   - Runs 166 tests (publishes results to Unify)
   - Builds 4 Docker images with kaniko (captures artifact UUIDs)
   - Deploys to Kubernetes via Helm
   - Registers 4 deployed artifacts in Unify
7. Access at: https://{domain}/{org}/{env}
```

### What Makes This Demo Valuable

**For CloudBees Unify:**
- ✅ Multi-service application (4 components in monorepo)
- ✅ Rich test reporting (166 tests across 3 Go APIs)
- ✅ Artifact tracking (4 artifacts with UUIDs)
- ✅ Feature Management integration (3 feature flags)
- ✅ Multi-environment deployment patterns
- ✅ Complete CI/CD workflows (build, test, deploy)
- ✅ Kubernetes deployment with Helm
- ✅ Path-based multi-tenancy (multiple demos simultaneously)

**For Sales/Demo Purposes:**
- Financial platform is relatable (banking, transactions, insights)
- Shows real-world microservices architecture
- Feature flags control visible functionality
- Tests demonstrate quality gates
- Deployment tracking shows artifact lifecycle

---

## Next Steps

1. **Review requirements** with Logan
2. **Decide on repository location** and access
3. **Determine secret automation** approach (KUBECONFIG, FM_KEY)
4. **Finalize infrastructure** (cluster, domain, ingress)
5. **Test end-to-end scenario** with real deployment
6. **Document post-setup steps** for demo users
7. **Add to official scenario pack** (if approved)

---

## Mimic Scenario File

Current scenario definition (`accountstack-demo.yaml`):

```yaml
id: accountstack-demo
name: AccountStack Financial Platform
summary: Multi-service financial demo platform (web + 3 APIs in monorepo)
details: |
  AccountStack is a complete financial platform demo showcasing:
  - React web application with feature management
  - 3 Go microservices (accounts, transactions, insights)
  - CloudBees CI/CD workflows with test reporting
  - Helm-based Kubernetes deployment
  - 166+ comprehensive tests for Unify test reporting demo

  POST-SETUP REQUIRED:
  After running this scenario, manually add these secrets:

  1. KUBECONFIG (Environment secret):
     - Go to CloudBees Unify → Environments → [your environment name]
     - Add secret: KUBECONFIG (paste your kubeconfig content)
     - Required for Kubernetes deployment

  2. FM_KEY (GitHub repository secret):
     - Go to GitHub → [your repo] → Settings → Secrets and variables → Actions
     - Add repository secret: FM_KEY (paste your Feature Management SDK key)
     - Required for feature flag functionality

parameter_schema:
  properties:
    project_name:
      type: string
      description: Project name (will be used for repo, component, and application)
      placeholder: accountstack
      pattern: "^[a-z0-9-]+$"
    target_org:
      type: string
      description: GitHub organization for repository
      placeholder: cb-demos
      pattern: "^[a-zA-Z0-9-]+$"
    environment:
      type: string
      description: Deployment environment
      enum: ["dev", "staging", "prod"]
      default: dev
  required:
    - project_name
    - target_org
    - environment

repositories:
  - source: CB-AccountStack/AccountStack
    target_org: "${target_org}"
    repo_name_template: "${project_name}"
    create_component: true
    replacements:
      "accountstack": "${project_name}"
      "dev": "${environment}"
    files_to_modify:
      - README.md
      - helm/accountstack/values.yaml
      - docker-compose.yaml

environments:
  - name: "${environment}"
    env:
      - name: ENVIRONMENT
        value: "${environment}"
    create_fm_token_var: true
    flags:
      - alerts-banner
      - transaction-insights
      - dark-mode

applications:
  - name: "${project_name}-app"
    repository: "${target_org}/${project_name}"
    components:
      - "${project_name}"
    environments:
      - "${environment}"

flags:
  - name: alerts-banner
    type: boolean
  - name: transaction-insights
    type: boolean
  - name: dark-mode
    type: boolean
```

---

## Contact

**Stu Brown**
GitHub: @stubrowncloudbees
CloudBees Org: stu-temp-tests (c7e80c5d-7f5c-4972-9683-1cf92e6ea6ce)
