-- Initial schema for CockroachDB
-- Table: system_logs
CREATE TABLE IF NOT EXISTS system_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: client_identification
CREATE TABLE IF NOT EXISTS client_identification (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    smbios_uuid UUID UNIQUE,
    motherboard_serial TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: network_interfaces
CREATE TABLE IF NOT EXISTS network_interfaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    mac_address TEXT UNIQUE NOT NULL,
    interface_name TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: network_chipsets
CREATE TABLE IF NOT EXISTS network_chipsets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    network_interface_id UUID REFERENCES network_interfaces(id) ON DELETE CASCADE,
    chipset TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: hardware_info
CREATE TABLE IF NOT EXISTS hardware_info (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    lshw_output TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: client_logs
CREATE TABLE IF NOT EXISTS client_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    timestamp TIMESTAMP DEFAULT current_timestamp,
    origin TEXT,
    description TEXT,
    name TEXT,
    result TEXT,
    event_type TEXT,
    files JSONB,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: client_status
CREATE TABLE IF NOT EXISTS client_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    progress INT CHECK (
        progress BETWEEN 0 AND 100
    ),
    message TEXT,
    updated_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: webhook_logs
CREATE TABLE IF NOT EXISTS webhook_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMP DEFAULT current_timestamp,
    request_data JSONB,
    response_data JSONB
);
-- Table: ipxe_configurations
CREATE TABLE IF NOT EXISTS ipxe_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    identifier TEXT UNIQUE NOT NULL,
    config TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: cloud_init_userdata
CREATE TABLE IF NOT EXISTS cloud_init_userdata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    mac_address TEXT UNIQUE NOT NULL,
    user_data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: cloud_init_metadata
CREATE TABLE IF NOT EXISTS cloud_init_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    mac_address TEXT UNIQUE NOT NULL,
    meta_data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: cloud_init_network
CREATE TABLE IF NOT EXISTS cloud_init_network (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    mac_address TEXT UNIQUE NOT NULL,
    network_config TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: cloud_init_history
CREATE TABLE IF NOT EXISTS cloud_init_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    mac_address TEXT NOT NULL,
    user_data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
-- Table: ipxe_history
CREATE TABLE IF NOT EXISTS ipxe_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES client_identification(id) ON DELETE CASCADE,
    config TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
