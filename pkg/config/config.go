package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type WebUIConfig struct {
	Enabled    bool   `json:"enabled"`
	ListenAddr string `json:"listen_addr"`
}

type NetworkConfig struct {
	GatewayIP  string `json:"gateway_ip"`
	GatewayMAC string `json:"gateway_mac"`
	SubnetMask string `json:"subnet_mask"`
}

type UDSFileConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type UDSAbstractConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type UDPConfig struct {
	Enabled    bool   `json:"enabled"`
	ListenAddr string `json:"listen_addr"`
}

type GRPCConfig struct {
	Enabled    bool   `json:"enabled"`
	ListenAddr string `json:"listen_addr"`
}

type TransportsConfig struct {
	UDSFile    UDSFileConfig    `json:"uds_file"`
	UDSAbstract UDSAbstractConfig `json:"uds_abstract"`
	UDP        UDPConfig        `json:"udp"`
	GRPC       GRPCConfig       `json:"grpc"`
}

type StorageConfig struct {
	LeasesFile string `json:"leases_file"`
}

type Config struct {
	WebUI      WebUIConfig     `json:"webui"`
	Network    NetworkConfig   `json:"network"`
	Transports TransportsConfig `json:"transports"`
	Storage    StorageConfig   `json:"storage"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		WebUI: WebUIConfig{
			Enabled:    true,
			ListenAddr: "0.0.0.0:9000",
		},
		Network: NetworkConfig{
			GatewayIP:  "10.0.0.254",
			GatewayMAC: "02:00:00:00:00:01",
			SubnetMask: "255.255.255.0",
		},
		Transports: TransportsConfig{
			UDSFile: UDSFileConfig{
				Enabled: true,
				Path:    "/data/data/com.termux/files/usr/tmp/vswitch.sock",
			},
			UDSAbstract: UDSAbstractConfig{
				Enabled: true,
				Path:    "@vswitch.sock",
			},
			UDP: UDPConfig{
				Enabled:    true,
				ListenAddr: "0.0.0.0:8080",
			},
			GRPC: GRPCConfig{
				Enabled:    true,
				ListenAddr: "0.0.0.0:50051",
			},
		},
		Storage: StorageConfig{
			LeasesFile: "leases.json",
		},
	}
}
