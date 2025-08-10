<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Webserver Microservice Architecture](#webserver-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Web Frontend](#web-frontend)
    - [Frontend Features](#frontend-features)
  - [API Services](#api-services)
    - [RESTful API](#restful-api)
    - [gRPC API](#grpc-api)
  - [Authentication and Authorization](#authentication-and-authorization)
    - [Authentication Methods](#authentication-methods)
    - [Role-Based Access Control (RBAC)](#role-based-access-control-rbac)
  - [Installation Client Support](#installation-client-support)
  - [Interface](#interface)
  - [Interactions with Other Services](#interactions-with-other-services)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Webserver Microservice Architecture

## Table of Contents

- [Webserver Microservice Architecture](#webserver-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Web Frontend](#web-frontend)
    - [Frontend Features](#frontend-features)
  - [API Services](#api-services)
    - [RESTful API](#restful-api)
    - [gRPC API](#grpc-api)
  - [Authentication and Authorization](#authentication-and-authorization)
    - [Authentication Methods](#authentication-methods)
    - [Role-Based Access Control (RBAC)](#role-based-access-control-rbac)
  - [Installation Client Support](#installation-client-support)
  - [Interface](#interface)
  - [Interactions with Other Services](#interactions-with-other-services)

## Overview

The Webserver microservice serves as the primary external interface for the
ubuntu-autoinstall-webhook system. It hosts both the web-based administration
interface and the API endpoints used by installation clients.

## Core Responsibilities

- Hosting the Angular-based web frontend
- Serving static installation resources (kernel, initrd, etc.)
- Providing RESTful and gRPC APIs
- Managing user authentication and authorization
- Handling installation client connections
- Serving cloud-init files for autoinstallation

## Web Frontend

- Serves an Angular single-page application (SPA)
- Implements responsive design for various devices
- Provides an intuitive administrative interface
- Supports theme customization

### Frontend Features

- Dashboard with system installation status
- Configuration management interface
- Template editor with validation
- User and role management
- System logs and monitoring views

## API Services

### RESTful API

- API endpoints for web frontend
- Documentation via OpenAPI/Swagger
- Versioned API for backwards compatibility
- Rate limiting and request validation

### gRPC API

- High-performance API for installation clients
- Streaming support for real-time updates
- Protocol buffer-based type safety
- Service discovery and health checking

## Authentication and Authorization

### Authentication Methods

- Username/password authentication
- API token-based authentication
- Certificate-based authentication for mutual TLS

### Role-Based Access Control (RBAC)

- Hierarchical role system
- Granular permissions model
- Role assignment and management
- Audit logging for access events

## Installation Client Support

- Handles registration of new installation clients
- Serves appropriate cloud-init configurations
- Receives installation progress updates
- Manages client certificate issuance
- Tracks installation status and metrics

## Interface

The Webserver service exposes:

- HTTP/HTTPS endpoints for the web frontend
- RESTful API for programmatic access
- gRPC API for installation clients
- WebSocket endpoints for real-time updates

## Interactions with Other Services

- Retrieves configurations from the Configuration service
- Accesses system information through the Database service
- Gets certificates from the Certificate Issuer service
- Receives file information from the File Editor service
- Obtains new system notifications from the DNSMasq Watcher service
