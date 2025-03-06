package db

import "time"

// ClientIdentification represents client hardware and unique identifiers.
type ClientIdentification struct {
	ID                string    `json:"id"`
	SMBIOSUUID        *string   `json:"smbios_uuid,omitempty"`
	MotherboardSerial string    `json:"motherboard_serial"`
	CreatedAt         time.Time `json:"created_at"`
}

// NetworkInterface represents a network adapter on a client.
type NetworkInterface struct {
	ID            string    `json:"id"`
	ClientID      string    `json:"client_id"`
	MacAddress    string    `json:"mac_address"`
	InterfaceName string    `json:"interface_name"`
	CreatedAt     time.Time `json:"created_at"`
}

// NetworkChipset represents the chipset used in a network adapter.
type NetworkChipset struct {
	ID                 string    `json:"id"`
	NetworkInterfaceID string    `json:"network_interface_id"`
	Chipset            string    `json:"chipset"`
	CreatedAt          time.Time `json:"created_at"`
}

// HardwareInfo stores lshw output for a client.
type HardwareInfo struct {
	ID         string    `json:"id"`
	ClientID   string    `json:"client_id"`
	LSHWOutput string    `json:"lshw_output"`
	CreatedAt  time.Time `json:"created_at"`
}

// ClientLog represents log messages sent by clients.
type ClientLog struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	Timestamp   time.Time `json:"timestamp"`
	Origin      string    `json:"origin"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Result      string    `json:"result"`
	EventType   string    `json:"event_type"`
	Files       string    `json:"files"` // JSONB
	CreatedAt   time.Time `json:"created_at"`
}

// ClientStatus represents the current status of an installation process.
type ClientStatus struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"` // Value between 0-100
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebhookLog represents logs of webhook requests.
type WebhookLog struct {
	ID           string    `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	RequestData  string    `json:"request_data"`  // JSONB
	ResponseData string    `json:"response_data"` // JSONB
}

// IPXEConfiguration represents a PXE boot configuration.
// Added Phase to differentiate between "install" and "post-install".
type IPXEConfiguration struct {
	ID         string    `json:"id"`
	ClientID   string    `json:"client_id"`
	Identifier string    `json:"identifier"`
	Config     string    `json:"config"`
	Phase      string    `json:"phase"` // New field: "install" or "post-install"
	CreatedAt  time.Time `json:"created_at"`
}

// CloudInitUserData represents the cloud-init user-data.
type CloudInitUserData struct {
	ID         string    `json:"id"`
	ClientID   string    `json:"client_id"`
	MacAddress string    `json:"mac_address"`
	UserData   string    `json:"user_data"`
	Phase      string    `json:"phase"` // New field: "install" or "post-install"
	CreatedAt  time.Time `json:"created_at"`
}

// CloudInitMetaData represents the cloud-init meta-data.
type CloudInitMetaData struct {
	ID         string    `json:"id"`
	ClientID   string    `json:"client_id"`
	MacAddress string    `json:"mac_address"`
	MetaData   string    `json:"meta_data"`
	CreatedAt  time.Time `json:"created_at"`
}

// CloudInitNetworkConfig represents the cloud-init network configuration.
type CloudInitNetworkConfig struct {
	ID            string    `json:"id"`
	ClientID      string    `json:"client_id"`
	MacAddress    string    `json:"mac_address"`
	NetworkConfig string    `json:"network_config"`
	CreatedAt     time.Time `json:"created_at"`
}

// CloudInitHistory stores the last five cloud-init configurations per client.
type CloudInitHistory struct {
	ID         string    `json:"id"`
	ClientID   string    `json:"client_id"`
	MacAddress string    `json:"mac_address"`
	UserData   string    `json:"user_data"`
	Phase      string    `json:"phase"` // New field: "install" or "post-install"
	CreatedAt  time.Time `json:"created_at"`
}

// IPXEHistory stores the last five iPXE configuration versions per client.
// Added Phase to allow separate config stacks for "install" and "post-install".
type IPXEHistory struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Config    string    `json:"config"`
	Phase     string    `json:"phase"` // New field: "install" or "post-install"
	CreatedAt time.Time `json:"created_at"`
}

// ServerLog represents a server log entry.
type ServerLog struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// ClientLogDetail provides detailed information for a client log.
type ClientLogDetail struct {
	ID          int
	ClientID    string
	Timestamp   time.Time
	Origin      string
	Description string
	Name        string
	Result      string
	EventType   string
	Files       string
	CreatedAt   time.Time
}

// IpxeConfig represents an iPXE configuration.
type IpxeConfig struct {
	ID        int
	ClientID  string
	Config    string
	CreatedAt time.Time
}

// CloudInitConfig represents a cloud-init configuration.
type CloudInitConfig struct {
	ID         int
	ClientID   string
	MacAddress string
	UserData   string
	CreatedAt  time.Time
}
