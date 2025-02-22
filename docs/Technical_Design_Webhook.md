# Technical Design Document for Ubuntu Autoinstall Webhook

## **1. Overview**

The Ubuntu Autoinstall Webhook is responsible for managing iPXE boot customization, tracking installation progress, and ensuring proper cloud-init execution. This document outlines the architecture, components, and technical details of the webhook service.

## **2. Architecture**

### **2.1. High-Level Architecture**

- **Clients** (PXE Booting Machines)
- **Webhook API** (Handles logs, hardware info, installation state tracking)
- **Database** (Stores client data, logs, and installation status)
- **iPXE Configuration Server** (Serves dynamic iPXE boot files)
- **Cloud-Init Configuration Server** (Provides correct cloud-init configurations)

### **2.2. Component Breakdown**

#### **2.2.1. Webhook API**

- **Endpoints:**
  - `POST /hardware-info` - Stores client hardware details
  - `POST /logs` - Receives logs from clients
  - `POST /status` - Updates installation progress
  - `GET /status/{mac_address}` - Returns installation state for a client
  - `POST /reinstall/{mac_address}` - Triggers reinstallation for a client

#### **2.2.2. Database**

- **Tables:**
  - `clients`: Stores MAC address, hardware details, and install status.
  - `logs`: Stores installation logs from clients.
  - `ipxe_configurations`: Stores dynamic iPXE boot files.

#### **2.2.3. iPXE Configuration Server**

- **Generates dynamic boot files:**
  - Example: `mac-<mac_address>.ipxe`
  - Points to the correct cloud-init configuration based on client status.

#### **2.2.4. Cloud-Init Configuration Server**

- Serves `user-data`, `meta-data`, `network-config` based on client MAC address.
- Adjusts settings dynamically based on installation progress.

## **3. Data Flow**

1. Client PXE boots and requests an iPXE script from `http://172.16.2.30/ipxe/boot.ipxe`.
2. iPXE loads a customized file for the client (`mac-<mac_address>.ipxe`).
3. Client loads a live Ubuntu image and fetches cloud-init config.
4. Installation process begins; logs and hardware info are sent to the webhook.
5. Webhook updates installation state and modifies iPXE files accordingly.
6. On first boot, the client receives a post-install cloud-init configuration.

## **4. Error Handling & Recovery**

- **Installation Failures:**
  - Logs errors and triggers corrective action.
  - Updates iPXE configuration for recovery boot.
- **Webhook API Errors:**
  - Validates JSON input.
  - Returns appropriate error codes (400, 403, 500).
- **Security Measures:**
  - Only known MAC addresses can trigger installation.
  - API requests require authentication.

## **5. Deployment & Scaling**

- **Containerized Deployment:**
  - Webhook runs as a Docker container.
  - Database runs separately for persistence.
- **Scalability Considerations:**
  - Load-balanced API instances.
  - Asynchronous processing for log storage.

## **6. Conclusion**

This document outlines the architecture and components of the Ubuntu Autoinstall Webhook, ensuring efficient and scalable operation for automated Ubuntu deployments.
