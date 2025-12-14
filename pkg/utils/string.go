package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// IsEmpty 检查字符串是否为空或仅包含空白字符
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 检查字符串是否非空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// TrimAll 去除字符串的所有空白字符
func TrimAll(s string) string {
	return strings.TrimSpace(s)
}

// ContainsIgnoreCase 检查字符串是否包含子串（忽略大小写）
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// EqualIgnoreCase 比较两个字符串是否相等（忽略大小写）
func EqualIgnoreCase(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

// SplitLines 将字符串分割为多行
func SplitLines(s string) []string {
	return strings.Split(s, "\n")
}

// JoinLines 将多行合并为单个字符串
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// RemoveEmptyLines 从字符串切片中移除空行
func RemoveEmptyLines(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if !IsEmpty(line) {
			result = append(result, line)
		}
	}
	return result
}

// Truncate 截断字符串到最大长度
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// PadLeft 在字符串左侧填充字符以达到最小长度
func PadLeft(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(pad), length-len(s)) + s
}

// PadRight 在字符串右侧填充字符以达到最小长度
func PadRight(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(string(pad), length-len(s))
}

// MD5Hash 返回字符串的MD5哈希值
func MD5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA256Hash 返回字符串的SHA256哈希值
func SHA256Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// IsValidIP 检查字符串是否为有效的IP地址（IPv4或IPv6）
func IsValidIP(ip string) bool {
	// IPv4的简单正则表达式
	ipv4Pattern := `^(\d{1,3}\.){3}\d{1,3}$`

	matched, _ := regexp.MatchString(ipv4Pattern, ip)
	if !matched {
		return false
	}

	// 验证每个八位组
	parts := strings.Split(ip, ".")
	for _, part := range parts {
		var num int
		if _, err := fmt.Sscanf(part, "%d", &num); err != nil {
			return false
		}
		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// IsValidPort 检查端口号是否有效（1-65535）
func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}

// IsValidPath 检查字符串是否为有效的文件路径
func IsValidPath(path string) bool {
	if IsEmpty(path) {
		return false
	}
	// 检查无效字符（基本检查）
	invalidChars := []string{"\x00", "\n", "\r"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return false
		}
	}
	return true
}

// SanitizePath 从路径中移除潜在危险字符
func SanitizePath(path string) string {
	// 移除空字节和换行符
	path = strings.ReplaceAll(path, "\x00", "")
	path = strings.ReplaceAll(path, "\n", "")
	path = strings.ReplaceAll(path, "\r", "")
	return path
}

// ExtractVariables 从模板字符串中提取变量
// 支持 ${VAR} 和 {{VAR}} 格式
func ExtractVariables(template string) []string {
	var vars []string
	seen := make(map[string]bool)

	// 提取 ${VAR} 格式
	re1 := regexp.MustCompile(`\$\{([^}]+)\}`)
	matches1 := re1.FindAllStringSubmatch(template, -1)
	for _, match := range matches1 {
		if len(match) > 1 && !seen[match[1]] {
			vars = append(vars, match[1])
			seen[match[1]] = true
		}
	}

	// 提取 {{VAR}} 格式
	re2 := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches2 := re2.FindAllStringSubmatch(template, -1)
	for _, match := range matches2 {
		if len(match) > 1 && !seen[match[1]] {
			vars = append(vars, match[1])
			seen[match[1]] = true
		}
	}

	return vars
}

// ReplaceMultiple 替换字符串中的多个子串
func ReplaceMultiple(s string, replacements map[string]string) string {
	result := s
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// IndentLines 为字符串的每一行添加给定的缩进前缀
func IndentLines(s string, indent string) string {
	lines := SplitLines(s)
	for i, line := range lines {
		if !IsEmpty(line) {
			lines[i] = indent + line
		}
	}
	return JoinLines(lines)
}

// RemoveComments 从字符串中移除shell风格的注释
func RemoveComments(s string) string {
	lines := SplitLines(s)
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		// 查找注释开始位置
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return JoinLines(result)
}

// ParseKeyValue 解析 key=value 字符串
func ParseKeyValue(s string) (key, value string, ok bool) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

// ParseKeyValueMap 将多行 key=value 解析为map
func ParseKeyValueMap(s string) map[string]string {
	result := make(map[string]string)
	lines := SplitLines(s)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if IsEmpty(line) || strings.HasPrefix(line, "#") {
			continue
		}
		if key, value, ok := ParseKeyValue(line); ok {
			result[key] = value
		}
	}
	return result
}
