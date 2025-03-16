<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Certificate Issuer Microservice Architecture](#certificate-issuer-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Certificate Authority Management](#certificate-authority-management)
    - [Root CA](#root-ca)
    - [Intermediate CA](#intermediate-ca)
    - [External CA Integration](#external-ca-integration)
  - [Certificate Lifecycle](#certificate-lifecycle)
    - [Certificate Issuance](#certificate-issuance)
    - [Certificate Renewal](#certificate-renewal)
    - [Certificate Revocation](#certificate-revocation)
  - [Authentication Methods](#authentication-methods)
  - [Interface](#interface)
  - [Security Considerations](#security-considerations)
  - [Interactions with Other Services](#interactions-with-other-services)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Certificate Issuer Microservice Architecture

## Table of Contents
- [Certificate Issuer Microservice Architecture](#certificate-issuer-microservice-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Core Responsibilities](#core-responsibilities)
  - [Certificate Authority Management](#certificate-authority-management)
    - [Root CA](#root-ca)
    - [Intermediate CA](#intermediate-ca)
    - [External CA Integration](#external-ca-integration)
  - [Certificate Lifecycle](#certificate-lifecycle)
    - [Certificate Issuance](#certificate-issuance)
    - [Certificate Renewal](#certificate-renewal)
    - [Certificate Revocation](#certificate-revocation)
  - [Authentication Methods](#authentication-methods)
  - [Interface](#interface)
  - [Security Considerations](#security-considerations)
  - [Interactions with Other Services](#interactions-with-other-services)

## Overview

The Certificate Issuer microservice manages the PKI infrastructure for the ubuntu-autoinstall-webhook system. It handles certificate requests, issuance, and revocation for secure communication between clients and services.

## Core Responsibilities

- Managing the internal certificate authority (CA)
- Processing certificate signing requests (CSRs)
- Issuing client and server certificates
- Maintaining certificate revocation lists (CRLs)
- Supporting mutual TLS authentication

## Certificate Authority Management

### Root CA
- Self-signed certificate with 20-year validity
- Generated on initial startup unless provided
- Stored securely in the database or filesystem
- Used only to sign intermediate CAs

### Intermediate CA
- Signed by the root CA with 2-year validity
- Used for signing all server and client certificates
- Can be rotated without affecting the root CA
- Full chain provided with issued certificates

### External CA Integration
- Optional integration with cert-manager
- Support for external PKI systems
- Ability to import existing CA certificates

## Certificate Lifecycle

### Certificate Issuance
- Processes CSRs from clients and services
- Validates CSR information against system records
- Issues certificates with appropriate validity periods
- Returns certificates with the full trust chain

### Certificate Renewal
- Monitors certificate expiration
- Supports automatic renewal workflows
- Handles grace periods for renewal

### Certificate Revocation
- Maintains CRLs for revoked certificates
- Supports OCSP for online revocation checking
- Handles emergency revocation procedures

## Authentication Methods

For initial certificate requests, the service supports:

- Pre-shared secrets
- IP address verification
- MAC address verification
- Existing client certificates
- External validation methods

## Interface

The Certificate Issuer service exposes a gRPC interface that provides:

- CSR submission
- Certificate issuance
- Revocation requests
- CA certificate retrieval
- Certificate validation

## Security Considerations

- Private keys never leave the system where they are generated
- Root CA keys are protected with additional security measures
- All certificate operations are logged for audit purposes
- Certificate templates enforce security best practices

## Interactions with Other Services

- Validates system identity with the Database service
- Provides certificates to the Webserver service for HTTPS
- Issues certificates to installation clients for mutual TLS
- Coordinates with the Configuration service for certificate policies
