# Troubleshooting Manual for Ubuntu Autoinstall Webhook

## **1. Introduction**

This manual provides solutions to common issues encountered while using the Ubuntu Autoinstall Webhook. It includes diagnostic steps, error messages, and recommended fixes.

## **2. Common Issues & Fixes**

### **2.1. Webhook API Not Responding**

#### **Symptoms:**

- Requests to the webhook return a timeout or connection refused.
- Logs do not show incoming requests.

#### **Diagnostic Steps:**

1. Check if the webhook service is running:

   ```bash
   $ systemctl status autoinstall-webhook
   ```

2. If using Docker, check container logs:

   ```bash
   $ docker logs autoinstall-webhook
   ```

3. Verify that the server firewall allows incoming traffic on the webhook port (default: 5000).

#### **Solution:**

- Restart the webhook service:

  ```bash
  $ systemctl restart autoinstall-webhook
  ```

- If using Docker:

  ```bash
  $ docker restart autoinstall-webhook
  ```

- Check and fix firewall rules:

  ```bash
  $ sudo ufw allow 5000/tcp
  ```

### **2.2. PXE Boot Fails**

#### **Symptoms:**

- PXE clients fail to retrieve `boot.ipxe`.
- Boot process stops at “TFTP Timeout” or “File Not Found.”

#### **Diagnostic Steps:**

1. Check if the PXE boot server is running:

   ```bash
   $ systemctl status tftpd-hpa
   ```

2. Verify that boot files exist in `/var/lib/tftpboot/`.
3. Test manual file retrieval from another machine:

   ```bash
   $ tftp <pxe-server-ip> -c get boot.ipxe
   ```

#### **Solution:**

- Restart the TFTP service:

  ```bash
  $ systemctl restart tftpd-hpa
  ```

- Ensure proper file permissions:

  ```bash
  $ sudo chmod -R 755 /var/lib/tftpboot/
  ```

### **2.3. Cloud-Init Fails to Fetch Configuration**

#### **Symptoms:**

- The installation halts, displaying `ds=nocloud` errors.
- Clients fail to retrieve `user-data`, `meta-data`, or `network-config`.

#### **Diagnostic Steps:**

1. Verify that the Cloud-Init configuration server is accessible:

   ```bash
   $ curl http://172.16.2.30/cloud-init/mac-<mac_address>/user-data
   ```

2. Check server logs for errors related to cloud-init retrieval.
3. Ensure the correct MAC address is being used in the request.

#### **Solution:**

- Restart the Cloud-Init configuration server:

  ```bash
  $ systemctl restart cloud-init-config
  ```

- Verify cloud-init logs on the client:

  ```bash
  $ cat /var/log/cloud-init.log
  ```

### **2.4. Webhook Does Not Receive Logs or Status Updates**

#### **Symptoms:**

- Clients do not appear in the webhook database.
- The `/status` or `/logs` API endpoints show no incoming data.

#### **Diagnostic Steps:**

1. Check webhook logs for incoming requests:

   ```bash
   $ tail -f /var/log/autoinstall-webhook.log
   ```

2. Verify that clients can reach the webhook API:

   ```bash
   $ curl -X POST http://webhook-server:5000/status -d '{"mac_address": "AA:BB:CC:DD:EE:FF", "status": "PXE Booted"}'
   ```

3. Confirm that webhook database entries exist:

   ```bash
   $ sqlite3 webhook.db "SELECT * FROM clients;"
   ```

#### **Solution:**

- Restart the webhook service.
- Check network configurations on both the client and server.

## **3. Additional Resources**

- Webhook logs: `/var/log/autoinstall-webhook.log`
- Cloud-Init logs: `/var/log/cloud-init.log`
- PXE boot logs: `/var/log/syslog`

## **4. Conclusion**

This manual provides troubleshooting guidance for common issues with the Ubuntu Autoinstall Webhook. If problems persist, refer to the logs and network diagnostics for further investigation.
