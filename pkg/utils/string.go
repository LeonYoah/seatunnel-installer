package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty checks if a string is not empty
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// TrimAll trims all whitespace from a string
func TrimAll(s string) string {
	return strings.TrimSpace(s)
}

// Contains checks if a string contains a substring (case-insensitive)
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// EqualIgnoreCase compares two strings ignoring case
func EqualIgnoreCase(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

// SplitLines splits a string into lines
func SplitLines(s string) []string {
	return strings.Split(s, "\n")
}

// JoinLines joins lines into a single string
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// RemoveEmptyLines removes empty lines from a slice of strings
func RemoveEmptyLines(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if !IsEmpty(line) {
			result = append(result, line)
		}
	}
	return result
}

// Truncate truncates a string to a maximum length
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// PadLeft pads a string on the left with a character to reach a minimum length
func PadLeft(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(pad), length-len(s)) + s
}

// PadRight pads a string on the right with a character to reach a minimum length
func PadRight(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(string(pad), length-len(s))
}

// MD5Hash returns the MD5 hash of a string
func MD5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA256Hash returns the SHA256 hash of a string
func SHA256Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// IsValidIP checks if a string is a valid IP address (IPv4 or IPv6)
func IsValidIP(ip string) bool {
	// Simple regex for IPv4
	ipv4Pattern := `^(\d{1,3}\.){3}\d{1,3}$`
	matched, _ := regexp.MatchString(ipv4Pattern, ip)
	if !matched {
		return false
	}

	// Validate each octet
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

// IsValidPort checks if a port number is valid (1-65535)
func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}

// IsValidPath checks if a string is a valid file path
func IsValidPath(path string) bool {
	if IsEmpty(path) {
		return false
	}
	// Check for invalid characters (basic check)
	invalidChars := []string{"\x00", "\n", "\r"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return false
		}
	}
	return true
}

// SanitizePath removes potentially dangerous characters from a path
func SanitizePath(path string) string {
	// Remove null bytes and newlines
	path = strings.ReplaceAll(path, "\x00", "")
	path = strings.ReplaceAll(path, "\n", "")
	path = strings.ReplaceAll(path, "\r", "")
	return path
}

// ExtractVariables extracts variables from a template string
// Supports ${VAR} and {{VAR}} formats
func ExtractVariables(template string) []string {
	var vars []string
	seen := make(map[string]bool)

	// Extract ${VAR} format
	re1 := regexp.MustCompile(`\$\{([^}]+)\}`)
	matches1 := re1.FindAllStringSubmatch(template, -1)
	for _, match := range matches1 {
		if len(match) > 1 && !seen[match[1]] {
			vars = append(vars, match[1])
			seen[match[1]] = true
		}
	}

	// Extract {{VAR}} format
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

// ReplaceMultiple replaces multiple substrings in a string
func ReplaceMultiple(s string, replacements map[string]string) string {
	result := s
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// IndentLines indents each line of a string with the given prefix
func IndentLines(s string, indent string) string {
	lines := SplitLines(s)
	for i, line := range lines {
		if !IsEmpty(line) {
			lines[i] = indent + line
		}
	}
	return JoinLines(lines)
}

// RemoveComments removes shell-style comments from a string
func RemoveComments(s string) string {
	lines := SplitLines(s)
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		// Find comment start
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

// ParseKeyValue parses a key=value string
func ParseKeyValue(s string) (key, value string, ok bool) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

// ParseKeyValueMap parses multiple key=value lines into a map
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
