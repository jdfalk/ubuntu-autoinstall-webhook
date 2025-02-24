package db

import (
	"encoding/json"
	"testing"
	"time"
)

func TestModelsJSONMarshalling(t *testing.T) {
	now := time.Now()
	uuid := "uuid-1234"
	tests := []struct {
		name  string
		model interface{}
	}{
		{
			name: "ClientIdentification",
			model: ClientIdentification{
				ID:                "client-1",
				SMBIOSUUID:        &uuid,
				MotherboardSerial: "mb-serial-123",
				CreatedAt:         now,
			},
		},
		{
			name: "NetworkInterface",
			model: NetworkInterface{
				ID:            "net-1",
				ClientID:      "client-1",
				MacAddress:    "aa:bb:cc:dd:ee:ff",
				InterfaceName: "eth0",
				CreatedAt:     now,
			},
		},
		{
			name: "NetworkChipset",
			model: NetworkChipset{
				ID:                 "chip-1",
				NetworkInterfaceID: "net-1",
				Chipset:            "chipset-model",
				CreatedAt:          now,
			},
		},
		{
			name: "HardwareInfo",
			model: HardwareInfo{
				ID:         "hw-1",
				ClientID:   "client-1",
				LSHWOutput: "lshw output sample",
				CreatedAt:  now,
			},
		},
		{
			name: "ClientLog",
			model: ClientLog{
				ID:          "log-1",
				ClientID:    "client-1",
				Timestamp:   now,
				Origin:      "origin-sample",
				Description: "description sample",
				Name:        "log name",
				Result:      "result sample",
				EventType:   "event-sample",
				Files:       "{}",
				CreatedAt:   now,
			},
		},
		{
			name: "ClientStatus",
			model: ClientStatus{
				ID:        "status-1",
				ClientID:  "client-1",
				Status:    "running",
				Progress:  50,
				Message:   "in progress",
				UpdatedAt: now,
			},
		},
		{
			name: "WebhookLog",
			model: WebhookLog{
				ID:           "wh-1",
				Timestamp:    now,
				RequestData:  `{"key":"value"}`,
				ResponseData: `{"key":"value"}`,
			},
		},
		{
			name: "IPXEConfiguration",
			model: IPXEConfiguration{
				ID:         "ipxe-1",
				ClientID:   "client-1",
				Identifier: "boot-identifier",
				Config:     "pxe config",
				CreatedAt:  now,
			},
		},
		{
			name: "CloudInitUserData",
			model: CloudInitUserData{
				ID:         "clouduser-1",
				ClientID:   "client-1",
				MacAddress: "aa:bb:cc:dd:ee:ff",
				UserData:   "user-data content",
				CreatedAt:  now,
			},
		},
		{
			name: "CloudInitMetaData",
			model: CloudInitMetaData{
				ID:         "cloudmeta-1",
				ClientID:   "client-1",
				MacAddress: "aa:bb:cc:dd:ee:ff",
				MetaData:   "meta-data content",
				CreatedAt:  now,
			},
		},
		{
			name: "CloudInitNetworkConfig",
			model: CloudInitNetworkConfig{
				ID:            "cloudnet-1",
				ClientID:      "client-1",
				MacAddress:    "aa:bb:cc:dd:ee:ff",
				NetworkConfig: "network configuration content",
				CreatedAt:     now,
			},
		},
		{
			name: "CloudInitHistory",
			model: CloudInitHistory{
				ID:         "cloudhist-1",
				ClientID:   "client-1",
				MacAddress: "aa:bb:cc:dd:ee:ff",
				UserData:   "historical user-data",
				CreatedAt:  now,
			},
		},
		{
			name: "ServerLog",
			model: ServerLog{
				ID:        "serverlog-1",
				Message:   "server started",
				CreatedAt: now,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.model)
			if err != nil {
				t.Errorf("Failed to marshal %s: %v", tc.name, err)
			}
			if len(data) == 0 {
				t.Errorf("Marshalled JSON for %s is empty", tc.name)
			}
		})
	}
}
