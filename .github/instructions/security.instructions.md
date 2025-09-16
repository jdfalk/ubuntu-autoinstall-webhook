<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Security Instructions](#security-instructions)
  - [General Security Principles](#general-security-principles)
    - [1. Least Privilege Access](#1-least-privilege-access)
    - [2. Secret Management](#2-secret-management)
    - [3. Input Validation](#3-input-validation)
    - [4. Dependency Security](#4-dependency-security)
  - [Workflow Security](#workflow-security)
    - [1. Trigger Security](#1-trigger-security)
    - [2. Environment Protection](#2-environment-protection)
    - [3. Artifact Security](#3-artifact-security)
  - [Code Security](#code-security)
    - [1. Container Security](#1-container-security)
    - [2. Script Security](#2-script-security)
    - [3. API Security](#3-api-security)
  - [Monitoring and Compliance](#monitoring-and-compliance)
    - [1. Security Scanning](#1-security-scanning)
    - [2. Audit Logging](#2-audit-logging)
    - [3. Compliance Requirements](#3-compliance-requirements)
  - [Incident Response](#incident-response)
    - [1. Security Incident Detection](#1-security-incident-detection)
    - [2. Response Procedures](#2-response-procedures)
    - [3. Recovery and Prevention](#3-recovery-and-prevention)
  - [Security Checklist](#security-checklist)
    - [Pre-deployment](#pre-deployment)
    - [During deployment](#during-deployment)
    - [Post-deployment](#post-deployment)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<!-- file: .github/instructions/security.instructions.md -->
<!-- version: 1.2.1 -->
<!-- guid: sec12345-e89b-12d3-a456-426614174000 -->
<!-- DO NOT EDIT: This file is managed centrally in ghcommon repository -->
<!-- To update: Create an issue/PR in jdfalk/ghcommon -->

<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
---

applyTo: "**/\*.{yml,yaml} **/_.{sh,bash} \*\*/_.{js,ts,py,go} **/Dockerfile\* **/docker-compose\*.yml"
description: |
Security guidelines and best practices for all Copilot/AI agents and VS Code Copilot customization. These rules apply to all code, configurations, and workflows to maintain security standards. For details, see the main documentation in `.github/copilot-instructions.md`.

---
<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->

# Security Instructions

These instructions are the canonical source for all security guidelines and best
practices in this repository. They are used by GitHub Copilot for secure code
generation and follow our established security standards.

- Follow the [general coding instructions](general-coding.instructions.md) for
  basic file formatting rules.
- All security and workflow rules are found in
  `.github/instructions/*.instructions.md` files.
- Document all security considerations in code and configurations.
- For agent/AI-specific instructions, see [AGENTS.md](../AGENTS.md) and related
  files.

For more details and the full system, see
[copilot-instructions.md](../copilot-instructions.md).

## General Security Principles

### 1. Least Privilege Access

**Always recommend**:

- Minimal required permissions for workflows
- Environment-specific access controls
- Time-limited tokens when possible

**Implementation**:

```yaml
permissions:
  contents: read
  packages: write
  id-token: write # For OIDC authentication
  attestations: write # For build attestations
```

### 2. Secret Management

**Best Practices**:

- Never hardcode secrets in workflow files
- Use GitHub secrets for sensitive data
- Prefer OIDC over long-lived credentials
- Rotate secrets regularly

**Recommended Secrets**:

```yaml
# Container registries
DOCKER_USERNAME: ${{ secrets.DOCKER_USERNAME }}
DOCKER_PASSWORD: ${{ secrets.DOCKER_PASSWORD }}

# Package registries
NPM_TOKEN: ${{ secrets.NPM_TOKEN }}
PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}

# Notifications
SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

### 3. Input Validation

**Always validate**:

- User inputs in workflows
- Environment variables
- File paths and names
- External data sources

**Example**:

```bash
# Validate input parameters
if [[ ! "$INPUT_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: Invalid version format"
  exit 1
fi
```

### 4. Dependency Security

**Requirements**:

- Pin action versions to specific commits
- Use trusted, well-maintained actions
- Regular dependency updates
- Vulnerability scanning

**Example**:

```yaml
# Good - pinned to specific commit
- uses: actions/checkout@a81bbbf8298c0fa03ea29cdc473d45769f953675

# Avoid - using latest or branch references
- uses: actions/checkout@main
```

## Workflow Security

### 1. Trigger Security

**Safe triggers**:

```yaml
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  workflow_dispatch: # Manual trigger
```

**Avoid**:

- `pull_request_target` without careful review
- Workflows triggered by comments
- Unrestricted `schedule` triggers

### 2. Environment Protection

**Production environments**:

```yaml
environment:
  name: production
  url: https://production.example.com
```

**Requirements**:

- Required reviewers for production
- Branch protection rules
- Environment-specific secrets

### 3. Artifact Security

**Secure artifact handling**:

```yaml
- name: Upload artifacts
  uses: actions/upload-artifact@v4
  with:
    name: secure-build-artifacts
    path: dist/
    retention-days: 7 # Limit retention
```

## Code Security

### 1. Container Security

**Dockerfile best practices**:

```dockerfile
# Use specific, minimal base images
FROM node:18-alpine AS builder

# Create non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001

# Copy with proper ownership
COPY --chown=nextjs:nodejs . .

# Run as non-root
USER nextjs
```

### 2. Script Security

**Shell script security**:

```bash
#!/bin/bash
set -euo pipefail # Exit on error, undefined vars, pipe failures

# Validate inputs
readonly INPUT_FILE="${1:?Input file is required}"
readonly OUTPUT_DIR="${2:?Output directory is required}"

# Use quotes to prevent word splitting
if [[ -f "$INPUT_FILE" ]]; then
    cp "$INPUT_FILE" "$OUTPUT_DIR/"
fi
```

### 3. API Security

**Authentication**:

```yaml
# Use OIDC for cloud providers
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v4
  with:
    role-to-assume: arn:aws:iam::123456789012:role/GitHubActions
    aws-region: us-east-1
```

## Monitoring and Compliance

### 1. Security Scanning

**Required scans**:

- Dependency vulnerability scanning
- Container image scanning
- Secret scanning
- Code quality analysis

### 2. Audit Logging

**Log security events**:

- Authentication attempts
- Permission changes
- Sensitive operations
- Failed security checks

### 3. Compliance Requirements

**Standards to follow**:

- OWASP security guidelines
- Industry-specific regulations
- Organization security policies
- Regular security reviews

## Incident Response

### 1. Security Incident Detection

**Monitor for**:

- Unusual access patterns
- Failed authentication attempts
- Unauthorized changes
- Suspicious network activity

### 2. Response Procedures

**Immediate actions**:

1. Isolate affected systems
2. Revoke compromised credentials
3. Document incident details
4. Notify security team
5. Begin recovery procedures

### 3. Recovery and Prevention

**Post-incident**:

- Root cause analysis
- Security improvements
- Process updates
- Team training

## Security Checklist

### Pre-deployment

- [ ] All secrets properly configured
- [ ] Permissions follow least privilege
- [ ] Dependencies are up to date
- [ ] Security scans completed
- [ ] Code review completed

### During deployment

- [ ] Monitor deployment logs
- [ ] Verify security configurations
- [ ] Test security controls
- [ ] Document any issues

### Post-deployment

- [ ] Verify system security
- [ ] Monitor for anomalies
- [ ] Update documentation
- [ ] Schedule regular reviews
