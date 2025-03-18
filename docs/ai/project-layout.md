<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Project Layout](#project-layout)
  - [Repository Structure](#repository-structure)
  - [Key Directories Explained](#key-directories-explained)
    - [`cmd/`](#cmd)
    - [`internal/`](#internal)
      - [Core Components in `internal/`:](#core-components-in-internal)
    - [`pkg/`](#pkg)
    - [`web/`](#web)
    - [`docs/`](#docs)
    - [`test/`](#test)
  - [Code Organization Principles](#code-organization-principles)
  - [Important Design Patterns](#important-design-patterns)
  - [Import Structure](#import-structure)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Project Layout

This document provides a structured overview of the Ubuntu Autoinstall Webhook repository organization to help AI assistants understand the codebase structure.

## Repository Structure

```
ubuntu-autoinstall-webhook/
├── cmd/                           # Command-line entry points
│   └── ubuntu-autoinstall-webhook/  # Main application binary
├── internal/                      # Private application code
│   ├── api/                       # API handlers and routes
│   ├── auth/                      # Authentication and authorization
│   ├── certificate/               # Certificate management
│   ├── cloud_init/                # Cloud-init template handling
│   ├── config/                    # Configuration processing
│   ├── database/                  # Database abstraction and models
│   │   ├── models/                # Database entity definitions
│   │   └── migrations/            # Database schema migrations
│   ├── dnsmasq/                   # DHCP/DNS integration
│   ├── installation/              # Installation process management
│   ├── ipxe/                      # iPXE script generation
│   ├── logging/                   # Logging infrastructure
│   ├── metrics/                   # Performance metrics collection
│   ├── notification/              # Alert and notification system
│   ├── security/                  # Security utilities
│   ├── server/                    # HTTP server implementation
│   ├── storage/                   # File system operations
│   ├── system/                    # System management logic
│   ├── template/                  # Installation template handling
│   └── utils/                     # Shared utilities
├── pkg/                           # Public libraries
│   ├── client/                    # API client library
│   └── webhook/                   # Core webhook functionality
├── web/                           # Web UI assets
│   ├── src/                       # Frontend source code
│   ├── static/                    # Static assets
│   └── templates/                 # HTML templates
├── build/                         # Build-related files
│   ├── docker/                    # Docker build files
│   └── package/                   # Packaging scripts
├── scripts/                       # Utility scripts
├── deployments/                   # Deployment configurations
│   ├── docker-compose/            # Docker Compose files
│   ├── kubernetes/                # Kubernetes manifests
│   └── terraform/                 # Terraform modules
├── docs/                          # Documentation
│   ├── architecture/              # Architecture documentation
│   │   ├── components/            # Component descriptions
│   │   └── diagrams/              # Architecture diagrams
│   ├── technical/                 # Technical specifications
│   ├── ai/                        # AI-specific documentation
│   ├── admin-guide.md             # Administrator guide
│   └── user-guide.md              # User guide
├── examples/                      # Example configurations
├── test/                          # Test files
│   ├── e2e/                       # End-to-end tests
│   ├── integration/               # Integration tests
│   └── fixtures/                  # Test data fixtures
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── Makefile                       # Build automation
├── .github/                       # GitHub specific files
│   └── workflows/                 # CI/CD workflows
└── README.md                      # Project readme
```

## Key Directories Explained

### `cmd/`
Contains the application entry points. Each subdirectory typically corresponds to a compiled binary.

### `internal/`
Contains application code that's private to this project. The Go compiler enforces that packages under `internal/` can only be imported by code in the parent directory.

#### Core Components in `internal/`:
- **api**: RESTful API implementation
- **certificate**: Certificate authority and management
- **database**: Database interactions and models
- **dnsmasq**: Integration with DNSMasq for DHCP monitoring
- **installation**: Manages OS installations
- **server**: HTTP server implementation
- **template**: Installation template processing

### `pkg/`
Contains code that may be used by external applications. These packages have stable APIs and are designed for reuse.

### `web/`
Contains web interface assets, including HTML templates, JavaScript, CSS, and other static files.

### `docs/`
Project documentation, organized by audience and purpose.

### `test/`
Contains test code, especially integration and end-to-end tests. Unit tests are typically co-located with the code they test.

## Code Organization Principles

1. **Domain-Driven Design**: Code is organized around business domains rather than technical function.

2. **Clean Architecture**: The codebase follows clean architecture principles with clear separation of:
   - Domain logic
   - Application business rules
   - Interface adapters
   - External frameworks and drivers

3. **Dependency Injection**: Dependencies are passed to components rather than created within them.

4. **Component Isolation**: Major system components are designed to work independently.

5. **Interface-Based Design**: Components interact through interfaces, allowing mock implementations during testing.

## Important Design Patterns

1. **Repository Pattern**: Used for database access, abstracting storage details.

2. **Service Layer**: Business logic is encapsulated in service objects.

3. **Middleware Pattern**: Used for HTTP request processing (logging, authentication, etc.).

4. **Factory Pattern**: Used for creating complex objects.

5. **Observer Pattern**: Used for event notification between components.

6. **Command Pattern**: Used for installation operations.

## Import Structure

The application follows a clear import hierarchy:

1. Standard library imports
2. Third-party imports
3. Project imports

Project imports maintain the following dependency direction:
- `cmd` → `internal/api` → `internal/service` → `internal/repository` → `internal/model`

This ensures clean separation and prevents circular dependencies.
