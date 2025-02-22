# Client Boot and Install Process Document

## **1. Overview**

This document outlines the step-by-step process for booting and installing Ubuntu using the autoinstall system. It covers PXE boot, iPXE customization, cloud-init execution, and post-install configuration.

## **2. Boot Process**

### **2.1. PXE Boot & iPXE Chainloading**

1. The client boots and requests an IP address via DHCP.
2. The DHCP server responds with the TFTP server location.
3. The client downloads and executes the network bootloader (iPXE).
4. iPXE fetches `http://172.16.2.30/ipxe/boot.ipxe` and executes it.
5. The boot script dynamically loads `mac-<mac_address>.ipxe`, customized per client.
6. The customized script loads `menu.ipxe` with adjusted defaults.

## **3. Installation Process**

### **3.1. Live Ubuntu Boot with Cloud-Init**

1. The client boots into a live Ubuntu instance.
2. The kernel command line includes `ds=nocloud;s=http://172.16.2.30/cloud-init/${mac:hexraw}/`.
3. The cloud-init service fetches `user-data`, `meta-data`, and `network-config` from the server.
4. The installation script formats the NVMe drive and installs Ubuntu.

### **3.2. Sending Logs & Status Updates**

1. The client periodically sends installation logs to the webhook API.
2. The webhook updates the installation state for tracking.
3. In case of an error, logs are stored for troubleshooting.

### **3.3. Post-Install Configuration**

1. Once installation completes, the client reboots.
2. The webhook updates `mac-<mac_address>.ipxe` to load a post-install cloud-init.
3. The final cloud-init script completes configuration (e.g., user setup, software install).
4. The client sends a final success status to the webhook.

## **4. Error Handling & Recovery**

- If a client fails at any step, it reboots and retries.
- The webhook can dynamically modify iPXE scripts to trigger a recovery mode.
- Admins can manually trigger a reinstall via the webhook API.

## **5. Conclusion**

This document outlines the client’s boot and installation process, ensuring a structured deployment for Ubuntu autoinstallation.
