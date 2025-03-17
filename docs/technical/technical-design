<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Ubuntu Autoinstall Webhook Technical Design Document](#ubuntu-autoinstall-webhook-technical-design-document)
  - [1. Introduction](#1-introduction)
    - [1.1. Purpose](#11-purpose)
    - [1.2. Scope](#12-scope)
    - [1.3. Definitions and Acronyms](#13-definitions-and-acronyms)
  - [2. System Architecture](#2-system-architecture)
    - [2.1. Architecture Overview](#21-architecture-overview)
    - [2.2. Design Principles](#22-design-principles)
    - [2.3. Component Interaction](#23-component-interaction)
  - [3. Technical Specifications](#3-technical-specifications)
    - [3.1. Programming Languages and Frameworks](#31-programming-languages-and-frameworks)
    - [3.2. Communication Protocols](#32-communication-protocols)
    - [3.3. Data Storage](#33-data-storage)
    - [3.4. Security Mechanisms](#34-security-mechanisms)
  - [4. Component Design](#4-component-design)
    - [4.1. File Editor Service](#41-file-editor-service)
    - [4.2. Database Service](#42-database-service)
    - [4.3. Configuration Service](#43-configuration-service)
    - [4.4. DNSMasq Watcher](#44-dnsmasq-watcher)
    - [4.5. Certificate Issuer](#45-certificate-issuer)
    - [4.6. Webserver](#46-webserver)
  - [5. Data Design](#5-data-design)
    - [5.1. Database Schema](#51-database-schema)
    - [5.2. File Structure](#52-file-structure)
    - [5.3. API Data Models](#53-api-data-models)
  - [6. Interface Design](#6-interface-design)
    - [6.1. API Endpoints](#61-api-endpoints)
    - [6.2. gRPC Service Definitions](#62-grpc-service-definitions)
    - [6.3. Web UI Structure](#63-web-ui-structure)
  - [7. Deployment Architecture](#7-deployment-architecture)
    - [7.1. Single-Node Deployment](#71-single-node-deployment)
    - [7.2. Kubernetes Deployment](#72-kubernetes-deployment)
    - [7.3. Hybrid Deployment](#73-hybrid-deployment)
  - [8. Performance Considerations](#8-performance-considerations)
    - [8.1. Scalability](#81-scalability)
    - [8.2. Availability](#82-availability)
    - [8.3. Resource Requirements](#83-resource-requirements)
  - [9. Security Considerations](#9-security-considerations)
    - [9.1. Authentication](#91-authentication)
    - [9.2. Authorization](#92-authorization)
    - [9.3. Data Protection](#93-data-protection)
  - [10. Development Considerations](#10-development-considerations)
    - [10.1. Development Environment](#101-development-environment)
    - [10.2. Build Process](#102-build-process)
    - [10.3. Testing Strategy](#103-testing-strategy)
  - [11. Operational Considerations](#11-operational-considerations)
    - [11.1. Monitoring](#111-monitoring)
    - [11.2. Logging](#112-logging)
    - [11.3. Backup and Recovery](#113-backup-and-recovery)
  - [12. References](#12-references)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Ubuntu Autoinstall Webhook Technical Design Document

## 1. Introduction

### 1.1. Purpose

This document provides a comprehensive technical design for the Ubuntu Autoinstall Webhook system, a solution designed to automate the installation of Ubuntu systems through PXE boot and cloud-init configurations. It serves as the primary reference for developers, architects, and technical staff implementing, maintaining, or extending the system.

### 1.2. Scope

This document covers:
- The system architecture and component design
- Technical specifications and protocols
- Data structures and storage models
- Interface definitions
- Deployment architectures
- Performance, security, and operational considerations

### 1.3. Definitions and Acronyms

- **PXE**: Preboot Execution Environment
- **iPXE**: Enhanced implementation of the PXE client
- **cloud-init**: Ubuntu's cloud instance initialization system
- **DHCP**: Dynamic Host Configuration Protocol
- **PKI**: Public Key Infrastructure
- **mTLS**: Mutual Transport Layer Security
- **RBAC**: Role-Based Access Control
- **CA**: Certificate Authority
- **CSR**: Certificate Signing Request
- **CRL**: Certificate Revocation List
- **gRPC**: Google Remote Procedure Call
- **REST**: Representational State Transfer
- **SPA**: Single-Page Application

## 2. System Architecture

### 2.1. Architecture Overview

The Ubuntu Autoinstall Webhook system follows a microservices architecture while being packaged as a single Go binary with multiple commands. This approach combines the benefits of microservices (separation of concerns, modular development) with the simplicity of deployment for a monolithic application.

The system consists of six core components:
1. **File Editor Service**: Manages all filesystem operations
2. **Database Service**: Provides persistent data storage
3. **Configuration Service**: Generates and validates system configurations
4. **DNSMasq Watcher**: Detects new systems on the network
5. **Certificate Issuer**: Manages the PKI infrastructure
6. **Webserver**: Hosts the web UI and API endpoints

### 2.2. Design Principles

The system architecture adheres to the following core principles:

1. **Separation of Concerns**: Each microservice has a distinct responsibility
2. **Interface-Based Design**: Components interact through well-defined interfaces
3. **Testability**: Components are designed for comprehensive testing
4. **Configurability**: All aspects of the system are configurable
5. **Security-First**: Security is built into the design, not added later
6. **Observability**: Comprehensive instrumentation for monitoring and debugging

### 2.3. Component Interaction

The microservices communicate primarily through gRPC, with the following key interactions:

**Installation Workflow**:
1. DNSMasq Watcher detects a new system via DHCP request
2. System information is stored in Database service
3. Configuration service generates cloud-init and iPXE files
4. Generated files are sent to File Editor for writing to disk
5. Client system boots using PXE and iPXE
6. Client retrieves cloud-init files from Webserver
7. Installation progress is reported back to Webserver

**Administrative Workflow**:
1. Admin accesses the web UI through Webserver
2. Authentication is verified against Database
3. Admin creates/modifies system configurations through web UI
4. Configuration service validates and stores changes in Database
5. File Editor updates filesystem with new configurations

## 3. Technical Specifications

### 3.1. Programming Languages and Frameworks

- **Backend**: Go programming language
- **Frontend**: Angular framework (TypeScript, HTML, CSS)
- **Configuration Templates**: YAML, Jinja2-style templating
- **Build System**: Make, Docker

### 3.2. Communication Protocols

- **gRPC**: For inter-service communication
- **HTTP/HTTPS**: For web UI and client interactions
- **WebSockets**: For real-time updates in the web UI
- **TLS/mTLS**: For encrypted communications

### 3.3. Data Storage

- **SQLite3**: Default database for single-node deployments
- **CockroachDB**: Optional database for distributed deployments
- **Filesystem**: For storing cloud-init and iPXE files
- **In-memory cache**: For frequently accessed configuration data

### 3.4. Security Mechanisms

- **PKI**: Self-managed certificate authority
- **mTLS**: For service-to-service authentication
- **RBAC**: For authorization in the web UI
- **TLS**: For encrypting all network traffic
- **Password hashing**: For user authentication

## 4. Component Design

### 4.1. File Editor Service

The File Editor service manages all filesystem operations within the system:

- **Leader Election**: Uses database locks to ensure only one instance writes to the filesystem
- **Directory Management**: Creates and maintains directory structures for cloud-init files
- **File Operations**: Provides CRUD operations for files with atomic write guarantees
- **Validation**: Ensures file content meets expected formats before writing

### 4.2. Database Service

The Database service provides persistent storage and abstracts database operations:

- **Backend Support**: SQLite3 (default) and CockroachDB
- **Schema Management**: Handles database migrations and schema evolution
- **Data Access Layer**: Provides a unified API for database operations
- **Transaction Management**: Ensures data consistency across operations

### 4.3. Configuration Service

The Configuration service manages all system configuration aspects:

- **Template System**: Manages templates for cloud-init and iPXE files
- **Configuration Generation**: Creates system-specific configurations from templates
- **Validation**: Ensures generated configurations are valid
- **Caching**: Implements efficient caching for frequently accessed configurations

### 4.4. DNSMasq Watcher

The DNSMasq Watcher monitors network activity to detect new systems:

- **Log Monitoring**: Watches dnsmasq logs for DHCP events
- **System Detection**: Identifies new systems based on MAC address
- **Hostname Generation**: Creates predictable hostnames for new systems
- **Notification**: Alerts other services about newly detected systems

### 4.5. Certificate Issuer

The Certificate Issuer manages the PKI infrastructure:

- **CA Management**: Maintains root and intermediate certificate authorities
- **Certificate Lifecycle**: Handles issuance, renewal, and revocation
- **Authentication**: Provides mechanisms for initial client authentication
- **Validation**: Ensures issued certificates meet security requirements

### 4.6. Webserver

The Webserver provides the external interface for the system:

- **Web UI**: Hosts the Angular-based administrative interface
- **RESTful API**: Provides endpoints for external integration
- **gRPC API**: Offers high-performance API for installation clients
- **Authentication**: Manages user sessions and access control
- **Static Files**: Serves installation resources like kernel and initrd

## 5. Data Design

### 5.1. Database Schema

The database schema consists of the following core tables:

- **Systems**: Stores information about systems being installed or managed
  - UUID (primary key)
  - MAC address
  - Hostname
  - IP address
  - Status
  - Created/Updated timestamps

- **Configurations**: Stores installation configurations and templates
  - UUID (primary key)
  - Name
  - Description
  - Configuration type
  - Template content
  - Created/Updated timestamps

- **Installation_Status**: Tracks installation progress for systems
  - UUID (primary key)
  - System UUID (foreign key)
  - Status code
  - Status message
  - Timestamp

- **Users**: Manages user authentication information
  - UUID (primary key)
  - Username
  - Password hash
  - Email
  - Created/Updated timestamps

- **Roles**: Defines user roles for authorization
  - UUID (primary key)
  - Name
  - Description
  - Permissions (JSON)

- **User_Roles**: Maps users to roles
  - User UUID (foreign key)
  - Role UUID (foreign key)
  - Primary key: (User UUID, Role UUID)

### 5.2. File Structure

The file system structure is organized as follows:
/var/www/html/
├── ipxe/
│ └── boot/
│ ├── mac-{MAC_ADDRESS}.ipxe
│ └── ...
├── cloud-init/
│ ├── {MAC_ADDRESS}/
│ │ ├── meta-data
│ │ ├── network-config
│ │ ├── user-data
│ │ └── variables.sh
│ ├── {MAC_ADDRESS}_install/
│ │ └── ... (installation-specific files)
│ ├── {HOSTNAME}/ -> {MAC_ADDRESS}/
│ └── {HOSTNAME}_install/ -> {MAC_ADDRESS}_install/
└── assets/
├── kernel/
├── initrd/
└── web-ui/


### 5.3. API Data Models

The system uses Protocol Buffers for defining gRPC service interfaces and data models. Key data models include:

- **SystemInfo**: Contains system identification and status information
- **Configuration**: Represents a system configuration or template
- **InstallationStatus**: Tracks the progress of an installation
- **User**: Represents a user account
- **Role**: Defines a set of permissions for authorization
- **FileContent**: Contains file data and metadata

## 6. Interface Design

### 6.1. API Endpoints

**RESTful API Endpoints**:

- `/api/v1/systems`: CRUD operations for systems
- `/api/v1/configurations`: CRUD operations for configurations
- `/api/v1/templates`: CRUD operations for templates
- `/api/v1/users`: User management
- `/api/v1/roles`: Role management
- `/api/v1/installations`: Installation status and control
- `/api/v1/certificates`: Certificate management
- `/api/v1/health`: System health check

**Installation Client Endpoints**:

- `/ipxe/boot/mac-{MAC_ADDRESS}.ipxe`: iPXE boot configuration
- `/cloud-init/{MAC_ADDRESS}/meta-data`: Cloud-init metadata
- `/cloud-init/{MAC_ADDRESS}/user-data`: Cloud-init user data
- `/cloud-init/{MAC_ADDRESS}/network-config`: Cloud-init network configuration
- `/cloud-init/{MAC_ADDRESS}_install/variables.sh`: Installation variables

### 6.2. gRPC Service Definitions

The system defines the following gRPC services:

- **FileEditorService**: File operations
  - `CreateFile(FileRequest) returns (FileResponse)`
  - `ReadFile(FileRequest) returns (FileResponse)`
  - `UpdateFile(FileRequest) returns (FileResponse)`
  - `DeleteFile(FileRequest) returns (StatusResponse)`
  - `CreateDirectory(DirectoryRequest) returns (StatusResponse)`
  - `CreateSymlink(SymlinkRequest) returns (StatusResponse)`

- **DatabaseService**: Data persistence
  - `Create(Entity) returns (Entity)`
  - `Read(Query) returns (Entity)`
  - `Update(Entity) returns (Entity)`
  - `Delete(Query) returns (StatusResponse)`
  - `List(Query) returns (EntityList)`
  - `ExecuteTransaction(Operations) returns (Results)`

- **ConfigurationService**: Configuration management
  - `GenerateConfiguration(SystemInfo) returns (ConfigurationFiles)`
  - `ValidateConfiguration(Configuration) returns (ValidationResult)`
  - `CreateTemplate(Template) returns (Template)`
  - `UpdateTemplate(Template) returns (Template)`
  - `DeleteTemplate(TemplateId) returns (StatusResponse)`
  - `ListTemplates(Query) returns (TemplateList)`

- **DNSMasqWatcherService**: System detection
  - `RegisterSystem(SystemInfo) returns (SystemInfo)`
  - `GetSystem(SystemQuery) returns (SystemInfo)`
  - `ListSystems(SystemQuery) returns (SystemList)`
  - `ConfigureWatcher(WatcherConfig) returns (StatusResponse)`

- **CertificateIssuerService**: Certificate management
  - `SubmitCSR(CSRRequest) returns (Certificate)`
  - `GetCACertificate(CARequest) returns (Certificate)`
  - `RevokeCertificate(RevocationRequest) returns (StatusResponse)`
  - `CheckCertificateStatus(CertificateQuery) returns (CertificateStatus)`

- **WebserverService**: External interface
  - `Authenticate(Credentials) returns (AuthToken)`
  - `GetSystemStatus(SystemQuery) returns (SystemStatus)`
  - `UpdateSystemStatus(SystemStatus) returns (StatusResponse)`
  - `GetConfiguration(ConfigQuery) returns (Configuration)`
  - `StreamInstallationProgress(SystemId) returns (stream ProgressUpdate)`

### 6.3. Web UI Structure

The web UI follows an Angular SPA architecture with the following main sections:

- **Dashboard**: Overview of system status
- **Systems**: Management of systems being installed
- **Configurations**: Template and configuration management
- **Users**: User and role management
- **Certificates**: Certificate management
- **Settings**: System-wide settings
- **Logs**: System logs and audit trails

## 7. Deployment Architecture

### 7.1. Single-Node Deployment

For small-scale deployments, the system runs as a single Go binary on one server:

- All services running on a single server
- SQLite3 as the database backend
- Suitable for environments with up to a few hundred nodes
- Simple backup and restore procedures
- Minimal resource requirements

### 7.2. Kubernetes Deployment

For large-scale or high-availability deployments, the system can be deployed on Kubernetes:

- Services deployed as separate containers
- CockroachDB for distributed database
- Horizontal scaling for stateless components
- Leader election for stateful components
- Leverages Kubernetes features for availability and scaling
- Suitable for environments with thousands of nodes

### 7.3. Hybrid Deployment

For geographically distributed environments:

- Core services deployed centrally
- DNSMasq Watcher deployed at network edge points
- File Editor with synchronized storage
- Federation capabilities for multi-region deployments
- Optimized for WAN connectivity between components

## 8. Performance Considerations

### 8.1. Scalability

- **Horizontal Scaling**: Stateless services can be scaled horizontally
- **Vertical Scaling**: Database and file storage benefit from vertical scaling
- **Caching**: Configuration service implements caching to reduce database load
- **Connection Pooling**: Database service manages connection pools efficiently

### 8.2. Availability

- **Service Redundancy**: Multiple instances of each service can be deployed
- **Database Replication**: CockroachDB provides high availability for data
- **Leader Election**: Ensures only one instance handles critical operations
- **Graceful Degradation**: System continues to function with reduced capability if some services are unavailable

### 8.3. Resource Requirements

Minimum system requirements for a single-node deployment:

- **CPU**: 2 cores
- **RAM**: 4 GB
- **Storage**: 20 GB
- **Network**: 1 Gbps Ethernet

Recommended system requirements for high-performance deployment:

- **CPU**: 8+ cores
- **RAM**: 16+ GB
- **Storage**: 100+ GB SSD
- **Network**: 10 Gbps Ethernet

## 9. Security Considerations

### 9.1. Authentication

- **Service-to-Service**: mTLS with certificate-based authentication
- **User-to-System**: Username/password or SSO integration
- **Client-to-System**: Initial authentication with pre-shared secrets or MAC verification

### 9.2. Authorization

- **RBAC**: Granular permissions for web UI access
- **Service Authorization**: Certificate-based verification of service identity
- **API Authorization**: Token-based authorization for API access

### 9.3. Data Protection

- **Encryption at Rest**: Sensitive data encrypted in the database
- **Encryption in Transit**: All network communication encrypted with TLS
- **Sensitive Data Handling**: Passwords and secrets stored securely
- **Certificate Security**: Private keys protected with appropriate permissions

## 10. Development Considerations

### 10.1. Development Environment

- **Local Development**: Docker Compose for local development environments
- **IDE Integration**: Go tools for code quality and consistency
- **Mock Services**: Mock implementations for testing
- **Database Migrations**: Tools for managing schema evolution

### 10.2. Build Process

- **Continuous Integration**: GitHub Actions or similar CI system
- **Build Toolchain**: Make for build automation
- **Containerization**: Docker for containerized deployments
- **Versioning**: Semantic versioning for releases

### 10.3. Testing Strategy

- **Unit Testing**: Test individual components in isolation
- **Integration Testing**: Test interactions between components
- **End-to-End Testing**: Test complete workflows
- **Performance Testing**: Test under load conditions
- **Security Testing**: Vulnerability scanning and penetration testing

## 11. Operational Considerations

### 11.1. Monitoring

- **Health Endpoints**: Each service exposes a health endpoint
- **Metrics**: Prometheus integration for metrics collection
- **Dashboards**: Grafana or similar for visualization
- **Alerts**: Configurable alerting for system events

### 11.2. Logging

- **Structured Logging**: JSON-formatted logs for machine processing
- **Log Levels**: Configurable log levels for each component
- **Log Aggregation**: Integration with centralized logging systems
- **Audit Logging**: Captures security-relevant events

### 11.3. Backup and Recovery

- **Database Backup**: Regular backups of database content
- **Configuration Backup**: Backup of all system configurations
- **Disaster Recovery**: Procedures for system restoration
- **High Availability**: Strategies for minimizing downtime

## 12. References

- [Architecture Overview](/docs/architecture/overview.md)
- [Component Documentation](/docs/architecture/components/)
- [API Documentation](/docs/api/)
- [Database Schema](/docs/database-schema.md)
- [Deployment Guide](/docs/deployment-guide.md)
