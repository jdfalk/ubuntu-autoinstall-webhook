<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Configuration Microservice Architecture](#configuration-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Configuration Management](#configuration-management)
    - [Global Configuration](#global-configuration)
    - [Per-System Configuration](#per-system-configuration)
  - [Template System](#template-system)
  - [Configuration Generation](#configuration-generation)
    - [Cloud-Init Files](#cloud-init-files)
    - [iPXE Boot Scripts](#ipxe-boot-scripts)
    - [Installation Variables](#installation-variables)
  - [Interface](#interface)
  - [Caching Strategy](#caching-strategy)
  - [Interactions with Other Services](#interactions-with-other-services)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Configuration Microservice Architecture

## Table of Contents

- [Configuration Microservice Architecture](#configuration-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Configuration Management](#configuration-management)
    - [Global Configuration](#global-configuration)
    - [Per-System Configuration](#per-system-configuration)
  - [Template System](#template-system)
  - [Configuration Generation](#configuration-generation)
    - [Cloud-Init Files](#cloud-init-files)
    - [iPXE Boot Scripts](#ipxe-boot-scripts)
    - [Installation Variables](#installation-variables)
  - [Interface](#interface)
  - [Caching Strategy](#caching-strategy)
  - [Interactions with Other Services](#interactions-with-other-services)

## Overview

The Configuration microservice manages all system configuration aspects within
the ubuntu-autoinstall-webhook system. It generates, validates, and distributes
configuration artifacts needed for Ubuntu autoinstallation.

## Core Responsibilities

- Creating and validating system-specific configuration files
- Managing configuration templates
- Generating cloud-init artifacts (user-data, meta-data, network-config)
- Creating iPXE boot scripts for systems
- Managing global system settings

## Configuration Management

### Global Configuration

- System-wide settings stored in the database
- Includes network settings, installation defaults, and server configurations
- Manages feature flags and optional component settings
- Handles environment-specific configurations

### Per-System Configuration

- Individual system configurations derived from templates and overrides
- Stored in the database with unique identifiers
- Contains both generated values and admin-specified settings
- Tracks configuration versions and changes

## Template System

- Manages templates for cloud-init files (user-data, meta-data, network-config)
- Supports variables and conditional logic within templates
- Allows for inheritance and extension of base templates
- Validates template syntax and structure

## Configuration Generation

For each system, the service generates:

### Cloud-Init Files

- `user-data`: Contains system provisioning instructions
- `meta-data`: System identity information
- `network-config`: Network configuration settings

### iPXE Boot Scripts

- Boot configuration specific to each system's MAC address
- Controls boot sequence and installation process
- Contains URLs for kernel, initrd, and cloud-init locations

### Installation Variables

- Generates `variables.sh` file with system-specific settings
- Provides environment variables used during the installation process

## Interface

The Configuration service exposes a gRPC interface that provides:

- Template management (CRUD operations)
- Configuration generation for specific systems
- Configuration validation
- Global settings management

## Caching Strategy

- Caches frequently accessed templates and configurations
- Implements a cache invalidation strategy on configuration changes
- Uses memory-efficient representation for cached items
- Synchronizes cache across multiple service instances

## Interactions with Other Services

- Uses Database service to store and retrieve configuration data
- Sends generated files to the File Editor service for persistence
- Receives system discovery information from the DNSMasq Watcher service
- Provides configuration information to the Webserver service for UI display
