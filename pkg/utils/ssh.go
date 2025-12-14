package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConfig holds SSH connection configuration
type SSHConfig struct {
	Host      string        // Host address (IP or hostname)
	Port      int           // SSH port (default: 22)
	User      string        // SSH username
	Password  string        // SSH password (optional if using key)
	KeyPath   string        // Path to private key file (optional if using password)
	Timeout   time.Duration // Connection timeout
	KeepAlive time.Duration // Keep-alive interval
}

// SSHClient wraps an SSH client connection
type SSHClient struct {
	config *SSHConfig
	client *ssh.Client
}

// NewSSHClient creates a new SSH client with the given configuration
func NewSSHClient(config *SSHConfig) (*SSHClient, error) {
	if config.Port == 0 {
		config.Port = 22
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.KeepAlive == 0 {
		config.KeepAlive = 10 * time.Second
	}

	return &SSHClient{
		config: config,
	}, nil
}

// Connect establishes an SSH connection
func (c *SSHClient) Connect() error {
	// Build SSH client config
	sshConfig := &ssh.ClientConfig{
		User:            c.config.User,
		Timeout:         c.config.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Add proper host key verification
	}

	// Add authentication methods
	var authMethods []ssh.AuthMethod

	// Try key-based authentication first
	if c.config.KeyPath != "" {
		key, err := os.ReadFile(c.config.KeyPath)
		if err != nil {
			return fmt.Errorf("failed to read private key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Add password authentication
	if c.config.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.config.Password))
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("no authentication method provided (password or key)")
	}

	sshConfig.Auth = authMethods

	// Connect to SSH server
	addr := net.JoinHostPort(c.config.Host, fmt.Sprintf("%d", c.config.Port))
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	c.client = client

	// Start keep-alive goroutine
	go c.keepAlive()

	return nil
}

// keepAlive sends periodic keep-alive messages
func (c *SSHClient) keepAlive() {
	if c.client == nil {
		return
	}

	ticker := time.NewTicker(c.config.KeepAlive)
	defer ticker.Stop()

	for range ticker.C {
		if c.client == nil {
			return
		}
		_, _, err := c.client.SendRequest("keepalive@openssh.com", true, nil)
		if err != nil {
			return
		}
	}
}

// Close closes the SSH connection
func (c *SSHClient) Close() error {
	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		return err
	}
	return nil
}

// ExecuteCommand executes a command on the remote host and returns stdout, stderr, and error
func (c *SSHClient) ExecuteCommand(command string) (stdout, stderr string, err error) {
	if c.client == nil {
		return "", "", fmt.Errorf("SSH client not connected")
	}

	// Create a new session
	session, err := c.client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	// Execute command
	err = session.Run(command)
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		return stdout, stderr, fmt.Errorf("command execution failed: %w", err)
	}

	return stdout, stderr, nil
}

// ExecuteCommandWithCallback executes a command and streams output to callback functions
func (c *SSHClient) ExecuteCommandWithCallback(command string, stdoutCallback, stderrCallback func(string)) error {
	if c.client == nil {
		return fmt.Errorf("SSH client not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up pipes for stdout and stderr
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start command
	if err := session.Start(command); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Read stdout
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdoutPipe.Read(buf)
			if n > 0 && stdoutCallback != nil {
				stdoutCallback(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// Read stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderrPipe.Read(buf)
			if n > 0 && stderrCallback != nil {
				stderrCallback(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// Wait for command to finish
	return session.Wait()
}

// UploadFile uploads a local file to the remote host using SCP
func (c *SSHClient) UploadFile(localPath, remotePath string) error {
	if c.client == nil {
		return fmt.Errorf("SSH client not connected")
	}

	// Read local file
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	// Get file info for permissions
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("failed to stat local file: %w", err)
	}

	// Create remote directory if needed
	remoteDir := filepath.Dir(remotePath)
	_, _, err = c.ExecuteCommand(fmt.Sprintf("mkdir -p %s", remoteDir))
	if err != nil {
		return fmt.Errorf("failed to create remote directory: %w", err)
	}

	// Create session for SCP
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up stdin pipe
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start SCP command
	go func() {
		defer stdin.Close()
		// Send file header
		fmt.Fprintf(stdin, "C%04o %d %s\n", fileInfo.Mode().Perm(), len(content), filepath.Base(remotePath))
		// Send file content
		stdin.Write(content)
		// Send termination byte
		fmt.Fprint(stdin, "\x00")
	}()

	// Execute SCP command
	if err := session.Run(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		return fmt.Errorf("SCP upload failed: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from the remote host to local path
func (c *SSHClient) DownloadFile(remotePath, localPath string) error {
	if c.client == nil {
		return fmt.Errorf("SSH client not connected")
	}

	// Create session
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Get stdout pipe
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start cat command to read file
	if err := session.Start(fmt.Sprintf("cat %s", remotePath)); err != nil {
		return fmt.Errorf("failed to start cat command: %w", err)
	}

	// Create local directory if needed
	localDir := filepath.Dir(localPath)
	if err := EnsureDir(localDir); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	// Copy content
	if _, err := io.Copy(localFile, stdout); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Wait for command to finish
	if err := session.Wait(); err != nil {
		return fmt.Errorf("cat command failed: %w", err)
	}

	return nil
}

// TestConnection tests if SSH connection can be established
func (c *SSHClient) TestConnection() error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	// Try to execute a simple command
	_, _, err := c.ExecuteCommand("echo test")
	return err
}
