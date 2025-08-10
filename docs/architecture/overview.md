<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Ubuntu Autoinstall Webhook System Architecture](#ubuntu-autoinstall-webhook-system-architecture)
  - [Table of Contents](#table-of-contents)
  - [System Overview](#system-overview)
    - [Core Components](#core-components)
  - [Architecture Principles](#architecture-principles)
  - [Component Interactions](#component-interactions)
    - [Installation Workflow](#installation-workflow)
    - [Administrative Workflow](#administrative-workflow)
  - [Deployment Models](#deployment-models)
    - [Single-Node Deployment](#single-node-deployment)
    - [Kubernetes Deployment](#kubernetes-deployment)
    - [Hybrid Deployment](#hybrid-deployment)
  - [Security Architecture](#security-architecture)
    - [Authentication Methods](#authentication-methods)
    - [Authorization](#authorization)
    - [Data Protection](#data-protection)
  - [Scalability and Availability](#scalability-and-availability)
    - [Stateless Services](#stateless-services)
    - [Stateful Services](#stateful-services)
  - [Data Flow](#data-flow)
    - [Installation Data Flow](#installation-data-flow)
    - [Configuration Data Flow](#configuration-data-flow)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Ubuntu Autoinstall Webhook System Architecture

## Table of Contents

- [System Overview](#system-overview)
- [Architecture Principles](#architecture-principles)
- [Component Interactions](#component-interactions)
- [Deployment Models](#deployment-models)
- [Security Architecture](#security-architecture)
- [Scalability and Availability](#scalability-and-availability)
- [Data Flow](#data-flow)

## System Overview

The Ubuntu Autoinstall Webhook system is a microservices-based application
designed to automate the installation of Ubuntu systems through PXE boot and
cloud-init. It provides both a web interface for administrators and a set of API
endpoints for automation.

Despite being implemented as microservices, the system is packaged as a single
Go binary with multiple commands, simplifying deployment while maintaining the
benefits of a microservices architecture.

### Core Components

1. **File Editor Service**: Manages the filesystem operations
2. **Database Service**: Handles persistent data storage
3. **Configuration Service**: Generates and validates system configurations
4. **DNSMasq Watcher**: Detects new systems on the network
5. **Certificate Issuer**: Manages the PKI infrastructure
6. **Webserver**: Hosts the web UI and API endpoints

## Architecture Principles

The system architecture follows these core principles:

1. **Separation of Concerns**: Each microservice has a distinct responsibility
2. **Interface-Based Design**: Components interact through well-defined
   interfaces
3. **Testability**: All components are designed for comprehensive testing
4. **Configurability**: All aspects of the system are configurable
5. **Security-First**: Security is built into the design, not added later
6. **Observability**: Comprehensive instrumentation for monitoring and debugging

## Component Interactions

The microservices communicate primarily through gRPC, with the following key
interactions:

### Installation Workflow

1. **System Discovery**:
   - DNSMasq Watcher detects a new system via DHCP request
   - System information is stored in Database service
   - Configuration service is notified of the new system

2. **Configuration Generation**:
   - Configuration service generates cloud-init and iPXE files
   - Generated files are sent to File Editor for writing to disk
   - System record in Database is updated with configuration status

3. **Installation Process**:
   - Client system boots using PXE and iPXE
   - Client retrieves cloud-init files from Webserver
   - Client establishes secure connection (mTLS) via Certificate Issuer
   - Installation progress is reported back to Webserver
   - Status updates are stored in Database

### Administrative Workflow

1. **User Authentication**:
   - Admin accesses the web UI through Webserver
   - Authentication is verified against Database
   - Session is established with appropriate RBAC permissions

2. **Configuration Management**:
   - Admin creates/modifies system configurations through web UI
   - Webserver forwards changes to Configuration service
   - Configuration service validates and stores changes in Database
   - File Editor updates filesystem with new configurations

## Deployment Models

The system supports multiple deployment models:

### Single-Node Deployment

- All services running on a single server
- Suitable for small-scale deployments
- Uses SQLite3 as the database backend

### Kubernetes Deployment

- Services deployed as separate containers
- Supports horizontal scaling for most components
- Uses CockroachDB for distributed database
- Leverages Kubernetes features for availability and scaling

### Hybrid Deployment

- Core services deployed centrally
- DNSMasq Watcher deployed at network edge points
- Installation clients across the network

## Security Architecture

### Authentication Methods

- Mutual TLS for service-to-service communication
- Pre-shared secrets for initial client authentication
- Username/password or SSO for web UI
- MAC address and IP-based verification for clients

### Authorization

- RBAC for web UI access control
- Service-to-service authorization via mTLS
- Granular permissions model for administrative functions

### Data Protection

- Encryption of sensitive data at rest
- TLS for all network communication
- Certificate lifecycle management

## Scalability and Availability

### Stateless Services

Most services are stateless and can be horizontally scaled:

- Configuration Service
- Certificate Issuer
- Webserver
- DNSMasq Watcher

### Stateful Services

Services with state have special considerations:

- Database Service: Uses CockroachDB for scalability when needed
- File Editor: Uses leader election to prevent conflicts

## Data Flow

### Installation Data Flow

1. Client system boots and obtains IP via DHCP
2. DNSMasq Watcher detects the system and registers it
3. Configuration is generated for the system
4. Client retrieves iPXE script and boots installer
5. Client downloads cloud-init files for configuration
6. Client reports installation progress
7. System is recorded as successfully installed

### Configuration Data Flow

1. Administrator creates or modifies configuration
2. Configuration is validated and stored in database
3. Configuration is applied to relevant systems
4. File configurations are updated on disk
5. Systems using the configuration are notified of changes
