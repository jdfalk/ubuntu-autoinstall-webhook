# Admin Guide for Ubuntu Autoinstall Webhook

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Ubuntu Autoinstall Webhook Admin Guide](#ubuntu-autoinstall-webhook-admin-guide)
  - [1. Installation](#1-installation)
    - [1.1. System Requirements](#11-system-requirements)
    - [1.2. Installation Methods](#12-installation-methods)
      - [1.2.1. Binary Installation](#121-binary-installation)
      - [1.2.2. Docker Installation](#122-docker-installation)
      - [1.2.3. Kubernetes Installation](#123-kubernetes-installation)
    - [1.3. First-Time Setup](#13-first-time-setup)
  - [2. System Configuration](#2-system-configuration)
    - [2.1. Configuration Files](#21-configuration-files)
    - [2.2. Environment Variables](#22-environment-variables)
    - [2.3. Command-Line Options](#23-command-line-options)
  - [3. Service Management](#3-service-management)
    - [3.1. Running as a Systemd Service](#31-running-as-a-systemd-service)
    - [3.2. Running in Docker](#32-running-in-docker)
    - [3.3. Running in Kubernetes](#33-running-in-kubernetes)
    - [3.4. Service Dependencies](#34-service-dependencies)
    - [3.5. Health Checks](#35-health-checks)
  - [4. User Management](#4-user-management)
    - [4.1. Authentication Methods](#41-authentication-methods)
    - [4.2. User Roles and Permissions](#42-user-roles-and-permissions)
    - [4.3. Managing Users](#43-managing-users)
    - [4.4. Password Policies](#44-password-policies)
    - [4.5. External Authentication](#45-external-authentication)
  - [5. Security Management](#5-security-management)
    - [5.1. PKI Infrastructure](#51-pki-infrastructure)
    - [5.2. Certificate Management](#52-certificate-management)
    - [5.3. TLS Configuration](#53-tls-configuration)
    - [5.4. API Security](#54-api-security)
    - [5.5. Network Security](#55-network-security)
    - [5.6. Data Security](#56-data-security)
    - [5.7. System Hardening](#57-system-hardening)
  - [6. Backup and Recovery](#6-backup-and-recovery)
    - [6.1. Database Backup](#61-database-backup)
      - [1.2.2. Docker Installation](#122-docker-installation-1)
      - [1.2.3. Kubernetes Installation](#123-kubernetes-installation-1)
    - [1.3. First-Time Setup](#13-first-time-setup-1)
  - [2. System Configuration](#2-system-configuration-1)
    - [2.1. Configuration Files](#21-configuration-files-1)
    - [2.2. Environment Variables](#22-environment-variables-1)
    - [2.3. Command-Line Options](#23-command-line-options-1)
  - [3. Service Management](#3-service-management-1)
    - [3.1. Running as a Systemd Service](#31-running-as-a-systemd-service-1)
    - [3.2. Running in Docker](#32-running-in-docker-1)
    - [3.3. Running in Kubernetes](#33-running-in-kubernetes-1)
    - [3.4. Service Dependencies](#34-service-dependencies-1)
    - [3.5. Health Checks](#35-health-checks-1)
  - [4. User Management](#4-user-management-1)
    - [4.1. Authentication Methods](#41-authentication-methods-1)
    - [4.2. User Roles and Permissions](#42-user-roles-and-permissions-1)
    - [4.3. Managing Users](#43-managing-users-1)
    - [4.4. Password Policies](#44-password-policies-1)
    - [4.5. External Authentication](#45-external-authentication-1)
  - [5. Security Management](#5-security-management-1)
    - [5.1. PKI Infrastructure](#51-pki-infrastructure-1)
    - [5.2. Certificate Management](#52-certificate-management-1)
    - [5.3. TLS Configuration](#53-tls-configuration-1)
    - [5.4. API Security](#54-api-security-1)
    - [5.5. Network Security](#55-network-security-1)
    - [5.6. Data Security](#56-data-security-1)
    - [5.7. System Hardening](#57-system-hardening-1)
  - [6. Backup and Recovery](#6-backup-and-recovery-1)
    - [6.1. Backup Strategy](#61-backup-strategy)
      - [6.1.1. Database Backup](#611-database-backup)
      - [6.1.2. Configuration Files Backup](#612-configuration-files-backup)
      - [6.1.3. PKI Infrastructure Backup](#613-pki-infrastructure-backup)
      - [6.1.4. Template Files Backup](#614-template-files-backup)
    - [6.2. Automated Backup Configuration](#62-automated-backup-configuration)
    - [6.3. Backup Verification](#63-backup-verification)
    - [6.4. Recovery Procedures](#64-recovery-procedures)
      - [6.4.1. Full System Recovery](#641-full-system-recovery)
      - [6.4.2. Database-Only Recovery](#642-database-only-recovery)
      - [6.4.3. Certificate Recovery](#643-certificate-recovery)
    - [6.5. Disaster Recovery Planning](#65-disaster-recovery-planning)
  - [7. Monitoring and Logging](#7-monitoring-and-logging)
    - [7.1. System Logging](#71-system-logging)
      - [7.1.1. Log File Locations](#711-log-file-locations)
      - [7.1.2. Log Configuration](#712-log-configuration)
      - [7.1.3. Log Rotation](#713-log-rotation)
    - [7.2. Monitoring System Health](#72-monitoring-system-health)
      - [7.2.1. Built-in Health Endpoints](#721-built-in-health-endpoints)
      - [7.2.2. Monitoring with Prometheus](#722-monitoring-with-prometheus)
      - [7.2.3. Grafana Dashboards](#723-grafana-dashboards)
    - [7.3. Alerting](#73-alerting)
      - [7.3.1. Configuring Alerts with Prometheus AlertManager](#731-configuring-alerts-with-prometheus-alertmanager)
      - [7.3.2. Email Notifications](#732-email-notifications)
      - [7.3.3. Webhook Notifications](#733-webhook-notifications)
    - [7.4. Log Analysis](#74-log-analysis)
      - [7.4.1. Common Log Patterns](#741-common-log-patterns)
      - [7.4.2. Using jq for Log Analysis](#742-using-jq-for-log-analysis)
      - [7.4.3. Centralized Logging](#743-centralized-logging)
    - [7.5. Audit Logging](#75-audit-logging)
      - [7.5.1. Audit Log Location](#751-audit-log-location)
      - [7.5.2. Audit Log Format](#752-audit-log-format)
      - [7.5.3. Audit Log Retention](#753-audit-log-retention)
      - [7.5.4. Audit Log Analysis](#754-audit-log-analysis)
  - [8. Performance Tuning](#8-performance-tuning)
    - [8.1. Resource Requirements](#81-resource-requirements)
    - [8.2. System Optimization](#82-system-optimization)
      - [8.2.1. Process Limits](#821-process-limits)
      - [8.2.2. Memory Management](#822-memory-management)
      - [8.2.3. Process Priority](#823-process-priority)
    - [8.3. Database Optimization](#83-database-optimization)
      - [8.3.1. SQLite Optimization](#831-sqlite-optimization)
      - [8.3.2. CockroachDB Optimization](#832-cockroachdb-optimization)
    - [8.4. Web Server Tuning](#84-web-server-tuning)
      - [8.4.1. HTTP Optimization](#841-http-optimization)
      - [8.4.2. Connection Limits](#842-connection-limits)
      - [8.4.3. Static File Serving](#843-static-file-serving)
      - [8.4.4. Load Balancing](#844-load-balancing)
    - [8.5. Filesystem Optimization](#85-filesystem-optimization)
      - [8.5.1. Filesystem Selection](#851-filesystem-selection)
      - [8.5.2. Mount Options](#852-mount-options)
      - [8.5.3. I/O Scheduling](#853-io-scheduling)
    - [8.6. Network Optimization](#86-network-optimization)
      - [8.6.1. TCP Tuning](#861-tcp-tuning)
      - [8.6.2. Network Interface Configuration](#862-network-interface-configuration)
    - [8.7. Caching Strategies](#87-caching-strategies)
      - [8.7.1. Application-Level Caching](#871-application-level-caching)
      - [8.7.2. External Caching](#872-external-caching)
      - [8.7.3. Content Delivery](#873-content-delivery)
    - [8.8. Monitoring Performance](#88-monitoring-performance)
      - [8.8.1. Key Performance Indicators](#881-key-performance-indicators)
      - [8.8.2. Performance Testing](#882-performance-testing)
    - [8.9. Scaling Strategies](#89-scaling-strategies)
      - [8.9.1. Vertical Scaling](#891-vertical-scaling)
      - [8.9.2. Horizontal Scaling](#892-horizontal-scaling)
      - [8.9.3. Service Decomposition](#893-service-decomposition)
    - [8.10. Hardware Recommendations](#810-hardware-recommendations)
      - [8.10.1. CPU](#8101-cpu)
      - [8.10.2. Memory](#8102-memory)
      - [8.10.3. Storage](#8103-storage)
      - [8.10.4. Network](#8104-network)
  - [9. Security Management](#9-security-management)
    - [9.1. Authentication and Authorization](#91-authentication-and-authorization)
      - [9.1.1. Authentication Methods](#911-authentication-methods)
      - [9.1.2. Role-Based Access Control (RBAC)](#912-role-based-access-control-rbac)
      - [9.1.3. Managing API Tokens](#913-managing-api-tokens)
    - [9.2. TLS Configuration](#92-tls-configuration)
      - [9.2.1. Web Interface TLS](#921-web-interface-tls)
      - [9.2.2. Certificate Management](#922-certificate-management)
      - [9.2.3. Mutual TLS (mTLS)](#923-mutual-tls-mtls)
    - [9.3. Network Security](#93-network-security)
      - [9.3.1. Firewall Configuration](#931-firewall-configuration)
      - [9.3.2. Network Isolation](#932-network-isolation)
      - [9.3.3. Traffic Encryption](#933-traffic-encryption)
    - [9.4. Secure Installation](#94-secure-installation)
      - [9.4.1. Secure Boot](#941-secure-boot)
      - [9.4.2. Installation Authentication](#942-installation-authentication)
      - [9.4.3. Installation Encryption](#943-installation-encryption)
    - [9.5. Vulnerability Management](#95-vulnerability-management)
      - [9.5.1. Security Updates](#951-security-updates)
      - [9.5.2. Regular Security Audits](#952-regular-security-audits)
      - [9.5.3. Security Hardening](#953-security-hardening)
    - [9.6. Data Protection](#96-data-protection)
      - [9.6.1. Sensitive Data Handling](#961-sensitive-data-handling)
      - [9.6.2. Data Retention](#962-data-retention)
      - [9.6.3. Secrets Management](#963-secrets-management)
  - [10. Troubleshooting](#10-troubleshooting)
    - [10.1. Diagnostic Tools](#101-diagnostic-tools)
      - [10.1.1. System Status Check](#1011-system-status-check)
      - [10.1.2. Log Analysis](#1012-log-analysis)
      - [10.1.3. Database Inspection](#1013-database-inspection)
      - [10.1.4. Network Diagnostics](#1014-network-diagnostics)
    - [10.2. Common Issues and Solutions](#102-common-issues-and-solutions)
      - [10.2.1. Installation Failures](#1021-installation-failures)
      - [10.2.2. Web Interface Issues](#1022-web-interface-issues)
      - [10.2.3. Database Connectivity Issues](#1023-database-connectivity-issues)
      - [10.2.4. Certificate Issues](#1024-certificate-issues)
    - [10.3. Service Specific Issues](#103-service-specific-issues)
      - [10.3.1. File Editor Service](#1031-file-editor-service)
      - [10.3.2. DNSMasq Watcher](#1032-dnsmasq-watcher)
      - [10.3.3. Certificate Issuer](#1033-certificate-issuer)
    - [10.4. Boot and Installation Debugging](#104-boot-and-installation-debugging)
      - [10.4.1. iPXE Debugging](#1041-ipxe-debugging)
      - [10.4.2. Cloud-Init Debugging](#1042-cloud-init-debugging)
      - [10.4.3. Live Installation Monitoring](#1043-live-installation-monitoring)
    - [10.5. Advanced Troubleshooting](#105-advanced-troubleshooting)
      - [10.5.1. Service Profiling](#1051-service-profiling)
      - [10.5.2. Database Query Analysis](#1052-database-query-analysis)
      - [10.5.3. Traffic Analysis](#1053-traffic-analysis)
      - [10.5.4. System Recovery](#1054-system-recovery)
    - [10.6. Getting Support](#106-getting-support)
      - [10.6.1. Generating Support Bundle](#1061-generating-support-bundle)
      - [10.6.2. Community Support](#1062-community-support)
      - [10.6.3. Commercial Support](#1063-commercial-support)
  - [11. Upgrading and Maintenance](#11-upgrading-and-maintenance)
    - [11.1. Version Upgrades](#111-version-upgrades)
      - [11.1.1. Before Upgrading](#1111-before-upgrading)
      - [11.1.2. Upgrade Procedures](#1112-upgrade-procedures)
      - [11.1.3. Post-Upgrade Tasks](#1113-post-upgrade-tasks)
    - [11.2. Database Maintenance](#112-database-maintenance)
      - [11.2.1. Database Migrations](#1121-database-migrations)
      - [11.2.2. Database Optimization](#1122-database-optimization)
      - [11.2.3. Data Cleanup](#1123-data-cleanup)
    - [11.3. Routine Maintenance Tasks](#113-routine-maintenance-tasks)
      - [11.3.1. Certificate Rotation](#1131-certificate-rotation)
      - [11.3.2. Log Rotation](#1132-log-rotation)
      - [11.3.3. File System Maintenance](#1133-file-system-maintenance)
      - [11.3.4. User Account Maintenance](#1134-user-account-maintenance)
    - [11.4. Configuration Management](#114-configuration-management)
      - [11.4.1. Configuration Backups](#1141-configuration-backups)
      - [11.4.2. Configuration Versioning](#1142-configuration-versioning)
      - [11.4.3. Configuration Validation](#1143-configuration-validation)
    - [11.5. Disaster Recovery Testing](#115-disaster-recovery-testing)
      - [11.5.1. Recovery Drills](#1151-recovery-drills)
      - [11.5.2. Failover Testing](#1152-failover-testing)
    - [11.6. System Monitoring](#116-system-monitoring)
      - [11.6.1. Monitoring Health Checks](#1161-monitoring-health-checks)
      - [11.6.2. Performance Baseline](#1162-performance-baseline)
    - [11.7. Planning for Major Upgrades](#117-planning-for-major-upgrades)
      - [11.7.1. Upgrade Impact Assessment](#1171-upgrade-impact-assessment)
      - [11.7.2. Rollback Planning](#1172-rollback-planning)
  - [12. Advanced Configuration](#12-advanced-configuration)
    - [12.1. Customizing Templates](#121-customizing-templates)
      - [12.1.1. Template Engine Overview](#1211-template-engine-overview)
      - [12.1.2. Creating Custom Template Functions](#1212-creating-custom-template-functions)
      - [12.1.3. Advanced Template Functions](#1213-advanced-template-functions)
      - [12.1.4. Template Inheritance](#1214-template-inheritance)
    - [12.2. API Customization](#122-api-customization)
      - [12.2.1. Custom API Endpoints](#1221-custom-api-endpoints)
      - [12.2.2. API Rate Limiting](#1222-api-rate-limiting)
      - [12.2.3. Custom API Authentication](#1223-custom-api-authentication)
    - [12.3. Advanced Networking](#123-advanced-networking)
      - [12.3.1. VLAN and Network Segmentation](#1231-vlan-and-network-segmentation)
      - [12.3.2. Network Bonding](#1232-network-bonding)
      - [12.3.3. IPv6 Configuration](#1233-ipv6-configuration)
    - [12.4. Clustering and High Availability](#124-clustering-and-high-availability)
      - [12.4.1. Basic Clustering Setup](#1241-basic-clustering-setup)
      - [12.4.2. Distributed File System](#1242-distributed-file-system)
      - [12.4.3. Load Balancer Configuration](#1243-load-balancer-configuration)
    - [12.5. Integration with External Systems](#125-integration-with-external-systems)
      - [12.5.1. CMDB Integration](#1251-cmdb-integration)
      - [12.5.2. Monitoring System Integration](#1252-monitoring-system-integration)
      - [12.5.3. Custom Webhook Notifications](#1253-custom-webhook-notifications)
    - [12.6. Advanced Storage Configuration](#126-advanced-storage-configuration)
      - [12.6.1. Object Storage for Installation Files](#1261-object-storage-for-installation-files)
      - [12.6.2. Database on External Storage](#1262-database-on-external-storage)
      - [12.6.3. Backup to Remote Storage](#1263-backup-to-remote-storage)
    - [12.7. Custom Authentication and Authorization](#127-custom-authentication-and-authorization)
      - [12.7.1. Custom LDAP Configuration](#1271-custom-ldap-configuration)
      - [12.7.2. OpenID Connect Configuration](#1272-openid-connect-configuration)
      - [12.7.3. Custom Authorization Rules](#1273-custom-authorization-rules)
    - [12.8. Advanced Logging](#128-advanced-logging)
      - [12.8.1. Structured Logging Configuration](#1281-structured-logging-configuration)
      - [12.8.2. Remote Logging Configuration](#1282-remote-logging-configuration)
      - [12.8.3. Log Correlation](#1283-log-correlation)
  - [13. Appendices](#13-appendices)
    - [13.1. Configuration Reference](#131-configuration-reference)
      - [13.1.1. Complete Configuration Schema](#1311-complete-configuration-schema)
      - [12.1.3. Advanced Template Functions](#1213-advanced-template-functions-1)
      - [12.1.4. Template Inheritance](#1214-template-inheritance-1)
    - [12.2. API Customization](#122-api-customization-1)
      - [12.2.1. Custom API Endpoints](#1221-custom-api-endpoints-1)
      - [12.2.2. API Rate Limiting](#1222-api-rate-limiting-1)
      - [12.2.3. Custom API Authentication](#1223-custom-api-authentication-1)
    - [12.3. Advanced Networking](#123-advanced-networking-1)
      - [12.3.1. VLAN and Network Segmentation](#1231-vlan-and-network-segmentation-1)
      - [12.3.2. Network Bonding](#1232-network-bonding-1)
      - [12.3.3. IPv6 Configuration](#1233-ipv6-configuration-1)
    - [12.4. Clustering and High Availability](#124-clustering-and-high-availability-1)
      - [12.4.1. Basic Clustering Setup](#1241-basic-clustering-setup-1)
      - [12.4.2. Distributed File System](#1242-distributed-file-system-1)
      - [12.4.3. Load Balancer Configuration](#1243-load-balancer-configuration-1)
    - [12.5. Integration with External Systems](#125-integration-with-external-systems-1)
      - [12.5.1. CMDB Integration](#1251-cmdb-integration-1)
      - [12.5.2. Monitoring System Integration](#1252-monitoring-system-integration-1)
      - [12.5.3. Custom Webhook Notifications](#1253-custom-webhook-notifications-1)
    - [12.6. Advanced Storage Configuration](#126-advanced-storage-configuration-1)
      - [12.6.1. Object Storage for Installation Files](#1261-object-storage-for-installation-files-1)
      - [12.6.2. Database on External Storage](#1262-database-on-external-storage-1)
      - [12.6.3. Backup to Remote Storage](#1263-backup-to-remote-storage-1)
    - [12.7. Custom Authentication and Authorization](#127-custom-authentication-and-authorization-1)
      - [12.7.1. Custom LDAP Configuration](#1271-custom-ldap-configuration-1)
      - [12.7.2. OpenID Connect Configuration](#1272-openid-connect-configuration-1)
      - [12.7.3. Custom Authorization Rules](#1273-custom-authorization-rules-1)
    - [12.8. Advanced Logging](#128-advanced-logging-1)
      - [12.8.1. Structured Logging Configuration](#1281-structured-logging-configuration-1)
      - [12.8.2. Remote Logging Configuration](#1282-remote-logging-configuration-1)
      - [12.8.3. Log Correlation](#1283-log-correlation-1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Ubuntu Autoinstall Webhook Admin Guide

## 1. Installation

### 1.1. System Requirements

**Minimum Requirements**:
- **CPU**: 2 cores
- **RAM**: 4 GB
- **Storage**: 20 GB
- **Network**: 1 Gbps Ethernet interface
- **Operating System**: Ubuntu 20.04 LTS or newer

**Recommended Requirements**:
- **CPU**: 4+ cores
- **RAM**: 8+ GB
- **Storage**: 50+ GB SSD
- **Network**: 10 Gbps Ethernet interface
- **Operating System**: Ubuntu 22.04 LTS

**Additional Requirements**:
- DHCP server with PXE boot support (can be the same machine)
- Network support for PXE booting
- Permissions to manage network infrastructure

### 1.2. Installation Methods

#### 1.2.1. Binary Installation

```bash
# Download the latest release
curl -L -o ubuntu-autoinstall-webhook.tar.gz https://github.com/jdfalk/ubuntu-autoinstall-webhook/releases/latest/download/ubuntu-autoinstall-webhook.tar.gz

# Extract the archive
tar -xzf ubuntu-autoinstall-webhook.tar.gz

# Move to /usr/local/bin
sudo mv ubuntu-autoinstall-webhook /usr/local/bin/

# Set permissions
sudo chmod +x /usr/local/bin/ubuntu-autoinstall-webhook

# Create configuration directory
sudo mkdir -p /etc/ubuntu-autoinstall-webhook

# Copy example configuration
sudo cp config.example.yaml /etc/ubuntu-autoinstall-webhook/config.yaml

# Create data directory
sudo mkdir -p /var/lib/ubuntu-autoinstall-webhook

# Set ownership
sudo chown -R nobody:nogroup /var/lib/ubuntu-autoinstall-webhook
```

#### 1.2.2. Docker Installation

```bash
# Pull the Docker image
docker pull ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest

# Create configuration directory
mkdir -p ./config

# Create data directory
mkdir -p ./data

# Create a docker-compose.yml file
cat > docker-compose.yml << 'EOF'
version: '3'
services:
  webhook:
    image: ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest
    ports:
      - "8443:8443"  # HTTPS
      - "8080:8080"  # HTTP
      - "69:69/udp"  # TFTP
    volumes:
      - ./config:/etc/ubuntu-autoinstall-webhook
      - ./data:/var/lib/ubuntu-autoinstall-webhook
    restart: unless-stopped
    environment:
      - UAW_LOG_LEVEL=info
EOF

# Start the service
docker-compose up -d
```

#### 1.2.3. Kubernetes Installation

```bash
# Create namespace
kubectl create namespace ubuntu-autoinstall-webhook

# Apply the deployment configurations
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/deployment.yaml
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/service.yaml
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/configmap.yaml
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/persistent-volume-claim.yaml
```

### 1.3. First-Time Setup

After installing the system, complete these initial configuration steps:

1. **Create the initial admin user**:
```bash
ubuntu-autoinstall-webhook create-admin --username admin --email admin@example.com
```

2. **Generate the root certificate**:
```bash
ubuntu-autoinstall-webhook cert-init --ca-name "Ubuntu Autoinstall Webhook CA"
```

3. **Verify the installation**:
```bash
ubuntu-autoinstall-webhook health-check
```

4. **Start the service**:
```bash
sudo systemctl start ubuntu-autoinstall-webhook
```

5. **Enable the service to start at boot**:
```bash
sudo systemctl enable ubuntu-autoinstall-webhook
```

## 2. System Configuration

### 2.1. Configuration Files

The system uses YAML configuration files located in the following directories:

- **Main configuration**: `/etc/ubuntu-autoinstall-webhook/config.yaml`
- **Database configuration**: `/etc/ubuntu-autoinstall-webhook/database.yaml`
- **Certificate configuration**: `/etc/ubuntu-autoinstall-webhook/certificates.yaml`
- **Template directory**: `/etc/ubuntu-autoinstall-webhook/templates/`

**Example main configuration file** (`config.yaml`):

```yaml
# Server configuration
server:
  host: 0.0.0.0
  port: 8443
  tls:
    enabled: true
    cert_file: /etc/ubuntu-autoinstall-webhook/certs/server.crt
    key_file: /etc/ubuntu-autoinstall-webhook/certs/server.key

# Database configuration
database:
  type: sqlite
  path: /var/lib/ubuntu-autoinstall-webhook/database.db

# File system configuration
filesystem:
  base_path: /var/lib/ubuntu-autoinstall-webhook
  ipxe_path: /var/lib/ubuntu-autoinstall-webhook/ipxe
  cloudinit_path: /var/lib/ubuntu-autoinstall-webhook/cloud-init

# Certificate authority configuration
certificates:
  ca_cert: /etc/ubuntu-autoinstall-webhook/certs/ca.crt
  ca_key: /etc/ubuntu-autoinstall-webhook/certs/ca.key
  validity_period_days: 365

# Logging configuration
logging:
  level: info
  format: json
  output: file
  file_path: /var/log/ubuntu-autoinstall-webhook/server.log

# DNSMasq watcher configuration
dnsmasq_watcher:
  enabled: true
  log_file: /var/log/dnsmasq.log
  poll_interval_seconds: 5

# Authentication configuration
authentication:
  method: local
  session_timeout_minutes: 60
  max_failed_attempts: 5
  lockout_time_minutes: 30
```

### 2.2. Environment Variables

The following environment variables can be used to override configuration values:

| Variable            | Description                | Example                          |
| ------------------- | -------------------------- | -------------------------------- |
| `UAW_SERVER_HOST`   | Server bind address        | `0.0.0.0`                        |
| `UAW_SERVER_PORT`   | Server port                | `8443`                           |
| `UAW_TLS_ENABLED`   | Enable TLS                 | `true`                           |
| `UAW_TLS_CERT_FILE` | TLS certificate path       | `/path/to/cert.crt`              |
| `UAW_TLS_KEY_FILE`  | TLS key path               | `/path/to/key.key`               |
| `UAW_DB_TYPE`       | Database type              | `sqlite`, `cockroachdb`          |
| `UAW_DB_PATH`       | SQLite database path       | `/path/to/db.db`                 |
| `UAW_DB_CONNECTION` | Database connection string | `host=localhost port=26257`      |
| `UAW_LOG_LEVEL`     | Log level                  | `debug`, `info`, `warn`, `error` |
| `UAW_LOG_FORMAT`    | Log format                 | `json`, `text`                   |
| `UAW_DNSMASQ_LOG`   | DNSMasq log path           | `/var/log/dnsmasq.log`           |

### 2.3. Command-Line Options

The application supports the following command-line options:

```
Usage: ubuntu-autoinstall-webhook [command] [options]

Commands:
  server               Start the webhook server
  cert-init            Initialize the certificate authority
  create-admin         Create an admin user
  import-template      Import installation templates
  backup               Backup system data
  restore              Restore system data
  health-check         Perform system health check
  migrate              Run database migrations
  version              Display version information
  help                 Display help information

Server options:
  --host               Server bind address (default: 0.0.0.0)
  --port               Server port (default: 8443)
  --config             Path to config file (default: /etc/ubuntu-autoinstall-webhook/config.yaml)
  --log-level          Log level: debug, info, warn, error (default: info)
  --log-format         Log format: text, json (default: text)
```

## 3. Service Management

### 3.1. Running as a Systemd Service

The system can be configured as a systemd service for automatic startup and management.

**Create systemd service file**:

```bash
sudo nano /etc/systemd/system/ubuntu-autoinstall-webhook.service
```

**Add the following content**:

```ini
[Unit]
Description=Ubuntu Autoinstall Webhook Server
After=network.target

[Service]
ExecStart=/usr/local/bin/ubuntu-autoinstall-webhook server --config /etc/ubuntu-autoinstall-webhook/config.yaml
Restart=always
RestartSec=5
User=nobody
Group=nogroup
Environment=UAW_LOG_LEVEL=info

[Install]
WantedBy=multi-user.target
```

**Enable and start the service**:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ubuntu-autoinstall-webhook
sudo systemctl start ubuntu-autoinstall-webhook
```

**Checking service status**:

```bash
sudo systemctl status ubuntu-autoinstall-webhook
```

**Viewing service logs**:

```bash
sudo journalctl -u ubuntu-autoinstall-webhook -f
```

### 3.2. Running in Docker

Running the service in Docker provides isolation and simplified deployment.

**Creating a persistent Docker setup**:

1. Create a directory structure:
```bash
mkdir -p ubuntu-autoinstall-webhook/{config,data,logs}
```

2. Create a configuration file:
```bash
nano ubuntu-autoinstall-webhook/config/config.yaml
```

3. Run with docker-compose:
```bash
docker-compose up -d
```

**Monitoring the Docker container**:

```bash
# View container logs
docker logs -f ubuntu-autoinstall-webhook

# Check container status
docker ps -a | grep ubuntu-autoinstall-webhook

# Restart the container
docker restart ubuntu-autoinstall-webhook
```

### 3.3. Running in Kubernetes

For larger deployments, Kubernetes provides scalability and high availability.

**Monitoring the Kubernetes deployment**:

```bash
# Get pod status
kubectl get pods -n ubuntu-autoinstall-webhook

# View pod logs
kubectl logs -f deployment/ubuntu-autoinstall-webhook -n ubuntu-autoinstall-webhook

# Describe the deployment
kubectl describe deployment ubuntu-autoinstall-webhook -n ubuntu-autoinstall-webhook
```

**Scaling the deployment**:

```bash
kubectl scale deployment ubuntu-autoinstall-webhook --replicas=3 -n ubuntu-autoinstall-webhook
```

**Updating the deployment**:

```bash
kubectl set image deployment/ubuntu-autoinstall-webhook ubuntu-autoinstall-webhook=ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest -n ubuntu-autoinstall-webhook
```

### 3.4. Service Dependencies

The Ubuntu Autoinstall Webhook system depends on the following external services:

1. **DHCP Server** (typically dnsmasq)
   - Required for handling PXE boot requests
   - Configuration must include PXE boot options

2. **TFTP Server** (built-in or external)
   - Required for serving boot files
   - Must be accessible to client systems

3. **Database** (SQLite by default, CockroachDB optional)
   - Stores system configuration and status
   - No external setup needed for SQLite

Example dnsmasq configuration for PXE booting:

```
# /etc/dnsmasq.conf
interface=eth0
domain=example.local
dhcp-range=192.168.1.50,192.168.1.150,12h
dhcp-boot=pxelinux.0
enable-tftp
tftp-root=/var/lib/tftpboot
```

### 3.5. Health Checks

The system provides health check endpoints for monitoring:

1. **HTTP Health Check**:
   - URL: `/api/v1/health`
   - Method: GET
   - Response: 200 OK if healthy

2. **Component Health Checks**:
   - URL: `/api/v1/health/components`
   - Method: GET
   - Response: Status of each component

Example health check script:

```bash
#!/bin/bash
# Health check script for automation tools

HEALTH_URL="https://localhost:8443/api/v1/health"

response=$(curl -sk $HEALTH_URL)
status=$?

if [ $status -ne 0 ]; then
  echo "ERROR: Could not connect to health endpoint"
  exit 1
fi

if [[ "$response" == *"healthy\":true"* ]]; then
  echo "System is healthy"
  exit 0
else
  echo "System is not healthy: $response"
  exit 1
fi
```

## 4. User Management

### 4.1. Authentication Methods

The system supports several authentication methods:

1. **Local Authentication**
   - Username and password stored in the system database
   - Passwords are securely hashed using bcrypt

2. **LDAP Authentication**
   - Integration with corporate LDAP directories
   - Configurable user attribute mapping

3. **OAuth2 Authentication**
   - Support for external identity providers
   - Configurable for providers like GitHub, Google, Azure AD

Example LDAP configuration:

```yaml
authentication:
  method: ldap
  ldap:
    url: "ldap://ldap.example.com:389"
    base_dn: "dc=example,dc=com"
    user_dn_pattern: "cn={0},ou=users,dc=example,dc=com"
    user_search_base: "ou=users"
    user_search_filter: "(uid={0})"
    group_search_base: "ou=groups"
    group_search_filter: "(member={0})"
    manager_dn: "cn=admin,dc=example,dc=com"
    manager_password: "password"
    user_attribute_mappings:
      username: "uid"
      email: "mail"
      display_name: "displayName"
```

### 4.2. User Roles and Permissions

The system implements Role-Based Access Control (RBAC) with these default roles:

1. **Administrator**
   - Full system access and control
   - Can manage users and roles
   - Can configure all system settings

2. **Operator**
   - Can manage systems and installations
   - Can view and modify templates
   - Cannot manage users or change system settings

3. **Viewer**
   - Can view systems and installation status
   - Cannot modify configurations or initiate installations

Custom roles can be created with specific permission sets:

```yaml
roles:
  - name: "DevOps"
    description: "Role for DevOps team members"
    permissions:
      - "systems:read"
      - "systems:create"
      - "systems:update"
      - "systems:install"
      - "templates:read"
      - "templates:use"
      - "installations:read"
      - "logs:read"
```

### 4.3. Managing Users

**Creating users with the CLI**:

```bash
# Create a new admin user
ubuntu-autoinstall-webhook create-admin --username admin --email admin@example.com

# Create a standard user with a specific role
ubuntu-autoinstall-webhook create-user --username operator --email operator@example.com --role operator
```

**Managing users via the Web UI**:

1. Navigate to "Settings" > "Users" in the web interface
2. Click "Add User" to create a new user
3. Fill in the required information:
   - Username
   - Email
   - Password
   - Role
4. Click "Save" to create the user

**User management API endpoints**:

```bash
# Create a new user
curl -X POST "https://localhost:8443/api/v1/users" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "SecurePassword123!",
    "role": "operator"
  }'

# Update a user
curl -X PUT "https://localhost:8443/api/v1/users/newuser" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "updated@example.com",
    "role": "viewer"
  }'

# Delete a user
curl -X DELETE "https://localhost:8443/api/v1/users/newuser" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### 4.4. Password Policies

The system supports configurable password policies:

```yaml
authentication:
  password_policy:
    min_length: 12
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_special_chars: true
    max_age_days: 90
    history_count: 5
    max_failed_attempts: 5
    lockout_time_minutes: 30
```

**Resetting user passwords**:

```bash
# Reset a user's password via CLI
ubuntu-autoinstall-webhook reset-password --username operator --generate
```

### 4.5. External Authentication

**OAuth2 Configuration Example**:

```yaml
authentication:
  method: oauth2
  oauth2:
    providers:
      github:
        enabled: true
        client_id: "your-client-id"
        client_secret: "your-client-secret"
        auth_url: "https://github.com/login/oauth/authorize"
        token_url: "https://github.com/login/oauth/access_token"
        user_info_url: "https://api.github.com/user"
        scope: "read:user,user:email"
        user_attribute_mappings:
          username: "login"
          email: "email"
          display_name: "name"
      google:
        enabled: false
        # Configuration for Google OAuth2
```

**SAML Configuration Example**:

```yaml
authentication:
  method: saml
  saml:
    idp_metadata_url: "https://idp.example.com/metadata"
    sp_entity_id: "https://ubuntu-autoinstall-webhook.example.com"
    assertion_consumer_service_url: "https://ubuntu-autoinstall-webhook.example.com/api/v1/auth/saml/callback"
    user_attribute_mappings:
      username: "NameID"
      email: "Email"
      display_name: "DisplayName"
```

## 5. Security Management

### 5.1. PKI Infrastructure

The system maintains its own Public Key Infrastructure (PKI) to secure communications:

1. **Root CA**: Self-signed certificate authority
2. **Intermediate CA**: Optional for larger deployments
3. **Service Certificates**: For internal services
4. **Client Certificates**: For mTLS authentication

**Initializing the PKI**:

```bash
# Initialize the root CA
ubuntu-autoinstall-webhook cert-init --ca-name "Ubuntu Autoinstall Webhook CA"

# Create an intermediate CA (optional)
ubuntu-autoinstall-webhook cert-create-intermediate --name "Intermediate CA"

# Issue service certificate
ubuntu-autoinstall-webhook cert-issue --name "webhook-server" --type server --dns webhook.example.com --ip 192.168.1.10
```

### 5.2. Certificate Management

**Certificate renewal**:

```bash
# Check for certificates nearing expiration
ubuntu-autoinstall-webhook cert-check-expiry

# Renew a specific certificate
ubuntu-autoinstall-webhook cert-renew --name "webhook-server"

# Renew all certificates expiring within 30 days
ubuntu-autoinstall-webhook cert-renew-all --days 30
```

**Certificate revocation**:

```bash
# Revoke a certificate
ubuntu-autoinstall-webhook cert-revoke --name "webhook-server" --reason "key-compromise"

# Generate CRL
ubuntu-autoinstall-webhook cert-gen-crl
```

### 5.3. TLS Configuration

The system uses secure TLS configurations by default. The following settings can be customized in `config.yaml`:

```yaml
server:
  tls:
    enabled: true
    cert_file: /etc/ubuntu-autoinstall-webhook/certs/server.crt
    key_file: /etc/ubuntu-autoinstall-webhook/certs/server.key
    min_version: "TLS1.2"
    cipher_suites:
      - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
      - TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
    prefer_server_cipher_suites: true
    client_auth: "request"  # options: "no", "request", "require"
    client_ca_file: /etc/ubuntu-autoinstall-webhook/certs/ca.crt
```

**Testing TLS Configuration**:

```bash
# Test server TLS configuration
openssl s_client -connect localhost:8443 -tls1_2

# Verify certificate chain
openssl verify -CAfile /etc/ubuntu-autoinstall-webhook/certs/ca.crt /etc/ubuntu-autoinstall-webhook/certs/server.crt
```

### 5.4. API Security

API access is secured through several mechanisms:

1. **API Keys**: Long-lived tokens for service integration
2. **JWT Tokens**: Short-lived tokens for user sessions
3. **Rate Limiting**: Prevents abuse of API endpoints
4. **IP Restrictions**: Optional IP-based access control

**API Key Management**:

```bash
# Generate a new API key
ubuntu-autoinstall-webhook generate-api-key --name "integration-service" --role operator

# List API keys
ubuntu-autoinstall-webhook list-api-keys

# Revoke an API key
ubuntu-autoinstall-webhook revoke-api-key --name "integration-service"
```

**API Security Configuration**:

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    burst: 10
  ip_restrictions:
    enabled: false
    allowed_ips:
      - "192.168.1.0/24"
      - "10.0.0.5"
  token_expiry:
    access_token_minutes: 15
    refresh_token_days: 7
```

### 5.5. Network Security

**Firewall Recommendations**:

The following ports should be open for system operation:

| Port | Protocol | Purpose                          |
| ---- | -------- | -------------------------------- |
| 8443 | TCP      | HTTPS API and Web UI             |
| 8080 | TCP      | HTTP (redirect to HTTPS)         |
| 69   | UDP      | TFTP Server (if using built-in)  |
| 67   | UDP      | DHCP Server (if running locally) |

**Example UFW configuration**:

```bash
sudo ufw allow 8443/tcp comment "Ubuntu Autoinstall Webhook HTTPS"
sudo ufw allow 8080/tcp comment "Ubuntu Autoinstall Webhook HTTP redirect"
sudo ufw allow 69/udp comment "TFTP Server"
```

**Securing network services**:

1. Use separate network segments for PXE booting when possible
2. Implement VLAN isolation for installation networks
3. Configure network ACLs to restrict access to the webhook server

### 5.6. Data Security

**Sensitive Data Handling**:

The system handles the following types of sensitive data:

1. User credentials
2. API keys
3. Private keys
4. Certificate authority keys

**Data Protection Measures**:

1. **Encryption at Rest**:
   - Database encryption options:
   ```yaml
   database:
     encryption:
       enabled: true
       key_file: /etc/ubuntu-autoinstall-webhook/keys/db-encryption-key
   ```

2. **Secure Storage of Secrets**:
   - Private keys stored with restricted permissions
   - Configuration option for external secret storage:
   ```yaml
   secrets:
     provider: vault
     vault:
       address: "https://vault.example.com:8200"
       token_file: "/etc/ubuntu-autoinstall-webhook/vault-token"
       path_prefix: "ubuntu-autoinstall-webhook"
   ```

### 5.7. System Hardening

**Operating System Hardening**:

1. Keep system updated with security patches
2. Use minimal installations for server environments
3. Implement AppArmor or SELinux profiles
4. Remove unnecessary packages and services

**Application Hardening**:

1. Run the service as a non-privileged user
2. Implement file permission restrictions
3. Configure proper umask settings
4. Use seccomp or capabilities for container deployments

**Example AppArmor Profile**:

```
# /etc/apparmor.d/usr.local.bin.ubuntu-autoinstall-webhook
#include <tunables/global>

/usr/local/bin/ubuntu-autoinstall-webhook {
  #include <abstractions/base>
  #include <abstractions/nameservice>
  #include <abstractions/ssl_certs>
  #include <abstractions/openssl>

  /usr/local/bin/ubuntu-autoinstall-webhook mr,
  /etc/ubuntu-autoinstall-webhook/** r,
  /etc/ubuntu-autoinstall-webhook/certs/* r,
  /var/lib/ubuntu-autoinstall-webhook/** rwk,
  /var/log/ubuntu-autoinstall-webhook/* w,
  /var/log/dnsmasq.log r,

  network tcp,
  network udp,
}
```

## 6. Backup and Recovery

### 6.1. Database Backup

**SQLite Database Backup**:

```bash
# Using the built-in backup command
ubuntu-autoinstall-webhook backup --component database --output /path/to/backup/database.db.bak

# Manual SQLite backup
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.db ".backup '/path/to/backup/database.db.bak'"
```

**CockroachDB Backup**:

```bash
# Using built-in backup command
ubuntu-autoinstall-webhook backup --component database --output /path/to/backup/cockroach-backup

# Manual CockroachDB backup
cockroach<!--
# Admin Guide for Ubuntu Autoinstall Webhook

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Ubuntu Autoinstall Webhook Admin Guide](#ubuntu-autoinstall-webhook-admin-guide)
  - [1. Installation](#1-installation)
    - [1.1. System Requirements](#11-system-requirements)
    - [1.2. Installation Methods](#12-installation-methods)
    - [1.3. First-Time Setup](#13-first-time-setup)
  - [2. System Configuration](#2-system-configuration)
    - [2.1. Configuration Files](#21-configuration-files)
    - [2.2. Environment Variables](#22-environment-variables)
    - [2.3. Command-Line Options](#23-command-line-options)
  - [3. Service Management](#3-service-management)
    - [3.1. Running as a Systemd Service](#31-running-as-a-systemd-service)
    - [3.2. Running in Docker](#32-running-in-docker)
    - [3.3. Running in Kubernetes](#33-running-in-kubernetes)
    - [3.4. Service Dependencies](#34-service-dependencies)
    - [3.5. Health Checks](#35-health-checks)
  - [4. User Management](#4-user-management)
    - [4.1. Authentication Methods](#41-authentication-methods)
    - [4.2. User Roles and Permissions](#42-user-roles-and-permissions)
    - [4.3. Managing Users](#43-managing-users)
    - [4.4. Password Policies](#44-password-policies)
    - [4.5. External Authentication](#45-external-authentication)
  - [5. Security Management](#5-security-management)
    - [5.1. PKI Infrastructure](#51-pki-infrastructure)
    - [5.2. Certificate Management](#52-certificate-management)
    - [5.3. TLS Configuration](#53-tls-configuration)
    - [5.4. API Security](#54-api-security)
    - [5.5. Network Security](#55-network-security)
    - [5.6. Data Security](#56-data-security)
    - [5.7. System Hardening](#57-system-hardening)
  - [6. Backup and Recovery](#6-backup-and-recovery)
    - [6.1. Database Backup](#61-database-backup)
    - [6.2. Certificate Backup](#62-certificate-backup)
    - [6.3. Configuration Backup](#63-configuration-backup)
    - [6.4. Automated Backup](#64-automated-backup)
    - [6.5. System Recovery](#65-system-recovery)
    - [6.6. Disaster Recovery](#66-disaster-recovery)
  - [7. Monitoring and Logging](#7-monitoring-and-logging)
    - [7.1. Log Files](#71-log-files)
    - [7.2. Log Rotation](#72-log-rotation)
    - [7.3. Audit Logging](#73-audit-logging)
    - [7.4. System Metrics](#74-system-metrics)
    - [7.5. Alerting](#75-alerting)
    - [7.6. Integration with External Monitoring](#76-integration-with-external-monitoring)
  - [8. Troubleshooting](#8-troubleshooting)
    - [8.1. Common Issues](#81-common-issues)
    - [8.2. Advanced Diagnostics](#82-advanced-diagnostics)
    - [8.3. Resource Management](#83-resource-management)
    - [8.4. Network Troubleshooting](#84-network-troubleshooting)
    - [8.5. Database Troubleshooting](#85-database-troubleshooting)
  - [9. Advanced Configuration](#9-advanced-configuration)
    - [9.1. DNSMasq Integration](#91-dnsmasq-integration)
    - [9.2. Custom Templates](#92-custom-templates)
    - [9.3. Multi-Environment Configuration](#93-multi-environment-configuration)
    - [9.4. Custom Certificate Authority](#94-custom-certificate-authority)
    - [9.5. High Availability Setup](#95-high-availability-setup)
  - [10. Upgrades and Maintenance](#10-upgrades-and-maintenance)
    - [10.1. Upgrade Process](#101-upgrade-process)
    - [10.2. Version Compatibility](#102-version-compatibility)
    - [10.3. Scheduled Maintenance](#103-scheduled-maintenance)
    - [10.4. Performance Tuning](#104-performance-tuning)
  - [11. Appendices](#11-appendices)
    - [11.1. Configuration Reference](#111-configuration-reference)
    - [11.2. Command-Line Reference](#112-command-line-reference)
    - [11.3. API Reference](#113-api-reference)
    - [11.4. Resource Requirements Sizing](#114-resource-requirements-sizing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Ubuntu Autoinstall Webhook Admin Guide

## 1. Installation

### 1.1. System Requirements

**Minimum Requirements**:
- **CPU**: 2 cores
- **RAM**: 4 GB
- **Storage**: 20 GB
- **Network**: 1 Gbps Ethernet interface
- **Operating System**: Ubuntu 20.04 LTS or newer

**Recommended Requirements**:
- **CPU**: 4+ cores
- **RAM**: 8+ GB
- **Storage**: 50+ GB SSD
- **Network**: 10 Gbps Ethernet interface
- **Operating System**: Ubuntu 22.04 LTS

**Additional Requirements**:
- DHCP server with PXE boot support (can be the same machine)
- Network support for PXE booting
- Permissions to manage network infrastructure

### 1.2. Installation Methods

#### 1.2.1. Binary Installation

```bash
# Download the latest release
curl -L -o ubuntu-autoinstall-webhook.tar.gz https://github.com/jdfalk/ubuntu-autoinstall-webhook/releases/latest/download/ubuntu-autoinstall-webhook.tar.gz

# Extract the archive
tar -xzf ubuntu-autoinstall-webhook.tar.gz

# Move to /usr/local/bin
sudo mv ubuntu-autoinstall-webhook /usr/local/bin/

# Set permissions
sudo chmod +x /usr/local/bin/ubuntu-autoinstall-webhook

# Create configuration directory
sudo mkdir -p /etc/ubuntu-autoinstall-webhook

# Copy example configuration
sudo cp config.example.yaml /etc/ubuntu-autoinstall-webhook/config.yaml

# Create data directory
sudo mkdir -p /var/lib/ubuntu-autoinstall-webhook

# Set ownership
sudo chown -R nobody:nogroup /var/lib/ubuntu-autoinstall-webhook
```

#### 1.2.2. Docker Installation

```bash
# Pull the Docker image
docker pull ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest

# Create configuration directory
mkdir -p ./config

# Create data directory
mkdir -p ./data

# Create a docker-compose.yml file
cat > docker-compose.yml << 'EOF'
version: '3'
services:
  webhook:
    image: ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest
    ports:
      - "8443:8443"  # HTTPS
      - "8080:8080"  # HTTP
      - "69:69/udp"  # TFTP
    volumes:
      - ./config:/etc/ubuntu-autoinstall-webhook
      - ./data:/var/lib/ubuntu-autoinstall-webhook
    restart: unless-stopped
    environment:
      - UAW_LOG_LEVEL=info
EOF

# Start the service
docker-compose up -d
```

#### 1.2.3. Kubernetes Installation

```bash
# Create namespace
kubectl create namespace ubuntu-autoinstall-webhook

# Apply the deployment configurations
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/deployment.yaml
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/service.yaml
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/configmap.yaml
kubectl apply -f https://github.com/jdfalk/ubuntu-autoinstall-webhook/deployments/kubernetes/persistent-volume-claim.yaml
```

### 1.3. First-Time Setup

After installing the system, complete these initial configuration steps:

1. **Create the initial admin user**:
```bash
ubuntu-autoinstall-webhook create-admin --username admin --email admin@example.com
```

2. **Generate the root certificate**:
```bash
ubuntu-autoinstall-webhook cert-init --ca-name "Ubuntu Autoinstall Webhook CA"
```

3. **Verify the installation**:
```bash
ubuntu-autoinstall-webhook health-check
```

4. **Start the service**:
```bash
sudo systemctl start ubuntu-autoinstall-webhook
```

5. **Enable the service to start at boot**:
```bash
sudo systemctl enable ubuntu-autoinstall-webhook
```

## 2. System Configuration

### 2.1. Configuration Files

The system uses YAML configuration files located in the following directories:

- **Main configuration**: `/etc/ubuntu-autoinstall-webhook/config.yaml`
- **Database configuration**: `/etc/ubuntu-autoinstall-webhook/database.yaml`
- **Certificate configuration**: `/etc/ubuntu-autoinstall-webhook/certificates.yaml`
- **Template directory**: `/etc/ubuntu-autoinstall-webhook/templates/`

**Example main configuration file** (`config.yaml`):

```yaml
# Server configuration
server:
  host: 0.0.0.0
  port: 8443
  tls:
    enabled: true
    cert_file: /etc/ubuntu-autoinstall-webhook/certs/server.crt
    key_file: /etc/ubuntu-autoinstall-webhook/certs/server.key

# Database configuration
database:
  type: sqlite
  path: /var/lib/ubuntu-autoinstall-webhook/database.db

# File system configuration
filesystem:
  base_path: /var/lib/ubuntu-autoinstall-webhook
  ipxe_path: /var/lib/ubuntu-autoinstall-webhook/ipxe
  cloudinit_path: /var/lib/ubuntu-autoinstall-webhook/cloud-init

# Certificate authority configuration
certificates:
  ca_cert: /etc/ubuntu-autoinstall-webhook/certs/ca.crt
  ca_key: /etc/ubuntu-autoinstall-webhook/certs/ca.key
  validity_period_days: 365

# Logging configuration
logging:
  level: info
  format: json
  output: file
  file_path: /var/log/ubuntu-autoinstall-webhook/server.log

# DNSMasq watcher configuration
dnsmasq_watcher:
  enabled: true
  log_file: /var/log/dnsmasq.log
  poll_interval_seconds: 5

# Authentication configuration
authentication:
  method: local
  session_timeout_minutes: 60
  max_failed_attempts: 5
  lockout_time_minutes: 30
```

### 2.2. Environment Variables

The following environment variables can be used to override configuration values:

| Variable            | Description                | Example                          |
| ------------------- | -------------------------- | -------------------------------- |
| `UAW_SERVER_HOST`   | Server bind address        | `0.0.0.0`                        |
| `UAW_SERVER_PORT`   | Server port                | `8443`                           |
| `UAW_TLS_ENABLED`   | Enable TLS                 | `true`                           |
| `UAW_TLS_CERT_FILE` | TLS certificate path       | `/path/to/cert.crt`              |
| `UAW_TLS_KEY_FILE`  | TLS key path               | `/path/to/key.key`               |
| `UAW_DB_TYPE`       | Database type              | `sqlite`, `cockroachdb`          |
| `UAW_DB_PATH`       | SQLite database path       | `/path/to/db.db`                 |
| `UAW_DB_CONNECTION` | Database connection string | `host=localhost port=26257`      |
| `UAW_LOG_LEVEL`     | Log level                  | `debug`, `info`, `warn`, `error` |
| `UAW_LOG_FORMAT`    | Log format                 | `json`, `text`                   |
| `UAW_DNSMASQ_LOG`   | DNSMasq log path           | `/var/log/dnsmasq.log`           |

### 2.3. Command-Line Options

The application supports the following command-line options:

```
Usage: ubuntu-autoinstall-webhook [command] [options]

Commands:
  server               Start the webhook server
  cert-init            Initialize the certificate authority
  create-admin         Create an admin user
  import-template      Import installation templates
  backup               Backup system data
  restore              Restore system data
  health-check         Perform system health check
  migrate              Run database migrations
  version              Display version information
  help                 Display help information

Server options:
  --host               Server bind address (default: 0.0.0.0)
  --port               Server port (default: 8443)
  --config             Path to config file (default: /etc/ubuntu-autoinstall-webhook/config.yaml)
  --log-level          Log level: debug, info, warn, error (default: info)
  --log-format         Log format: text, json (default: text)
```

## 3. Service Management

### 3.1. Running as a Systemd Service

The system can be configured as a systemd service for automatic startup and management.

**Create systemd service file**:

```bash
sudo nano /etc/systemd/system/ubuntu-autoinstall-webhook.service
```

**Add the following content**:

```ini
[Unit]
Description=Ubuntu Autoinstall Webhook Server
After=network.target

[Service]
ExecStart=/usr/local/bin/ubuntu-autoinstall-webhook server --config /etc/ubuntu-autoinstall-webhook/config.yaml
Restart=always
RestartSec=5
User=nobody
Group=nogroup
Environment=UAW_LOG_LEVEL=info

[Install]
WantedBy=multi-user.target
```

**Enable and start the service**:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ubuntu-autoinstall-webhook
sudo systemctl start ubuntu-autoinstall-webhook
```

**Checking service status**:

```bash
sudo systemctl status ubuntu-autoinstall-webhook
```

**Viewing service logs**:

```bash
sudo journalctl -u ubuntu-autoinstall-webhook -f
```

### 3.2. Running in Docker

Running the service in Docker provides isolation and simplified deployment.

**Creating a persistent Docker setup**:

1. Create a directory structure:
```bash
mkdir -p ubuntu-autoinstall-webhook/{config,data,logs}
```

2. Create a configuration file:
```bash
nano ubuntu-autoinstall-webhook/config/config.yaml
```

3. Run with docker-compose:
```bash
docker-compose up -d
```

**Monitoring the Docker container**:

```bash
# View container logs
docker logs -f ubuntu-autoinstall-webhook

# Check container status
docker ps -a | grep ubuntu-autoinstall-webhook

# Restart the container
docker restart ubuntu-autoinstall-webhook
```

### 3.3. Running in Kubernetes

For larger deployments, Kubernetes provides scalability and high availability.

**Monitoring the Kubernetes deployment**:

```bash
# Get pod status
kubectl get pods -n ubuntu-autoinstall-webhook

# View pod logs
kubectl logs -f deployment/ubuntu-autoinstall-webhook -n ubuntu-autoinstall-webhook

# Describe the deployment
kubectl describe deployment ubuntu-autoinstall-webhook -n ubuntu-autoinstall-webhook
```

**Scaling the deployment**:

```bash
kubectl scale deployment ubuntu-autoinstall-webhook --replicas=3 -n ubuntu-autoinstall-webhook
```

**Updating the deployment**:

```bash
kubectl set image deployment/ubuntu-autoinstall-webhook ubuntu-autoinstall-webhook=ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest -n ubuntu-autoinstall-webhook
```

### 3.4. Service Dependencies

The Ubuntu Autoinstall Webhook system depends on the following external services:

1. **DHCP Server** (typically dnsmasq)
   - Required for handling PXE boot requests
   - Configuration must include PXE boot options

2. **TFTP Server** (built-in or external)
   - Required for serving boot files
   - Must be accessible to client systems

3. **Database** (SQLite by default, CockroachDB optional)
   - Stores system configuration and status
   - No external setup needed for SQLite

Example dnsmasq configuration for PXE booting:

```
# /etc/dnsmasq.conf
interface=eth0
domain=example.local
dhcp-range=192.168.1.50,192.168.1.150,12h
dhcp-boot=pxelinux.0
enable-tftp
tftp-root=/var/lib/tftpboot
```

### 3.5. Health Checks

The system provides health check endpoints for monitoring:

1. **HTTP Health Check**:
   - URL: `/api/v1/health`
   - Method: GET
   - Response: 200 OK if healthy

2. **Component Health Checks**:
   - URL: `/api/v1/health/components`
   - Method: GET
   - Response: Status of each component

Example health check script:

```bash
#!/bin/bash
# Health check script for automation tools

HEALTH_URL="https://localhost:8443/api/v1/health"

response=$(curl -sk $HEALTH_URL)
status=$?

if [ $status -ne 0 ]; then
  echo "ERROR: Could not connect to health endpoint"
  exit 1
fi

if [[ "$response" == *"healthy\":true"* ]]; then
  echo "System is healthy"
  exit 0
else
  echo "System is not healthy: $response"
  exit 1
fi
```

## 4. User Management

### 4.1. Authentication Methods

The system supports several authentication methods:

1. **Local Authentication**
   - Username and password stored in the system database
   - Passwords are securely hashed using bcrypt

2. **LDAP Authentication**
   - Integration with corporate LDAP directories
   - Configurable user attribute mapping

3. **OAuth2 Authentication**
   - Support for external identity providers
   - Configurable for providers like GitHub, Google, Azure AD

Example LDAP configuration:

```yaml
authentication:
  method: ldap
  ldap:
    url: "ldap://ldap.example.com:389"
    base_dn: "dc=example,dc=com"
    user_dn_pattern: "cn={0},ou=users,dc=example,dc=com"
    user_search_base: "ou=users"
    user_search_filter: "(uid={0})"
    group_search_base: "ou=groups"
    group_search_filter: "(member={0})"
    manager_dn: "cn=admin,dc=example,dc=com"
    manager_password: "password"
    user_attribute_mappings:
      username: "uid"
      email: "mail"
      display_name: "displayName"
```

### 4.2. User Roles and Permissions

The system implements Role-Based Access Control (RBAC) with these default roles:

1. **Administrator**
   - Full system access and control
   - Can manage users and roles
   - Can configure all system settings

2. **Operator**
   - Can manage systems and installations
   - Can view and modify templates
   - Cannot manage users or change system settings

3. **Viewer**
   - Can view systems and installation status
   - Cannot modify configurations or initiate installations

Custom roles can be created with specific permission sets:

```yaml
roles:
  - name: "DevOps"
    description: "Role for DevOps team members"
    permissions:
      - "systems:read"
      - "systems:create"
      - "systems:update"
      - "systems:install"
      - "templates:read"
      - "templates:use"
      - "installations:read"
      - "logs:read"
```

### 4.3. Managing Users

**Creating users with the CLI**:

```bash
# Create a new admin user
ubuntu-autoinstall-webhook create-admin --username admin --email admin@example.com

# Create a standard user with a specific role
ubuntu-autoinstall-webhook create-user --username operator --email operator@example.com --role operator
```

**Managing users via the Web UI**:

1. Navigate to "Settings" > "Users" in the web interface
2. Click "Add User" to create a new user
3. Fill in the required information:
   - Username
   - Email
   - Password
   - Role
4. Click "Save" to create the user

**User management API endpoints**:

```bash
# Create a new user
curl -X POST "https://localhost:8443/api/v1/users" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "SecurePassword123!",
    "role": "operator"
  }'

# Update a user
curl -X PUT "https://localhost:8443/api/v1/users/newuser" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "updated@example.com",
    "role": "viewer"
  }'

# Delete a user
curl -X DELETE "https://localhost:8443/api/v1/users/newuser" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### 4.4. Password Policies

The system supports configurable password policies:

```yaml
authentication:
  password_policy:
    min_length: 12
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_special_chars: true
    max_age_days: 90
    history_count: 5
    max_failed_attempts: 5
    lockout_time_minutes: 30
```

**Resetting user passwords**:

```bash
# Reset a user's password via CLI
ubuntu-autoinstall-webhook reset-password --username operator --generate
```

### 4.5. External Authentication

**OAuth2 Configuration Example**:

```yaml
authentication:
  method: oauth2
  oauth2:
    providers:
      github:
        enabled: true
        client_id: "your-client-id"
        client_secret: "your-client-secret"
        auth_url: "https://github.com/login/oauth/authorize"
        token_url: "https://github.com/login/oauth/access_token"
        user_info_url: "https://api.github.com/user"
        scope: "read:user,user:email"
        user_attribute_mappings:
          username: "login"
          email: "email"
          display_name: "name"
      google:
        enabled: false
        # Configuration for Google OAuth2
```

**SAML Configuration Example**:

```yaml
authentication:
  method: saml
  saml:
    idp_metadata_url: "https://idp.example.com/metadata"
    sp_entity_id: "https://ubuntu-autoinstall-webhook.example.com"
    assertion_consumer_service_url: "https://ubuntu-autoinstall-webhook.example.com/api/v1/auth/saml/callback"
    user_attribute_mappings:
      username: "NameID"
      email: "Email"
      display_name: "DisplayName"
```

## 5. Security Management

### 5.1. PKI Infrastructure

The system maintains its own Public Key Infrastructure (PKI) to secure communications:

1. **Root CA**: Self-signed certificate authority
2. **Intermediate CA**: Optional for larger deployments
3. **Service Certificates**: For internal services
4. **Client Certificates**: For mTLS authentication

**Initializing the PKI**:

```bash
# Initialize the root CA
ubuntu-autoinstall-webhook cert-init --ca-name "Ubuntu Autoinstall Webhook CA"

# Create an intermediate CA (optional)
ubuntu-autoinstall-webhook cert-create-intermediate --name "Intermediate CA"

# Issue service certificate
ubuntu-autoinstall-webhook cert-issue --name "webhook-server" --type server --dns webhook.example.com --ip 192.168.1.10
```

### 5.2. Certificate Management

**Certificate renewal**:

```bash
# Check for certificates nearing expiration
ubuntu-autoinstall-webhook cert-check-expiry

# Renew a specific certificate
ubuntu-autoinstall-webhook cert-renew --name "webhook-server"

# Renew all certificates expiring within 30 days
ubuntu-autoinstall-webhook cert-renew-all --days 30
```

**Certificate revocation**:

```bash
# Revoke a certificate
ubuntu-autoinstall-webhook cert-revoke --name "webhook-server" --reason "key-compromise"

# Generate CRL
ubuntu-autoinstall-webhook cert-gen-crl
```

### 5.3. TLS Configuration

The system uses secure TLS configurations by default. The following settings can be customized in `config.yaml`:

```yaml
server:
  tls:
    enabled: true
    cert_file: /etc/ubuntu-autoinstall-webhook/certs/server.crt
    key_file: /etc/ubuntu-autoinstall-webhook/certs/server.key
    min_version: "TLS1.2"
    cipher_suites:
      - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
      - TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
    prefer_server_cipher_suites: true
    client_auth: "request"  # options: "no", "request", "require"
    client_ca_file: /etc/ubuntu-autoinstall-webhook/certs/ca.crt
```

**Testing TLS Configuration**:

```bash
# Test server TLS configuration
openssl s_client -connect localhost:8443 -tls1_2

# Verify certificate chain
openssl verify -CAfile /etc/ubuntu-autoinstall-webhook/certs/ca.crt /etc/ubuntu-autoinstall-webhook/certs/server.crt
```

### 5.4. API Security

API access is secured through several mechanisms:

1. **API Keys**: Long-lived tokens for service integration
2. **JWT Tokens**: Short-lived tokens for user sessions
3. **Rate Limiting**: Prevents abuse of API endpoints
4. **IP Restrictions**: Optional IP-based access control

**API Key Management**:

```bash
# Generate a new API key
ubuntu-autoinstall-webhook generate-api-key --name "integration-service" --role operator

# List API keys
ubuntu-autoinstall-webhook list-api-keys

# Revoke an API key
ubuntu-autoinstall-webhook revoke-api-key --name "integration-service"
```

**API Security Configuration**:

```yaml
api:
  rate_limit:
    enabled: true
    requests_per_minute: 60
    burst: 10
  ip_restrictions:
    enabled: false
    allowed_ips:
      - "192.168.1.0/24"
      - "10.0.0.5"
  token_expiry:
    access_token_minutes: 15
    refresh_token_days: 7
```

### 5.5. Network Security

**Firewall Recommendations**:

The following ports should be open for system operation:

| Port | Protocol | Purpose                          |
| ---- | -------- | -------------------------------- |
| 8443 | TCP      | HTTPS API and Web UI             |
| 8080 | TCP      | HTTP (redirect to HTTPS)         |
| 69   | UDP      | TFTP Server (if using built-in)  |
| 67   | UDP      | DHCP Server (if running locally) |

**Example UFW configuration**:

```bash
sudo ufw allow 8443/tcp comment "Ubuntu Autoinstall Webhook HTTPS"
sudo ufw allow 8080/tcp comment "Ubuntu Autoinstall Webhook HTTP redirect"
sudo ufw allow 69/udp comment "TFTP Server"
```

**Securing network services**:

1. Use separate network segments for PXE booting when possible
2. Implement VLAN isolation for installation networks
3. Configure network ACLs to restrict access to the webhook server

### 5.6. Data Security

**Sensitive Data Handling**:

The system handles the following types of sensitive data:

1. User credentials
2. API keys
3. Private keys
4. Certificate authority keys

**Data Protection Measures**:

1. **Encryption at Rest**:
   - Database encryption options:
   ```yaml
   database:
     encryption:
       enabled: true
       key_file: /etc/ubuntu-autoinstall-webhook/keys/db-encryption-key
   ```

2. **Secure Storage of Secrets**:
   - Private keys stored with restricted permissions
   - Configuration option for external secret storage:
   ```yaml
   secrets:
     provider: vault
     vault:
       address: "https://vault.example.com:8200"
       token_file: "/etc/ubuntu-autoinstall-webhook/vault-token"
       path_prefix: "ubuntu-autoinstall-webhook"
   ```

### 5.7. System Hardening

**Operating System Hardening**:

1. Keep system updated with security patches
2. Use minimal installations for server environments
3. Implement AppArmor or SELinux profiles
4. Remove unnecessary packages and services

**Application Hardening**:

1. Run the service as a non-privileged user
2. Implement file permission restrictions
3. Configure proper umask settings
4. Use seccomp or capabilities for container deployments

**Example AppArmor Profile**:

```
# /etc/apparmor.d/usr.local.bin.ubuntu-autoinstall-webhook
#include <tunables/global>

/usr/local/bin/ubuntu-autoinstall-webhook {
  #include <abstractions/base>
  #include <abstractions/nameservice>
  #include <abstractions/ssl_certs>
  #include <abstractions/openssl>

  /usr/local/bin/ubuntu-autoinstall-webhook mr,
  /etc/ubuntu-autoinstall-webhook/** r,
  /etc/ubuntu-autoinstall-webhook/certs/* r,
  /var/lib/ubuntu-autoinstall-webhook/** rwk,
  /var/log/ubuntu-autoinstall-webhook/* w,
  /var/log/dnsmasq.log r,

  network tcp,
  network udp,
}
```

## 6. Backup and Recovery

Maintaining regular backups is critical for any production system. The Ubuntu Autoinstall Webhook system stores important data in multiple locations that need to be backed up consistently.

### 6.1. Backup Strategy

A comprehensive backup strategy for the system should include:

#### 6.1.1. Database Backup

The database contains critical system data including system records, templates, and installation status information.

**SQLite3 Backup (Default Database)**

For systems using the default SQLite3 database:

```bash
# Stop the service to ensure database consistency
systemctl stop ubuntu-autoinstall-webhook

# Backup the database
cp /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 /path/to/backup/database.sqlite3.bak

# Restart the service
systemctl start ubuntu-autoinstall-webhook
```

For more advanced SQLite backups with point-in-time recovery:

```bash
# Using sqlite3 command to dump the database
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 ".backup '/path/to/backup/database.sqlite3.bak'"
```

**CockroachDB Backup (Distributed Deployment)**

For systems using CockroachDB:

```bash
# Create a backup to a specific location
cockroach sql --execute="BACKUP DATABASE ubuntu_autoinstall TO 'nodelocal://1/backups/$(date +%Y-%m-%d)';"

# For cloud storage (AWS S3 example)
cockroach sql --execute="BACKUP DATABASE ubuntu_autoinstall TO 's3://bucket-name/backups/$(date +%Y-%m-%d)?AWS_ACCESS_KEY_ID=key&AWS_SECRET_ACCESS_KEY=secret';"
```

#### 6.1.2. Configuration Files Backup

Essential configuration files are stored in the `/etc/ubuntu-autoinstall-webhook/` directory:

```bash
# Create a backup of configuration files
tar -czf /path/to/backup/config-$(date +%Y-%m-%d).tar.gz /etc/ubuntu-autoinstall-webhook/
```

#### 6.1.3. PKI Infrastructure Backup

The Certificate Authority private keys and certificates are critical for secure system operation:

```bash
# Backup certificate store
tar -czf /path/to/backup/certificates-$(date +%Y-%m-%d).tar.gz /var/lib/ubuntu-autoinstall-webhook/certificates/
```

#### 6.1.4. Template Files Backup

Installation templates and custom files:

```bash
# Backup template files
tar -czf /path/to/backup/templates-$(date +%Y-%m-%d).tar.gz /var/lib/ubuntu-autoinstall-webhook/templates/
```

### 6.2. Automated Backup Configuration

To automate the backup process, create a backup script and schedule it with cron:

```bash
#!/bin/bash
# ubuntu-autoinstall-webhook-backup.sh

# Set backup directory
BACKUP_DIR="/path/to/backup"
TIMESTAMP=$(date +%Y-%m-%d-%H%M)

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Stop service temporarily
systemctl stop ubuntu-autoinstall-webhook

# Database backup
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 ".backup '$BACKUP_DIR/database-$TIMESTAMP.sqlite3'"

# Configuration files backup
tar -czf $BACKUP_DIR/config-$TIMESTAMP.tar.gz /etc/ubuntu-autoinstall-webhook/

# Certificate store backup
tar -czf $BACKUP_DIR/certificates-$TIMESTAMP.tar.gz /var/lib/ubuntu-autoinstall-webhook/certificates/

# Template files backup
tar -czf $BACKUP_DIR/templates-$TIMESTAMP.tar.gz /var/lib/ubuntu-autoinstall-webhook/templates/

# Start service
systemctl start ubuntu-autoinstall-webhook

# Clean up old backups (keep last 7 days)
find $BACKUP_DIR -name "*.sqlite3" -type f -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -type f -mtime +7 -delete

echo "Backup completed at $TIMESTAMP"
```

Make the script executable and schedule it:

```bash
chmod +x /usr/local/bin/ubuntu-autoinstall-webhook-backup.sh

# Add to crontab to run daily at 2:00 AM
(crontab -l 2>/dev/null; echo "0 2 * * * /usr/local/bin/ubuntu-autoinstall-webhook-backup.sh") | crontab -
```

### 6.3. Backup Verification

Regularly test your backups to ensure they can be used for system recovery:

1. Create a test environment mimicking your production setup
2. Restore backups to the test environment
3. Verify the system functions correctly with restored data
4. Document any issues encountered during the verification process

### 6.4. Recovery Procedures

#### 6.4.1. Full System Recovery

To restore the entire system after a catastrophic failure:

1. Install the Ubuntu Autoinstall Webhook software on a new server
2. Stop the service: `systemctl stop ubuntu-autoinstall-webhook`
3. Restore the database:

   ```bash
   # For SQLite3
   cp /path/to/backup/database.sqlite3.bak /var/lib/ubuntu-autoinstall-webhook/database.sqlite3

   # For CockroachDB
   cockroach sql --execute="RESTORE DATABASE ubuntu_autoinstall FROM 's3://bucket-name/backups/2023-01-01';"
   ```

4. Restore configuration files:

   ```bash
   tar -xzf /path/to/backup/config.tar.gz -C /
   ```

5. Restore certificate store:

   ```bash
   tar -xzf /path/to/backup/certificates.tar.gz -C /
   ```

6. Restore template files:

   ```bash
   tar -xzf /path/to/backup/templates.tar.gz -C /
   ```

7. Set correct permissions:

   ```bash
   chown -R ubuntu-autoinstall:ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/
   chmod -R 750 /var/lib/ubuntu-autoinstall-webhook/certificates/
   ```

8. Start the service: `systemctl start ubuntu-autoinstall-webhook`

#### 6.4.2. Database-Only Recovery

To recover just the database:

1. Stop the service: `systemctl stop ubuntu-autoinstall-webhook`
2. Restore the database file:

   ```bash
   # For SQLite3
   cp /path/to/backup/database.sqlite3.bak /var/lib/ubuntu-autoinstall-webhook/database.sqlite3
   chown ubuntu-autoinstall:ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/database.sqlite3
   chmod 640 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3
   ```

3. Start the service: `systemctl start ubuntu-autoinstall-webhook`

#### 6.4.3. Certificate Recovery

To recover just the certificate store:

1. Stop the service: `systemctl stop ubuntu-autoinstall-webhook`
2. Restore the certificate files:

   ```bash
   tar -xzf /path/to/backup/certificates.tar.gz -C /
   chown -R ubuntu-autoinstall:ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/certificates/
   chmod -R 750 /var/lib/ubuntu-autoinstall-webhook/certificates/
   ```

3. Start the service: `systemctl start ubuntu-autoinstall-webhook`

### 6.5. Disaster Recovery Planning

1. **Document Your Environment**:
   - Keep an updated inventory of all servers, network configurations, and service dependencies
   - Document all custom configurations and modifications

2. **Define Recovery Objectives**:
   - Recovery Time Objective (RTO): Maximum acceptable time to restore service
   - Recovery Point Objective (RPO): Maximum acceptable data loss period

3. **Create Recovery Runbooks**:
   - Step-by-step procedures for various failure scenarios
   - Clear assignment of responsibilities during recovery

4. **Off-site Backup Storage**:
   - Store backups in a geographically separate location
   - Consider secure cloud storage for backup redundancy

5. **Regular Testing**:
   - Conduct disaster recovery drills at least quarterly
   - Update procedures based on test results and system changes

## 7. Monitoring and Logging

Effective monitoring and logging are essential for maintaining the health, security, and performance of your Ubuntu Autoinstall Webhook deployment.

### 7.1. System Logging

#### 7.1.1. Log File Locations

The Ubuntu Autoinstall Webhook system writes logs to several locations:

- **Service Logs**: `/var/log/ubuntu-autoinstall-webhook/*.log`
- **Installation Logs**: `/var/log/ubuntu-autoinstall-webhook/installations/`
- **Web Access Logs**: `/var/log/ubuntu-autoinstall-webhook/webserver/access.log`
- **Web Error Logs**: `/var/log/ubuntu-autoinstall-webhook/webserver/error.log`
- **System Journal**: Service logs are also sent to systemd journal

#### 7.1.2. Log Configuration

Log settings can be adjusted in the main configuration file at `/etc/ubuntu-autoinstall-webhook/config.yaml`:

```yaml
logging:
  # Log levels: debug, info, warn, error
  level: info
  # Format options: text, json
  format: json
  # Output options: file, stdout, both
  output: both
  # File rotation
  rotate:
    max_size_mb: 100
    max_backups: 10
    max_age_days: 30
    compress: true
```

The log level can be adjusted dynamically without a service restart:

```bash
# Set log level to debug temporarily
curl -X PUT http://localhost:8080/api/v1/config/loglevel \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"level": "debug"}'
```

#### 7.1.3. Log Rotation

Logs are automatically rotated using logrotate. The default configuration rotates logs daily and keeps 14 days of history:

```bash
# /etc/logrotate.d/ubuntu-autoinstall-webhook
/var/log/ubuntu-autoinstall-webhook/*.log {
  daily
  rotate 14
  compress
  delaycompress
  missingok
  notifempty
  create 0640 ubuntu-autoinstall ubuntu-autoinstall
  postrotate
    systemctl kill -s HUP ubuntu-autoinstall-webhook >/dev/null 2>&1 || true
  endscript
}
```

You can customize this configuration to meet your specific needs.

### 7.2. Monitoring System Health

#### 7.2.1. Built-in Health Endpoints

The system provides health check endpoints for monitoring:

- **Basic Health**: `http://localhost:8080/health`
- **Detailed Health**: `http://localhost:8080/health/detailed`
- **Component Status**: `http://localhost:8080/health/components`

Example health check with curl:

```bash
curl -s http://localhost:8080/health | jq
```

Sample output:
```json
{
  "status": "healthy",
  "version": "1.2.3",
  "timestamp": "2023-01-15T12:34:56Z",
  "uptime": "3d 5h 12m"
}
```

For detailed component status:

```bash
curl -s http://localhost:8080/health/components | jq
```

Sample output:
```json
{
  "status": "healthy",
  "components": {
    "database": {
      "status": "healthy",
      "details": "Connected to SQLite database"
    },
    "file_editor": {
      "status": "healthy",
      "details": "Leader election active"
    },
    "webserver": {
      "status": "healthy",
      "details": "Listening on ports 8080, 8443"
    },
    "dnsmasq_watcher": {
      "status": "healthy",
      "details": "Monitoring /var/log/dnsmasq.log"
    },
    "cert_issuer": {
      "status": "healthy",
      "details": "CA certificates valid"
    }
  }
}
```

#### 7.2.2. Monitoring with Prometheus

The system exposes metrics in Prometheus format at `http://localhost:8080/metrics`. These metrics can be collected by a Prometheus server for monitoring and alerting.

Example Prometheus configuration to scrape metrics:

```yaml
scrape_configs:
  - job_name: 'ubuntu-autoinstall-webhook'
    scrape_interval: 15s
    static_configs:
      - targets: ['webhook-server:8080']
```

Key metrics exposed include:

- **`webhook_http_requests_total`**: Total HTTP requests processed
- **`webhook_http_request_duration_seconds`**: HTTP request latency histogram
- **`webhook_installations_total`**: Total installation attempts
- **`webhook_installations_active`**: Currently active installations
- **`webhook_installations_completed`**: Successfully completed installations
- **`webhook_installations_failed`**: Failed installations
- **`webhook_database_queries_total`**: Total database queries executed
- **`webhook_file_operations_total`**: Total filesystem operations
- **`webhook_certificate_operations_total`**: Certificate operations performed
- **`webhook_memory_usage_bytes`**: Memory usage of the application
- **`webhook_goroutines`**: Number of active goroutines

#### 7.2.3. Grafana Dashboards

A sample Grafana dashboard is available to visualize system metrics. Import the dashboard JSON from `/usr/share/ubuntu-autoinstall-webhook/dashboards/system-overview.json` into your Grafana instance.

The dashboard includes panels for:
- Request rates and latencies
- Installation success/failure metrics
- Resource utilization
- Component health status
- Error rates

### 7.3. Alerting

#### 7.3.1. Configuring Alerts with Prometheus AlertManager

Create alert rules for common failure scenarios:

```yaml
# /etc/prometheus/rules/ubuntu-autoinstall-webhook.yml
groups:
- name: ubuntu-autoinstall-webhook
  rules:
  - alert: HighErrorRate
    expr: rate(webhook_http_requests_total{status_code=~"5.."}[5m]) > 0.01
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value | humanizePercentage }} for the past 5 minutes"

  - alert: ServiceDown
    expr: up{job="ubuntu-autoinstall-webhook"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Service is down"
      description: "The Ubuntu Autoinstall Webhook service is not responding"

  - alert: ComponentUnhealthy
    expr: webhook_component_health{status="unhealthy"} > 0
    for: 3m
    labels:
      severity: critical
    annotations:
      summary: "Component is unhealthy"
      description: "The {{ $labels.component }} component is reporting unhealthy status"

  - alert: HighInstallationFailureRate
    expr: rate(webhook_installations_failed[30m]) / rate(webhook_installations_total[30m]) > 0.2
    for: 10m
    labels:
      severity: warning
    annotations:
      summary: "High installation failure rate"
      description: "Installation failure rate is {{ $value | humanizePercentage }} in the past 30 minutes"
```

#### 7.3.2. Email Notifications

Configure AlertManager to send email notifications:

```yaml
# /etc/alertmanager/alertmanager.yml
global:
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alertmanager@example.com'
  smtp_auth_username: 'alertmanager'
  smtp_auth_password: 'password'
  smtp_require_tls: true

route:
  receiver: 'team-email'
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  routes:
  - match:
      severity: critical
    receiver: 'team-pager'
    repeat_interval: 1h

receivers:
- name: 'team-email'
  email_configs:
  - to: 'team@example.com'

- name: 'team-pager'
  email_configs:
  - to: 'oncall@example.com'
  - to: 'backup-oncall@example.com'
```

#### 7.3.3. Webhook Notifications

For integration with services like Slack, Microsoft Teams, or PagerDuty:

```yaml
receivers:
- name: 'slack-notifications'
  slack_configs:
  - api_url: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXX'
    channel: '#alerts'
    text: "{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}\n{{ end }}"
    send_resolved: true

- name: 'pagerduty'
  pagerduty_configs:
  - service_key: '1234567890abcdef'
    description: '{{ .CommonAnnotations.summary }}'
    client: 'AlertManager'
    client_url: 'https://alertmanager.example.com'
    details:
      firing: '{{ .Alerts.Firing | len }}'
      description: '{{ .CommonAnnotations.description }}'
```

### 7.4. Log Analysis

#### 7.4.1. Common Log Patterns

Understanding common log patterns helps with troubleshooting:

**Successful Installation Start**:
```
{"level":"info","timestamp":"2023-01-15T12:34:56Z","message":"Installation started","system":"server1","mac":"00:11:22:33:44:55","template":"minimal-server"}
```

**Installation Completion**:
```
{"level":"info","timestamp":"2023-01-15T13:00:00Z","message":"Installation completed successfully","system":"server1","mac":"00:11:22:33:44:55","duration_seconds":1504}
```

**Installation Failure**:
```
{"level":"error","timestamp":"2023-01-15T12:40:00Z","message":"Installation failed","system":"server2","mac":"AA:BB:CC:DD:EE:FF","error":"Disk partitioning failed","stage":"storage_configuration"}
```

**System Discovery**:
```
{"level":"info","timestamp":"2023-01-15T12:30:00Z","message":"New system discovered","mac":"00:11:22:33:44:55","ip":"192.168.1.100"}
```

**Service Error**:
```
{"level":"error","timestamp":"2023-01-15T14:00:00Z","message":"Service error","component":"database","error":"connection refused","attempt":3}
```

#### 7.4.2. Using jq for Log Analysis

The JSON log format allows for powerful filtering and analysis with tools like `jq`:

```bash
# Find all failed installations
cat /var/log/ubuntu-autoinstall-webhook/webhook.log | jq -c 'select(.message == "Installation failed")'

# Count installations by status
cat /var/log/ubuntu-autoinstall-webhook/webhook.log | jq -c 'select(.message | test("Installation (started|completed|failed)"))' | \
  jq -r '.message' | sort | uniq -c

# Find all errors for a specific system
cat /var/log/ubuntu-autoinstall-webhook/webhook.log | jq -c 'select(.level == "error" and .system == "server1")'

# Track installation duration
cat /var/log/ubuntu-autoinstall-webhook/webhook.log | jq -c 'select(.message == "Installation completed successfully")' | \
  jq -r '[.timestamp, .system, .duration_seconds] | @csv' > installation_durations.csv
```

#### 7.4.3. Centralized Logging

For production environments, sending logs to a centralized logging system is recommended:

**Configuring rsyslog forwarding**:

Add the following to `/etc/rsyslog.d/10-ubuntu-autoinstall-webhook.conf`:

```
# Forward Ubuntu Autoinstall Webhook logs to a central syslog server
if $programname == 'ubuntu-autoinstall-webhook' then @@logserver.example.com:514
```

**ELK Stack Integration**:

Configure Filebeat to ship logs to Elasticsearch:

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/ubuntu-autoinstall-webhook/*.log
  json.keys_under_root: true
  json.add_error_key: true
  fields:
    application: ubuntu-autoinstall-webhook
  fields_under_root: true

output.elasticsearch:
  hosts: ["elasticsearch.example.com:9200"]
```

### 7.5. Audit Logging

The system maintains separate audit logs for security-relevant events:

#### 7.5.1. Audit Log Location

Audit logs are stored in `/var/log/ubuntu-autoinstall-webhook/audit.log` and contain records of:

- User authentication attempts (successful and failed)
- Administrative actions (user/role changes, system configuration)
- Template modifications
- Installation initiations
- Certificate operations

#### 7.5.2. Audit Log Format

Audit logs use a structured JSON format for easier parsing:

```json
{
  "timestamp": "2023-01-15T12:34:56Z",
  "event_type": "USER_LOGIN",
  "user": "admin",
  "result": "success",
  "client_ip": "192.168.1.100",
  "details": {
    "method": "password"
  }
}
```

#### 7.5.3. Audit Log Retention

Audit logs have a separate retention policy from regular logs. By default, they are kept for 1 year:

```
# /etc/logrotate.d/ubuntu-autoinstall-webhook-audit
/var/log/ubuntu-autoinstall-webhook/audit.log {
  monthly
  rotate 12
  compress
  delaycompress
  missingok
  notifempty
  create 0640 ubuntu-autoinstall ubuntu-autoinstall
}
```

#### 7.5.4. Audit Log Analysis

Query the audit logs for security analysis:

```bash
# Find all failed login attempts
cat /var/log/ubuntu-autoinstall-webhook/audit.log | jq -c 'select(.event_type == "USER_LOGIN" and .result == "failure")'

# Count actions by user
cat /var/log/ubuntu-autoinstall-webhook/audit.log | jq -r '.user' | sort | uniq -c | sort -nr

# Identify suspicious activity (multiple failures)
cat /var/log/ubuntu-autoinstall-webhook/audit.log | \
  jq -c 'select(.event_type == "USER_LOGIN" and .result == "failure")' | \
  jq -r '.client_ip' | sort | uniq -c | sort -nr | head -10
```

## 8. Performance Tuning

As your Ubuntu Autoinstall Webhook deployment grows, you may need to optimize performance to support larger environments and more concurrent installations.

### 8.1. Resource Requirements

The system's resource requirements depend on the scale of your deployment:

| Scale      | Systems  | Concurrent Installations | CPU Cores | RAM    | Storage |
| ---------- | -------- | ------------------------ | --------- | ------ | ------- |
| Small      | < 100    | 1-5                      | 2         | 4 GB   | 20 GB   |
| Medium     | 100-500  | 5-20                     | 4         | 8 GB   | 50 GB   |
| Large      | 500-2000 | 20-50                    | 8         | 16 GB  | 100 GB  |
| Enterprise | 2000+    | 50+                      | 16+       | 32+ GB | 200+ GB |

### 8.2. System Optimization

#### 8.2.1. Process Limits

Adjust system limits in `/etc/security/limits.conf`:

```
ubuntu-autoinstall    soft    nofile      65536
ubuntu-autoinstall    hard    nofile      65536
ubuntu-autoinstall    soft    nproc       16384
ubuntu-autoinstall    hard    nproc       16384
```

Update systemd service limits in `/etc/systemd/system/ubuntu-autoinstall-webhook.service.d/override.conf`:

```
[Service]
LimitNOFILE=65536
LimitNPROC=16384
```

Apply changes:

```bash
systemctl daemon-reload
systemctl restart ubuntu-autoinstall-webhook
```

#### 8.2.2. Memory Management

Configure Go runtime parameters in the service environment:

```bash
# /etc/default/ubuntu-autoinstall-webhook
GOGC=100
GOMEMLIMIT=0
GOMAXPROCS=0
```

For memory-constrained environments:

```bash
# /etc/default/ubuntu-autoinstall-webhook
GOGC=50
GOMEMLIMIT=1024MiB
GOMAXPROCS=2
```

#### 8.2.3. Process Priority

For shared systems, adjust the service nice level:

```
# /etc/systemd/system/ubuntu-autoinstall-webhook.service.d/override.conf
[Service]
Nice=10
```

### 8.3. Database Optimization

#### 8.3.1. SQLite Optimization

For SQLite deployments:

1. Enable WAL mode:

```bash
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 "PRAGMA journal_mode=WAL;"
```

2. Optimize other SQLite parameters:

```bash
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 << EOF
PRAGMA synchronous = NORMAL;
PRAGMA temp_store = MEMORY;
PRAGMA mmap_size = 30000000000;
PRAGMA cache_size = -64000;
EOF
```

3. Regular maintenance:

```bash
# Add to cron to run weekly
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 "VACUUM;"
```

#### 8.3.2. CockroachDB Optimization

For CockroachDB deployments:

1. Indexing strategy:

```sql
-- Add indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_systems_status ON systems(status);
CREATE INDEX IF NOT EXISTS idx_installations_status ON installations(status);
CREATE INDEX IF NOT EXISTS idx_systems_mac_address ON systems(mac_address);
```

2. SQL query optimization:

```sql
-- Example optimized query
EXPLAIN ANALYZE SELECT * FROM systems
WHERE status = 'installing'
AND last_seen > now() - INTERVAL '6 hours';
```

3. Connection pooling:

Configure in `config.yaml`:

```yaml
database:
  cockroach:
    max_open_conns: 25
    max_idle_conns: 10
    conn_max_lifetime_minutes: 60
```

### 8.4. Web Server Tuning

#### 8.4.1. HTTP Optimization

Adjust HTTP server parameters in `config.yaml`:

```yaml
webserver:
  read_timeout_seconds: 30
  write_timeout_seconds: 60
  idle_timeout_seconds: 120
  max_header_bytes: 1048576  # 1MB
  shutdown_timeout_seconds: 30
```

#### 8.4.2. Connection Limits

```yaml
webserver:
  max_connections: 1000
  keep_alive_seconds: 60
```

#### 8.4.3. Static File Serving

Optimize static file serving:

```yaml
webserver:
  static_files:
    cache_control: "public, max-age=86400"
    gzip: true
    min_compress_size: 1024  # 1KB
```

#### 8.4.4. Load Balancing

For multi-instance deployments, use a load balancer:

```
+----------------+
| Load Balancer  |
+----------------+
        |
        +--------------------------+
        |                          |
+----------------+       +----------------+
| Instance 1     |       | Instance 2     |
+----------------+       +----------------+
        |                          |
        +--------------------------+
        |
+----------------+
| Database       |
+----------------+
```

### 8.5. Filesystem Optimization

#### 8.5.1. Filesystem Selection

For production deployments, choose an appropriate filesystem:

- **ext4**: Good general-purpose choice
- **XFS**: Better for large files and high concurrency
- **ZFS**: Advanced features like compression and snapshots

#### 8.5.2. Mount Options

Optimize mount options in `/etc/fstab`:

```
# Example for an SSD
/dev/sda1 /var/lib/ubuntu-autoinstall-webhook ext4 defaults,noatime,discard 0 2
```

#### 8.5.3. I/O Scheduling

For SSDs, use the `none` or `mq-deadline` scheduler:

```bash
echo "mq-deadline" > /sys/block/sda/queue/scheduler
```

Add to `/etc/sysfs.conf` for persistence:

```
block/sda/queue/scheduler = mq-deadline
```

### 8.6. Network Optimization

#### 8.6.1. TCP Tuning

Optimize kernel parameters in `/etc/sysctl.conf`:

```
# Increase TCP max buffer size
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216

# Increase Linux autotuning TCP buffer limits
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216

# Enable TCP fast open
net.ipv4.tcp_fastopen = 3

# Other optimizations
net.ipv4.tcp_slow_start_after_idle = 0
net.ipv4.tcp_no_metrics_save = 1
```

Apply changes:

```bash
sysctl -p
```

#### 8.6.2. Network Interface Configuration

For high-performance networks:

```bash
# Set NIC parameters
ethtool -G eth0 rx 4096 tx 4096
ethtool -K eth0 tso on gso on gro on
```

### 8.7. Caching Strategies

#### 8.7.1. Application-Level Caching

Configure in-memory cache parameters in `config.yaml`:

```yaml
cache:
  # Cache size in MB
  size_mb: 256
  # TTL in seconds
  ttl_seconds: 300
  # Cache types
  enable_template_cache: true
  enable_config_cache: true
  enable_system_cache: true
```

#### 8.7.2. External Caching

For large deployments, add Redis caching:

```yaml
cache:
  type: "redis"
  redis:
    address: "redis:6379"
    password: ""
    db: 0
    pool_size: 10
```

#### 8.7.3. Content Delivery

For geographically distributed deployments, consider using a CDN for static boot files.

### 8.8. Monitoring Performance

#### 8.8.1. Key Performance Indicators

Monitor these key metrics:

- **Response Time**: API endpoint response times
- **Throughput**: Requests per second
- **Error Rate**: Percentage of failed requests
- **Resource Utilization**: CPU, memory, disk, network usage
- **Installation Time**: End-to-end installation duration
- **Concurrency**: Number of simultaneous installations

#### 8.8.2. Performance Testing

Conduct regular performance tests:

```bash
# Example using vegeta for HTTP load testing
echo "GET http://localhost:8080/api/v1/systems" | \
  vegeta attack -rate=50 -duration=30s | \
  vegeta report
```

### 8.9. Scaling Strategies

#### 8.9.1. Vertical Scaling

- Increase resources (CPU, RAM) on the existing server
- Improve storage with faster disks (NVMe SSDs)
- Upgrade network interfaces to higher bandwidth

#### 8.9.2. Horizontal Scaling

- Deploy multiple instances behind a load balancer
- Use shared storage for installation files
- Implement database replication or clustering
- Configure leader election for coordinated operations

#### 8.9.3. Service Decomposition

For very large deployments, consider separating services:

```
+----------------+   +----------------+   +----------------+
| Web Frontend   |   | API Server     |   | File Editor    |
+----------------+   +----------------+   +----------------+
        |                    |                    |
+------------------------------------------------------+
|                      Message Bus                      |
+------------------------------------------------------+
        |                    |                    |
+----------------+   +----------------+   +----------------+
| DNSMasq Watcher |   | Cert Issuer    |   | Config Gen    |
+----------------+   +----------------+   +----------------+
```

### 8.10. Hardware Recommendations

#### 8.10.1. CPU

- Prefer higher clock speeds over more cores for small deployments
- Balance core count and clock speed for medium deployments
- Prioritize core count for large deployments with many concurrent installations

#### 8.10.2. Memory

- Minimum: 4GB RAM
- Recommended: 16GB RAM for medium deployments
- Consider 32GB+ for large-scale environments

#### 8.10.3. Storage

- Use SSDs for all deployments
- NVMe SSDs recommended for database and log storage
- RAID configurations for redundancy in production environments
- Separate volumes for logs, database, and installation files

#### 8.10.4. Network

- Minimum: 1 Gbps Ethernet
- Recommended: 10 Gbps Ethernet for medium to large deployments
- Consider redundant NICs for high-availability configurations

## 9. Security Management

Security is a critical aspect of the Ubuntu Autoinstall Webhook system, especially as it manages the installation of operating systems across your infrastructure. This section covers security best practices, configuration, and management.

### 9.1. Authentication and Authorization

#### 9.1.1. Authentication Methods

The system supports multiple authentication methods:

1. **Local Authentication**: Username and password stored in the database
2. **LDAP/Active Directory**: Integration with enterprise directory services
3. **OAuth2/OpenID Connect**: Integration with identity providers
4. **API Token Authentication**: For programmatic access

Configure authentication methods in `/etc/ubuntu-autoinstall-webhook/config.yaml`:

```yaml
auth:
  # Local authentication settings
  local:
    enabled: true
    password_min_length: 12
    password_require_mixed_case: true
    password_require_number: true
    password_require_special: true
    password_max_age_days: 90

  # LDAP configuration
  ldap:
    enabled: false
    server: "ldap.example.com"
    port: 636
    use_ssl: true
    bind_dn: "cn=service-account,ou=users,dc=example,dc=com"
    bind_password: "secret"
    search_base: "ou=users,dc=example,dc=com"
    search_filter: "(sAMAccountName=%s)"
    group_search_base: "ou=groups,dc=example,dc=com"
    group_search_filter: "(member=%s)"
    admin_group: "ubuntu-autoinstall-admins"

  # OAuth2 configuration
  oauth2:
    enabled: false
    provider: "github"
    client_id: "your-client-id"
    client_secret: "your-client-secret"
    redirect_url: "https://webhook.example.com/auth/callback"
    scopes: ["user:email"]
```

#### 9.1.2. Role-Based Access Control (RBAC)

The system implements RBAC with predefined roles:

1. **Admin**: Full access to all features
2. **Operator**: Can manage systems and installations but not users or settings
3. **Installer**: Can only initiate and monitor installations
4. **Viewer**: Read-only access to system status and logs

Custom roles can be defined in the configuration:

```yaml
rbac:
  custom_roles:
    - name: "TemplateManager"
      description: "Can create and edit templates"
      permissions:
        - "templates:read"
        - "templates:write"
        - "systems:read"

    - name: "SecurityAuditor"
      description: "Audit security settings and logs"
      permissions:
        - "logs:read"
        - "certificates:read"
        - "users:read"
        - "audit:read"
```

Assign roles to users through the web interface or API.

#### 9.1.3. Managing API Tokens

API tokens provide programmatic access to the system. Manage tokens securely:

1. Generate tokens with appropriate scopes:

```bash
# Using the built-in CLI
ubuntu-autoinstall-webhook tokens create --name="automation-token" \
  --scopes="systems:read,installations:write" \
  --expires-in="720h"
```

2. Rotate tokens regularly:

```bash
# List existing tokens
ubuntu-autoinstall-webhook tokens list

# Revoke a token
ubuntu-autoinstall-webhook tokens revoke --id="token-id"
```

3. Configure token settings:

```yaml
api:
  tokens:
    max_lifetime_hours: 8760  # 1 year
    default_lifetime_hours: 720  # 30 days
    inactive_timeout_hours: 72
```

### 9.2. TLS Configuration

#### 9.2.1. Web Interface TLS

Configure TLS for the web interface in `/etc/ubuntu-autoinstall-webhook/config.yaml`:

```yaml
webserver:
  tls:
    enabled: true
    cert_file: "/etc/ubuntu-autoinstall-webhook/certs/server.crt"
    key_file: "/etc/ubuntu-autoinstall-webhook/certs/server.key"
    min_version: "1.2"  # TLS 1.2
    preferred_cipher_suites:
      - "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
      - "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
      - "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
    hsts_max_age_seconds: 31536000  # 1 year
    hsts_include_subdomains: true
```

#### 9.2.2. Certificate Management

The system includes tools for managing certificates:

1. Generate a new server certificate:

```bash
ubuntu-autoinstall-webhook cert generate --type=server \
  --cn="webhook.example.com" \
  --org="Example Corp" \
  --out-cert="/etc/ubuntu-autoinstall-webhook/certs/server.crt" \
  --out-key="/etc/ubuntu-autoinstall-webhook/certs/server.key"
```

2. Install a commercial or Let's Encrypt certificate:

```bash
# Copy your certificates
cp /etc/letsencrypt/live/webhook.example.com/fullchain.pem \
   /etc/ubuntu-autoinstall-webhook/certs/server.crt

cp /etc/letsencrypt/live/webhook.example.com/privkey.pem \
   /etc/ubuntu-autoinstall-webhook/certs/server.key

# Set permissions
chown ubuntu-autoinstall:ubuntu-autoinstall /etc/ubuntu-autoinstall-webhook/certs/server.*
chmod 640 /etc/ubuntu-autoinstall-webhook/certs/server.*

# Restart service
systemctl restart ubuntu-autoinstall-webhook
```

3. Certificate renewal automation:

```bash
# For Let's Encrypt certificates
echo '#!/bin/bash
cp /etc/letsencrypt/live/webhook.example.com/fullchain.pem /etc/ubuntu-autoinstall-webhook/certs/server.crt
cp /etc/letsencrypt/live/webhook.example.com/privkey.pem /etc/ubuntu-autoinstall-webhook/certs/server.key
chown ubuntu-autoinstall:ubuntu-autoinstall /etc/ubuntu-autoinstall-webhook/certs/server.*
chmod 640 /etc/ubuntu-autoinstall-webhook/certs/server.*
systemctl reload ubuntu-autoinstall-webhook
' > /etc/letsencrypt/renewal-hooks/post/ubuntu-autoinstall-webhook.sh

chmod +x /etc/letsencrypt/renewal-hooks/post/ubuntu-autoinstall-webhook.sh
```

#### 9.2.3. Mutual TLS (mTLS)

Configure mutual TLS for client authentication:

```yaml
webserver:
  tls:
    # Enable mutual TLS
    mutual_tls:
      enabled: true
      client_ca_file: "/etc/ubuntu-autoinstall-webhook/certs/client-ca.crt"
      client_verification: "require"  # Options: require, optional

  # Exempt paths from mTLS (for initial access)
  tls_exempt_paths:
    - "/api/v1/auth/login"
    - "/api/v1/cert/request"
```

### 9.3. Network Security

#### 9.3.1. Firewall Configuration

Recommended firewall rules:

```bash
# Allow web interface access
ufw allow 8443/tcp comment 'HTTPS Web Interface'

# Allow PXE/installation services
ufw allow 8080/tcp comment 'HTTP Installation Files'
ufw allow 69/udp comment 'TFTP'

# Optional - only if using external database
ufw allow out to 192.168.1.10 port 5432 proto tcp comment 'PostgreSQL'

# Apply rules
ufw enable
```

#### 9.3.2. Network Isolation

For production environments, implement network segmentation:

1. **Management Network**: For administrative access
2. **Provisioning Network**: For PXE boot and installations
3. **Storage Network**: For accessing installation files

Example network configuration:

```
+-------------------+
| Web/Admin Traffic |
|    (10.1.0.0/24)  |
+-------------------+
         |
+-------------------+
| Ubuntu Autoinstall|
| Webhook Server    |
+-------------------+
         |
+-------------------+
| Provisioning/PXE  |
|   (10.2.0.0/24)   |
+-------------------+
```

Configure multiple network interfaces in `/etc/ubuntu-autoinstall-webhook/config.yaml`:

```yaml
webserver:
  interfaces:
    admin: "10.1.0.5"  # Management interface
    provision: "10.2.0.5"  # Provisioning interface

  # Bind services to specific interfaces
  bind:
    web_interface: "admin"
    installation_files: "provision"
    tftp: "provision"
```

#### 9.3.3. Traffic Encryption

Ensure all traffic is encrypted:

1. Enable HTTPS for web interface
2. Configure iPXE to use HTTPS for downloads
3. Use WS-TLS for WebSocket connections

```yaml
webserver:
  secure_boot_files: true
  ipxe:
    force_https: true
```

### 9.4. Secure Installation

#### 9.4.1. Secure Boot

Configure the system to support UEFI Secure Boot:

```yaml
installation:
  secure_boot:
    enabled: true
    keys_path: "/etc/ubuntu-autoinstall-webhook/secure-boot-keys/"
    require_signed_kernels: true
```

#### 9.4.2. Installation Authentication

Secure the installation process with authentication:

```yaml
installation:
  authentication:
    method: "token"  # Options: none, token, certificate
    token_lifetime_minutes: 30
    max_attempts: 3
```

#### 9.4.3. Installation Encryption

Configure disk encryption for installations:

```yaml
installation:
  encryption:
    enabled: true
    method: "luks"  # Options: luks, zfs
    key_management: "tpm"  # Options: passphrase, tpm, network
```

### 9.5. Vulnerability Management

#### 9.5.1. Security Updates

Keep the system updated with security patches:

```bash
# Update the package
apt-get update
apt-get install --only-upgrade ubuntu-autoinstall-webhook

# Restart the service
systemctl restart ubuntu-autoinstall-webhook
```

For automatic updates, configure `unattended-upgrades`:

```bash
apt-get install unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades
```

Edit `/etc/apt/apt.conf.d/50unattended-upgrades`:

```
Unattended-Upgrade::Allowed-Origins {
    "Ubuntu focal-security";
    "Ubuntu focal-updates";
    "UbuntuAutoinstallWebhook:focal";
};

Unattended-Upgrade::Package-Blacklist {
};
```

#### 9.5.2. Regular Security Audits

Perform regular security audits:

1. Review system logs for unauthorized access attempts
2. Audit user accounts and permissions
3. Verify TLS certificate validity and configurations
4. Check for outdated software components
5. Scan for vulnerabilities in dependencies

#### 9.5.3. Security Hardening

Additional security hardening steps:

1. Enable AppArmor profiles:

```bash
aa-enforce /etc/apparmor.d/usr.bin.ubuntu-autoinstall-webhook
systemctl reload apparmor
```

2. Restrict filesystem access:

```bash
# Ensure proper permissions
find /etc/ubuntu-autoinstall-webhook -type f -exec chmod 640 {} \;
find /etc/ubuntu-autoinstall-webhook -type d -exec chmod 750 {} \;
```

3. Implement IP-based access restrictions:

```yaml
webserver:
  security:
    allowed_ip_ranges:
      - "10.0.0.0/8"
      - "172.16.0.0/12"
      - "192.168.0.0/16"
```

### 9.6. Data Protection

#### 9.6.1. Sensitive Data Handling

Configure how sensitive data is handled:

```yaml
security:
  sensitive_data:
    encryption_key_file: "/etc/ubuntu-autoinstall-webhook/keys/data-encryption.key"
    encrypt_passwords: true
    encrypt_ssh_keys: true
    encrypt_api_tokens: true
    hide_sensitive_logs: true
```

#### 9.6.2. Data Retention

Configure data retention policies:

```yaml
security:
  data_retention:
    installation_logs_days: 90
    completed_installations_days: 180
    failed_installations_days: 30
    audit_logs_days: 365
```

#### 9.6.3. Secrets Management

For managing secrets like passwords and keys:

1. Generate a new encryption key:

```bash
ubuntu-autoinstall-webhook keys generate \
  --type=encryption \
  --out=/etc/ubuntu-autoinstall-webhook/keys/data-encryption.key
```

2. Rotate encryption keys:

```bash
# Generate new key
ubuntu-autoinstall-webhook keys generate \
  --type=encryption \
  --out=/etc/ubuntu-autoinstall-webhook/keys/data-encryption-new.key

# Re-encrypt data with new key
ubuntu-autoinstall-webhook maintenance reencrypt-data \
  --old-key=/etc/ubuntu-autoinstall-webhook/keys/data-encryption.key \
  --new-key=/etc/ubuntu-autoinstall-webhook/keys/data-encryption-new.key

# Replace old key with new key
mv /etc/ubuntu-autoinstall-webhook/keys/data-encryption-new.key \
   /etc/ubuntu-autoinstall-webhook/keys/data-encryption.key
```

## 10. Troubleshooting

This section provides guidance for identifying, diagnosing, and resolving common issues with the Ubuntu Autoinstall Webhook system.

### 10.1. Diagnostic Tools

#### 10.1.1. System Status Check

The system includes a comprehensive status check tool:

```bash
ubuntu-autoinstall-webhook status check
```

This command checks:
- Service status for all components
- Database connectivity
- File system permissions
- Certificate validity
- Network connectivity
- External service dependencies

For a more detailed report:

```bash
ubuntu-autoinstall-webhook status check --verbose
```

#### 10.1.2. Log Analysis

Quickly search logs for specific issues:

```bash
# Search for errors in the main log
grep -i error /var/log/ubuntu-autoinstall-webhook/webhook.log

# Show recent installation failures
grep -i "installation failed" /var/log/ubuntu-autoinstall-webhook/webhook.log | tail -n 20

# Check for authentication failures
grep -i "authentication failed" /var/log/ubuntu-autoinstall-webhook/audit.log
```

For structured analysis of JSON logs:

```bash
# Count errors by component
cat /var/log/ubuntu-autoinstall-webhook/webhook.log | \
  jq -r 'select(.level=="error") | .component' | \
  sort | uniq -c | sort -nr

# Find systems with most installation failures
cat /var/log/ubuntu-autoinstall-webhook/webhook.log | \
  jq -r 'select(.message=="Installation failed") | .system' | \
  sort | uniq -c | sort -nr | head -10
```

#### 10.1.3. Database Inspection

Check database health and contents:

```bash
# For SQLite3
ubuntu-autoinstall-webhook db check --fix-integrity

# Query specific information
ubuntu-autoinstall-webhook db query "SELECT mac_address, hostname, status FROM systems WHERE status='failed'"

# Export data for analysis
ubuntu-autoinstall-webhook db export --tables=systems,installations --format=csv --output=/tmp/export
```

#### 10.1.4. Network Diagnostics

Test network connectivity for installation services:

```bash
# Test if TFTP is accessible
ubuntu-autoinstall-webhook network test --service=tftp --interface=eth0

# Check DNS resolution
ubuntu-autoinstall-webhook network test --dns=webhook.example.com

# Validate port availability
ubuntu-autoinstall-webhook network test --ports=8080,8443,69
```

### 10.2. Common Issues and Solutions

#### 10.2.1. Installation Failures

**Symptom**: Systems fail during installation phase

**Possible Causes and Solutions**:

1. **PXE Boot Issues**:
   - Verify DHCP is providing correct next-server and filename options
   - Check that TFTP service is running and accessible
   - Ensure boot files (initrd, kernel) are present and readable

   ```bash
   # Check DHCP options
   tcpdump -i eth0 -n port 67 or port 68

   # Verify TFTP service
   systemctl status tftpd-hpa

   # Test TFTP file retrieval
   tftp 192.168.1.1 -c get pxelinux.0
   ```

2. **Network Configuration**:
   - Confirm target system has network connectivity to the server
   - Verify IP addressing in templates is correct
   - Check for network restrictions (firewalls, VLANs)

3. **Storage Issues**:
   - Verify disk meets minimum size requirements
   - Check for unsupported RAID controllers
   - Ensure partitioning scheme in template is valid

   ```bash
   # Review partitioning errors in logs
   grep -i "partition\|storage\|disk" /var/log/ubuntu-autoinstall-webhook/installations/failed-server.log
   ```

4. **Template Problems**:
   - Validate template syntax
   - Check for missing or invalid variables
   - Test template generation with debug mode:

   ```bash
   ubuntu-autoinstall-webhook template debug --id=template-uuid --mac=00:11:22:33:44:55
   ```

#### 10.2.2. Web Interface Issues

**Symptom**: Web interface is inaccessible or shows errors

**Possible Causes and Solutions**:

1. **Service Not Running**:
   ```bash
   systemctl status ubuntu-autoinstall-webhook
   # If not running
   systemctl start ubuntu-autoinstall-webhook
   # Check for errors
   journalctl -u ubuntu-autoinstall-webhook -n 50
   ```

2. **Network/Firewall Issues**:
   ```bash
   # Check if port is listening
   ss -tulpn | grep 8443

   # Test firewall rules
   ufw status
   ```

3. **Certificate Problems**:
   - Verify certificates are valid and not expired
   - Check for certificate path issues

   ```bash
   # Verify certificate
   openssl x509 -in /etc/ubuntu-autoinstall-webhook/certs/server.crt -text -noout

   # Test TLS connection
   openssl s_client -connect localhost:8443
   ```

4. **File Permission Issues**:
   ```bash
   # Check permissions on web files
   find /var/www/ubuntu-autoinstall-webhook -type f -exec ls -l {} \;

   # Fix permissions if needed
   chown -R ubuntu-autoinstall:ubuntu-autoinstall /var/www/ubuntu-autoinstall-webhook
   ```

#### 10.2.3. Database Connectivity Issues

**Symptom**: System reports database connection errors

**Possible Causes and Solutions**:

1. **Database Service Down**:
   ```bash
   # For SQLite
   ls -la /var/lib/ubuntu-autoinstall-webhook/database.sqlite3

   # For CockroachDB
   systemctl status cockroachdb
   ```

2. **Permission Problems**:
   ```bash
   # Check ownership and permissions
   ls -la /var/lib/ubuntu-autoinstall-webhook/

   # Fix if needed
   chown ubuntu-autoinstall:ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/database.sqlite3
   chmod 640 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3
   ```

3. **Connection Configuration**:
   - Verify database settings in config.yaml
   - Check network connectivity to external database

   ```bash
   # Test connection to external database
   nc -zv cockroachdb-host 26257
   ```

4. **Database Corruption**:
   ```bash
   # Check and repair SQLite database
   sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 "PRAGMA integrity_check;"

   # If corrupted, restore from backup
   cp /path/to/backup/database.sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3
   ```

#### 10.2.4. Certificate Issues

**Symptom**: Certificate errors in logs or during installation

**Possible Causes and Solutions**:

1. **Expired Certificates**:
   ```bash
   # Check certificate expiration
   ubuntu-autoinstall-webhook cert list --show-expiry

   # Renew expired certificates
   ubuntu-autoinstall-webhook cert renew --all
   ```

2. **Trust Chain Problems**:
   - Verify CA certificates are properly installed
   - Check that client systems trust the CA

   ```bash
   # Export CA certificate for clients
   ubuntu-autoinstall-webhook cert export --ca --format=pem --output=/tmp/ca.crt
   ```

3. **Certificate Mismatch**:
   - Ensure hostname in certificate matches server name
   - Check Subject Alternative Names (SANs)

   ```bash
   # View certificate details
   openssl x509 -in /etc/ubuntu-autoinstall-webhook/certs/server.crt -text -noout | grep -A1 "Subject Alternative Name"
   ```

### 10.3. Service Specific Issues

#### 10.3.1. File Editor Service

**Common Issues**:

1. **Permission Denied**:
   ```bash
   # Check ownership and permissions
   ls -la /var/www/html/ipxe/boot/

   # Fix permissions
   chown -R ubuntu-autoinstall:ubuntu-autoinstall /var/www/html/ipxe/
   chmod -R 755 /var/www/html/ipxe/
   ```

2. **Leader Election Conflicts**:
   ```bash
   # Check leader status
   ubuntu-autoinstall-webhook service status --component=file-editor

   # Force leader election
   ubuntu-autoinstall-webhook service election --component=file-editor --force
   ```

#### 10.3.2. DNSMasq Watcher

**Common Issues**:

1. **Log File Access**:
   ```bash
   # Check log file accessibility
   ls -la /var/log/dnsmasq.log

   # Grant access if needed
   usermod -a -G adm ubuntu-autoinstall
   ```

2. **DNSMasq Not Logging**:
   ```bash
   # Verify DNSMasq logging is enabled
   grep "log-dhcp" /etc/dnsmasq.conf

   # Add if missing
   echo "log-dhcp" >> /etc/dnsmasq.conf

   # Restart DNSMasq
   systemctl restart dnsmasq
   ```

3. **Missing DHCP Events**:
   - Verify DNSMasq is configured as DHCP server
   - Check that systems send DHCP requests

   ```bash
   # Capture DHCP traffic
   tcpdump -i eth0 -n port 67 or port 68
   ```

#### 10.3.3. Certificate Issuer

**Common Issues**:

1. **CA Not Initialized**:
   ```bash
   # Check CA status
   ubuntu-autoinstall-webhook cert ca-status

   # Initialize CA if needed
   ubuntu-autoinstall-webhook cert init-ca
   ```

2. **CSR Failures**:
   ```bash
   # Review recent CSR failures
   grep "CSR validation failed" /var/log/ubuntu-autoinstall-webhook/cert-issuer.log
   ```

3. **Key Permission Issues**:
   ```bash
   # Check key permissions
   ls -la /var/lib/ubuntu-autoinstall-webhook/certificates/private/

   # Fix if needed
   chmod 600 /var/lib/ubuntu-autoinstall-webhook/certificates/private/*.key
   ```

### 10.4. Boot and Installation Debugging

#### 10.4.1. iPXE Debugging

Enable verbose logging in iPXE scripts:

```
# Add to your iPXE template
set debug:dhcp 1
set debug:proto 1
set console:boot
```

Capture iPXE console output:

```bash
# Setup a netconsole server
nc -u -l 6666 | tee ipxe-debug.log

# In iPXE script
set netconsole:server 192.168.1.100
set netconsole:port 6666
```

#### 10.4.2. Cloud-Init Debugging

Enable verbose logging in cloud-init by adding to user-data:

```yaml
debug:
  verbose: true
output: {all: '| tee -a /var/log/cloud-init-output.log'}
```

Configure cloud-init to report back to webhook server:

```yaml
phone_home:
  url: http://webhook.example.com:8080/api/v1/installations/phone-home
  tries: 10
  post:
    - hostname
    - instance_id
    - pub_key_dsa
    - pub_key_rsa
    - pub_key_ecdsa
```

#### 10.4.3. Live Installation Monitoring

```bash
# Track installation logs in real-time
ubuntu-autoinstall-webhook installation watch --mac=00:11:22:33:44:55

# Capture installation metrics
ubuntu-autoinstall-webhook installation stats --last=24h
```

### 10.5. Advanced Troubleshooting

#### 10.5.1. Service Profiling

Analyze service performance:

```bash
# Enable profiling
ubuntu-autoinstall-webhook debug enable-profiling

# Access profiling data (default: http://localhost:6060/debug/pprof/)
go tool pprof http://localhost:6060/debug/pprof/heap

# Disable profiling
ubuntu-autoinstall-webhook debug disable-profiling
```

#### 10.5.2. Database Query Analysis

```bash
# For CockroachDB
echo "EXPLAIN ANALYZE SELECT * FROM systems WHERE status = 'installing';" | \
  cockroach sql --host=localhost --insecure

# For SQLite with query plan
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 \
  "EXPLAIN QUERY PLAN SELECT * FROM systems WHERE status = 'installing';"
```

#### 10.5.3. Traffic Analysis

Capture and analyze network traffic:

```bash
# Capture PXE boot traffic
tcpdump -i eth0 -n port 67 or port 68 or port 69 or port 4011 -w /tmp/pxe-capture.pcap

# Analyze HTTP installation traffic
tcpdump -i eth0 -n port 8080 -A | grep -i "user-agent\|host:\|get\|post"
```

#### 10.5.4. System Recovery

In case of severe system issues, a recovery mode is available:

```bash
# Start in recovery mode
systemctl stop ubuntu-autoinstall-webhook
ubuntu-autoinstall-webhook start --recovery-mode

# Repair database
ubuntu-autoinstall-webhook db repair

# Reset to factory defaults
ubuntu-autoinstall-webhook recovery factory-reset --keep-certificates

# Exit recovery mode
CTRL+C
systemctl start ubuntu-autoinstall-webhook
```

### 10.6. Getting Support

#### 10.6.1. Generating Support Bundle

Create a comprehensive support bundle for analysis:

```bash
ubuntu-autoinstall-webhook support bundle \
  --include-logs \
  --include-configs \
  --include-db-schema \
  --redact-sensitive \
  --output=/tmp/support-bundle.zip
```

The bundle includes:
- System information
- Service configuration (with sensitive data redacted)
- Recent logs
- Database schema (without data)
- Component status

#### 10.6.2. Community Support

For community support options:
- GitHub Issues: [https://github.com/jdfalk/ubuntu-autoinstall-webhook/issues](https://github.com/jdfalk/ubuntu-autoinstall-webhook/issues)
- Community Forums: [https://discourse.ubuntu.com/c/server/autoinstall](https://discourse.ubuntu.com/c/server/autoinstall)

#### 10.6.3. Commercial Support

No support this software is offered as-is.

## 11. Upgrading and Maintenance

Regular upgrading and maintenance are crucial for keeping your Ubuntu Autoinstall Webhook system secure, performant, and reliable. This section covers procedures for upgrades, migrations, and routine maintenance tasks.

### 11.1. Version Upgrades

#### 11.1.1. Before Upgrading

Before performing any upgrade, follow these preparatory steps:

1. **Create a backup**:
   ```bash
   # Back up all system data
   ubuntu-autoinstall-webhook backup create --full --output=/path/to/backups/
   ```

2. **Review release notes**:
   Check the release documentation for breaking changes, new features, and migration requirements.

3. **Verify system health**:
   ```bash
   # Check current system status
   ubuntu-autoinstall-webhook status check --verbose

   # Resolve any existing issues before upgrading
   ```

4. **Schedule maintenance window**:
   - Schedule downtime during low-usage periods
   - Notify users of planned downtime
   - Consider impact on active installations

#### 11.1.2. Upgrade Procedures

**Method 1: Package Manager (Recommended)**

For systems installed via APT:

```bash
# Update package information
apt update

# View what will be upgraded
apt list --upgradable | grep ubuntu-autoinstall-webhook

# Perform the upgrade
apt upgrade ubuntu-autoinstall-webhook

# Check status after upgrade
systemctl status ubuntu-autoinstall-webhook
```

**Method 2: Docker Container**

For Docker-based deployments:

```bash
# Pull the new image
docker pull ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest

# Stop the current container
docker stop ubuntu-autoinstall-webhook

# Remove the current container (keeping volumes)
docker rm ubuntu-autoinstall-webhook

# Start a new container with the updated image
docker run -d --name ubuntu-autoinstall-webhook \
  --restart unless-stopped \
  -v webhook-data:/var/lib/ubuntu-autoinstall-webhook \
  -v webhook-config:/etc/ubuntu-autoinstall-webhook \
  -p 8080:8080 -p 8443:8443 -p 69:69/udp \
  ghcr.io/jdfalk/ubuntu-autoinstall-webhook:latest
```

**Method 3: Manual Binary Installation**

For systems with manual binary installations:

```bash
# Download the latest release
wget https://github.com/jdfalk/ubuntu-autoinstall-webhook/releases/download/v1.2.3/ubuntu-autoinstall-webhook-1.2.3-linux-amd64.tar.gz

# Extract files
tar -xzf ubuntu-autoinstall-webhook-1.2.3-linux-amd64.tar.gz

# Stop the service
systemctl stop ubuntu-autoinstall-webhook

# Back up the existing binary
cp /usr/local/bin/ubuntu-autoinstall-webhook /usr/local/bin/ubuntu-autoinstall-webhook.bak

# Install the new binary
cp ubuntu-autoinstall-webhook /usr/local/bin/

# Restore permissions
chmod 755 /usr/local/bin/ubuntu-autoinstall-webhook

# Start the service
systemctl start ubuntu-autoinstall-webhook
```

#### 11.1.3. Post-Upgrade Tasks

After upgrading, perform these verification steps:

1. **Check service status**:
   ```bash
   systemctl status ubuntu-autoinstall-webhook
   ```

2. **Verify version**:
   ```bash
   ubuntu-autoinstall-webhook --version
   ```

3. **Run database migrations** (if needed):
   ```bash
   ubuntu-autoinstall-webhook db migrate
   ```

4. **Check logs for errors**:
   ```bash
   journalctl -u ubuntu-autoinstall-webhook -n 50
   ```

5. **Verify web interface functionality**:
   - Log in to the web interface
   - Check dashboard and key features
   - Verify template generation works

6. **Test a basic installation** (if possible):
   - Start a test installation on a non-production system
   - Verify the complete workflow functions as expected

### 11.2. Database Maintenance

#### 11.2.1. Database Migrations

When upgrading between versions, database schema migrations may be required:

```bash
# Check migration status
ubuntu-autoinstall-webhook db migration-status

# Run pending migrations
ubuntu-autoinstall-webhook db migrate
```

For manual control over migrations:

```bash
# Apply specific migration
ubuntu-autoinstall-webhook db migrate --to-version=5

# Roll back a migration (if needed)
ubuntu-autoinstall-webhook db rollback --to-version=4
```

#### 11.2.2. Database Optimization

Perform regular database optimization for better performance:

**For SQLite**:
```bash
# Optimize SQLite database
ubuntu-autoinstall-webhook db optimize

# Or manually
sqlite3 /var/lib/ubuntu-autoinstall-webhook/database.sqlite3 "VACUUM; ANALYZE;"
```

**For CockroachDB**:
```bash
# Run statistics update
cockroach sql --execute="SET CLUSTER SETTING sql.stats.automatic_collection.enabled = true;"

# Manual statistics update
cockroach sql --execute="CREATE STATISTICS systems_stats FROM systems;"
cockroach sql --execute="CREATE STATISTICS installations_stats FROM installations;"
```

#### 11.2.3. Data Cleanup

Implement regular data cleanup to prevent database growth:

```bash
# Clean up old installation data
ubuntu-autoinstall-webhook maintenance cleanup --older-than=90d --status=completed

# Remove expired tokens
ubuntu-autoinstall-webhook maintenance cleanup-tokens --expired-only

# Purge audit logs (if regulatory compliance allows)
ubuntu-autoinstall-webhook maintenance cleanup-audit-logs --older-than=365d
```

### 11.3. Routine Maintenance Tasks

#### 11.3.1. Certificate Rotation

Regularly rotate certificates to maintain security:

```bash
# Check certificates nearing expiration
ubuntu-autoinstall-webhook cert list --expiring-within=30d

# Renew expiring certificates
ubuntu-autoinstall-webhook cert renew --expiring-within=30d

# Force renewal of specific certificate
ubuntu-autoinstall-webhook cert renew --id=cert-uuid
```

#### 11.3.2. Log Rotation

Though the system configures logrotate automatically, verify the configuration:

```bash
cat /etc/logrotate.d/ubuntu-autoinstall-webhook
```

To manually rotate logs:

```bash
logrotate -f /etc/logrotate.d/ubuntu-autoinstall-webhook
```

#### 11.3.3. File System Maintenance

Maintain the file system organization:

```bash
# Clean up temporary files
ubuntu-autoinstall-webhook maintenance cleanup-temp-files

# Verify file permissions
ubuntu-autoinstall-webhook maintenance verify-permissions --fix

# Check for and remove orphaned files
ubuntu-autoinstall-webhook maintenance find-orphaned-files --remove
```

#### 11.3.4. User Account Maintenance

Regularly audit user accounts:

```bash
# List all users
ubuntu-autoinstall-webhook users list

# Find inactive users
ubuntu-autoinstall-webhook users list --inactive-days=90

# Disable inactive accounts
ubuntu-autoinstall-webhook users disable --inactive-days=90

# Reset failed login counters
ubuntu-autoinstall-webhook users reset-login-attempts --all
```

### 11.4. Configuration Management

#### 11.4.1. Configuration Backups

Back up configuration before making changes:

```bash
# Export current configuration
ubuntu-autoinstall-webhook config export --output=/path/to/config-backup.yaml
```

#### 11.4.2. Configuration Versioning

Track configuration changes with version control:

```bash
# Initialize git repository for configuration
cd /etc/ubuntu-autoinstall-webhook
git init
git add .
git commit -m "Initial configuration"

# After changes
git add -u
git commit -m "Updated network settings"
```

#### 11.4.3. Configuration Validation

Validate configuration changes before applying:

```bash
# Validate configuration file
ubuntu-autoinstall-webhook config validate --file=/path/to/new-config.yaml

# Apply validated configuration
ubuntu-autoinstall-webhook config apply --file=/path/to/new-config.yaml
```

### 11.5. Disaster Recovery Testing

#### 11.5.1. Recovery Drills

Periodically test your disaster recovery procedures:

```bash
# Create a test environment
ubuntu-autoinstall-webhook test setup-dr-environment

# Restore from backup to test environment
ubuntu-autoinstall-webhook backup restore \
  --file=/path/to/backup.tar.gz \
  --target-dir=/tmp/dr-test

# Validate restored data
ubuntu-autoinstall-webhook test validate-dr-restore
```

#### 11.5.2. Failover Testing

For high-availability setups, test failover mechanisms:

```bash
# Test leader election by stopping the current leader
systemctl stop ubuntu-autoinstall-webhook@leader

# Verify standby instance takes over
ubuntu-autoinstall-webhook cluster status

# Restore normal operation
systemctl start ubuntu-autoinstall-webhook@leader
```

### 11.6. System Monitoring

#### 11.6.1. Monitoring Health Checks

Configure regular health checks:

```bash
# Create a health check script
cat > /usr/local/bin/webhook-health-check.sh << 'EOF'
#!/bin/bash
curl -s http://localhost:8080/health | grep -q '"status":"healthy"'
exit $?
EOF

chmod +x /usr/local/bin/webhook-health-check.sh

# Add to crontab
(crontab -l 2>/dev/null; echo "*/5 * * * * /usr/local/bin/webhook-health-check.sh || systemctl restart ubuntu-autoinstall-webhook") | crontab -
```

#### 11.6.2. Performance Baseline

Establish performance baselines for future comparison:

```bash
# Capture baseline metrics
ubuntu-autoinstall-webhook benchmark run --output=/var/lib/ubuntu-autoinstall-webhook/baselines/

# Compare current performance with baseline
ubuntu-autoinstall-webhook benchmark compare \
  --baseline=/var/lib/ubuntu-autoinstall-webhook/baselines/baseline-2023-01-15.json \
  --current
```

### 11.7. Planning for Major Upgrades

#### 11.7.1. Upgrade Impact Assessment

Before major version upgrades, assess potential impacts:

1. Review the breaking changes documentation
2. Analyze custom templates for compatibility
3. Check API integrations for deprecated endpoints
4. Plan for database schema changes
5. Test the upgrade process in a staging environment

#### 11.7.2. Rollback Planning

Always have a rollback plan for major upgrades:

1. Document the rollback procedure specific to the version
2. Test the rollback procedure in a staging environment
3. Ensure database backups are backward compatible
4. Prepare to revert configuration changes if needed
5. Plan for data migration rollback (if applicable)

## 12. Advanced Configuration

This section covers advanced configuration options and customizations for the Ubuntu Autoinstall Webhook system.

### 12.1. Customizing Templates

#### 12.1.1. Template Engine Overview

The Ubuntu Autoinstall Webhook system uses a powerful templating engine to generate installation configurations. The template engine:

- Supports multiple output formats (cloud-init user-data, network-config, meta-data, iPXE scripts)
- Uses Go templates as the base syntax
- Provides custom functions for system-specific operations
- Supports inheritance and composition
- Allows for complex logic and conditionals

#### 12.1.2. Creating Custom Template Functions

You can extend the template engine with custom functions:

```bash
# Create a plugin directory if it doesn't exist
mkdir -p /etc/ubuntu-autoinstall-webhook/plugins/template-functions

# Create a simple template function plugin
cat > /etc/ubuntu-autoinstall-webhook/plugins/template-functions/custom.go << 'EOF'
package main

import (
    "strings"
)

// CustomUppercase - Example custom function that converts text to uppercase
func CustomUppercase(input string) string {
    return strings.ToUpper(input)
}

// Export is required - maps function names to implementations
var Export = map[string]interface{}{
    "customUpper": CustomUppercase,
}
EOF

# Build the plugin
cd /etc/ubuntu-autoinstall-webhook/plugins/template-functions
go build -buildmode=plugin -o custom.so custom.go

# Restart the service to load the new plugin
systemctl restart ubuntu-autoinstall-webhook
```

Now you can use your custom function in templates:

```
Hostname: {{ customUpper .System.Hostname }}
```

#### 12.1.3. Advanced Template Functions

The system provides several advanced template functions:

**Network Functions**:
```
{{ cidrHost "192.168.1.0/24" 5 }}           // Returns: 192.168.1.5
{{ cidrSubnet "10.0.0.0/16" 8 10 }}         // Returns: 10.0.10.0/24
{{ cidrNetmask "10.0.0.0/24" }}             // Returns: 255.255.255.0
```

**Cryptographic Functions**:
```
{{ sha256sum "data" }}                      // Returns SHA256 hash
{{ genPassword 16 }}                         // Generates random 16-char password
{{ genSSHKey "rsa" }}                        // Generates SSH key pair
```

**Data Manipulation**:
```
{{ toJSON .System }}                         // Convert to JSON
{{ fromYAML $yamlData }}                     // Parse YAML to object
{{ indent 4 $content }}                      // Indent content by 4 spaces
```

**Conditional Logic**:
```
{{ if eq .System.Role "webserver" }}
packages:
  - apache2
  - php
{{ else if eq .System.Role "database" }}
packages:
  - mysql-server
{{ else }}
packages:
  - basic-utils
{{ end }}
```

#### 12.1.4. Template Inheritance

Templates can inherit from other templates using the `extends` directive:

```yaml
# Base template (base.yaml)
#extends: none
apt:
  primary:
    - arches: [default]
      uri: http://archive.ubuntu.com/ubuntu
packages:
  - openssh-server
  - cloud-init

# Derived template (webserver.yaml)
#extends: base.yaml
packages:
  - openssh-server
  - cloud-init
  - apache2
  - php
```

When a template extends another, it inherits all settings, with the child template's settings overriding the parent's when conflicts occur.

### 12.2. API Customization

#### 12.2.1. Custom API Endpoints

Create custom API endpoints for specific needs:

```bash
# Create a plugin directory if it doesn't exist
mkdir -p /etc/ubuntu-autoinstall-webhook/plugins/api-endpoints

# Create a custom API endpoint plugin
cat > /etc/ubuntu-autoinstall-webhook/plugins/api-endpoints/custom-status.go << 'EOF'
package main

import (
    "encoding/json"
    "net/http"
    "time"
)

// StatusResponse defines the structure of our response
type StatusResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Custom    string    `json:"custom"`
}

// Handler is the entrypoint for our custom endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
    // Only allow GET requests
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Create our response
    resp := StatusResponse{
        Status:    "operational",
        Timestamp: time.Now(),
        Custom:    "Custom endpoint is working!",
    }

    // Set content type and encode response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// Route tells the system where to mount this handler
var Route = "/api/v1/custom/status"
EOF

# Build the plugin
cd /etc/ubuntu-autoinstall-webhook/plugins/api-endpoints
go build -buildmode=plugin -o custom-status.so custom-status.go

# Restart the service to load the new endpoint
systemctl restart ubuntu-autoinstall-webhook
```

Access your custom endpoint at `http://localhost:8080/api/v1/custom/status`

#### 12.2.2. API Rate Limiting

Configure custom rate limiting for API endpoints:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
api:
  rate_limiting:
    enabled: true
    requests_per_minute: 60
    burst: 20

    # Per-endpoint overrides
    endpoint_overrides:
      - path: "/api/v1/systems"
        method: "GET"
        requests_per_minute: 120

      - path: "/api/v1/installations"
        method: "POST"
        requests_per_minute: 30
```

#### 12.2.3. Custom API Authentication

Implement specialized API authentication methods:

```bash
# Create a plugin directory if it doesn't exist
mkdir -p /etc/ubuntu-autoinstall-webhook/plugins/auth-providers

# Create a custom authentication plugin
cat > /etc/ubuntu-autoinstall-webhook/plugins/auth-providers/hmac-auth.go << 'EOF'
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "strings"
    "time"
)

// Authenticate implements HMAC-based authentication
func Authenticate(r *http.Request) (bool, string, map[string]interface{}) {
    // Extract authorization header
    authHeader := r.Header.Get("Authorization")
    if !strings.HasPrefix(authHeader, "HMAC ") {
        return false, "", nil
    }

    parts := strings.Split(strings.TrimPrefix(authHeader, "HMAC "), ":")
    if len(parts) != 2 {
        return false, "", nil
    }

    apiKeyID := parts[0]
    signatureProvided := parts[1]

    // In a real implementation, look up the secret based on apiKeyID
    secret := "your-secret-key" // This would come from a secure store

    // Create message to sign (timestamp + method + path)
    timestamp := r.Header.Get("X-Timestamp")
    message := timestamp + r.Method + r.URL.Path

    // Calculate HMAC signature
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(message))
    signatureExpected := hex.EncodeToString(h.Sum(nil))

    // Compare signatures
    if signatureProvided == signatureExpected {
        // Authentication successful, return user ID and any custom claims
        return true, apiKeyID, map[string]interface{}{
            "authenticated_at": time.Now().Unix(),
            "method": "hmac",
        }
    }

    return false, "", nil
}

// Name identifies this authentication provider
var Name = "hmac"
EOF

# Build the plugin
cd /etc/ubuntu-autoinstall-webhook/plugins/auth-providers
go build -buildmode=plugin -o hmac-auth.so hmac-auth.go

# Configure the system to use the plugin
cat >> /etc/ubuntu-autoinstall-webhook/config.yaml << 'EOF'
auth:
  providers:
    - name: local
      enabled: true
    - name: hmac
      enabled: true
      priority: 10
EOF

# Restart the service to load the new auth provider
systemctl restart ubuntu-autoinstall-webhook
```

### 12.3. Advanced Networking

#### 12.3.1. VLAN and Network Segmentation

Configure the system to work with VLANs:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
networking:
  vlans:
    enabled: true
    interfaces:
      - name: "eth0.100"
        parent: "eth0"
        id: 100
        address: "10.100.0.5/24"
        gateway: "10.100.0.1"
        service_bindings: ["web"]

      - name: "eth0.200"
        parent: "eth0"
        id: 200
        address: "10.200.0.5/24"
        gateway: "10.200.0.1"
        service_bindings: ["pxe"]
```

Create corresponding VLAN interfaces on the system:

```bash
# Load 8021q kernel module if not already loaded
modprobe 8021q

# Create VLAN interfaces
ip link add link eth0 name eth0.100 type vlan id 100
ip link add link eth0 name eth0.200 type vlan id 200

# Configure addresses
ip addr add 10.100.0.5/24 dev eth0.100
ip addr add 10.200.0.5/24 dev eth0.200

# Bring interfaces up
ip link set eth0.100 up
ip link set eth0.200 up

# Configure default routes if needed
ip route add 10.100.0.0/24 via 10.100.0.1 dev eth0.100
ip route add 10.200.0.0/24 via 10.200.0.1 dev eth0.200
```

Make the configuration persistent with Netplan:

```yaml
# /etc/netplan/60-vlans.yaml
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: true
  vlans:
    eth0.100:
      id: 100
      link: eth0
      addresses: [10.100.0.5/24]
      routes:
        - to: 10.100.0.0/24
          via: 10.100.0.1
    eth0.200:
      id: 200
      link: eth0
      addresses: [10.200.0.5/24]
      routes:
        - to: 10.200.0.0/24
          via: 10.200.0.1
```

Apply the configuration:

```bash
netplan apply
```

#### 12.3.2. Network Bonding

Configure network bonding for redundancy:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
networking:
  bonding:
    enabled: true
    interfaces:
      - name: "bond0"
        mode: "802.3ad"
        slaves: ["eth0", "eth1"]
        address: "192.168.1.5/24"
        gateway: "192.168.1.1"
        options:
          miimon: 100
          lacp_rate: fast
```

Set up bonding on the system:

```bash
# Load bonding kernel module
modprobe bonding

# Create bonding configuration
cat > /etc/modprobe.d/bonding.conf << 'EOF'
alias bond0 bonding
options bond0 mode=802.3ad miimon=100 lacp_rate=1
EOF

# Configure with Netplan
cat > /etc/netplan/50-bonding.yaml << 'EOF'
network:
  version: 2
  bonds:
    bond0:
      interfaces: [eth0, eth1]
      parameters:
        mode: 802.3ad
        mii-monitor-interval: 100
        lacp-rate: fast
      addresses: [192.168.1.5/24]
      routes:
        - to: default
          via: 192.168.1.1
EOF

# Apply configuration
netplan apply
```

#### 12.3.3. IPv6 Configuration

Enable IPv6 support:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
networking:
  ipv6:
    enabled: true
    address: "2001:db8::5/64"
    gateway: "2001:db8::1"
```

Configure the system for IPv6:

```bash
# Add IPv6 address to interface
ip -6 addr add 2001:db8::5/64 dev eth0

# Add default route
ip -6 route add default via 2001:db8::1 dev eth0

# Make persistent with Netplan
cat > /etc/netplan/70-ipv6.yaml << 'EOF'
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: true
      addresses: [2001:db8::5/64]
      routes:
        - to: ::/0
          via: 2001:db8::1
EOF

netplan apply
```

### 12.4. Clustering and High Availability

#### 12.4.1. Basic Clustering Setup

Configure a basic cluster with multiple instances:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml on each node
cluster:
  enabled: true
  node_name: "node1"  # Unique for each node
  nodes:
    - name: "node1"
      address: "192.168.1.10"
      role: "primary"
    - name: "node2"
      address: "192.168.1.11"
      role: "secondary"
    - name: "node3"
      address: "192.168.1.12"
      role: "secondary"

  quorum:
    required: true
    min_nodes: 2
```

#### 12.4.2. Distributed File System

Configure a shared storage for installation files:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
storage:
  type: "shared"
  shared:
    type: "nfs"
    mount_point: "/var/lib/ubuntu-autoinstall-webhook/shared"
    server: "nfs.example.com"
    export: "/exports/ubuntu-autoinstall"
    options: "rw,sync,no_subtree_check"
```

Mount the NFS share:

```bash
# Create mount point
mkdir -p /var/lib/ubuntu-autoinstall-webhook/shared

# Mount NFS share
mount -t nfs nfs.example.com:/exports/ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/shared

# Make persistent by adding to /etc/fstab
echo "nfs.example.com:/exports/ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/shared nfs rw,sync,no_subtree_check 0 0" >> /etc/fstab
```

#### 12.4.3. Load Balancer Configuration

Example NGINX load balancer configuration:

```nginx
# /etc/nginx/sites-available/webhook-lb
upstream webhook_backend {
    # ip_hash ensures the same client always hits the same server
    ip_hash;

    server 192.168.1.10:8443;
    server 192.168.1.11:8443;
    server 192.168.1.12:8443;

    # Health check parameters
    check interval=5000 rise=2 fall=3 timeout=4000;
}

server {
    listen 443 ssl;
    server_name webhook.example.com;

    ssl_certificate /etc/nginx/ssl/webhook.crt;
    ssl_certificate_key /etc/nginx/ssl/webhook.key;
    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass https://webhook_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable and restart NGINX:

```bash
ln -s /etc/nginx/sites-available/webhook-lb /etc/nginx/sites-enabled/
systemctl restart nginx
```

### 12.5. Integration with External Systems

#### 12.5.1. CMDB Integration

Configure integration with a Configuration Management Database:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
integrations:
  cmdb:
    enabled: true
    type: "servicenow"
    url: "https://example.service-now.com/api"
    auth:
      type: "basic"
      username: "integration-user"
      password: "integration-password"

    # Field mappings
    mappings:
      system_id: "sys_id"
      hostname: "name"
      mac_address: "mac_address"
      ip_address: "ip_address"
      status: "status"

    # Synchronization settings
    sync:
      enabled: true
      interval_minutes: 60
      direction: "bidirectional"
```

#### 12.5.2. Monitoring System Integration

Configure integration with monitoring systems:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
integrations:
  monitoring:
    - type: "prometheus"
      enabled: true
      endpoint: "/metrics"
      labels:
        environment: "production"
        service: "ubuntu-autoinstall"

    - type: "nagios"
      enabled: true
      nrdp_url: "https://nagios.example.com/nrdp/"
      token: "your-nagios-token"
      host_name: "ubuntu-autoinstall"
      interval_minutes: 5
```

#### 12.5.3. Custom Webhook Notifications

Configure custom webhook notifications:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
notifications:
  webhooks:
    - name: "slack-alerts"
      url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXX"
      events:
        - system.discovered
        - installation.started
        - installation.completed
        - installation.failed
      headers:
        Content-Type: "application/json"
      template: |
        {
          "text": "Event: {{ .Event.Type }}",
          "blocks": [
            {
              "type": "section",
              "text": {
                "type": "mrkdwn",
                "text": "*{{ .Event.Type }}*\nSystem: {{ .System.Hostname }}\nMAC: {{ .System.MACAddress }}\nTimestamp: {{ .Event.Timestamp }}"
              }
            }
          ]
        }

    - name: "teams-notification"
      url: "https://outlook.office.com/webhook/..."
      events:
        - installation.failed
      headers:
        Content-Type: "application/json"
      template: |
        {
          "@type": "MessageCard",
          "@context": "https://schema.org/extensions",
          "summary": "Installation Failed",
          "sections": [
            {
              "activityTitle": "Installation Failed for {{ .System.Hostname }}",
              "facts": [
                { "name": "System", "value": "{{ .System.Hostname }}" },
                { "name": "MAC Address", "value": "{{ .System.MACAddress }}" },
                { "name": "Error", "value": "{{ .Event.Details.Error }}" },
                { "name": "Time", "value": "{{ .Event.Timestamp }}" }
              ]
            }
          ]
        }
```

### 12.6. Advanced Storage Configuration

#### 12.6.1. Object Storage for Installation Files

Configure S3-compatible object storage:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
storage:
  installation_files:
    type: "s3"
    s3:
      endpoint: "https://s3.amazonaws.com"
      region: "us-east-1"
      bucket: "ubuntu-autoinstall-files"
      prefix: "files/"
      access_key: "YOUR_ACCESS_KEY"
      secret_key: "YOUR_SECRET_KEY"
      use_path_style: false
```

#### 12.6.2. Database on External Storage

Configure external database storage:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
database:
  type: "cockroach"
  cockroach:
    host: "db.example.com"
    port: 26257
    database: "ubuntu_autoinstall"
    user: "ubuntu_autoinstall"
    password: "your_password"
    ssl:
      mode: "verify-full"
      ca_cert: "/etc/ubuntu-autoinstall-webhook/certs/ca.crt"
      client_cert: "/etc/ubuntu-autoinstall-webhook/certs/client.crt"
      client_key: "/etc/ubuntu-autoinstall-webhook/certs/client.key"
```

#### 12.6.3. Backup to Remote Storage

Configure backups to remote storage:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
backup:
  enabled: true
  schedule: "0 2 * * *"  # 2 AM daily
  retention:
    days: 30
    count: 10

  storage:
    type: "sftp"
    sftp:
      host: "backup.example.com"
      port: 22
      user: "backup-user"
      key_file: "/etc/ubuntu-autoinstall-webhook/keys/backup-key"
      path: "/backups/ubuntu-autoinstall/"
```

### 12.7. Custom Authentication and Authorization

#### 12.7.1. Custom LDAP Configuration

Configure advanced LDAP integration:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
auth:
  ldap:
    enabled: true
    servers:
      - uri: "ldaps://ldap1.example.com"
        priority: 1
      - uri: "ldaps://ldap2.example.com"
        priority: 2

    bind_dn: "cn=service-account,ou=ServiceAccounts,dc=example,dc=com"
    bind_password: "your_password"

    user:
      base_dn: "ou=Users,dc=example,dc=com"
      filter: "(&(objectClass=person)(sAMAccountName=%s))"
      username_attribute: "sAMAccountName"
      name_attribute: "displayName"
      email_attribute: "mail"

    group:
      base_dn: "ou=Groups,dc=example,dc=com"
      filter: "(&(objectClass=group)(member=%s))"
      name_attribute: "cn"

    role_mappings:
      - ldap_group: "AutoinstallAdmins"
        role: "admin"
      - ldap_group: "AutoinstallOperators"
        role: "operator"
      - ldap_group: "AutoinstallUsers"
        role: "installer"
```

#### 12.7.2. OpenID Connect Configuration

Configure authentication with OpenID Connect providers:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
auth:
  oidc:
    enabled: true
    providers:
      - name: "google"
        client_id: "your_client_id"
        client_secret: "your_client_secret"
        discovery_url: "https://accounts.google.com/.well-known/openid-configuration"
        scopes: ["openid", "profile", "email"]
        claims:
          username: "email"
          name: "name"
          email: "email"

      - name: "azure_ad"
        client_id: "your_client_id"
        client_secret: "your_client_secret"
        discovery_url: "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration"
        scopes: ["openid", "profile", "email"]
        claims:
          username: "preferred_username"
          name: "name"
          email: "email"

    role_mappings:
      - claim: "groups"
        values:
          - value: "autoinstall_admins"
            role: "admin"
          - value: "autoinstall_operators"
            role: "operator"
```

#### 12.7.3. Custom Authorization Rules

Configure fine-grained authorization rules:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
authorization:
  rules:
    - name: "ReadOnlyForNonProduction"
      description: "Read-only access to non-production systems"
      match:
        roles:
          - "viewer"
        resources:
          - type: "system"
            attributes:
              environment: "non-production"
      permissions:
        - "read"

    - name: "FullAccessToOwnedSystems"
      description: "Full access to systems created by the user"
      match:
        resources:
          - type: "system"
            attributes:
              created_by: "${user.id}"
      permissions:
        - "read"
        - "write"
        - "delete"
        - "install"
```

### 12.8. Advanced Logging

#### 12.8.1. Structured Logging Configuration

Configure advanced structured logging:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
logging:
  structured:
    enabled: true
    format: "json"

  levels:
    default: "info"
    components:
      webserver: "info"
      database: "warn"
      file_editor: "info"
      dnsmasq_watcher: "debug"
      certificate_issuer: "info"

  output:
    console:
      enabled: true
      format: "text"
    file:
      enabled: true
      format: "json"
      path: "/var/log/ubuntu-autoinstall-webhook/webhook.log"
    syslog:
      enabled: true
      format: "json"
      facility: "local0"
```

#### 12.8.2. Remote Logging Configuration

Configure logging to remote systems:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
logging:
  remote:
    - type: "elasticsearch"
      enabled: true
      url: "https://elasticsearch.example.com:9200"
      index: "ubuntu-autoinstall-logs"
      username: "elastic"
      password: "your_password"
      batch_size: 100
      flush_interval_seconds: 5

    - type: "fluentd"
      enabled: true
      host: "fluentd.example.com"
      port: 24224
      tag: "ubuntu-autoinstall"

    - type: "graylog"
      enabled: true
      server: "graylog.example.com"
      port: 12201
      protocol: "tcp"
      tls:
        enabled: true
        skip_verify: false
```

#### 12.8.3. Log Correlation

Configure distributed tracing:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
logging:
  tracing:
    enabled: true
    type: "opentelemetry"
    service_name: "ubuntu-autoinstall-webhook"

    opentelemetry:
      endpoint: "http://jaeger.example.com:14268/api/traces"
      propagation: "w3c"
      sample_ratio: 0.1
```

## 13. Appendices

### 13.1. Configuration Reference

#### 13.1.1. Complete Configuration Schema

Below is a reference for all configuration options available in `/etc/ubuntu-autoinstall-webhook/config.yaml`:

```yaml
# General settings
general:
  # Application name - for display in UI and logs
  name: "Ubuntu Autoinstall Webhook"
  # Default environment name
  environment: "production"
  # Data directory path
  data_dir: "/var/lib/ubuntu-autoinstall-webhook"
  # Temporary file directory
  temp_dir: "/tmp/ubuntu-autoinstall-webhook"

# Webserver settings
webserver:
  # HTTP service settings
  http:
    # Enable HTTP service
    enabled: true
    # Listening address (blank for all interfaces)
    address: ""
    # Listening port
    port: 8080
    # Redirect to HTTPS
    redirect_to_https: true

  # HTTPS service settings
  https:
    # Enable HTTPS service
    enabled: true
    # Listening address (blank for all interfaces)
    address: ""
    # Listening port
    port: 8443
    # TLS certificate file
    cert_file: "/etc/ubuntu-autoinstall-webhook/certs/server.crt"
    # TLS key file
    key_file: "/etc/ubuntu-autoinstall-webhook/certs/server.key"

  # TFTP service settings (for PXE)
  tftp:
    # Enable TFTP service
    enabled: true
    # Listening address (blank for all interfaces)
    address: ""
    # Root directory for TFTP files
    root_dir: "/var/lib/ubuntu-autoinstall-webhook/tftp"

  # Web UI settings
  ui:
    # Enable web UI
    enabled: true
    # Page title
    title: "Ubuntu Autoinstall Webhook"
    # Logo URL (relative to webroot)
    logo: "/static/logo.png"
    # Theme (light, dark, auto)
    theme: "auto"
    # Session timeout in minutes
    session_timeout_minutes: 60

# Database settings
database:
  # Database type (sqlite3, cockroach)
  type: "sqlite3"

  # SQLite settings
  sqlite3:
    # Database file path
    path: "/var/lib/ubuntu-autoinstall-webhook/database.sqlite3"
    # Use WAL mode for better performance
    wal: true

  # CockroachDB settings
  cockroach:
    # Hosts (comma-separated for multiple)
    hosts: "localhost"
    # Port
    port: 26257
    # Database name
    name: "ubuntu_autoinstall"
    # Username
    user: "ubuntu_autoinstall"
    # Password
    password: "change-me"
    # SSL mode (disable, require, verify-ca, verify-full)
    ssl_mode: "verify-full"
    # CA certificate for verification
    ca_cert: "/etc/ubuntu-autoinstall-webhook/certs/ca.crt"
    # Client certificate
    client_cert: "/etc/ubuntu-autoinstall-webhook/certs/client.crt"
    # Client key
    client_key: "/etc/ubuntu-autoinstall-webhook/certs/client.key"
    # Connection pool settings
    max_open_conns: 25
    max_idle_conns: 10
    conn_max_lifetime_minutes: 60

# File editor service settings
file_editor:
  # Root directory for installation files
  root_dir: "/var/www/html/ipxe"
  # Enable file system write operations
  write_enabled: true
  # Temporary file directory
  temp_dir: "/tmp/ubuntu-autoinstall-webhook/file-editor"
  # File permissions
  permissions:<!--
## 12. Advanced Configuration

This section covers advanced configuration options and customizations for the Ubuntu Autoinstall Webhook system.

### 12.1. Customizing Templates

#### 12.1.1. Template Engine Overview

The Ubuntu Autoinstall Webhook system uses a powerful templating engine to generate installation configurations. The template engine:

- Supports multiple output formats (cloud-init user-data, network-config, meta-data, iPXE scripts)
- Uses Go templates as the base syntax
- Provides custom functions for system-specific operations
- Supports inheritance and composition
- Allows for complex logic and conditionals

#### 12.1.2. Creating Custom Template Functions

You can extend the template engine with custom functions:

```bash
# Create a plugin directory if it doesn't exist
mkdir -p /etc/ubuntu-autoinstall-webhook/plugins/template-functions

# Create a simple template function plugin
cat > /etc/ubuntu-autoinstall-webhook/plugins/template-functions/custom.go << 'EOF'
package main

import (
    "strings"
)

// CustomUppercase - Example custom function that converts text to uppercase
func CustomUppercase(input string) string {
    return strings.ToUpper(input)
}

// Export is required - maps function names to implementations
var Export = map[string]interface{}{
    "customUpper": CustomUppercase,
}
EOF

# Build the plugin
cd /etc/ubuntu-autoinstall-webhook/plugins/template-functions
go build -buildmode=plugin -o custom.so custom.go

# Restart the service to load the new plugin
systemctl restart ubuntu-autoinstall-webhook
```

Now you can use your custom function in templates:

```
Hostname: {{ customUpper .System.Hostname }}
```

#### 12.1.3. Advanced Template Functions

The system provides several advanced template functions:

**Network Functions**:
```
{{ cidrHost "192.168.1.0/24" 5 }}           // Returns: 192.168.1.5
{{ cidrSubnet "10.0.0.0/16" 8 10 }}         // Returns: 10.0.10.0/24
{{ cidrNetmask "10.0.0.0/24" }}             // Returns: 255.255.255.0
```

**Cryptographic Functions**:
```
{{ sha256sum "data" }}                      // Returns SHA256 hash
{{ genPassword 16 }}                         // Generates random 16-char password
{{ genSSHKey "rsa" }}                        // Generates SSH key pair
```

**Data Manipulation**:
```
{{ toJSON .System }}                         // Convert to JSON
{{ fromYAML $yamlData }}                     // Parse YAML to object
{{ indent 4 $content }}                      // Indent content by 4 spaces
```

**Conditional Logic**:
```
{{ if eq .System.Role "webserver" }}
packages:
  - apache2
  - php
{{ else if eq .System.Role "database" }}
packages:
  - mysql-server
{{ else }}
packages:
  - basic-utils
{{ end }}
```

#### 12.1.4. Template Inheritance

Templates can inherit from other templates using the `extends` directive:

```yaml
# Base template (base.yaml)
#extends: none
apt:
  primary:
    - arches: [default]
      uri: http://archive.ubuntu.com/ubuntu
packages:
  - openssh-server
  - cloud-init

# Derived template (webserver.yaml)
#extends: base.yaml
packages:
  - openssh-server
  - cloud-init
  - apache2
  - php
```

When a template extends another, it inherits all settings, with the child template's settings overriding the parent's when conflicts occur.

### 12.2. API Customization

#### 12.2.1. Custom API Endpoints

Create custom API endpoints for specific needs:

```bash
# Create a plugin directory if it doesn't exist
mkdir -p /etc/ubuntu-autoinstall-webhook/plugins/api-endpoints

# Create a custom API endpoint plugin
cat > /etc/ubuntu-autoinstall-webhook/plugins/api-endpoints/custom-status.go << 'EOF'
package main

import (
    "encoding/json"
    "net/http"
    "time"
)

// StatusResponse defines the structure of our response
type StatusResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Custom    string    `json:"custom"`
}

// Handler is the entrypoint for our custom endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
    // Only allow GET requests
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Create our response
    resp := StatusResponse{
        Status:    "operational",
        Timestamp: time.Now(),
        Custom:    "Custom endpoint is working!",
    }

    // Set content type and encode response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// Route tells the system where to mount this handler
var Route = "/api/v1/custom/status"
EOF

# Build the plugin
cd /etc/ubuntu-autoinstall-webhook/plugins/api-endpoints
go build -buildmode=plugin -o custom-status.so custom-status.go

# Restart the service to load the new endpoint
systemctl restart ubuntu-autoinstall-webhook
```

Access your custom endpoint at `http://localhost:8080/api/v1/custom/status`

#### 12.2.2. API Rate Limiting

Configure custom rate limiting for API endpoints:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
api:
  rate_limiting:
    enabled: true
    requests_per_minute: 60
    burst: 20

    # Per-endpoint overrides
    endpoint_overrides:
      - path: "/api/v1/systems"
        method: "GET"
        requests_per_minute: 120

      - path: "/api/v1/installations"
        method: "POST"
        requests_per_minute: 30
```

#### 12.2.3. Custom API Authentication

Implement specialized API authentication methods:

```bash
# Create a plugin directory if it doesn't exist
mkdir -p /etc/ubuntu-autoinstall-webhook/plugins/auth-providers

# Create a custom authentication plugin
cat > /etc/ubuntu-autoinstall-webhook/plugins/auth-providers/hmac-auth.go << 'EOF'
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "strings"
    "time"
)

// Authenticate implements HMAC-based authentication
func Authenticate(r *http.Request) (bool, string, map[string]interface{}) {
    // Extract authorization header
    authHeader := r.Header.Get("Authorization")
    if !strings.HasPrefix(authHeader, "HMAC ") {
        return false, "", nil
    }

    parts := strings.Split(strings.TrimPrefix(authHeader, "HMAC "), ":")
    if len(parts) != 2 {
        return false, "", nil
    }

    apiKeyID := parts[0]
    signatureProvided := parts[1]

    // In a real implementation, look up the secret based on apiKeyID
    secret := "your-secret-key" // This would come from a secure store

    // Create message to sign (timestamp + method + path)
    timestamp := r.Header.Get("X-Timestamp")
    message := timestamp + r.Method + r.URL.Path

    // Calculate HMAC signature
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(message))
    signatureExpected := hex.EncodeToString(h.Sum(nil))

    // Compare signatures
    if signatureProvided == signatureExpected {
        // Authentication successful, return user ID and any custom claims
        return true, apiKeyID, map[string]interface{}{
            "authenticated_at": time.Now().Unix(),
            "method": "hmac",
        }
    }

    return false, "", nil
}

// Name identifies this authentication provider
var Name = "hmac"
EOF

# Build the plugin
cd /etc/ubuntu-autoinstall-webhook/plugins/auth-providers
go build -buildmode=plugin -o hmac-auth.so hmac-auth.go

# Configure the system to use the plugin
cat >> /etc/ubuntu-autoinstall-webhook/config.yaml << 'EOF'
auth:
  providers:
    - name: local
      enabled: true
    - name: hmac
      enabled: true
      priority: 10
EOF

# Restart the service to load the new auth provider
systemctl restart ubuntu-autoinstall-webhook
```

### 12.3. Advanced Networking

#### 12.3.1. VLAN and Network Segmentation

Configure the system to work with VLANs:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
networking:
  vlans:
    enabled: true
    interfaces:
      - name: "eth0.100"
        parent: "eth0"
        id: 100
        address: "10.100.0.5/24"
        gateway: "10.100.0.1"
        service_bindings: ["web"]

      - name: "eth0.200"
        parent: "eth0"
        id: 200
        address: "10.200.0.5/24"
        gateway: "10.200.0.1"
        service_bindings: ["pxe"]
```

Create corresponding VLAN interfaces on the system:

```bash
# Load 8021q kernel module if not already loaded
modprobe 8021q

# Create VLAN interfaces
ip link add link eth0 name eth0.100 type vlan id 100
ip link add link eth0 name eth0.200 type vlan id 200

# Configure addresses
ip addr add 10.100.0.5/24 dev eth0.100
ip addr add 10.200.0.5/24 dev eth0.200

# Bring interfaces up
ip link set eth0.100 up
ip link set eth0.200 up

# Configure default routes if needed
ip route add 10.100.0.0/24 via 10.100.0.1 dev eth0.100
ip route add 10.200.0.0/24 via 10.200.0.1 dev eth0.200
```

Make the configuration persistent with Netplan:

```yaml
# /etc/netplan/60-vlans.yaml
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: true
  vlans:
    eth0.100:
      id: 100
      link: eth0
      addresses: [10.100.0.5/24]
      routes:
        - to: 10.100.0.0/24
          via: 10.100.0.1
    eth0.200:
      id: 200
      link: eth0
      addresses: [10.200.0.5/24]
      routes:
        - to: 10.200.0.0/24
          via: 10.200.0.1
```

Apply the configuration:

```bash
netplan apply
```

#### 12.3.2. Network Bonding

Configure network bonding for redundancy:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
networking:
  bonding:
    enabled: true
    interfaces:
      - name: "bond0"
        mode: "802.3ad"
        slaves: ["eth0", "eth1"]
        address: "192.168.1.5/24"
        gateway: "192.168.1.1"
        options:
          miimon: 100
          lacp_rate: fast
```

Set up bonding on the system:

```bash
# Load bonding kernel module
modprobe bonding

# Create bonding configuration
cat > /etc/modprobe.d/bonding.conf << 'EOF'
alias bond0 bonding
options bond0 mode=802.3ad miimon=100 lacp_rate=1
EOF

# Configure with Netplan
cat > /etc/netplan/50-bonding.yaml << 'EOF'
network:
  version: 2
  bonds:
    bond0:
      interfaces: [eth0, eth1]
      parameters:
        mode: 802.3ad
        mii-monitor-interval: 100
        lacp-rate: fast
      addresses: [192.168.1.5/24]
      routes:
        - to: default
          via: 192.168.1.1
EOF

# Apply configuration
netplan apply
```

#### 12.3.3. IPv6 Configuration

Enable IPv6 support:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
networking:
  ipv6:
    enabled: true
    address: "2001:db8::5/64"
    gateway: "2001:db8::1"
```

Configure the system for IPv6:

```bash
# Add IPv6 address to interface
ip -6 addr add 2001:db8::5/64 dev eth0

# Add default route
ip -6 route add default via 2001:db8::1 dev eth0

# Make persistent with Netplan
cat > /etc/netplan/70-ipv6.yaml << 'EOF'
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: true
      addresses: [2001:db8::5/64]
      routes:
        - to: ::/0
          via: 2001:db8::1
EOF

netplan apply
```

### 12.4. Clustering and High Availability

#### 12.4.1. Basic Clustering Setup

Configure a basic cluster with multiple instances:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml on each node
cluster:
  enabled: true
  node_name: "node1"  # Unique for each node
  nodes:
    - name: "node1"
      address: "192.168.1.10"
      role: "primary"
    - name: "node2"
      address: "192.168.1.11"
      role: "secondary"
    - name: "node3"
      address: "192.168.1.12"
      role: "secondary"

  quorum:
    required: true
    min_nodes: 2
```

#### 12.4.2. Distributed File System

Configure a shared storage for installation files:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
storage:
  type: "shared"
  shared:
    type: "nfs"
    mount_point: "/var/lib/ubuntu-autoinstall-webhook/shared"
    server: "nfs.example.com"
    export: "/exports/ubuntu-autoinstall"
    options: "rw,sync,no_subtree_check"
```

Mount the NFS share:

```bash
# Create mount point
mkdir -p /var/lib/ubuntu-autoinstall-webhook/shared

# Mount NFS share
mount -t nfs nfs.example.com:/exports/ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/shared

# Make persistent by adding to /etc/fstab
echo "nfs.example.com:/exports/ubuntu-autoinstall /var/lib/ubuntu-autoinstall-webhook/shared nfs rw,sync,no_subtree_check 0 0" >> /etc/fstab
```

#### 12.4.3. Load Balancer Configuration

Example NGINX load balancer configuration:

```nginx
# /etc/nginx/sites-available/webhook-lb
upstream webhook_backend {
    # ip_hash ensures the same client always hits the same server
    ip_hash;

    server 192.168.1.10:8443;
    server 192.168.1.11:8443;
    server 192.168.1.12:8443;

    # Health check parameters
    check interval=5000 rise=2 fall=3 timeout=4000;
}

server {
    listen 443 ssl;
    server_name webhook.example.com;

    ssl_certificate /etc/nginx/ssl/webhook.crt;
    ssl_certificate_key /etc/nginx/ssl/webhook.key;
    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass https://webhook_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable and restart NGINX:

```bash
ln -s /etc/nginx/sites-available/webhook-lb /etc/nginx/sites-enabled/
systemctl restart nginx
```

### 12.5. Integration with External Systems

#### 12.5.1. CMDB Integration

Configure integration with a Configuration Management Database:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
integrations:
  cmdb:
    enabled: true
    type: "servicenow"
    url: "https://example.service-now.com/api"
    auth:
      type: "basic"
      username: "integration-user"
      password: "integration-password"

    # Field mappings
    mappings:
      system_id: "sys_id"
      hostname: "name"
      mac_address: "mac_address"
      ip_address: "ip_address"
      status: "status"

    # Synchronization settings
    sync:
      enabled: true
      interval_minutes: 60
      direction: "bidirectional"
```

#### 12.5.2. Monitoring System Integration

Configure integration with monitoring systems:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
integrations:
  monitoring:
    - type: "prometheus"
      enabled: true
      endpoint: "/metrics"
      labels:
        environment: "production"
        service: "ubuntu-autoinstall"

    - type: "nagios"
      enabled: true
      nrdp_url: "https://nagios.example.com/nrdp/"
      token: "your-nagios-token"
      host_name: "ubuntu-autoinstall"
      interval_minutes: 5
```

#### 12.5.3. Custom Webhook Notifications

Configure custom webhook notifications:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
notifications:
  webhooks:
    - name: "slack-alerts"
      url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXX"
      events:
        - system.discovered
        - installation.started
        - installation.completed
        - installation.failed
      headers:
        Content-Type: "application/json"
      template: |
        {
          "text": "Event: {{ .Event.Type }}",
          "blocks": [
            {
              "type": "section",
              "text": {
                "type": "mrkdwn",
                "text": "*{{ .Event.Type }}*\nSystem: {{ .System.Hostname }}\nMAC: {{ .System.MACAddress }}\nTimestamp: {{ .Event.Timestamp }}"
              }
            }
          ]
        }

    - name: "teams-notification"
      url: "https://outlook.office.com/webhook/..."
      events:
        - installation.failed
      headers:
        Content-Type: "application/json"
      template: |
        {
          "@type": "MessageCard",
          "@context": "https://schema.org/extensions",
          "summary": "Installation Failed",
          "sections": [
            {
              "activityTitle": "Installation Failed for {{ .System.Hostname }}",
              "facts": [
                { "name": "System", "value": "{{ .System.Hostname }}" },
                { "name": "MAC Address", "value": "{{ .System.MACAddress }}" },
                { "name": "Error", "value": "{{ .Event.Details.Error }}" },
                { "name": "Time", "value": "{{ .Event.Timestamp }}" }
              ]
            }
          ]
        }
```

### 12.6. Advanced Storage Configuration

#### 12.6.1. Object Storage for Installation Files

Configure S3-compatible object storage:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
storage:
  installation_files:
    type: "s3"
    s3:
      endpoint: "https://s3.amazonaws.com"
      region: "us-east-1"
      bucket: "ubuntu-autoinstall-files"
      prefix: "files/"
      access_key: "YOUR_ACCESS_KEY"
      secret_key: "YOUR_SECRET_KEY"
      use_path_style: false
```

#### 12.6.2. Database on External Storage

Configure external database storage:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
database:
  type: "cockroach"
  cockroach:
    host: "db.example.com"
    port: 26257
    database: "ubuntu_autoinstall"
    user: "ubuntu_autoinstall"
    password: "your_password"
    ssl:
      mode: "verify-full"
      ca_cert: "/etc/ubuntu-autoinstall-webhook/certs/ca.crt"
      client_cert: "/etc/ubuntu-autoinstall-webhook/certs/client.crt"
      client_key: "/etc/ubuntu-autoinstall-webhook/certs/client.key"
```

#### 12.6.3. Backup to Remote Storage

Configure backups to remote storage:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
backup:
  enabled: true
  schedule: "0 2 * * *"  # 2 AM daily
  retention:
    days: 30
    count: 10

  storage:
    type: "sftp"
    sftp:
      host: "backup.example.com"
      port: 22
      user: "backup-user"
      key_file: "/etc/ubuntu-autoinstall-webhook/keys/backup-key"
      path: "/backups/ubuntu-autoinstall/"
```

### 12.7. Custom Authentication and Authorization

#### 12.7.1. Custom LDAP Configuration

Configure advanced LDAP integration:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
auth:
  ldap:
    enabled: true
    servers:
      - uri: "ldaps://ldap1.example.com"
        priority: 1
      - uri: "ldaps://ldap2.example.com"
        priority: 2

    bind_dn: "cn=service-account,ou=ServiceAccounts,dc=example,dc=com"
    bind_password: "your_password"

    user:
      base_dn: "ou=Users,dc=example,dc=com"
      filter: "(&(objectClass=person)(sAMAccountName=%s))"
      username_attribute: "sAMAccountName"
      name_attribute: "displayName"
      email_attribute: "mail"

    group:
      base_dn: "ou=Groups,dc=example,dc=com"
      filter: "(&(objectClass=group)(member=%s))"
      name_attribute: "cn"

    role_mappings:
      - ldap_group: "AutoinstallAdmins"
        role: "admin"
      - ldap_group: "AutoinstallOperators"
        role: "operator"
      - ldap_group: "AutoinstallUsers"
        role: "installer"
```

#### 12.7.2. OpenID Connect Configuration

Configure authentication with OpenID Connect providers:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
auth:
  oidc:
    enabled: true
    providers:
      - name: "google"
        client_id: "your_client_id"
        client_secret: "your_client_secret"
        discovery_url: "https://accounts.google.com/.well-known/openid-configuration"
        scopes: ["openid", "profile", "email"]
        claims:
          username: "email"
          name: "name"
          email: "email"

      - name: "azure_ad"
        client_id: "your_client_id"
        client_secret: "your_client_secret"
        discovery_url: "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration"
        scopes: ["openid", "profile", "email"]
        claims:
          username: "preferred_username"
          name: "name"
          email: "email"

    role_mappings:
      - claim: "groups"
        values:
          - value: "autoinstall_admins"
            role: "admin"
          - value: "autoinstall_operators"
            role: "operator"
```

#### 12.7.3. Custom Authorization Rules

Configure fine-grained authorization rules:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
authorization:
  rules:
    - name: "ReadOnlyForNonProduction"
      description: "Read-only access to non-production systems"
      match:
        roles:
          - "viewer"
        resources:
          - type: "system"
            attributes:
              environment: "non-production"
      permissions:
        - "read"

    - name: "FullAccessToOwnedSystems"
      description: "Full access to systems created by the user"
      match:
        resources:
          - type: "system"
            attributes:
              created_by: "${user.id}"
      permissions:
        - "read"
        - "write"
        - "delete"
        - "install"
```

### 12.8. Advanced Logging

#### 12.8.1. Structured Logging Configuration

Configure advanced structured logging:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
logging:
  structured:
    enabled: true
    format: "json"

  levels:
    default: "info"
    components:
      webserver: "info"
      database: "warn"
      file_editor: "info"
      dnsmasq_watcher: "debug"
      certificate_issuer: "info"

  output:
    console:
      enabled: true
      format: "text"
    file:
      enabled: true
      format: "json"
      path: "/var/log/ubuntu-autoinstall-webhook/webhook.log"
    syslog:
      enabled: true
      format: "json"
      facility: "local0"
```

#### 12.8.2. Remote Logging Configuration

Configure logging to remote systems:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
logging:
  remote:
    - type: "elasticsearch"
      enabled: true
      url: "https://elasticsearch.example.com:9200"
      index: "ubuntu-autoinstall-logs"
      username: "elastic"
      password: "your_password"
      batch_size: 100
      flush_interval_seconds: 5

    - type: "fluentd"
      enabled: true
      host: "fluentd.example.com"
      port: 24224
      tag: "ubuntu-autoinstall"

    - type: "graylog"
      enabled: true
      server: "graylog.example.com"
      port: 12201
      protocol: "tcp"
      tls:
        enabled: true
        skip_verify: false
```

#### 12.8.3. Log Correlation

Configure distributed tracing:

```yaml
# In /etc/ubuntu-autoinstall-webhook/config.yaml
logging:
  tracing:
    enabled: true
    type: "opentelemetry"
    service_name: "ubuntu-autoinstall-webhook"

    opentelemetry:
      endpoint: "http://jaeger.example.com:14268/api/traces"
      propagation: "w3c"
      sample_ratio: 0.1
```
