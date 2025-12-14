package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileExists 检查文件或目录是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir 检查路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir 如果目录不存在则创建
func EnsureDir(path string) error {
	if FileExists(path) {
		if !IsDir(path) {
			return fmt.Errorf("路径存在但不是目录: %s", path)
		}
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// ReadFile 读取整个文件内容
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile 将内容写入文件，如需要则创建父目录
func WriteFile(path string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	return os.WriteFile(path, content, perm)
}

// CopyFile 从src复制文件到dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 获取源文件信息以保留权限
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	// 如需要则创建目标目录
	dstDir := filepath.Dir(dst)
	if err := EnsureDir(dstDir); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	// 复制内容
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 设置权限
	if err := os.Chmod(dst, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("设置文件权限失败: %w", err)
	}

	return nil
}

// CopyDir 递归复制目录从src到dst
func CopyDir(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源目录信息失败: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("源路径不是目录: %s", src)
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 读取目录条目
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	// 复制每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveAll 递归删除文件或目录
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Chmod 修改文件权限
func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// Chown 修改文件所有者和组
func Chown(path string, uid, gid int) error {
	return os.Chown(path, uid, gid)
}

// GetFileSize 返回文件大小（字节）
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ListFiles 返回目录中的所有文件（非递归）
func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	return files, nil
}

// ListDirs 返回目录中的所有子目录（非递归）
func ListDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(dir, entry.Name()))
		}
	}
	return dirs, nil
}
