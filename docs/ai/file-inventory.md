<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [File Inventory](#file-inventory)
  - [Entry Points](#entry-points)
  - [Configuration](#configuration)
  - [Core Components](#core-components)
  - [Models and Data Structures](#models-and-data-structures)
  - [API Handlers](#api-handlers)
  - [Business Logic](#business-logic)
  - [Authentication and Security](#authentication-and-security)
  - [Web Interface](#web-interface)
  - [Utility Functions](#utility-functions)
  - [Tests](#tests)
  - [Build and Deployment](#build-and-deployment)
  - [Configuration Files](#configuration-files)
  - [Documentation](#documentation)
  - [File Dependencies and Relationships](#file-dependencies-and-relationships)
    - [Core Component Dependencies](#core-component-dependencies)
    - [API Handler Dependencies](#api-handler-dependencies)
    - [Business Logic Dependencies](#business-logic-dependencies)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# File Inventory

This document provides a comprehensive catalog of key files in the Ubuntu
Autoinstall Webhook project with descriptions of their purpose and
relationships. This information is specifically organized to help AI assistants
understand the codebase structure and functionality.

## Entry Points

| File Path                                          | Description                                                           |
| -------------------------------------------------- | --------------------------------------------------------------------- |
| `cmd/ubuntu-autoinstall-webhook/main.go`           | Main application entry point that initializes and starts all services |
| `cmd/ubuntu-autoinstall-webhook/commands/`         | Command-line interface implementation using Cobra                     |
| `cmd/ubuntu-autoinstall-webhook/commands/root.go`  | Root command definition and global flags                              |
| `cmd/ubuntu-autoinstall-webhook/commands/serve.go` | Web server command implementation                                     |

## Configuration

| File Path                       | Description                                 |
| ------------------------------- | ------------------------------------------- |
| `internal/config/config.go`     | Configuration definition and loading logic  |
| `internal/config/defaults.go`   | Default configuration values                |
| `internal/config/validation.go` | Configuration validation rules              |
| `internal/config/file.go`       | File-based configuration handling           |
| `internal/config/env.go`        | Environment variable configuration handling |

## Core Components

| File Path                         | Description                                      |
| --------------------------------- | ------------------------------------------------ |
| `internal/server/server.go`       | HTTP server implementation                       |
| `internal/database/database.go`   | Database interface and factory                   |
| `internal/database/sqlite.go`     | SQLite implementation of database interface      |
| `internal/database/cockroach.go`  | CockroachDB implementation of database interface |
| `internal/certificate/ca.go`      | Certificate Authority implementation             |
| `internal/certificate/issuer.go`  | Certificate issuance service                     |
| `internal/dnsmasq/watcher.go`     | DNSMasq log monitoring for DHCP events           |
| `internal/storage/file_editor.go` | File system operations for template management   |

## Models and Data Structures

| File Path                                  | Description                                              |
| ------------------------------------------ | -------------------------------------------------------- |
| `internal/database/models/system.go`       | System entity representing a physical or virtual machine |
| `internal/database/models/template.go`     | Template entity for installation configurations          |
| `internal/database/models/installation.go` | Installation entity tracking installation processes      |
| `internal/database/models/user.go`         | User entity for authentication and authorization         |
| `internal/database/models/certificate.go`  | Certificate entity for managing TLS certificates         |

## API Handlers

| File Path                       | Description                                |
| ------------------------------- | ------------------------------------------ |
| `internal/api/router.go`        | API route definitions and middleware setup |
| `internal/api/systems.go`       | Systems API endpoints                      |
| `internal/api/templates.go`     | Templates API endpoints                    |
| `internal/api/installations.go` | Installations API endpoints                |
| `internal/api/users.go`         | User management API endpoints              |
| `internal/api/certificates.go`  | Certificate management API endpoints       |
| `internal/api/health.go`        | Health check and monitoring endpoints      |
| `internal/api/metrics.go`       | Prometheus metrics endpoints               |

## Business Logic

| File Path                          | Description                         |
| ---------------------------------- | ----------------------------------- |
| `internal/system/service.go`       | System management business logic    |
| `internal/template/service.go`     | Template processing and management  |
| `internal/installation/service.go` | Installation process management     |
| `internal/installation/tracker.go` | Installation status tracking        |
| `internal/cloud_init/generator.go` | Cloud-init configuration generation |
| `internal/ipxe/generator.go`       | iPXE script generation              |

## Authentication and Security

| File Path                         | Description                          |
| --------------------------------- | ------------------------------------ |
| `internal/auth/auth.go`           | Authentication service interface     |
| `internal/auth/local.go`          | Local authentication implementation  |
| `internal/auth/ldap.go`           | LDAP authentication implementation   |
| `internal/auth/oauth.go`          | OAuth2 authentication implementation |
| `internal/auth/jwt.go`            | JWT token generation and validation  |
| `internal/auth/middleware.go`     | Authentication middleware            |
| `internal/security/encryption.go` | Data encryption utilities            |
| `internal/security/password.go`   | Password hashing and validation      |

## Web Interface

| File Path                          | Description                            |
| ---------------------------------- | -------------------------------------- |
| `web/templates/layout.html`        | Main layout template for web interface |
| `web/templates/dashboard.html`     | Dashboard view template                |
| `web/templates/systems.html`       | Systems management view                |
| `web/templates/templates.html`     | Template management view               |
| `web/templates/installations.html` | Installation monitoring view           |
| `web/static/js/main.js`            | Main JavaScript for web interface      |
| `web/static/css/styles.css`        | Main stylesheet for web interface      |

## Utility Functions

| File Path                      | Description                          |
| ------------------------------ | ------------------------------------ |
| `internal/utils/network.go`    | Network-related utility functions    |
| `internal/utils/crypto.go`     | Cryptographic utility functions      |
| `internal/utils/filesystem.go` | File system utility functions        |
| `internal/utils/validation.go` | Input validation utility functions   |
| `internal/utils/template.go`   | Template rendering utility functions |

## Tests

| File Path                               | Description                               |
| --------------------------------------- | ----------------------------------------- |
| `internal/server/server_test.go`        | Tests for HTTP server functionality       |
| `internal/database/database_test.go`    | Tests for database layer                  |
| `internal/api/systems_test.go`          | Tests for systems API endpoints           |
| `internal/installation/service_test.go` | Tests for installation service            |
| `test/integration/api_test.go`          | Integration tests for API functionality   |
| `test/e2e/installation_test.go`         | End-to-end tests for installation process |
| `test/fixtures/templates/`              | Test template fixtures                    |

## Build and Deployment

| File Path                                       | Description                          |
| ----------------------------------------------- | ------------------------------------ |
| `Makefile`                                      | Build automation rules               |
| `build/docker/Dockerfile`                       | Docker image definition              |
| `build/package/deb/`                            | Debian package build scripts         |
| `build/package/rpm/`                            | RPM package build scripts            |
| `deployments/docker-compose/docker-compose.yml` | Docker Compose deployment definition |
| `deployments/kubernetes/deployment.yaml`        | Kubernetes deployment manifest       |

## Configuration Files

| File Path                                 | Description                      |
| ----------------------------------------- | -------------------------------- |
| `examples/config.yaml`                    | Example configuration file       |
| `examples/templates/base.yaml`            | Example base template            |
| `examples/templates/web-server.yaml`      | Example web server template      |
| `examples/templates/database-server.yaml` | Example database server template |

## Documentation

| File Path                            | Description                        |
| ------------------------------------ | ---------------------------------- |
| `docs/admin-guide.md`                | Administrator guide                |
| `docs/user-guide.md`                 | End-user guide                     |
| `docs/architecture/overview.md`      | Architecture overview              |
| `docs/architecture/components/`      | Individual component documentation |
| `docs/technical/technical-design.md` | Technical design specifications    |
| `docs/technical/test-design.md`      | Test design documentation          |
| `docs/ai/`                           | AI-specific documentation          |

## File Dependencies and Relationships

### Core Component Dependencies

```
server.go
├── api/router.go
├── config/config.go
├── database/database.go
├── auth/middleware.go
└── logging/logging.go

database.go
├── models/*.go
├── config/config.go
└── migrations/*.sql

dnsmasq/watcher.go
├── system/service.go
├── config/config.go
└── logging/logging.go

certificate/issuer.go
├── certificate/ca.go
├── database/models/certificate.go
└── config/config.go
```

### API Handler Dependencies

```
api/systems.go
├── database/models/system.go
├── system/service.go
└── auth/middleware.go

api/templates.go
├── database/models/template.go
├── template/service.go
└── auth/middleware.go

api/installations.go
├── database/models/installation.go
├── installation/service.go
└── auth/middleware.go
```

### Business Logic Dependencies

```
system/service.go
├── database/models/system.go
├── installation/service.go
└── template/service.go

installation/service.go
├── database/models/installation.go
├── cloud_init/generator.go
├── ipxe/generator.go
└── storage/file_editor.go

template/service.go
├── database/models/template.go
├── cloud_init/generator.go
└── storage/file_editor.go
```

This file inventory provides a foundation for understanding the key components
of the project. Each file serves a specific purpose within the overall
architecture, and understanding these relationships is crucial for navigating
and modifying the codebase.
