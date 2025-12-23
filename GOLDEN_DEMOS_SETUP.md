# Golden Demos Setup Progress

This document tracks the progress of adding InsuranceStack and AccountStack to CloudBees Golden Demos (Mimic).

## Overview

**Goal**: Integrate InsuranceStack and AccountStack into CloudBees Golden Demos platform to enable one-command deployment for demos.

**Applications**:
- InsuranceStack (this repo)
- AccountStack (CB-AccountStack/AccountStack)

## Progress Log

### 2025-12-23

#### KUBECONFIG Secret Migration ✅
- **Status**: **SUCCESS** - Deployment working with org-level property
- **Change**: Migrated KUBECONFIG from environment-level secret to organization-level property
- **Rationale**: Allows reuse across multiple environments and projects
- **Key Finding**: CloudBees Unify resolves org-level properties even when Helm charts reference environment-level secrets
  - Helm chart expects secret at environment level
  - Works seamlessly with org-level property instead
  - Simplifies golden demos setup (one KUBECONFIG for all demos)
- **Testing**: ✅ Verified in InsuranceStack deployment
- **Next Steps**:
  - Apply same pattern to AccountStack
  - Consider moving FM_KEY to org-level property as well
  - Document this secret resolution behavior for Mimic integration

## Configuration Requirements

### Secrets & Properties

| Name | Type | Level | Status | Notes |
|------|------|-------|--------|-------|
| KUBECONFIG | Property | Organization | ✅ Working | Org-level property works despite Helm expecting env-level |
| FM_KEY | Secret | Repository | Configured | CloudBees Feature Management API key |

### Deployment Pattern

Both applications use:
- **Architecture**: Full-stack monorepo with multiple components
- **Database**: In-memory (no external database required)
- **Ingress**: Path-based routing for multi-tenancy
- **Feature Management**: Reactive SSE for instant flag updates

## Mimic Integration Plan

### Phase 1: Configuration Setup
- [x] Validate KUBECONFIG as org-level property
- [x] Document required secrets and properties
- [x] Test deployment workflows with new configuration

### Phase 2: Mimic Scenario Creation
- [ ] Create `insurancestack-demo` scenario
- [ ] Create `accountstack-demo` scenario
- [ ] Test end-to-end deployment via Mimic
- [ ] Document post-setup manual steps

### Phase 3: Documentation
- [ ] Update Confluence pages with Mimic instructions
- [ ] Create demo scripts and talking points
- [ ] Record demo videos

## Open Questions

1. Should FM_KEY also be moved to org-level property?
2. Any additional configuration needed for Mimic automation?
3. Deployment target cluster configuration?

## Related Resources

- [InsuranceStack Confluence](https://cloudbees.atlassian.net/wiki/spaces/UDEMO/pages/5809700869)
- [AccountStack Confluence](https://cloudbees.atlassian.net/wiki/spaces/UDEMO/pages/5781618696)
- [Golden Demos Overview](https://cloudbees.atlassian.net/wiki/spaces/UDEMO/pages/5779816519)
- [Mimic GitHub Repository](https://github.com/cb-demos/mimic)
