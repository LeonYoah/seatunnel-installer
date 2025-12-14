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

// SSHConfig SSH连接配置
type SSHConfig struct {
	Host      string        // 主机地址（IP或主机名）
	Port      int           // SSH端口（默认：22）
	User      string        // SSH用户名
	Password  string        // SSH密码（使用密钥时可选）
	KeyPath   string        // 私钥文件路径（使用密码时可选）
	Timeout   time.Duration // 连接超时时间
	KeepAlive time.Duration // 保持连接间隔
}

// SSHClient 封装SSH客户端连接
type SSHClient struct {
	config *SSHConfig
	client *ssh.Client
}

// NewSSHClient 使用给定配置创建新的SSH客户端
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

// Connect 建立SSH连接
func (c *SSHClient) Connect() error {
	// 构建SSH客户端配置
	sshConfig := &ssh.ClientConfig{
		User:            c.config.User,
		Timeout:         c.config.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: 添加适当的主机密钥验证
	}

	// 添加认证方法
	var authMethods []ssh.AuthMethod

	// 首先尝试基于密钥的认证
	if c.config.KeyPath != "" {
		key, err := os.ReadFile(c.config.KeyPath)
		if err != nil {
			return fmt.Errorf("读取私钥失败: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("解析私钥失败: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// 添加密码认证
	if c.config.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.config.Password))
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("未提供认证方法（密码或密钥）")
	}

	sshConfig.Auth = authMethods

	// 连接到SSH服务器
	addr := net.JoinHostPort(c.config.Host, fmt.Sprintf("%d", c.config.Port))
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("连接SSH服务器失败: %w", err)
	}

	c.client = client

	// 启动保持连接的goroutine
	go c.keepAlive()

	return nil
}

// keepAlive 定期发送保持连接消息
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

// Close 关闭SSH连接
func (c *SSHClient) Close() error {
	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		return err
	}
	return nil
}

// ExecuteCommand 在远程主机上执行命令并返回标准输出、标准错误和错误
func (c *SSHClient) ExecuteCommand(command string) (stdout, stderr string, err error) {
	if c.client == nil {
		return "", "", fmt.Errorf("SSH客户端未连接")
	}

	// 创建新会话
	session, err := c.client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("创建会话失败: %w", err)
	}
	defer session.Close()

	// 捕获标准输出和标准错误
	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	// 执行命令
	err = session.Run(command)
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		return stdout, stderr, fmt.Errorf("命令执行失败: %w", err)
	}

	return stdout, stderr, nil
}

// ExecuteCommandWithCallback 执行命令并将输出流式传输到回调函数
func (c *SSHClient) ExecuteCommandWithCallback(command string, stdoutCallback, stderrCallback func(string)) error {
	if c.client == nil {
		return fmt.Errorf("SSH客户端未连接")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}
	defer session.Close()

	// 设置标准输出和标准错误的管道
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建标准输出管道失败: %w", err)
	}

	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建标准错误输出管道失败: %w", err)
	}

	// 启动命令
	if err := session.Start(command); err != nil {
		return fmt.Errorf("启动命令失败: %w", err)
	}

	// 读取标准输出
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

	// 读取标准错误输出
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

	// 等待命令完成
	return session.Wait()
}

// UploadFile 使用SCP将本地文件上传到远程主机
func (c *SSHClient) UploadFile(localPath, remotePath string) error {
	if c.client == nil {
		return fmt.Errorf("SSH客户端未连接")
	}

	// 读取本地文件
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("读取本地文件失败: %w", err)
	}

	// 获取文件信息以保留权限
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("获取本地文件信息失败: %w", err)
	}

	// 如需要则创建远程目录
	remoteDir := filepath.Dir(remotePath)
	_, _, err = c.ExecuteCommand(fmt.Sprintf("mkdir -p %s", remoteDir))
	if err != nil {
		return fmt.Errorf("创建远程目录失败: %w", err)
	}

	// 为SCP创建会话
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}
	defer session.Close()

	// 设置标准输入管道
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建标准输入管道失败: %w", err)
	}

	// 启动SCP命令
	go func() {
		defer stdin.Close()
		// 发送文件头
		fmt.Fprintf(stdin, "C%04o %d %s\n", fileInfo.Mode().Perm(), len(content), filepath.Base(remotePath))
		// 发送文件内容
		stdin.Write(content)
		// 发送终止字节
		fmt.Fprint(stdin, "\x00")
	}()

	// 执行SCP命令
	if err := session.Run(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		return fmt.Errorf("SCP上传失败: %w", err)
	}

	return nil
}

// DownloadFile 从远程主机下载文件到本地路径
func (c *SSHClient) DownloadFile(remotePath, localPath string) error {
	if c.client == nil {
		return fmt.Errorf("SSH客户端未连接")
	}

	// 创建会话
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}
	defer session.Close()

	// 获取标准输出管道
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建标准输出管道失败: %w", err)
	}

	// 启动cat命令读取文件
	if err := session.Start(fmt.Sprintf("cat %s", remotePath)); err != nil {
		return fmt.Errorf("启动cat命令失败: %w", err)
	}

	// 如需要则创建本地目录
	localDir := filepath.Dir(localPath)
	if err := EnsureDir(localDir); err != nil {
		return fmt.Errorf("创建本地目录失败: %w", err)
	}

	// 创建本地文件
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %w", err)
	}
	defer localFile.Close()

	// 复制内容
	if _, err := io.Copy(localFile, stdout); err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 等待命令完成
	if err := session.Wait(); err != nil {
		return fmt.Errorf("cat命令失败: %w", err)
	}

	return nil
}

// TestConnection 测试是否可以建立SSH连接
func (c *SSHClient) TestConnection() error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	// 尝试执行简单命令
	_, _, err := c.ExecuteCommand("echo test")
	return err
}
