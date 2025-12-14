package config

import (
	"fmt"
	"net"
	"regexp"
)

// InstallConfig represents the installation configuration from web form
type InstallConfig struct {
	// Basic parameters (required)
	Version     string   `json:"version" binding:"required"`      // SeaTunnel version
	InstallDir  string   `json:"install_dir" binding:"required"`  // Installation directory
	DeployMode  string   `json:"deploy_mode" binding:"required"`  // hybrid or separated
	NodeIPs     []string `json:"node_ips" binding:"required"`     // Node IP list
	InstallMode string   `json:"install_mode" binding:"required"` // online or offline

	// Advanced parameters (optional, with defaults)
	Ports      PortConfig       `json:"ports"`
	Memory     MemoryConfig     `json:"memory"`
	Checkpoint CheckpointConfig `json:"checkpoint"`
	SSH        SSHConfig        `json:"ssh"`
	User       UserConfig       `json:"user"`
	Firewall   FirewallConfig   `json:"firewall"`

	// Plugin selection
	Plugins []string `json:"plugins"`
}

// PortConfig represents port configuration
type PortConfig struct {
	HybridPort     int `json:"hybrid_port"`      // Default: 5801
	MasterPort     int `json:"master_port"`      // Default: 5801
	WorkerPort     int `json:"worker_port"`      // Default: 5802
	MasterHTTPPort int `json:"master_http_port"` // Default: 8080
}

// MemoryConfig represents JVM memory configuration
type MemoryConfig struct {
	HybridHeapSize int `json:"hybrid_heap_size"` // GB, Default: 3
	MasterHeapSize int `json:"master_heap_size"` // GB, Default: 2
	WorkerHeapSize int `json:"worker_heap_size"` // GB, Default: 2
}

// CheckpointConfig represents checkpoint storage configuration
type CheckpointConfig struct {
	StorageType string `json:"storage_type"` // LOCAL_FILE, HDFS, OSS, S3
	Namespace   string `json:"namespace"`    // Storage path

	// HDFS specific
	HDFSNameNode string `json:"hdfs_namenode,omitempty"`
	HDFSPort     int    `json:"hdfs_port,omitempty"`

	// OSS/S3 specific
	Endpoint  string `json:"endpoint,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
}

// SSHConfig represents SSH configuration
type SSHConfig struct {
	Port int `json:"port"` // Default: 22
}

// UserConfig represents user configuration
type UserConfig struct {
	User  string `json:"user"`  // Default: root
	Group string `json:"group"` // Default: root
}

// FirewallConfig represents firewall check configuration
type FirewallConfig struct {
	Check  bool   `json:"check"`  // Default: true
	Action string `json:"action"` // error or warn, Default: error
}

// SetDefaults sets default values for optional parameters
func (c *InstallConfig) SetDefaults() {
	// Port defaults
	if c.Ports.HybridPort == 0 {
		c.Ports.HybridPort = 5801
	}
	if c.Ports.MasterPort == 0 {
		c.Ports.MasterPort = 5801
	}
	if c.Ports.WorkerPort == 0 {
		c.Ports.WorkerPort = 5802
	}
	if c.Ports.MasterHTTPPort == 0 {
		c.Ports.MasterHTTPPort = 8080
	}

	// Memory defaults
	if c.Memory.HybridHeapSize == 0 {
		c.Memory.HybridHeapSize = 3
	}
	if c.Memory.MasterHeapSize == 0 {
		c.Memory.MasterHeapSize = 2
	}
	if c.Memory.WorkerHeapSize == 0 {
		c.Memory.WorkerHeapSize = 2
	}

	// Checkpoint defaults
	if c.Checkpoint.StorageType == "" {
		c.Checkpoint.StorageType = "LOCAL_FILE"
	}

	// SSH defaults
	if c.SSH.Port == 0 {
		c.SSH.Port = 22
	}

	// User defaults
	if c.User.User == "" {
		c.User.User = "root"
	}
	if c.User.Group == "" {
		c.User.Group = "root"
	}

	// Firewall defaults
	if c.Firewall.Action == "" {
		c.Firewall.Action = "error"
	}
}

// Validate validates the installation configuration
func (c *InstallConfig) Validate() error {
	// Validate version
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}
	versionRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !versionRegex.MatchString(c.Version) {
		return fmt.Errorf("invalid version format: %s (expected: x.y.z)", c.Version)
	}

	// Validate install directory
	if c.InstallDir == "" {
		return fmt.Errorf("install_dir is required")
	}

	// Validate deploy mode
	if c.DeployMode != "hybrid" && c.DeployMode != "separated" {
		return fmt.Errorf("invalid deploy_mode: %s (must be 'hybrid' or 'separated')", c.DeployMode)
	}

	// Validate node IPs
	if len(c.NodeIPs) == 0 {
		return fmt.Errorf("node_ips is required")
	}
	for _, ip := range c.NodeIPs {
		if net.ParseIP(ip) == nil {
			return fmt.Errorf("invalid IP address: %s", ip)
		}
	}

	// Validate install mode
	if c.InstallMode != "online" && c.InstallMode != "offline" {
		return fmt.Errorf("invalid install_mode: %s (must be 'online' or 'offline')", c.InstallMode)
	}

	// Validate ports
	if err := validatePort(c.Ports.HybridPort, "hybrid_port"); err != nil {
		return err
	}
	if err := validatePort(c.Ports.MasterPort, "master_port"); err != nil {
		return err
	}
	if err := validatePort(c.Ports.WorkerPort, "worker_port"); err != nil {
		return err
	}
	if err := validatePort(c.Ports.MasterHTTPPort, "master_http_port"); err != nil {
		return err
	}

	// Validate memory
	if c.Memory.HybridHeapSize < 1 {
		return fmt.Errorf("hybrid_heap_size must be at least 1 GB")
	}
	if c.Memory.MasterHeapSize < 1 {
		return fmt.Errorf("master_heap_size must be at least 1 GB")
	}
	if c.Memory.WorkerHeapSize < 1 {
		return fmt.Errorf("worker_heap_size must be at least 1 GB")
	}

	// Validate checkpoint storage
	validStorageTypes := map[string]bool{
		"LOCAL_FILE": true,
		"HDFS":       true,
		"OSS":        true,
		"S3":         true,
	}
	if !validStorageTypes[c.Checkpoint.StorageType] {
		return fmt.Errorf("invalid checkpoint storage_type: %s", c.Checkpoint.StorageType)
	}

	// Validate firewall action
	if c.Firewall.Action != "error" && c.Firewall.Action != "warn" {
		return fmt.Errorf("invalid firewall action: %s (must be 'error' or 'warn')", c.Firewall.Action)
	}

	return nil
}

// validatePort validates a port number
func validatePort(port int, name string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid %s: %d (must be between 1 and 65535)", name, port)
	}
	return nil
}
