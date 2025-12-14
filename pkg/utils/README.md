# 工具包

此包为 SeaTunnel 企业平台提供了一套全面的实用工具函数。

## 模块

### 文件操作 (`file.go`)

提供文件和目录操作工具：

- **FileExists(path)** - 检查文件或目录是否存在
- **IsDir(path)** - 检查路径是否为目录
- **EnsureDir(path)** - 如果目录不存在则创建
- **ReadFile(path)** - 读取整个文件内容
- **WriteFile(path, content, perm)** - 将内容写入文件
- **CopyFile(src, dst)** - 复制文件及其权限
- **CopyDir(src, dst)** - 递归复制目录
- **RemoveAll(path)** - 递归删除文件或目录
- **Chmod(path, mode)** - 更改文件权限
- **Chown(path, uid, gid)** - 更改文件所有者
- **GetFileSize(path)** - 获取文件大小（字节）
- **ListFiles(dir)** - 列出目录中的所有文件
- **ListDirs(dir)** - 列出所有子目录

### SSH 操作 (`ssh.go`)

提供 SSH 连接和远程命令执行：
```go
// 创建 SSH 客户端
config := &SSHConfig{
    Host:     "192.168.1.100",
    Port:     22,
    User:     "admin",
    Password: "password",
    // 或使用 KeyPath: "/path/to/private/key"
}
client, err := NewSSHClient(config)

// 连接
err = client.Connect()
defer client.Close()

// 执行命令
stdout, stderr, err := client.ExecuteCommand("ls -la")

// 上传文件
err = client.UploadFile("/local/path", "/remote/path")

// 下载文件
err = client.DownloadFile("/remote/path", "/local/path")
```
