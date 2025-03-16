<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [DNSMasq Watcher Microservice Architecture](#dnsmasq-watcher-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Log Monitoring](#log-monitoring)
  - [System Detection](#system-detection)
    - [Identification Process](#identification-process)
    - [Deduplication Strategy](#deduplication-strategy)
  - [Interface](#interface)
  - [Hostname Generation](#hostname-generation)
  - [Interactions with Other Services](#interactions-with-other-services)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# DNSMasq Watcher Microservice Architecture

## Table of Contents
- [DNSMasq Watcher Microservice Architecture](#dnsmasq-watcher-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Log Monitoring](#log-monitoring)
  - [System Detection](#system-detection)
    - [Identification Process](#identification-process)
    - [Deduplication Strategy](#deduplication-strategy)
  - [Interface](#interface)
  - [Hostname Generation](#hostname-generation)
  - [Interactions with Other Services](#interactions-with-other-services)

## Overview

The DNSMasq Watcher microservice monitors dnsmasq logs or API to detect new systems on the network. It extracts MAC addresses, IP addresses, and hostnames from DHCP requests and registers these systems for potential automatic installation.

## Core Responsibilities

- Monitoring dnsmasq logs or API for DHCP events
- Detecting new systems based on MAC addresses
- Extracting system information (IP address, hostname)
- Generating default hostnames when not provided
- Registering discovered systems with the Configuration service

## Log Monitoring

- Supports multiple methods for monitoring dnsmasq:
  - Log file tailing
  - Syslog reception
  - Journal monitoring
  - Direct API calls if available
- Parses log entries to extract DHCP lease information
- Handles log rotation and reconnection scenarios
- Filters relevant DHCP events from other log entries

## System Detection

### Identification Process
- Detects DHCP DISCOVER, REQUEST, and ACK messages
- Extracts MAC address as primary system identifier
- Captures assigned IP address from DHCP ACK
- Notes client-provided hostname if available
- Cross-references with existing system records

### Deduplication Strategy
- Uses MAC address as unique identifier
- Updates existing system records rather than creating duplicates
- Tracks IP address changes for the same MAC
- Resolves hostname conflicts according to policy

## Interface

The DNSMasq Watcher service exposes a gRPC interface that provides:

- Manual system registration
- Query capabilities for discovered systems
- Configuration of monitoring parameters
- Status information about the monitoring service

## Hostname Generation

When a client doesn't provide a hostname:

- Generates a predictable hostname based on MAC address
- Ensures uniqueness within the system database
- Follows configurable naming patterns (e.g., prefix + MAC)
- Handles collision detection and resolution

## Interactions with Other Services

- Notifies the Configuration service of newly discovered systems
- Uses the Database service to check for existing system records
- Updates system records with new information (IP, hostname)
- Retrieves naming policies from the Configuration service
