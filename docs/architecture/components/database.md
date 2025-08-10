<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Database Microservice Architecture](#database-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Database Backend Support](#database-backend-support)
    - [SQLite3 (Default)](#sqlite3-default)
    - [CockroachDB](#cockroachdb)
  - [Schema Management](#schema-management)
    - [Systems](#systems)
    - [Configurations](#configurations)
    - [Installation Status](#installation-status)
    - [Users and Roles](#users-and-roles)
  - [Data Access Patterns](#data-access-patterns)
  - [Interface](#interface)
  - [Scalability and Redundancy](#scalability-and-redundancy)
  - [Interactions with Other Services](#interactions-with-other-services)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Database Microservice Architecture

## Table of Contents

- [Database Microservice Architecture](#database-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Database Backend Support](#database-backend-support)
    - [SQLite3 (Default)](#sqlite3-default)
    - [CockroachDB](#cockroachdb)
  - [Schema Management](#schema-management)
    - [Systems](#systems)
    - [Configurations](#configurations)
    - [Installation Status](#installation-status)
    - [Users and Roles](#users-and-roles)
  - [Data Access Patterns](#data-access-patterns)
  - [Interface](#interface)
  - [Scalability and Redundancy](#scalability-and-redundancy)
  - [Interactions with Other Services](#interactions-with-other-services)

## Overview

The Database microservice is the central persistence layer for the
ubuntu-autoinstall-webhook system. It abstracts database operations and provides
a consistent API for data storage and retrieval regardless of the underlying
database technology.

## Core Responsibilities

- Providing a unified data access layer via gRPC
- Managing database connections and transaction boundaries
- Handling schema creation and migrations
- Implementing data validation before storage
- Supporting data queries with filtering and pagination

## Database Backend Support

The Database microservice supports two database backends:

### SQLite3 (Default)

- File-based storage suitable for single-instance deployments
- Simple setup with minimal configuration
- Lower concurrency capabilities
- Stores UUID primary keys as strings

### CockroachDB

- Distributed SQL database for high availability
- Supports multiple service instances and horizontal scaling
- Higher performance for large-scale deployments
- Native UUID support

## Schema Management

The service automatically manages the following core tables:

### Systems

- Stores information about systems being installed or managed
- Primary key: UUID
- Tracks MAC address, hostname, IP address, status

### Configurations

- Stores installation configurations and templates
- Primary key: UUID
- Relates to Systems through foreign keys

### Installation Status

- Tracks installation progress for systems
- Primary key: UUID
- Stores timestamps for installation phases

### Users and Roles

- Manages authentication and authorization information
- Primary key: UUID
- Supports role hierarchies for the web frontend

## Data Access Patterns

- Implements transactional boundaries for multi-operation consistency
- Provides caching for frequently accessed data
- Uses prepared statements to prevent SQL injection
- Implements appropriate indexes for query optimization

## Interface

The Database service exposes a gRPC interface that provides:

- CRUD operations for all schema entities
- Query capabilities with filtering and pagination
- Transaction management
- Schema information and validation

## Scalability and Redundancy

- Single-instance operation when using SQLite3
- Multi-instance support with CockroachDB
- Read replicas when appropriate for the backend
- Connection pooling for efficient resource utilization

## Interactions with Other Services

- Serves as persistence layer for all other microservices
- Provides system information to the Configuration service
- Stores file metadata for the File Editor service
- Maintains user and role information for the Webserver service
- Stores certificate information for the Cert-Issuer service
