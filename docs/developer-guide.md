# Developer Guide for Ubuntu Autoinstall Webhook

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

  - [1. Introduction](#1-introduction)
    - [1.1. Project Overview](#11-project-overview)
    - [1.2. Development Philosophy](#12-development-philosophy)
    - [1.3. Repository Overview](#13-repository-overview)
  - [2. Development Environment Setup](#2-development-environment-setup)
    - [2.1. Prerequisites](#21-prerequisites)
    - [2.2. Cloning the Repository](#22-cloning-the-repository)
    - [2.3. Setting Up the Development Environment](#23-setting-up-the-development-environment)
      - [Basic Setup](#basic-setup)
      - [Development Configuration](#development-configuration)
      - [Frontend Development Setup](#frontend-development-setup)
    - [2.4. IDE Configuration](#24-ide-configuration)
      - [Visual Studio Code](#visual-studio-code)
      - [GoLand](#goland)
    - [2.5. Docker Development Environment](#25-docker-development-environment)
  - [3. Project Structure and Code Organization](#3-project-structure-and-code-organization)
    - [3.1. Detailed Directory Structure](#31-detailed-directory-structure)
    - [3.2. Component Architecture](#32-component-architecture)
    - [3.3. Code Organization Principles](#33-code-organization-principles)
      - [Dependency Direction](#dependency-direction)
      - [Interface Definitions](#interface-definitions)
      - [Error Handling](#error-handling)
    - [3.4. Key Design Patterns](#34-key-design-patterns)
  - [4. Code Standards and Guidelines](#4-code-standards-and-guidelines)
    - [4.1. Go Code Style](#41-go-code-style)
      - [Naming Conventions](#naming-conventions)
      - [File Organization](#file-organization)
    - [4.2. Documentation Standards](#42-documentation-standards)
      - [Code Comments](#code-comments)
      - [Package Documentation](#package-documentation)
    - [4.3. Error Handling](#43-error-handling)
    - [4.4. Logging Guidelines](#44-logging-guidelines)
    - [4.5. Testing Standards](#45-testing-standards)
  - [5. Building and Testing](#5-building-and-testing)
    - [5.1. Building the Project](#51-building-the-project)
    - [5.2. Running Tests](#52-running-tests)
    - [5.3. Linting and Code Quality](#53-linting-and-code-quality)
    - [5.4. Building Docker Images](#54-building-docker-images)
    - [5.5. Continuous Integration](#55-continuous-integration)
  - [6. Contribution Workflow](#6-contribution-workflow)
    - [6.1. Getting Started](#61-getting-started)
    - [6.2. Creating a Feature Branch](#62-creating-a-feature-branch)
    - [6.3. Making Changes](#63-making-changes)
    - [6.4. Commit Guidelines](#64-commit-guidelines)
    - [6.5. Submitting a Pull Request](#65-submitting-a-pull-request)
    - [6.6. Code Review Process](#66-code-review-process)
    - [6.7. After Merge](#67-after-merge)
  - [7. API Development](#7-api-development)
    - [7.1. API Design Principles](#71-api-design-principles)
    - [7.2. Adding New Endpoints](#72-adding-new-endpoints)
    - [7.3. Error Handling in API](#73-error-handling-in-api)
- [Developer Guide for Ubuntu Autoinstall Webhook](#developer-guide-for-ubuntu-autoinstall-webhook)
  - [1. Introduction](#1-introduction-1)
    - [1.1. Project Overview](#11-project-overview-1)
    - [1.2. Development Philosophy](#12-development-philosophy-1)
    - [1.3. Repository Overview](#13-repository-overview-1)
  - [2. Development Environment Setup](#2-development-environment-setup-1)
    - [2.1. Prerequisites](#21-prerequisites-1)
    - [2.2. Cloning the Repository](#22-cloning-the-repository-1)
    - [2.3. Setting Up the Development Environment](#23-setting-up-the-development-environment-1)
      - [Basic Setup](#basic-setup-1)
      - [Development Configuration](#development-configuration-1)
      - [Frontend Development Setup](#frontend-development-setup-1)
    - [2.4. IDE Configuration](#24-ide-configuration-1)
      - [Visual Studio Code](#visual-studio-code-1)
      - [GoLand](#goland-1)
    - [2.5. Docker Development Environment](#25-docker-development-environment-1)
  - [3. Project Structure and Code Organization](#3-project-structure-and-code-organization-1)
    - [3.1. Detailed Directory Structure](#31-detailed-directory-structure-1)
    - [3.2. Component Architecture](#32-component-architecture-1)
    - [3.3. Code Organization Principles](#33-code-organization-principles-1)
      - [Dependency Direction](#dependency-direction-1)
      - [Interface Definitions](#interface-definitions-1)
      - [Error Handling](#error-handling-1)
    - [3.4. Key Design Patterns](#34-key-design-patterns-1)
  - [4. Code Standards and Guidelines](#4-code-standards-and-guidelines-1)
    - [4.1. Go Code Style](#41-go-code-style-1)
      - [Naming Conventions](#naming-conventions-1)
      - [File Organization](#file-organization-1)
    - [4.2. Documentation Standards](#42-documentation-standards-1)
      - [Code Comments](#code-comments-1)
      - [Package Documentation](#package-documentation-1)
    - [4.3. Error Handling](#43-error-handling-1)
    - [4.4. Logging Guidelines](#44-logging-guidelines-1)
    - [4.5. Testing Standards](#45-testing-standards-1)
  - [5. Building and Testing](#5-building-and-testing-1)
    - [5.1. Building the Project](#51-building-the-project-1)
    - [5.2. Running Tests](#52-running-tests-1)
    - [5.3. Linting and Code Quality](#53-linting-and-code-quality-1)
    - [5.4. Building Docker Images](#54-building-docker-images-1)
    - [5.5. Continuous Integration](#55-continuous-integration-1)
  - [6. Contribution Workflow](#6-contribution-workflow-1)
    - [6.1. Getting Started](#61-getting-started-1)
    - [6.2. Creating a Feature Branch](#62-creating-a-feature-branch-1)
    - [6.3. Making Changes](#63-making-changes-1)
    - [6.4. Commit Guidelines](#64-commit-guidelines-1)
    - [6.5. Submitting a Pull Request](#65-submitting-a-pull-request-1)
    - [6.6. Code Review Process](#66-code-review-process-1)
    - [6.7. After Merge](#67-after-merge-1)
  - [7. API Development](#7-api-development-1)
    - [7.1. API Design Principles](#71-api-design-principles-1)
    - [7.2. Adding New Endpoints](#72-adding-new-endpoints-1)
    - [7.2. Endpoint Implementation](#72-endpoint-implementation)
    - [7.3. Error Handling in API](#73-error-handling-in-api-1)
      - [7.3.1. Error Response Structure](#731-error-response-structure)
      - [7.3.2. Standard Error Codes](#732-standard-error-codes)
      - [7.3.3. Implementing Error Handling](#733-implementing-error-handling)
      - [7.3.4. Error Helper Functions](#734-error-helper-functions)
      - [7.3.5. Validation Errors](#735-validation-errors)
      - [7.3.6. Domain Errors vs. API Errors](#736-domain-errors-vs-api-errors)
      - [7.3.7. Error Logging](#737-error-logging)
      - [7.3.8. Security Considerations](#738-security-considerations)
      - [7.3.9. Testing Error Conditions](#739-testing-error-conditions)
    - [7.4. API Authentication](#74-api-authentication)
      - [7.4.1. Authentication Methods](#741-authentication-methods)
      - [7.4.2. API Token Implementation](#742-api-token-implementation)
      - [7.4.3. JWT Authentication](#743-jwt-authentication)
      - [7.4.4. Authentication Middleware](#744-authentication-middleware)
      - [7.4.5. Authorization Middleware](#745-authorization-middleware)
      - [7.4.6. Securing API Routes](#746-securing-api-routes)
      - [7.4.7. Testing Authentication](#747-testing-authentication)
    - [7.5. API Versioning](#75-api-versioning)
      - [7.5.1. Versioning Strategy](#751-versioning-strategy)
      - [7.5.2. Implementing Versioned Routes](#752-implementing-versioned-routes)
      - [7.5.3. Version Compatibility](#753-version-compatibility)
      - [7.5.4. Supporting Multiple Versions](#754-supporting-multiple-versions)
    - [7.6. API Documentation](#76-api-documentation)
      - [7.6.1. OpenAPI Specification](#761-openapi-specification)
      - [7.6.2. Documenting Endpoints](#762-documenting-endpoints)
      - [7.6.3. Schema Definitions](#763-schema-definitions)
      - [7.6.4. Generating API Documentation](#764-generating-api-documentation)
      - [7.6.5. Keeping Documentation in Sync](#765-keeping-documentation-in-sync)
    - [7.7. Rate Limiting](#77-rate-limiting)
      - [7.7.1. Rate Limiting Strategy](#771-rate-limiting-strategy)
      - [7.7.2. Implementing Rate Limiting](#772-implementing-rate-limiting)
      - [7.7.3. Different Rate Limits by Role](#773-different-rate-limits-by-role)
      - [7.7.4. Rate Limit Headers](#774-rate-limit-headers)
  - [8. Database Development](#8-database-development)
    - [8.1. Database Architecture](#81-database-architecture)
      - [8.1.1. Database Abstraction Layer](#811-database-abstraction-layer)
      - [8.1.2. Supported Database Backends](#812-supported-database-backends)
      - [8.1.3. Component Structure](#813-component-structure)
      - [8.1.4. Database Interface](#814-database-interface)
    - [8.2. Model Definition](#82-model-definition)
      - [8.2.1. Basic Model Structure](#821-basic-model-structure)
      - [8.2.2. Entity Model Example](#822-entity-model-example)
      - [8.2.3. Custom Data Types](#823-custom-data-types)
      - [8.2.4. Model Relationships](#824-model-relationships)
    - [8.3. Repository Pattern](#83-repository-pattern)
      - [8.3.1. Repository Interfaces](#831-repository-interfaces)
      - [8.3.2. SQLite Implementation](#832-sqlite-implementation)
      - [8.3.3. CockroachDB Implementation](#833-cockroachdb-implementation)
    - [8.4. Database Migrations](#84-database-migrations)
      - [8.4.1. Migration Structure](#841-migration-structure)
      - [8.4.2. Migration Implementation](#842-migration-implementation)
      - [8.4.3. Migration Example](#843-migration-example)
    - [8.5. Query Optimization](#85-query-optimization)
      - [8.5.1. Indexing Strategy](#851-indexing-strategy)
      - [8.5.2. Query Optimization Techniques](#852-query-optimization-techniques)
      - [8.5.3. Database-Specific Optimizations](#853-database-specific-optimizations)
    - [8.6. Testing Database Code](#86-testing-database-code)
      - [8.6.1. Repository Testing with SQLite](#861-repository-testing-with-sqlite)
      - [8.6.2. Mock Repositories for Service Tests](#862-mock-repositories-for-service-tests)
      - [8.6.3. Integration Tests with Test Database](#863-integration-tests-with-test-database)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 1. Introduction

This developer guide provides comprehensive information for developers who want
to understand, modify, or contribute to the Ubuntu Autoinstall Webhook project.
Whether you're fixing bugs, adding features, or just exploring the codebase,
this document will help you navigate the project effectively.

### 1.1. Project Overview

The Ubuntu Autoinstall Webhook system is designed to automate the deployment of
Ubuntu Server systems at scale. It provides:

- A centralized service for managing Ubuntu installations
- Template-based configuration
- Integration with PXE boot and network services
- A REST API for programmatic control
- A web interface for human operators

The system follows a modular design with clear separation between components,
making it easier to understand, test, and extend individual parts without
affecting the whole system.

### 1.2. Development Philosophy

This project follows several key principles:

1. **Separation of Concerns** Each component has a well-defined responsibility
   and minimal dependencies.

2. **Interface-Based Design** Components interact through interfaces, allowing
   for flexible implementation and easier testing.

3. **Security First** Security is integrated into the design from the beginning,
   not added as an afterthought.

4. **Test-Driven Development** Comprehensive testing is a core part of the
   development process.

5. **Maintainable Code** Code clarity and maintainability are prioritized over
   clever optimizations.

### 1.3. Repository Overview

The project is organized as a Go module with the following high-level structure:

```
ubuntu-autoinstall-webhook/
├── cmd/                  # Command-line entry points
├── internal/             # Private application code
├── pkg/                  # Public libraries
├── web/                  # Web interface assets
├── docs/                 # Documentation
├── test/                 # Test files
├── build/                # Build-related files
└── scripts/              # Utility scripts
```

For a more detailed breakdown of the repository structure, see
[Project Layout](./ai/project-layout.md).

## 2. Development Environment Setup

### 2.1. Prerequisites

Before you begin, ensure you have the following installed:

- **Go** (version 1.21 or newer)
- **Git** for version control
- **Make** for build automation
- **Docker** and **Docker Compose** for containerized development
- **Node.js** and **npm** for web frontend development (if working on the UI)

For Ubuntu/Debian systems, you can install these dependencies with:

```bash
# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install other dependencies
sudo apt update
sudo apt install -y git make docker.io docker-compose nodejs npm
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER  # Log out and back in after this
```

### 2.2. Cloning the Repository

Clone the repository and navigate into it:

```bash
git clone https://github.com/jdfalk/ubuntu-autoinstall-webhook.git
cd ubuntu-autoinstall-webhook
```

### 2.3. Setting Up the Development Environment

#### Basic Setup

1. Install Go dependencies:

```bash
go mod download
```

2. Install development tools:

```bash
make install-tools
```

This will install:

- `golint` for linting
- `goimports` for import formatting
- `errcheck` for error checking
- `staticcheck` for static analysis
- `mockgen` for generating mock interfaces

#### Development Configuration

Create a development configuration file:

```bash
cp examples/config.yaml ./dev-config.yaml
```

Edit `dev-config.yaml` to set appropriate values for your environment.

#### Frontend Development Setup

If you're working on the web interface:

```bash
cd web
npm install
```

### 2.4. IDE Configuration

#### Visual Studio Code

For VS Code users, we recommend the following extensions:

- Go (ms-vscode.go)
- EditorConfig for VS Code
- Docker
- YAML

Add the following settings to your `.vscode/settings.json`:

```json
{
  "go.lintTool": "golint",
  "go.formatTool": "goimports",
  "go.useLanguageServer": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  },
  "go.testFlags": ["-v"],
  "go.testTimeout": "300s"
}
```

#### GoLand

For GoLand users:

1. Enable Go modules integration (Settings -> Go -> Go Modules)
2. Configure the Go formatter to use goimports
3. Enable code inspections for Go

### 2.5. Docker Development Environment

For a containerized development environment:

```bash
# Build development container
docker-compose -f build/docker-compose/dev.yml build

# Start development environment
docker-compose -f build/docker-compose/dev.yml up -d

# Access development shell
docker-compose -f build/docker-compose/dev.yml exec app bash
```

## 3. Project Structure and Code Organization

### 3.1. Detailed Directory Structure

```
ubuntu-autoinstall-webhook/
├── cmd/
│   └── ubuntu-autoinstall-webhook/
│       ├── main.go                 # Application entry point
│       └── commands/               # CLI command definitions
├── internal/
│   ├── api/                       # API handlers
│   ├── auth/                      # Authentication & authorization
│   ├── certificate/               # Certificate management
│   ├── cloud_init/                # Cloud-init configuration
│   ├── config/                    # Configuration management
│   ├── database/                  # Database abstraction
│   │   ├── models/                # Data models
│   │   └── migrations/            # Database migrations
│   ├── dnsmasq/                   # DNSMasq integration
│   ├── installation/              # Installation management
│   ├── ipxe/                      # iPXE script generation
│   ├── server/                    # HTTP server
│   ├── system/                    # System management
│   ├── template/                  # Template handling
│   └── utils/                     # Shared utilities
├── pkg/
│   ├── client/                    # API client library
│   └── webhook/                   # Core webhook functionality
├── web/
│   ├── src/                       # Frontend source code
│   ├── static/                    # Static assets
│   └── templates/                 # HTML templates
├── docs/
│   ├── architecture/              # Architecture documentation
│   ├── technical/                 # Technical specifications
│   └── ai/                        # AI-specific documentation
├── test/
│   ├── e2e/                       # End-to-end tests
│   ├── integration/               # Integration tests
│   └── fixtures/                  # Test fixtures
└── build/
    ├── docker/                    # Docker build files
    └── package/                   # Packaging scripts
```

### 3.2. Component Architecture

The system is built around these core components:

1. **Web Server**: Handles HTTP requests for both API and web interface
2. **Certificate Issuer**: Manages TLS certificates
3. **File Editor**: Handles file system operations
4. **DNSMasq Watcher**: Monitors DHCP events
5. **Database**: Stores system data and configurations
6. **Configuration**: Manages system settings

Each component is designed with clear interfaces to separate the component's
contract from its implementation details.

### 3.3. Code Organization Principles

#### Dependency Direction

The code follows a clear dependency direction:

```
cmd → api → services → repositories → models
```

This ensures that higher-level components depend on lower-level ones, but not
vice versa.

#### Interface Definitions

Interfaces are typically defined in the package where they are used, not where
they are implemented. For example:

```go
// internal/system/service.go
type Repository interface {
    GetSystem(id string) (*model.System, error)
    // Other methods...
}

// The implementation lives in internal/database/system_repository.go
```

#### Error Handling

Errors are propagated up the call stack and handled at appropriate levels:

- Low-level functions return specific errors
- Mid-level functions add context to errors
- High-level functions (API handlers, etc.) translate errors to appropriate
  responses

### 3.4. Key Design Patterns

The codebase uses several design patterns consistently:

1. **Repository Pattern**: Abstracts data access behind interfaces
2. **Dependency Injection**: Components receive their dependencies
3. **Factory Pattern**: For creating complex objects
4. **Middleware**: For processing HTTP requests
5. **Adapter Pattern**: For integrating external systems

## 4. Code Standards and Guidelines

### 4.1. Go Code Style

This project follows the standard Go code style and conventions:

- Use `gofmt` or `goimports` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Adhere to
  [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

#### Naming Conventions

- Use **CamelCase** for exported names (visible outside the package)
- Use **lowerCamelCase** for non-exported names
- Choose descriptive names; avoid abbreviations except for common ones
- Keep package names short, concise, and lowercase

#### File Organization

- One file per logical component
- Related files in the same package
- Test files next to the code they test (`foo.go` and `foo_test.go`)

### 4.2. Documentation Standards

#### Code Comments

- Every exported function, type, and constant must have a comment
- Comments should explain the "why" not just the "what"
- Use complete sentences starting with the function name

Example:

```go
// GetSystem retrieves a system by its ID. It returns nil and ErrSystemNotFound
// if the system doesn't exist, or another error if the database query fails.
func GetSystem(id string) (*System, error) {
    // ...
}
```

#### Package Documentation

Each package should have a package comment in one of its files:

```go
// Package system provides functionality for managing system entities,
// including creation, retrieval, updating, and deletion operations.
package system
```

### 4.3. Error Handling

- Always check errors
- Don't use `panic` in production code
- Add context to errors with `fmt.Errorf("doing X: %w", err)`
- Use custom error types for expected error conditions

```go
var (
    // ErrSystemNotFound is returned when a system with the specified ID does not exist.
    ErrSystemNotFound = errors.New("system not found")
)
```

### 4.4. Logging Guidelines

- Use the project's logging package (`internal/logging`)
- Log at appropriate levels (debug, info, warning, error, fatal)
- Include relevant context in log messages
- Avoid logging sensitive information

```go
// Incorrect:
log.Infof("User %s logged in with password %s", username, password)

// Correct:
log.Infof("User %s logged in successfully", username)
```

### 4.5. Testing Standards

- All packages should have tests
- Aim for at least 70% test coverage
- Use table-driven tests for multiple test cases
- Use mocks for external dependencies
- Test both success and error paths

```go
func TestGetSystem(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        wantErr bool
    }{
        {"existing system", "valid-id", false},
        {"non-existent system", "invalid-id", true},
        {"empty id", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 5. Building and Testing

### 5.1. Building the Project

The project uses Make for build automation. Common Make targets:

```bash
# Build the main binary
make build

# Build for all supported platforms
make build-all

# Build and install locally
make install

# Clean build artifacts
make clean
```

### 5.2. Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific tests
go test ./internal/system/...

# Run tests with verbose output
go test -v ./...
```

### 5.3. Linting and Code Quality

```bash
# Run all linters
make lint

# Format code
make fmt

# Check for common errors
make errcheck

# Run static analysis
make staticcheck
```

### 5.4. Building Docker Images

```bash
# Build development image
make docker-build-dev

# Build production image
make docker-build-prod

# Run using Docker
make docker-run
```

### 5.5. Continuous Integration

The project uses GitHub Actions for CI/CD. Each pull request triggers:

- Linting
- Unit tests
- Integration tests
- Build verification

See `.github/workflows/` for CI configuration.

## 6. Contribution Workflow

### 6.1. Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Add the upstream repository as a remote

```bash
git remote add upstream https://github.com/jdfalk/ubuntu-autoinstall-webhook.git
```

### 6.2. Creating a Feature Branch

Create a new branch for your feature or bugfix:

```bash
# Update your main branch
git checkout main
git pull upstream main

# Create a new branch
git checkout -b feature/your-feature-name
```

### 6.3. Making Changes

1. Make your changes in the feature branch
2. Add tests for your changes
3. Ensure all tests pass
4. Update documentation if needed

### 6.4. Commit Guidelines

- Use descriptive commit messages
- Start with a verb in the present tense (Add, Fix, Update, etc.)
- Reference issue numbers if applicable
- Keep changes focused and atomic

Example:

```
Add system discovery feature

- Implement DHCP event monitoring
- Add system registration endpoint
- Create database schema for discovered systems

Fixes #123
```

### 6.5. Submitting a Pull Request

1. Push your changes to your fork
2. Create a pull request against the upstream repository
3. Describe your changes in the PR description
4. Address any feedback from code reviews

### 6.6. Code Review Process

- All PRs require at least one review
- Address all review comments
- Ensure CI tests pass
- Maintain a respectful and collaborative tone

### 6.7. After Merge

After your PR is merged:

1. Update your local main:
   ```bash
   git checkout main
   git pull upstream main
   ```
2. Delete your feature branch:
   ```bash
   git branch -d feature/your-feature-name
   ```

## 7. API Development

### 7.1. API Design Principles

The API follows RESTful principles:

- Resource-based URLs
- Standard HTTP methods (GET, POST, PUT, DELETE)
- JSON request and response bodies
- Consistent error responses
- Stateless operations

### 7.2. Adding New Endpoints

To add a new API endpoint:

1. Define the endpoint in the appropriate file in `internal/api/`
2. Create handler functions for the endpoint
3. Register the endpoint in `internal/api/router.go`
4. Add authentication/authorization middleware if needed
5. Document the endpoint

Example:

```go
// internal/api/systems.go

// GetSystem handles GET requests to retrieve a system by ID.
func (h *SystemsHandler) GetSystem(w http.ResponseWriter, r *http.Request) {
    // Extract ID from request
    id := mux.Vars(r)["id"]

    // Call service
    system, err := h.systemService.GetSystem(id)
    if err != nil {
        if errors.Is(err, ErrSystemNotFound) {
            api.RespondNotFound(w, "System not found")
            return
        }
        api.RespondError(w, http.StatusInternalServerError, "Failed to get system")
        return
    }

    // Return response
    api.RespondJSON(w, http.StatusOK, system)
}

// RegisterRoutes registers all system API routes.
func (h *SystemsHandler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/systems/{id}", h.GetSystem).Methods("GET")
    // Other routes...
}
```

### 7.3. Error Handling in API

Use<!-- filepath: /Users/jdfalk/repos/github.com/jdfalk/ubuntu-autoinstall-webhook/docs/developer-guide.md -->

# Developer Guide for Ubuntu Autoinstall Webhook

<!-- START doctoc -->
<!-- END doctoc -->

## 1. Introduction

This developer guide provides comprehensive information for developers who want
to understand, modify, or contribute to the Ubuntu Autoinstall Webhook project.
Whether you're fixing bugs, adding features, or just exploring the codebase,
this document will help you navigate the project effectively.

### 1.1. Project Overview

The Ubuntu Autoinstall Webhook system is designed to automate the deployment of
Ubuntu Server systems at scale. It provides:

- A centralized service for managing Ubuntu installations
- Template-based configuration
- Integration with PXE boot and network services
- A REST API for programmatic control
- A web interface for human operators

The system follows a modular design with clear separation between components,
making it easier to understand, test, and extend individual parts without
affecting the whole system.

### 1.2. Development Philosophy

This project follows several key principles:

1. **Separation of Concerns** Each component has a well-defined responsibility
   and minimal dependencies.

2. **Interface-Based Design** Components interact through interfaces, allowing
   for flexible implementation and easier testing.

3. **Security First** Security is integrated into the design from the beginning,
   not added as an afterthought.

4. **Test-Driven Development** Comprehensive testing is a core part of the
   development process.

5. **Maintainable Code** Code clarity and maintainability are prioritized over
   clever optimizations.

### 1.3. Repository Overview

The project is organized as a Go module with the following high-level structure:

```
ubuntu-autoinstall-webhook/
├── cmd/                  # Command-line entry points
├── internal/             # Private application code
├── pkg/                  # Public libraries
├── web/                  # Web interface assets
├── docs/                 # Documentation
├── test/                 # Test files
├── build/                # Build-related files
└── scripts/              # Utility scripts
```

For a more detailed breakdown of the repository structure, see
[Project Layout](./ai/project-layout.md).

## 2. Development Environment Setup

### 2.1. Prerequisites

Before you begin, ensure you have the following installed:

- **Go** (version 1.21 or newer)
- **Git** for version control
- **Make** for build automation
- **Docker** and **Docker Compose** for containerized development
- **Node.js** and **npm** for web frontend development (if working on the UI)

For Ubuntu/Debian systems, you can install these dependencies with:

```bash
# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install other dependencies
sudo apt update
sudo apt install -y git make docker.io docker-compose nodejs npm
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER  # Log out and back in after this
```

### 2.2. Cloning the Repository

Clone the repository and navigate into it:

```bash
git clone https://github.com/jdfalk/ubuntu-autoinstall-webhook.git
cd ubuntu-autoinstall-webhook
```

### 2.3. Setting Up the Development Environment

#### Basic Setup

1. Install Go dependencies:

```bash
go mod download
```

2. Install development tools:

```bash
make install-tools
```

This will install:

- `golint` for linting
- `goimports` for import formatting
- `errcheck` for error checking
- `staticcheck` for static analysis
- `mockgen` for generating mock interfaces

#### Development Configuration

Create a development configuration file:

```bash
cp examples/config.yaml ./dev-config.yaml
```

Edit `dev-config.yaml` to set appropriate values for your environment.

#### Frontend Development Setup

If you're working on the web interface:

```bash
cd web
npm install
```

### 2.4. IDE Configuration

#### Visual Studio Code

For VS Code users, we recommend the following extensions:

- Go (ms-vscode.go)
- EditorConfig for VS Code
- Docker
- YAML

Add the following settings to your `.vscode/settings.json`:

```json
{
  "go.lintTool": "golint",
  "go.formatTool": "goimports",
  "go.useLanguageServer": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  },
  "go.testFlags": ["-v"],
  "go.testTimeout": "300s"
}
```

#### GoLand

For GoLand users:

1. Enable Go modules integration (Settings -> Go -> Go Modules)
2. Configure the Go formatter to use goimports
3. Enable code inspections for Go

### 2.5. Docker Development Environment

For a containerized development environment:

```bash
# Build development container
docker-compose -f build/docker-compose/dev.yml build

# Start development environment
docker-compose -f build/docker-compose/dev.yml up -d

# Access development shell
docker-compose -f build/docker-compose/dev.yml exec app bash
```

## 3. Project Structure and Code Organization

### 3.1. Detailed Directory Structure

```
ubuntu-autoinstall-webhook/
├── cmd/
│   └── ubuntu-autoinstall-webhook/
│       ├── main.go                 # Application entry point
│       └── commands/               # CLI command definitions
├── internal/
│   ├── api/                       # API handlers
│   ├── auth/                      # Authentication & authorization
│   ├── certificate/               # Certificate management
│   ├── cloud_init/                # Cloud-init configuration
│   ├── config/                    # Configuration management
│   ├── database/                  # Database abstraction
│   │   ├── models/                # Data models
│   │   └── migrations/            # Database migrations
│   ├── dnsmasq/                   # DNSMasq integration
│   ├── installation/              # Installation management
│   ├── ipxe/                      # iPXE script generation
│   ├── server/                    # HTTP server
│   ├── system/                    # System management
│   ├── template/                  # Template handling
│   └── utils/                     # Shared utilities
├── pkg/
│   ├── client/                    # API client library
│   └── webhook/                   # Core webhook functionality
├── web/
│   ├── src/                       # Frontend source code
│   ├── static/                    # Static assets
│   └── templates/                 # HTML templates
├── docs/
│   ├── architecture/              # Architecture documentation
│   ├── technical/                 # Technical specifications
│   └── ai/                        # AI-specific documentation
├── test/
│   ├── e2e/                       # End-to-end tests
│   ├── integration/               # Integration tests
│   └── fixtures/                  # Test fixtures
└── build/
    ├── docker/                    # Docker build files
    └── package/                   # Packaging scripts
```

### 3.2. Component Architecture

The system is built around these core components:

1. **Web Server**: Handles HTTP requests for both API and web interface
2. **Certificate Issuer**: Manages TLS certificates
3. **File Editor**: Handles file system operations
4. **DNSMasq Watcher**: Monitors DHCP events
5. **Database**: Stores system data and configurations
6. **Configuration**: Manages system settings

Each component is designed with clear interfaces to separate the component's
contract from its implementation details.

### 3.3. Code Organization Principles

#### Dependency Direction

The code follows a clear dependency direction:

```
cmd → api → services → repositories → models
```

This ensures that higher-level components depend on lower-level ones, but not
vice versa.

#### Interface Definitions

Interfaces are typically defined in the package where they are used, not where
they are implemented. For example:

```go
// internal/system/service.go
type Repository interface {
    GetSystem(id string) (*model.System, error)
    // Other methods...
}

// The implementation lives in internal/database/system_repository.go
```

#### Error Handling

Errors are propagated up the call stack and handled at appropriate levels:

- Low-level functions return specific errors
- Mid-level functions add context to errors
- High-level functions (API handlers, etc.) translate errors to appropriate
  responses

### 3.4. Key Design Patterns

The codebase uses several design patterns consistently:

1. **Repository Pattern**: Abstracts data access behind interfaces
2. **Dependency Injection**: Components receive their dependencies
3. **Factory Pattern**: For creating complex objects
4. **Middleware**: For processing HTTP requests
5. **Adapter Pattern**: For integrating external systems

## 4. Code Standards and Guidelines

### 4.1. Go Code Style

This project follows the standard Go code style and conventions:

- Use `gofmt` or `goimports` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Adhere to
  [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

#### Naming Conventions

- Use **CamelCase** for exported names (visible outside the package)
- Use **lowerCamelCase** for non-exported names
- Choose descriptive names; avoid abbreviations except for common ones
- Keep package names short, concise, and lowercase

#### File Organization

- One file per logical component
- Related files in the same package
- Test files next to the code they test (`foo.go` and `foo_test.go`)

### 4.2. Documentation Standards

#### Code Comments

- Every exported function, type, and constant must have a comment
- Comments should explain the "why" not just the "what"
- Use complete sentences starting with the function name

Example:

```go
// GetSystem retrieves a system by its ID. It returns nil and ErrSystemNotFound
// if the system doesn't exist, or another error if the database query fails.
func GetSystem(id string) (*System, error) {
    // ...
}
```

#### Package Documentation

Each package should have a package comment in one of its files:

```go
// Package system provides functionality for managing system entities,
// including creation, retrieval, updating, and deletion operations.
package system
```

### 4.3. Error Handling

- Always check errors
- Don't use `panic` in production code
- Add context to errors with `fmt.Errorf("doing X: %w", err)`
- Use custom error types for expected error conditions

```go
var (
    // ErrSystemNotFound is returned when a system with the specified ID does not exist.
    ErrSystemNotFound = errors.New("system not found")
)
```

### 4.4. Logging Guidelines

- Use the project's logging package (`internal/logging`)
- Log at appropriate levels (debug, info, warning, error, fatal)
- Include relevant context in log messages
- Avoid logging sensitive information

```go
// Incorrect:
log.Infof("User %s logged in with password %s", username, password)

// Correct:
log.Infof("User %s logged in successfully", username)
```

### 4.5. Testing Standards

- All packages should have tests
- Aim for at least 70% test coverage
- Use table-driven tests for multiple test cases
- Use mocks for external dependencies
- Test both success and error paths

```go
func TestGetSystem(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        wantErr bool
    }{
        {"existing system", "valid-id", false},
        {"non-existent system", "invalid-id", true},
        {"empty id", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 5. Building and Testing

### 5.1. Building the Project

The project uses Make for build automation. Common Make targets:

```bash
# Build the main binary
make build

# Build for all supported platforms
make build-all

# Build and install locally
make install

# Clean build artifacts
make clean
```

### 5.2. Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific tests
go test ./internal/system/...

# Run tests with verbose output
go test -v ./...
```

### 5.3. Linting and Code Quality

```bash
# Run all linters
make lint

# Format code
make fmt

# Check for common errors
make errcheck

# Run static analysis
make staticcheck
```

### 5.4. Building Docker Images

```bash
# Build development image
make docker-build-dev

# Build production image
make docker-build-prod

# Run using Docker
make docker-run
```

### 5.5. Continuous Integration

The project uses GitHub Actions for CI/CD. Each pull request triggers:

- Linting
- Unit tests
- Integration tests
- Build verification

See `.github/workflows/` for CI configuration.

## 6. Contribution Workflow

### 6.1. Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Add the upstream repository as a remote

```bash
git remote add upstream https://github.com/jdfalk/ubuntu-autoinstall-webhook.git
```

### 6.2. Creating a Feature Branch

Create a new branch for your feature or bugfix:

```bash
# Update your main branch
git checkout main
git pull upstream main

# Create a new branch
git checkout -b feature/your-feature-name
```

### 6.3. Making Changes

1. Make your changes in the feature branch
2. Add tests for your changes
3. Ensure all tests pass
4. Update documentation if needed

### 6.4. Commit Guidelines

- Use descriptive commit messages
- Start with a verb in the present tense (Add, Fix, Update, etc.)
- Reference issue numbers if applicable
- Keep changes focused and atomic

Example:

```
Add system discovery feature

- Implement DHCP event monitoring
- Add system registration endpoint
- Create database schema for discovered systems

Fixes #123
```

### 6.5. Submitting a Pull Request

1. Push your changes to your fork
2. Create a pull request against the upstream repository
3. Describe your changes in the PR description
4. Address any feedback from code reviews

### 6.6. Code Review Process

- All PRs require at least one review
- Address all review comments
- Ensure CI tests pass
- Maintain a respectful and collaborative tone

### 6.7. After Merge

After your PR is merged:

1. Update your local main:
   ```bash
   git checkout main
   git pull upstream main
   ```
2. Delete your feature branch:
   ```bash
   git branch -d feature/your-feature-name
   ```

## 7. API Development

### 7.1. API Design Principles

The API follows RESTful principles:

- Resource-based URLs
- Standard HTTP methods (GET, POST, PUT, DELETE)
- JSON request and response bodies
- Consistent error responses
- Stateless operations

### 7.2. Adding New Endpoints

To add a new API endpoint:

1. Define the endpoint in the appropriate file in `internal/api/`
2. Create handler functions for the endpoint
3. Register the endpoint in `internal/api/router.go`
4. Add authentication/authorization middleware if needed
5. Document the endpoint

Example:

```go
// internal/api/systems.go

// GetSystem handles GET requests to retrieve a system by ID.
func (h *SystemsHandler) GetSystem(w http.ResponseWriter, r *http.Request) {
    // Extract ID from request
    id := mux.Vars(r)["id"]

    // Call service
    system, err := h.systemService.GetSystem(id)
    if err != nil {
        if errors.Is(err, ErrSystemNotFound) {
            api.RespondNotFound(w, "System not found")
            return
        }
        api.RespondError(w, http.StatusInternalServerError, "Failed to get system")
        return
    }

    // Return response
    api.RespondJSON(w, http.StatusOK, system)
}

// RegisterRoutes registers all system API routes.
func (h *SystemsHandler) RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/systems/{id}", h.GetSystem).Methods("GET")
    // Other routes...
}
```

### 7.2. Endpoint Implementation

<!-- This section would cover endpoint implementation -->

### 7.3. Error Handling in API

Consistent error handling is crucial for creating a reliable and user-friendly
API. This section outlines the error handling patterns and best practices for
the Ubuntu Autoinstall Webhook API.

#### 7.3.1. Error Response Structure

All API error responses follow a standard JSON format:

```json
{
  "status": "error",
  "error": {
    "code": "resource_not_found",
    "message": "The requested system was not found",
    "details": {
      "resource_id": "550e8400-e29b-41d4-a716-446655440000",
      "resource_type": "system"
    }
  }
}
```

Key components:

- `status`: Always "error" for error responses
- `error.code`: A machine-readable error code (snake_case)
- `error.message`: A human-readable error message
- `error.details`: Optional object with additional error context

#### 7.3.2. Standard Error Codes

Use these standard error codes consistently across the API:

| HTTP Status | Error Code                | Usage                            |
| ----------- | ------------------------- | -------------------------------- |
| 400         | `invalid_request`         | General validation error         |
| 400         | `missing_required_field`  | A required field is missing      |
| 400         | `invalid_field_format`    | A field has invalid format       |
| 401         | `authentication_required` | Authentication is missing        |
| 401         | `invalid_credentials`     | Authentication failed            |
| 401         | `token_expired`           | Auth token has expired           |
| 403         | `permission_denied`       | User lacks permission            |
| 404         | `resource_not_found`      | Requested resource doesn't exist |
| 409         | `resource_conflict`       | Resource state conflict          |
| 422         | `validation_failed`       | Semantic validation failed       |
| 429         | `too_many_requests`       | Rate limit exceeded              |
| 500         | `server_error`            | Unexpected server error          |
| 503         | `service_unavailable`     | Service temporarily unavailable  |

#### 7.3.3. Implementing Error Handling

Use the provided error handling utilities in `internal/api/errors.go`:

```go
// Example of returning a standard error from a handler
func (h *SystemsHandler) GetSystem(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    system, err := h.systemService.GetSystem(id)
    if err != nil {
        // Check for specific error types
        if errors.Is(err, ErrSystemNotFound) {
            // Return standard 404 error
            api.WriteNotFoundError(w, "system_not_found", "System not found", map[string]interface{}{
                "system_id": id,
            })
            return
        }

        // Log unexpected errors and return 500
        h.logger.WithError(err).Error("Failed to retrieve system")
        api.WriteInternalError(w, "Failed to retrieve system")
        return
    }

    api.WriteJSON(w, http.StatusOK, api.SuccessResponse(system))
}
```

#### 7.3.4. Error Helper Functions

Use these helper functions from the `api` package:

```go
// Write a 400 Bad Request error
api.WriteBadRequestError(w, code, message, details)

// Write a 401 Unauthorized error
api.WriteUnauthorizedError(w, code, message, details)

// Write a 403 Forbidden error
api.WriteForbiddenError(w, code, message, details)

// Write a 404 Not Found error
api.WriteNotFoundError(w, code, message, details)

// Write a 409 Conflict error
api.WriteConflictError(w, code, message, details)

// Write a 422 Unprocessable Entity error
api.WriteValidationError(w, code, message, details)

// Write a 429 Too Many Requests error
api.WriteRateLimitError(w, code, message, details)

// Write a 500 Internal Server Error
api.WriteInternalError(w, message)
```

#### 7.3.5. Validation Errors

For request validation errors, provide detailed information:

```go
func (h *SystemsHandler) CreateSystem(w http.ResponseWriter, r *http.Request) {
    var req CreateSystemRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        api.WriteBadRequestError(w, "invalid_json", "Invalid JSON in request body", nil)
        return
    }

    // Validate request fields
    errors := validateCreateSystemRequest(req)
    if len(errors) > 0 {
        api.WriteValidationError(w, "validation_failed", "Request validation failed", map[string]interface{}{
            "validation_errors": errors,
        })
        return
    }

    // Proceed with creating the system...
}

func validateCreateSystemRequest(req CreateSystemRequest) []map[string]string {
    var errors []map[string]string

    if req.Hostname == "" {
        errors = append(errors, map[string]string{
            "field": "hostname",
            "error": "Hostname is required",
        })
    }

    if req.MacAddress == "" {
        errors = append(errors, map[string]string{
            "field": "mac_address",
            "error": "MAC address is required",
        })
    } else if !isValidMacAddress(req.MacAddress) {
        errors = append(errors, map[string]string{
            "field": "mac_address",
            "error": "Invalid MAC address format",
        })
    }

    return errors
}
```

#### 7.3.6. Domain Errors vs. API Errors

Distinguish between domain errors and API presentation errors:

1. **Domain Errors**: Generated in service and repository layers
2. **API Errors**: How errors are presented to API clients

Map domain errors to API errors in your handlers:

```go
// In service layer: Domain-specific error
var ErrSystemNotFound = errors.New("system not found")

// In handler: Map domain error to API error
if errors.Is(err, services.ErrSystemNotFound) {
    api.WriteNotFoundError(w, "system_not_found", "System not found", map[string]interface{}{
        "system_id": id,
    })
    return
}
```

#### 7.3.7. Error Logging

Follow these guidelines for error logging:

1. **Log Unexpected Errors**: Always log unexpected errors at error level
2. **Include Context**: Add relevant request context to logs
3. **Don't Log Expected Errors**: No need to log 404s and validation errors as
   errors
4. **Include Stack Trace**: For internal errors, include stack traces

Example:

```go
if err != nil {
    // Expected error - don't log as error
    if errors.Is(err, ErrSystemNotFound) {
        api.WriteNotFoundError(w, "system_not_found", "System not found", nil)
        return
    }

    // Unexpected error - log with context
    h.logger.WithFields(log.Fields{
        "system_id": id,
        "client_ip": r.RemoteAddr,
        "error": err.Error(),
    }).Error("Failed to retrieve system")

    api.WriteInternalError(w, "Failed to retrieve system")
    return
}
```

#### 7.3.8. Security Considerations

1. **Don't Leak Implementation Details**: Internal error messages should not be
   exposed to clients
2. **Be Careful with Error Messages**: Error messages should not reveal
   sensitive information
3. **Use Generic Messages for 500 Errors**: Don't include detailed error
   information in 500 responses
4. **Log Sensitive Errors**: Log the full error internally for debugging

#### 7.3.9. Testing Error Conditions

Always write tests for error conditions:

```go
func TestGetSystem_NotFound(t *testing.T) {
    // Arrange
    mockService := &MockSystemService{}
    mockService.On("GetSystem", "not-found-id").Return(nil, ErrSystemNotFound)

    handler := NewSystemsHandler(mockService, mockLogger)
    req := httptest.NewRequest("GET", "/api/v1/systems/not-found-id", nil)
    res := httptest.NewRecorder()

    // Act
    handler.GetSystem(res, req)

    // Assert
    assert.Equal(t, http.StatusNotFound, res.Code)

    var response api.ErrorResponse
    err := json.Unmarshal(res.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "error", response.Status)
    assert.Equal(t, "system_not_found", response.Error.Code)
}
```

### 7.4. API Authentication

The Ubuntu Autoinstall Webhook API uses a robust authentication system to secure
endpoints. This section covers implementation details for API authentication
mechanisms.

#### 7.4.1. Authentication Methods

The API supports three authentication methods:

1. **API Token Authentication**: Long-lived tokens for service accounts and
   automation
2. **JWT Authentication**: Short-lived JSON Web Tokens for user sessions
3. **Basic Authentication**: Username/password for initial authentication only

#### 7.4.2. API Token Implementation

API tokens provide a simple way for services to authenticate with the API.

**Token Format**:

- 64-character randomly generated string
- Prefixed with `uaw_` to identify token type
- Example: `uaw_a1b2c3d4e5f6...`

**Implementation**:

```go
// In internal/auth/token.go
func GenerateAPIToken() (string, error) {
    // Generate 32 bytes of random data
    randomBytes := make([]byte, 32)
    if _, err := rand.Read(randomBytes); err != nil {
        return "", fmt.Errorf("failed to generate random bytes: %w", err)
    }

    // Encode to hex and add prefix
    token := fmt.Sprintf("uaw_%x", randomBytes)
    return token, nil
}
```

**Token Storage**:

- Store only the hashed value of the token in the database
- Use bcrypt for hashing tokens

```go
// Hash a token for storage
func HashToken(token string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash token: %w", err)
    }
    return string(hash), nil
}

// Verify a token against its hash
func VerifyToken(token, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
    return err == nil
}
```

#### 7.4.3. JWT Authentication

JWT tokens are used for authenticating user sessions.

**JWT Claims**:

```go
type Claims struct {
    UserID    string   `json:"uid"`
    Username  string   `json:"username"`
    Email     string   `json:"email"`
    Roles     []string `json:"roles"`
    jwt.RegisteredClaims
}
```

**Generating JWT Tokens**:

```go
// In internal/auth/jwt.go
func GenerateJWT(user *models.User) (string, error) {
    // Set expiration time
    expirationTime := time.Now().Add(1 * time.Hour)

    // Create claims
    claims := &Claims{
        UserID:    user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Roles:     user.Roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "ubuntu-autoinstall-webhook",
        },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with the secret key
    tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
    if err != nil {
        return "", fmt.Errorf("failed to sign JWT: %w", err)
    }

    return tokenString, nil
}
```

**Validating JWT Tokens**:

```go
func ValidateJWT(tokenString string) (*Claims, error) {
    // Parse token
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            // Validate signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }

            return []byte(config.GetJWTSecret()), nil
        },
    )

    if err != nil {
        return nil, err
    }

    // Extract claims
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
```

#### 7.4.4. Authentication Middleware

The API uses middleware to authenticate requests:

```go
// In internal/api/middleware/auth.go
func TokenAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            api.WriteUnauthorizedError(w, "authentication_required",
                "Authentication required", nil)
            return
        }

        // Check for Bearer prefix
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            api.WriteUnauthorizedError(w, "invalid_auth_format",
                "Invalid authorization format", nil)
            return
        }

        token := parts[1]

        // Check if it's an API token (prefixed with uaw_)
        if strings.HasPrefix(token, "uaw_") {
            // Validate API token
            valid, userContext := auth.ValidateAPIToken(token)
            if !valid {
                api.WriteUnauthorizedError(w, "invalid_token",
                    "Invalid API token", nil)
                return
            }

            // Add user context to request
            ctx := context.WithValue(r.Context(), auth.ContextKeyUser, userContext)
            next.ServeHTTP(w, r.WithContext(ctx))
            return
        }

        // Assume JWT token
        claims, err := auth.ValidateJWT(token)
        if err != nil {
            if errors.Is(err, jwt.ErrTokenExpired) {
                api.WriteUnauthorizedError(w, "token_expired",
                    "Token has expired", nil)
            } else {
                api.WriteUnauthorizedError(w, "invalid_token",
                    "Invalid authentication token", nil)
            }
            return
        }

        // Create user context from JWT claims
        userContext := &auth.UserContext{
            UserID:   claims.UserID,
            Username: claims.Username,
            Email:    claims.Email,
            Roles:    claims.Roles,
        }

        // Add user context to request
        ctx := context.WithValue(r.Context(), auth.ContextKeyUser, userContext)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

#### 7.4.5. Authorization Middleware

In addition to authentication, implement authorization to check if authenticated
users have the necessary permissions:

```go
// In internal/api/middleware/auth.go
func RequireRoleMiddleware(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get user from context
            userContext, ok := r.Context().Value(auth.ContextKeyUser).(*auth.UserContext)
            if !ok {
                api.WriteUnauthorizedError(w, "authentication_required",
                    "Authentication required", nil)
                return
            }

            // Check if user has any of the required roles
            hasRole := false
            for _, requiredRole := range roles {
                for _, userRole := range userContext.Roles {
                    if requiredRole == userRole {
                        hasRole = true
                        break
                    }
                }
                if hasRole {
                    break
                }
            }

            if !hasRole {
                api.WriteForbiddenError(w, "permission_denied",
                    "You don't have permission to perform this action", nil)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

#### 7.4.6. Securing API Routes

Apply authentication and authorization middleware to API routes:

```go
// In internal/api/router.go
func SetupRouter(deps *Dependencies) *mux.Router {
    router := mux.NewRouter()

    // Public endpoints (no auth required)
    router.HandleFunc("/api/v1/health", handlers.HealthCheckHandler).Methods("GET")
    router.HandleFunc("/api/v1/auth/login", deps.AuthHandler.Login).Methods("POST")

    // Create API subrouter with authentication
    apiRouter := router.PathPrefix("/api/v1").Subrouter()
    apiRouter.Use(middleware.TokenAuthMiddleware)

    // Routes for authenticated users
    apiRouter.HandleFunc("/systems", deps.SystemsHandler.ListSystems).Methods("GET")
    apiRouter.HandleFunc("/systems/{id}", deps.SystemsHandler.GetSystem).Methods("GET")

    // Routes requiring admin role
    adminRouter := apiRouter.PathPrefix("/").Subrouter()
    adminRouter.Use(middleware.RequireRoleMiddleware("admin"))

    adminRouter.HandleFunc("/users", deps.UsersHandler.ListUsers).Methods("GET")
    adminRouter.HandleFunc("/users", deps.UsersHandler.CreateUser).Methods("POST")

    return router
}
```

#### 7.4.7. Testing Authentication

Write tests for authentication logic:

```go
func TestTokenAuthMiddleware(t *testing.T) {
    // Create mock token validator
    mockAuth := &MockAuthService{}
    mockAuth.On("ValidateAPIToken", "uaw_validtoken123").Return(true, &auth.UserContext{
        UserID:   "123",
        Username: "testuser",
        Roles:    []string{"user"},
    })
    mockAuth.On("ValidateAPIToken", "uaw_invalidtoken").Return(false, nil)

    // Create test handler
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if user context is set correctly
        user := r.Context().Value(auth.ContextKeyUser).(*auth.UserContext)
        assert.Equal(t, "testuser", user.Username)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("success"))
    })

    // Create middleware with mock auth
    middleware := middleware.CreateTokenAuthMiddleware(mockAuth)
    handlerToTest := middleware(nextHandler)

    // Test valid token
    req := httptest.NewRequest("GET", "/api/v1/systems", nil)
    req.Header.Set("Authorization", "Bearer uaw_validtoken123")
    res := httptest.NewRecorder()

    handlerToTest.ServeHTTP(res, req)

    assert.Equal(t, http.StatusOK, res.Code)
    assert.Equal(t, "success", res.Body.String())

    // Test invalid token
    req = httptest.NewRequest("GET", "/api/v1/systems", nil)
    req.Header.Set("Authorization", "Bearer uaw_invalidtoken")
    res = httptest.NewRecorder()

    handlerToTest.ServeHTTP(res, req)

    assert.Equal(t, http.StatusUnauthorized, res.Code)
}
```

### 7.5. API Versioning

API versioning ensures backward compatibility as the API evolves. The Ubuntu
Autoinstall Webhook uses URL-based versioning.

#### 7.5.1. Versioning Strategy

The API version is specified in the URL path:

- Current version: `/api/v1/...`
- Future versions: `/api/v2/...`, `/api/v3/...`

#### 7.5.2. Implementing Versioned Routes

Organize route handlers by version:

```go
// In internal/api/router.go
func SetupRouter(deps *Dependencies) *mux.Router {
    router := mux.NewRouter()

    // Set up v1 routes
    v1Router := router.PathPrefix("/api/v1").Subrouter()
    setupV1Routes(v1Router, deps)

    // In the future, add v2 routes
    // v2Router := router.PathPrefix("/api/v2").Subrouter()
    // setupV2Routes(v2Router, deps)

    return router
}

func setupV1Routes(router *mux.Router, deps *Dependencies) {
    // Health check
    router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

    // Auth routes
    router.HandleFunc("/auth/login", deps.AuthHandler.Login).Methods("POST")
    router.HandleFunc("/auth/refresh", deps.AuthHandler.RefreshToken).Methods("POST")

    // Add authentication middleware for protected routes
    protectedRouter := router.PathPrefix("/").Subrouter()
    protectedRouter.Use(middleware.TokenAuthMiddleware)

    // Systems routes
    protectedRouter.HandleFunc("/systems", deps.SystemsHandler.ListSystems).Methods("GET")
    protectedRouter.HandleFunc("/systems", deps.SystemsHandler.CreateSystem).Methods("POST")
    protectedRouter.HandleFunc("/systems/{id}", deps.SystemsHandler.GetSystem).Methods("GET")
    // Other routes...
}
```

#### 7.5.3. Version Compatibility

When making changes to the API:

1. **Non-Breaking Changes**: Add to existing version
   - Adding new endpoints
   - Adding optional request fields
   - Adding new response fields

2. **Breaking Changes**: Require a new version
   - Removing endpoints
   - Removing or renaming fields
   - Changing field types
   - Changing response structure

#### 7.5.4. Supporting Multiple Versions

When supporting multiple API versions:

1. Create separate route handlers for each version
2. Share business logic between versions when possible
3. Use interface adapters to transform between versions

```go
// V1 handler
func (h *SystemsHandlerV1) GetSystem(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    // Call common service
    system, err := h.systemService.GetSystem(id)
    if err != nil {
        // Handle error...
        return
    }

    // Transform to V1 response format
    responseV1 := transformToV1SystemResponse(system)
    api.WriteJSON(w, http.StatusOK, api.SuccessResponse(responseV1))
}

// V2 handler
func (h *SystemsHandlerV2) GetSystem(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    // Call the same service
    system, err := h.systemService.GetSystem(id)
    if err != nil {
        // Handle error...
        return
    }

    // Transform to V2 response format (with additional fields)
    responseV2 := transformToV2SystemResponse(system)
    api.WriteJSON(w, http.StatusOK, api.SuccessResponse(responseV2))
}
```

### 7.6. API Documentation

Comprehensive API documentation is essential for developers using the API. This
section covers how to document the API using OpenAPI/Swagger.

#### 7.6.1. OpenAPI Specification

The Ubuntu Autoinstall Webhook API is documented using OpenAPI 3.0. The
specification is maintained in YAML format.

**Basic Structure**:

```yaml
# in api/openapi/v1.yaml
openapi: 3.0.3
info:
  title: Ubuntu Autoinstall Webhook API
  description: API for managing Ubuntu automated installations
  version: 1.0.0
  contact:
    name: Ubuntu Autoinstall Webhook Team
    url: https://github.com/jdfalk/ubuntu-autoinstall-webhook
servers:
  - url: https://webhook.example.com/api/v1
    description: Production API
  - url: http://localhost:8080/api/v1
    description: Development API
paths:
  /systems:
    get:
      summary: List all systems
      description: Returns a list of all registered systems
      operationId: listSystems
      parameters:
        - name: status
          in: query
          description: Filter systems by status
          required: false
          schema:
            type: string
            enum: [discovered, ready, installing, completed, failed]
      responses:
        '200':
          description: Successful operation
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SystemList'
        '401':
          $ref: '#/components/responses/Unauthorized'
      security:
        - BearerAuth: []
    # Other methods...
```

#### 7.6.2. Documenting Endpoints

Document each endpoint with:

1. **Summary & Description**: Brief and detailed explanations
2. **Parameters**: Path, query, and header parameters
3. **Request Body**: Schema for request data
4. **Responses**: All possible response codes and their schemas
5. **Security**: Authentication requirements

```yaml
/systems/{id}:
  get:
    summary: Get a system by ID
    description: Returns detailed information about a specific system
    operationId: getSystem
    parameters:
      - name: id
        in: path
        description: System ID
        required: true
        schema:
          type: string
          format: uuid
    responses:
      '200':
        description: Successful operation
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/SystemResponse'
      '404':
        description: System not found
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ErrorResponse'
            example:
              status: error
              error:
                code: resource_not_found
                message: The requested system was not found
                details:
                  system_id: 550e8400-e29b-41d4-a716-446655440000
    security:
      - BearerAuth: []
```

#### 7.6.3. Schema Definitions

Define all data models in the components/schemas section:

```yaml
components:
  schemas:
    System:
      type: object
      properties:
        id:
          type: string
          format: uuid
          description: Unique system identifier
        hostname:
          type: string
          description: System hostname
        mac_address:
          type: string
          description: System MAC address
          pattern: '^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$'
        ip_address:
          type: string
          format: ipv4
          description: System IP address
        status:
          type: string
          enum: [discovered, ready, installing, completed, failed]
          description: Current system status
        template_id:
          type: string
          format: uuid
          description: ID of the assigned installation template
        created_at:
          type: string
          format: date-time
          description: Timestamp when the system was created
        updated_at:
          type: string
          format: date-time
          description: Timestamp when the system was last updated
      required:
        - id
        - hostname
        - mac_address
        - status
```

#### 7.6.4. Generating API Documentation

Use the OpenAPI specification to generate interactive documentation:

1. Install Swagger UI as part of the API:

```go
// In internal/api/docs.go
import (
    "net/http"

    "github.com/go-openapi/runtime/middleware"
)

// SetupSwaggerUI configures Swagger UI endpoints
func SetupSwaggerUI(router *mux.Router) {
    // Serve OpenAPI specification
    router.PathPrefix("/api/swagger.json").Handler(http.FileServer(http.Dir("./api/openapi")))

    // Serve Swagger UI
    opts := middleware.SwaggerUIOpts{
        SpecURL: "/api/swagger.json",
        Path:    "/api/docs",
    }
    router.PathPrefix("/api/docs").Handler(middleware.SwaggerUI(opts, nil))
}
```

2. Add Swagger UI to the main router:

```go
func SetupRouter(deps *Dependencies) *mux.Router {
    router := mux.NewRouter()

    // Setup API routes
    v1Router := router.PathPrefix("/api/v1").Subrouter()
    setupV1Routes(v1Router, deps)

    // Setup Swagger UI
    SetupSwaggerUI(router)

    return router
}
```

#### 7.6.5. Keeping Documentation in Sync

Ensure API documentation stays in sync with implementation:

1. **Review Process**: Include OpenAPI spec updates in code reviews
2. **Testing**: Test that API responses match documented schemas
3. **Automation**: Consider tools to validate API responses against OpenAPI
   schemas

### 7.7. Rate Limiting

Implement rate limiting to protect the API from abuse and ensure fair resource
usage.

#### 7.7.1. Rate Limiting Strategy

The API uses a token bucket algorithm for rate limiting with different limits
based on:

- Authentication status
- User role
- Endpoint sensitivity

#### 7.7.2. Implementing Rate Limiting

Use the `golang.org/x/time/rate` package to implement rate limiting:

```go
// In internal/api/middleware/ratelimit.go
import (
    "net/http"
    "sync"

    "golang.org/x/time/rate"
)

// RateLimiter stores limiters for different clients
type RateLimiter struct {
    visitors map[string]*rate.Limiter
    mu       sync.Mutex
    limit    rate.Limit
    burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
    return &RateLimiter{
        visitors: make(map[string]*rate.Limiter),
        limit:    rate.Limit(rps),
        burst:    burst,
    }
}

// GetLimiter returns a limiter for a specific client
func (rl *RateLimiter) GetLimiter(clientIP string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    limiter, exists := rl.visitors[clientIP]
    if !exists {
        limiter = rate.NewLimiter(rl.limit, rl.burst)
        rl.visitors[clientIP] = limiter
    }

    return limiter
}

// RateLimitMiddleware creates middleware that limits request rate
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get client IP
            clientIP := r.RemoteAddr

            // Get rate limiter for this client
            limiter := limiter.GetLimiter(clientIP)

            // Check if request is allowed
            if !limiter.Allow() {
                api.WriteRateLimitError(w, "too_many_requests",
                    "Rate limit exceeded", nil)
                return
            }

            // Continue processing the request
            next.ServeHTTP(w, r)
        })
    }
}
```

#### 7.7.3. Different Rate Limits by Role

Implement different rate limits based on authentication status and roles:

```go
func SetupRateLimiters() map[string]*RateLimiter {
    limiters := make(map[string]*RateLimiter)

    // Default limiter for unauthenticated users (5 requests per second, burst of 10)
    limiters["default"] = NewRateLimiter(5, 10)

    // Higher limits for authenticated users (20 requests per second, burst of 50)
    limiters["authenticated"] = NewRateLimiter(20, 50)

    // Even higher limits for admin users (50 requests per second, burst of 100)
    limiters["admin"] = NewRateLimiter(50, 100)

    return limiters
}

func DynamicRateLimitMiddleware(limiters map[string]*RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get client IP
            clientIP := r.RemoteAddr

            // Determine which limiter to use based on authentication
            limiterKey := "default"

            // If user is authenticated, get their context
            if userContext, ok := r.Context().Value(auth.ContextKeyUser).(*auth.UserContext); ok {
                limiterKey = "authenticated"

                // Check for admin role
                for _, role := range userContext.Roles {
                    if role == "admin" {
                        limiterKey = "admin"
                        break
                    }
                }
            }

            // Create composite key from IP and limiter type
            key := clientIP + ":" + limiterKey

            // Get appropriate limiter
            limiter := limiters[limiterKey].GetLimiter(key)

            // Check if request is allowed
            if !limiter.Allow() {
                api.WriteRateLimitError(w, "too_many_requests",
                    "Rate limit exceeded", nil)
                return
            }

            // Continue processing the request
            next.ServeHTTP(w, r)
        })
    }
}
```

#### 7.7.4. Rate Limit Headers

Include rate limit information in response headers:

```go
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            clientIP := r.RemoteAddr
            limiter := limiter.GetLimiter(clientIP)

            // Check if request is allowed
            if !limiter.Allow() {
                api.WriteRateLimitError(w, "too_many_requests",
                    "Rate limit exceeded", nil)
                return
            }

            // Add rate limit headers
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Burst()))
            w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Burst()-1)) // Approximate
            w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))

            next.ServeHTTP(w, r)
        })
    }
}
```

## 8. Database Development

### 8.1. Database Architecture

The Ubuntu Autoinstall Webhook uses a flexible database architecture that
supports multiple database backends while maintaining consistent access
patterns.

#### 8.1.1. Database Abstraction Layer

The system follows a layered approach to database access:

1. **Models Layer**: Domain entities and data structures
2. **Repository Layer**: Interface for data access operations
3. **Database Implementation Layer**: Concrete implementations for different
   database backends

This architecture allows for:

- Easy switching between database backends
- Consistent access patterns across the application
- Simplified testing with mock repositories

#### 8.1.2. Supported Database Backends

The system primarily supports:

1. **SQLite**: Default for development and small deployments
2. **CockroachDB**: For distributed, high-availability deployments

#### 8.1.3. Component Structure

```
internal/database/
├── models/           # Domain model definitions
│   ├── system.go
│   ├── template.go
│   └── ...
├── database.go       # Database interface definition
├── sqlite.go         # SQLite implementation
├── cockroach.go      # CockroachDB implementation
└── migrations/       # Database schema migrations
    ├── sqlite/
    └── cockroach/
```

#### 8.1.4. Database Interface

The core database interface provides a factory for repositories:

```go
// In internal/database/database.go
type Database interface {
    // Get repositories
    SystemRepository() SystemRepository
    TemplateRepository() TemplateRepository
    InstallationRepository() InstallationRepository
    UserRepository() UserRepository

    // Database management
    Connect() error
    Close() error
    Migrate() error
    Ping() error
}

// NewDatabase creates a database instance based on configuration
func NewDatabase(config *config.Config) (Database, error) {
    switch config.Database.Type {
    case "sqlite":
        return NewSQLiteDatabase(config)
    case "cockroach":
        return NewCockroachDatabase(config)
    default:
        return nil, fmt.Errorf("unsupported database type: %s", config.Database.Type)
    }
}
```

### 8.2. Model Definition

Models represent the core data entities in the system. They serve as both
database records and domain objects.

#### 8.2.1. Basic Model Structure

All models follow a consistent structure with common fields:

```go
// In internal/database/models/base.go
type BaseModel struct {
    ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Generate a new ID for models
func (m *BaseModel) GenerateID() {
    if m.ID == "" {
        m.ID = uuid.New().String()
    }
}
```

#### 8.2.2. Entity Model Example

Example of a specific entity model:

```go
// In internal/database/models/system.go
type System struct {
    BaseModel
    Hostname    string     `json:"hostname" gorm:"uniqueIndex"`
    MacAddress  string     `json:"mac_address" gorm:"uniqueIndex"`
    IPAddress   string     `json:"ip_address"`
    Status      string     `json:"status" gorm:"index"`
    TemplateID  *string    `json:"template_id"`
    LastSeen    *time.Time `json:"last_seen"`
    Description string     `json:"description"`
    Metadata    JSON       `json:"metadata" gorm:"type:jsonb"`
}

// Validate performs validation on the system model
func (s *System) Validate() error {
    if s.Hostname == "" {
        return errors.New("hostname is required")
    }

    if s.MacAddress == "" {
        return errors.New("mac address is required")
    }

    // Validate MAC address format
    if !IsValidMacAddress(s.MacAddress) {
        return errors.New("invalid mac address format")
    }

    return nil
}

// Status constants
const (
    SystemStatusDiscovered = "discovered"
    SystemStatusReady      = "ready"
    SystemStatusInstalling = "installing"
    SystemStatusCompleted  = "completed"
    SystemStatusFailed     = "failed"
)
```

#### 8.2.3. Custom Data Types

Define custom types for special data formats:

```go
// In internal/database/models/json.go
import (
    "database/sql/driver"
    "encoding/json"
    "errors"
)

// JSON is a custom type for handling JSON data in database
type JSON map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
    if j == nil {
        return nil, nil
    }
    return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
    if value == nil {
        *j = nil
        return nil
    }

    var data []byte
    switch v := value.(type) {
    case []byte:
        data = v
    case string:
        data = []byte(v)
    default:
        return errors.New("type assertion to []byte failed")
    }

    return json.Unmarshal(data, &j)
}
```

#### 8.2.4. Model Relationships

Define relationships between models:

```go
// In internal/database/models/installation.go
type Installation struct {
    BaseModel
    SystemID    string       `json:"system_id" gorm:"index"`
    TemplateID  string       `json:"template_id" gorm:"index"`
    Status      string       `json:"status" gorm:"index"`
    StartedAt   time.Time    `json:"started_at"`
    CompletedAt *time.Time   `json:"completed_at"`
    LogFile     string       `json:"log_file"`
    ErrorReason *string      `json:"error_reason"`
    Variables   JSON         `json:"variables" gorm:"type:jsonb"`

    // Relationships
    System      *System      `json:"system" gorm:"foreignKey:SystemID"`
    Template    *Template    `json:"template" gorm:"foreignKey:TemplateID"`
}
```

### 8.3. Repository Pattern

The repository pattern provides a clean abstraction for database operations
while keeping business logic separate from data access.

#### 8.3.1. Repository Interfaces

Define repository interfaces for each entity:

```go
// In internal/database/repositories.go
type SystemRepository interface {
    // Basic CRUD operations
    Create(system *models.System) error
    GetByID(id string) (*models.System, error)
    GetByMacAddress(mac string) (*models.System, error)
    Update(system *models.System) error
    Delete(id string) error

    // Query operations
    List(filter SystemFilter) ([]*models.System, error)
    Count(filter SystemFilter) (int64, error)

    // Specialized operations
    UpdateStatus(id string, status string) error
}

type SystemFilter struct {
    Status     string
    Hostname   string
    MacAddress string
    IPAddress  string
    TemplateID string
    Limit      int
    Offset     int
    SortBy     string
    SortDesc   bool
}
```

#### 8.3.2. SQLite Implementation

Implement the repository interface for SQLite:

```go
// In internal/database/sqlite_system_repository.go
type sqliteSystemRepository struct {
    db *gorm.DB
}

// Ensure interface compliance
var _ SystemRepository = (*sqliteSystemRepository)(nil)

func NewSQLiteSystemRepository(db *gorm.DB) SystemRepository {
    return &sqliteSystemRepository{db: db}
}

func (r *sqliteSystemRepository) Create(system *models.System) error {
    // Generate ID if not provided
    system.GenerateID()

    // Perform validation
    if err := system.Validate(); err != nil {
        return err
    }

    // Insert into database
    result := r.db.Create(system)
    return result.Error
}

func (r *sqliteSystemRepository) GetByID(id string) (*models.System, error) {
    var system models.System
    result := r.db.Where("id = ?", id).First(&system)

    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, ErrSystemNotFound
        }
        return nil, result.Error
    }

    return &system, nil
}

func (r *sqliteSystemRepository) List(filter SystemFilter) ([]*models.System, error) {
    var systems []*models.System
    query := r.db.Model(&models.System{})

    // Apply filters
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }

    if filter.Hostname != "" {
        query = query.Where("hostname LIKE ?", "%"+filter.Hostname+"%")
    }

    if filter.MacAddress != "" {
        query = query.Where("mac_address = ?", filter.MacAddress)
    }

    // Apply sorting
    if filter.SortBy != "" {
        order := filter.SortBy
        if filter.SortDesc {
            order += " DESC"
        } else {
            order += " ASC"
        }
        query = query.Order(order)
    } else {
        // Default sorting by created_at
        query = query.Order("created_at DESC")
    }

    // Apply pagination
    if filter.Limit > 0 {
        query = query.Limit(filter.Limit)
    }

    if filter.Offset > 0 {
        query = query.Offset(filter.Offset)
    }

    result := query.Find(&systems)
    return systems, result.Error
}

// Implement other methods...
```

#### 8.3.3. CockroachDB Implementation

For CockroachDB, leverage GORM's compatibility while handling
CockroachDB-specific features:

```go
// In internal/database/cockroach_system_repository.go
type cockroachSystemRepository struct {
    db *gorm.DB
}

func NewCockroachSystemRepository(db *gorm.DB) SystemRepository {
    return &cockroachSystemRepository{db: db}
}

// Most methods can be similar to SQLite, but handle CockroachDB specifics
func (r *cockroachSystemRepository) Create(system *models.System) error {
    // Generate ID if not provided
    system.GenerateID()

    // Perform validation
    if err := system.Validate(); err != nil {
        return err
    }

    // Use a transaction for better consistency
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Check for unique constraint violations explicitly
        var existing models.System
        result := tx.Where("hostname = ? OR mac_address = ?",
                          system.Hostname, system.MacAddress).
                   First(&existing)

        if result.Error == nil {
            // Record exists, check which constraint would be violated
            if existing.Hostname == system.Hostname {
                return ErrDuplicateHostname
            }
            return ErrDuplicateMacAddress
        } else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
            // Unexpected error
            return result.Error
        }

        // Create the record
        result = tx.Create(system)
        return result.Error
    })
}

// Implement other methods following same pattern...
```

### 8.4. Database Migrations

Database migrations enable schema evolution while preserving data integrity.

#### 8.4.1. Migration Structure

Organize migrations by database type and version:

```
internal/database/migrations/
├── sqlite/
│   ├── 001_initial_schema.sql
│   ├── 002_add_metadata_column.sql
│   └── 003_add_foreign_keys.sql
└── cockroach/
    ├── 001_initial_schema.sql
    ├── 002_add_metadata_column.sql
    └── 003_add_foreign_keys.sql
```

#### 8.4.2. Migration Implementation

Implement migration logic in the database interface:

```go
// In internal/database/sqlite.go
func (db *sqliteDatabase) Migrate() error {
    // Create migrations table if it doesn't exist
    if err := db.createMigrationsTable(); err != nil {
        return err
    }

    // Get applied migrations
    appliedMigrations, err := db.getAppliedMigrations()
    if err != nil {
        return err
    }

    // Get available migrations
    migrations, err := db.getAvailableMigrations()
    if err != nil {
        return err
    }

    // Apply missing migrations in order
    for _, migration := range migrations {
        if !contains(appliedMigrations, migration.Name) {
            log.Infof("Applying migration: %s", migration.Name)

            if err := db.applyMigration(migration); err != nil {
                return fmt.Errorf("failed to apply migration %s: %w", migration.Name, err)
            }

            log.Infof("Migration applied: %s", migration.Name)
        }
    }

    return nil
}

// createMigrationsTable creates the migrations tracking table
func (db *sqliteDatabase) createMigrationsTable() error {
    query := `
        CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `

    _, err := db.db.Exec(query)
    return err
}

// Migration represents a database migration file
type Migration struct {
    Name     string
    Contents string
}

// applyMigration applies a single migration and records it
func (db *sqliteDatabase) applyMigration(migration Migration) error {
    // Use a transaction to ensure atomicity
    tx, err := db.db.Begin()
    if err != nil {
        return err
    }

    // Apply the migration SQL
    if _, err := tx.Exec(migration.Contents); err != nil {
        tx.Rollback()
        return err
    }

    // Record the migration
    if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", migration.Name); err != nil {
        tx.Rollback()
        return err
    }

    // Commit the transaction
    return tx.Commit()
}

// getAvailableMigrations reads migration files from the filesystem
func (db *sqliteDatabase) getAvailableMigrations() ([]Migration, error) {
    migrationDir := "internal/database/migrations/sqlite"
    files, err := os.ReadDir(migrationDir)
    if err != nil {
        return nil, err
    }

    var migrations []Migration
    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
            continue
        }

        content, err := os.ReadFile(filepath.Join(migrationDir, file.Name()))
        if err != nil {
            return nil, err
        }

        migrations = append(migrations, Migration{
            Name:     file.Name(),
            Contents: string(content),
        })
    }

    // Sort migrations by name (assuming numeric prefix)
    sort.Slice(migrations, func(i, j int) bool {
        return migrations[i].Name < migrations[j].Name
    })

    return migrations, nil
}

// getAppliedMigrations fetches previously applied migrations
func (db *sqliteDatabase) getAppliedMigrations() ([]string, error) {
    rows, err := db.db.Query("SELECT name FROM migrations ORDER BY id")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var migrations []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, err
        }
        migrations = append(migrations, name)
    }

    return migrations, rows.Err()
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

#### 8.4.3. Migration Example

Example of an actual migration file:

```sql
-- internal/database/migrations/sqlite/001_initial_schema.sql

-- Systems table
CREATE TABLE systems (
    id TEXT PRIMARY KEY,
    hostname TEXT UNIQUE NOT NULL,
    mac_address TEXT UNIQUE NOT NULL,
    ip_address TEXT,
    status TEXT NOT NULL DEFAULT 'discovered',
    template_id TEXT,
    last_seen TIMESTAMP,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_systems_status ON systems(status);
CREATE INDEX idx_systems_template_id ON systems(template_id);

-- Templates table
CREATE TABLE templates (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    ubuntu_version TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Installations table
CREATE TABLE installations (
    id TEXT PRIMARY KEY,
    system_id TEXT NOT NULL,
    template_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    log_file TEXT,
    error_reason TEXT,
    variables TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (system_id) REFERENCES systems(id),
    FOREIGN KEY (template_id) REFERENCES templates(id)
);

CREATE INDEX idx_installations_system_id ON installations(system_id);
CREATE INDEX idx_installations_template_id ON installations(template_id);
CREATE INDEX idx_installations_status ON installations(status);
```

### 8.5. Query Optimization

Optimize database queries for performance and efficiency.

#### 8.5.1. Indexing Strategy

Create appropriate indexes based on query patterns:

1. **Primary Key Indexes**: Automatically created (e.g., `id`)
2. **Foreign Key Indexes**: Create for all foreign keys (e.g., `system_id`)
3. **Filter Indexes**: Create for columns frequently used in WHERE clauses
   (e.g., `status`)
4. **Unique Indexes**: For columns with uniqueness constraints (e.g.,
   `hostname`, `mac_address`)

Example of adding indexes:

```sql
-- Add index for frequent filtering by status
CREATE INDEX idx_systems_status ON systems(status);

-- Add composite index for common query pattern
CREATE INDEX idx_installations_system_status ON installations(system_id, status);
```

#### 8.5.2. Query Optimization Techniques

1. **Select Only Needed Columns**:

```go
// Bad: Fetching all columns
db.Find(&systems)

// Good: Fetch only needed columns
db.Select("id, hostname, status").Find(&systems)
```

2. **Use Prepared Statements**:

```go
// Prepare statement once, use multiple times
stmt, err := db.Prepare("SELECT * FROM systems WHERE status = ?")
if err != nil {
    return err
}
defer stmt.Close()

// Execute for different statuses
rows, err := stmt.Query("ready")
```

3. **Pagination for Large Result Sets**:

```go
func (r *sqliteSystemRepository) List(filter SystemFilter) ([]*models.System, error) {
    // Always use pagination for lists
    query := r.db.Model(&models.System{})

    // Apply filters...

    // Default limit if not specified
    limit := filter.Limit
    if limit <= 0 {
        limit = 100 // Default page size
    }

    return query.Limit(limit).Offset(filter.Offset).Find(&systems)
}
```

4. **Batch Operations**:

```go
// Instead of multiple individual inserts
func (r *sqliteSystemRepository) BulkCreate(systems []*models.System) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        for _, system := range systems {
            system.GenerateID()
            if err := system.Validate(); err != nil {
                return err
            }

            if err := tx.Create(system).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

#### 8.5.3. Database-Specific Optimizations

**SQLite Optimizations**:

```go
// In internal/database/sqlite.go
func optimizeSQLiteConnection(db *sql.DB) error {
    // Enable WAL mode for better concurrency
    if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
        return err
    }

    // Set synchronous mode for better performance
    if _, err := db.Exec("PRAGMA synchronous=NORMAL"); err != nil {
        return err
    }

    // Set cache size (in pages)
    if _, err := db.Exec("PRAGMA cache_size=-16000"); err != nil {
        return err
    }

    // Enable foreign keys
    if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
        return err
    }

    return nil
}
```

**CockroachDB Optimizations**:

```go
// In internal/database/cockroach.go
func optimizeCockroachConnection(db *sql.DB) error {
    // Set appropriate connection pooling
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(30 * time.Minute)
    db.SetConnMaxIdleTime(10 * time.Minute)

    return nil
}
```

### 8.6. Testing Database Code

Effective testing strategies for database-related code.

#### 8.6.1. Repository Testing with SQLite

Use in-memory SQLite for fast repository tests:

```go
// In internal/database/system_repository_test.go
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
    // Use in-memory SQLite database
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect to test database: %v", err)
    }

    // Migrate schema
    if err := db.AutoMigrate(&models.System{}, &models.Template{}, &models.Installation{}); err != nil {
        t.Fatalf("failed to migrate test database: %v", err)
    }

    // Return cleanup function
    cleanup := func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }

    return db, cleanup
}

func TestSystemRepository_Create(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := NewSQLiteSystemRepository(db)

    // Test create operation
    system := &models.System{
        Hostname:   "test-host",
        MacAddress: "00:11:22:33:44:55",
        Status:     models.SystemStatusDiscovered,
    }

    err := repo.Create(system)
    assert.NoError(t, err)
    assert.NotEmpty(t, system.ID)

    // Verify system was created
    saved, err := repo.GetByID(system.ID)
    assert.NoError(t, err)
    assert.Equal(t, system.Hostname, saved.Hostname)
    assert.Equal(t, system.MacAddress, saved.MacAddress)
}
```

#### 8.6.2. Mock Repositories for Service Tests

Use mock repositories when testing services:

```go
// In internal/system/service_test.go
type MockSystemRepository struct {
    mock.Mock
}

func (m *MockSystemRepository) Create(system *models.System) error {
    args := m.Called(system)
    return args.Error(0)
}

func (m *MockSystemRepository) GetByID(id string) (*models.System, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.System), args.Error(1)
}

// Implement other required methods...

func TestSystemService_CreateSystem(t *testing.T) {
    mockRepo := &MockSystemRepository{}
    service := NewSystemService(mockRepo)

    // Setup mock expectations
    system := &models.System{
        Hostname:   "test-host",
        MacAddress: "00:11:22:33:44:55",
    }

    mockRepo.On("Create", mock.MatchedBy(func(s *models.System) bool {
        return s.Hostname == system.Hostname && s.MacAddress == system.MacAddress
    })).Return(nil)

    // Test the service method
    err := service.CreateSystem(system)
    assert.NoError(t, err)

    // Verify expectations
    mockRepo.AssertExpectations(t)
}
```

#### 8.6.3. Integration Tests with Test Database

For comprehensive testing, use a dedicated test database:

```go
// In test/integration/database_test.go
func setupIntegrationTestDB() (*gorm.DB, func()) {
    // Use environment variable to control test database
    dbPath := os.Getenv("TEST_DB_PATH")
    if dbPath == "" {
        dbPath = "file:test.db?mode=memory&cache=shared"
    }

    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        panic(fmt.Sprintf("failed to connect to test database: %v", err))
    }

    // Run migrations
    err = db.AutoMigrate(
        &models.System{},
        &models.Template{},
        &models.Installation{},
        &models.User{},
    )
    if err != nil {
        panic(fmt.Sprintf("failed to migrate test database: %v", err))
    }

    // Return cleanup function
    cleanup := func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()

        // If using a file-based database, remove it
        if !strings.Contains(dbPath, ":memory:") {
            os.Remove(strings.TrimPrefix(dbPath, "file:"))
        }
    }

    return db, cleanup
}

func TestDatabaseIntegration(t *testing.T) {
    // Skip in short mode
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    db, cleanup := setupIntegrationTestDB()
    defer cleanup()

    // Create repositories
    systemRepo := NewSQLiteSystemRepository(db)
    templateRepo := NewSQLiteTemplateRepository(db)
    installationRepo := NewSQLiteInstallationRepository(db)

    // Run integration tests...
    t.Run("SystemCRUD", func(t *testing.T) {
        // Test system CRUD operations
        system := &models.System{
            Hostname:   "integration-test",
            MacAddress: "00:11:22:33:44:55",
            Status:     models.SystemStatusDiscovered,
        }

        // Test create
        err := systemRepo.Create(system)
        assert.NoError(t, err)

        // Test get
        saved, err := systemRepo.GetByID(system.ID)
        assert.NoError(t, err)
        assert.Equal(t, system.Hostname, saved.Hostname)

        // Test update
        system.Status = models.SystemStatusReady
        err = systemRepo.Update(system)
        assert.NoError(t, err)

        updated, _ := systemRepo.GetByID(system.ID)
        assert.Equal(t, models.SystemStatusReady, updated.Status)

        // Test delete
        err = systemRepo.Delete(system.ID)
        assert.NoError(t, err)

        _, err = systemRepo.GetByID(system.ID)
        assert.Error(t, err)
    })
}
```
