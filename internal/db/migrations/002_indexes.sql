-- Indexes for fast retrieval


CREATE INDEX IF NOT EXISTS idx_client_logs_client_id ON client_logs(client_id);
CREATE INDEX IF NOT EXISTS idx_client_status_client_id ON client_status(client_id);
CREATE INDEX IF NOT EXISTS idx_webhook_logs_timestamp ON webhook_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_ipxe_config_identifier ON ipxe_configurations(identifier);
CREATE INDEX IF NOT EXISTS idx_cloud_init_userdata_mac ON cloud_init_userdata(mac_address);
CREATE INDEX IF NOT EXISTS idx_cloud_init_metadata_mac ON cloud_init_metadata(mac_address);
CREATE INDEX IF NOT EXISTS idx_cloud_init_network_mac ON cloud_init_network(mac_address);


-- Add more indexes as needed


-- Index for faster driver lookups in network interfaces
CREATE INDEX IF NOT EXISTS idx_network_interfaces_driver ON network_interfaces(driver);

-- Index for retrieving cloud-init history by MAC address
CREATE INDEX IF NOT EXISTS idx_cloud_init_history_mac ON cloud_init_history(mac_address);
