<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [File Editor Microservice Architecture](#file-editor-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Leader Election](#leader-election)
  - [Directory and File Management](#directory-and-file-management)
    - [iPXE Boot Configuration Files](#ipxe-boot-configuration-files)
    - [Cloud-Init Directory Structure](#cloud-init-directory-structure)
    - [Core Configuration Files](#core-configuration-files)
  - [File Validation](#file-validation)
  - [Interface](#interface)
  - [Error Handling](#error-handling)
  - [Interactions with Other Services](#interactions-with-other-services)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# File Editor Microservice Architecture

## Table of Contents

- [File Editor Microservice Architecture](#file-editor-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Leader Election](#leader-election)
  - [Directory and File Management](#directory-and-file-management)
    - [iPXE Boot Configuration Files](#ipxe-boot-configuration-files)
    - [Cloud-Init Directory Structure](#cloud-init-directory-structure)
    - [Core Configuration Files](#core-configuration-files)
  - [File Validation](#file-validation)
  - [Interface](#interface)
  - [Error Handling](#error-handling)
  - [Interactions with Other Services](#interactions-with-other-services)

## Overview

The File Editor microservice is responsible for all filesystem operations within
the ubuntu-autoinstall-webhook system. It provides a safe interface for reading,
writing, and validating files that are critical to the PXE boot process and
Ubuntu autoinstallation.

## Core Responsibilities

- Managing iPXE boot configuration files
- Creating and maintaining cloud-init directory structures
- Managing symbolic links between MAC address and hostname-based directories
- Validating configuration files before writing them to disk
- Ensuring file system operations are atomic and consistent

## Leader Election

To prevent file conflicts and ensure data integrity, only one instance of the
File Editor can be active at any time:

- Uses a database lock or Kubernetes lease for leader election
- Other instances remain in standby mode until leader fails
- When using SQLite3, multi-instance deployment is not supported
- Handles graceful failover if leader becomes unavailable

## Directory and File Management

### iPXE Boot Configuration Files

- Writes and manages files in `/var/www/html/ipxe/boot/`
- Uses consistent naming convention:
  `/var/www/html/ipxe/boot/mac-{MAC_ADDRESS}.ipxe`
- Updates files when installation status changes

### Cloud-Init Directory Structure

- Creates per-system directories using MAC address:
  `/var/www/html/cloud-init/{MAC_ADDRESS}/`
- Creates installation directories:
  `/var/www/html/cloud-init/{MAC_ADDRESS}_install/`
- Maintains symbolic links from hostname to MAC address:
  - `/var/www/html/cloud-init/{HOSTNAME}/` →
    `/var/www/html/cloud-init/{MAC_ADDRESS}/`
  - `/var/www/html/cloud-init/{HOSTNAME}_install/` →
    `/var/www/html/cloud-init/{MAC_ADDRESS}_install/`
- Prevents hostname conflicts through validation

### Core Configuration Files

For each system, manages the following files:

- `meta-data`: System metadata (instance-id, local-hostname)
- `network-config`: Network configuration (typically minimal)
- `user-data`: Cloud-init instructions for system configuration
- `variables.sh`: System-specific variables for the installation process

## File Validation

- Validates iPXE boot configurations for syntax and correctness
- Uses cloud-init libraries to validate cloud-init files
- Ensures variables.sh contains all required fields
- Performs permissions checking to ensure files are readable by web server

## Interface

The File Editor exposes a gRPC interface that allows other services to:

- Request file creation or modification
- Request file reads
- Check if a file exists
- Create directory structures
- Manage symbolic links

## Error Handling

- Reports detailed error information when file operations fail
- Implements retry mechanisms for transient failures
- Logs all file operations for audit and debugging
- Provides status information about the filesystem state

## Interactions with Other Services

- Receives file creation/modification requests from Configuration service
- Notifies Webserver service when files are updated
- Uses Database service to track file metadata and state
- Reports metrics and traces to observability systems
