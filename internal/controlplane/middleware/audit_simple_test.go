package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditHelperFunctions(t *testing.T) {
	t.Run("TestGetActionFromMethod", func(t *testing.T) {
		tests := []struct {
			method   string
			expected string
		}{
			{"POST", "create"},
			{"PUT", "update"},
			{"PATCH", "update"},
			{"DELETE", "delete"},
			{"GET", "unknown"},
			{"HEAD", "unknown"},
		}

		for _, test := range tests {
			result := getActionFromMethod(test.method)
			assert.Equal(t, test.expected, result, "Method %s should return %s", test.method, test.expected)
		}
	})

	t.Run("TestGetResourceFromPath", func(t *testing.T) {
		tests := []struct {
			path     string
			expected string
		}{
			{"/api/v1/hosts", "hosts"},
			{"/api/v1/clusters/123", "clusters"},
			{"/api/v1/tasks/456/runs", "tasks"},
			{"/hosts", "hosts"},
			{"/", "unknown"},
			{"", "unknown"},
		}

		for _, test := range tests {
			result := getResourceFromPath(test.path)
			assert.Equal(t, test.expected, result, "Path %s should return %s", test.path, test.expected)
		}
	})

	t.Run("TestExtractResourceID", func(t *testing.T) {
		tests := []struct {
			path     string
			expected string
		}{
			{"/api/v1/hosts/550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440000"},
			{"/api/v1/clusters/123", ""},
			{"/api/v1/tasks", ""},
			{"/hosts/550e8400-e29b-41d4-a716-446655440001/status", "550e8400-e29b-41d4-a716-446655440001"},
		}

		for _, test := range tests {
			result := extractResourceID(test.path)
			assert.Equal(t, test.expected, result, "Path %s should return %s", test.path, test.expected)
		}
	})

	t.Run("TestIsWriteOperation", func(t *testing.T) {
		tests := []struct {
			method   string
			expected bool
		}{
			{"POST", true},
			{"PUT", true},
			{"PATCH", true},
			{"DELETE", true},
			{"GET", false},
			{"HEAD", false},
			{"OPTIONS", false},
		}

		for _, test := range tests {
			result := isWriteOperation(test.method)
			assert.Equal(t, test.expected, result, "Method %s should return %t", test.method, test.expected)
		}
	})

	t.Run("TestContainsSensitiveData", func(t *testing.T) {
		tests := []struct {
			body     string
			expected bool
		}{
			{`{"username": "test", "password": "secret"}`, true},
			{`{"name": "test", "value": "normal"}`, false},
			{`{"api_key": "secret123"}`, true},
			{`{"token": "abc123"}`, true},
			{`{"credential": "secret"}`, true},
			{`{"normal": "data"}`, false},
		}

		for _, test := range tests {
			result := containsSensitiveData([]byte(test.body))
			assert.Equal(t, test.expected, result, "Body %s should return %t", test.body, test.expected)
		}
	})
}
