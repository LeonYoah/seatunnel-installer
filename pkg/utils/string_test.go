package utils

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"hello", false},
		{" hello ", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsEmpty(tt.input); got != tt.want {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"/"+tt.substr, func(t *testing.T) {
			if got := ContainsIgnoreCase(tt.s, tt.substr); got != tt.want {
				t.Errorf("ContainsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello..."},
		{"test", 4, "test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Truncate(tt.input, tt.maxLen); got != tt.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		input  string
		length int
		pad    rune
		want   string
	}{
		{"5", 3, '0', "005"},
		{"hello", 3, ' ', "hello"},
		{"x", 5, '-', "----x"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := PadLeft(tt.input, tt.length, tt.pad); got != tt.want {
				t.Errorf("PadLeft(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.pad, got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		input  string
		length int
		pad    rune
		want   string
	}{
		{"5", 3, '0', "500"},
		{"hello", 3, ' ', "hello"},
		{"x", 5, '-', "x----"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := PadRight(tt.input, tt.length, tt.pad); got != tt.want {
				t.Errorf("PadRight(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.pad, got, tt.want)
			}
		})
	}
}

func TestMD5Hash(t *testing.T) {
	input := "hello"
	hash := MD5Hash(input)

	// MD5 hash should be 32 characters (hex)
	if len(hash) != 32 {
		t.Errorf("MD5Hash length = %d, want 32", len(hash))
	}

	// Same input should produce same hash
	hash2 := MD5Hash(input)
	if hash != hash2 {
		t.Error("MD5Hash should be deterministic")
	}
}

func TestSHA256Hash(t *testing.T) {
	input := "hello"
	hash := SHA256Hash(input)

	// SHA256 hash should be 64 characters (hex)
	if len(hash) != 64 {
		t.Errorf("SHA256Hash length = %d, want 64", len(hash))
	}

	// Same input should produce same hash
	hash2 := SHA256Hash(input)
	if hash != hash2 {
		t.Error("SHA256Hash should be deterministic")
	}
}

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"256.1.1.1", false},
		{"192.168.1", false},
		{"192.168.1.1.1", false},
		{"abc.def.ghi.jkl", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			if got := IsValidIP(tt.ip); got != tt.want {
				t.Errorf("IsValidIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsValidPort(t *testing.T) {
	tests := []struct {
		port int
		want bool
	}{
		{80, true},
		{443, true},
		{8080, true},
		{65535, true},
		{0, false},
		{-1, false},
		{65536, false},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.port)), func(t *testing.T) {
			if got := IsValidPort(tt.port); got != tt.want {
				t.Errorf("IsValidPort(%d) = %v, want %v", tt.port, got, tt.want)
			}
		})
	}
}

func TestParseKeyValue(t *testing.T) {
	tests := []struct {
		input     string
		wantKey   string
		wantValue string
		wantOk    bool
	}{
		{"key=value", "key", "value", true},
		{"name = John", "name", "John", true},
		{"invalid", "", "", false},
		{"a=b=c", "a", "b=c", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, value, ok := ParseKeyValue(tt.input)
			if ok != tt.wantOk {
				t.Errorf("ParseKeyValue(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if key != tt.wantKey {
				t.Errorf("ParseKeyValue(%q) key = %q, want %q", tt.input, key, tt.wantKey)
			}
			if value != tt.wantValue {
				t.Errorf("ParseKeyValue(%q) value = %q, want %q", tt.input, value, tt.wantValue)
			}
		})
	}
}

func TestParseKeyValueMap(t *testing.T) {
	input := `
# Comment line
key1=value1
key2 = value2

key3=value3
`
	result := ParseKeyValueMap(input)

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	if len(result) != len(expected) {
		t.Errorf("Map length = %d, want %d", len(result), len(expected))
	}

	for key, want := range expected {
		if got, ok := result[key]; !ok || got != want {
			t.Errorf("result[%q] = %q, want %q", key, got, want)
		}
	}
}

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		template string
		want     []string
	}{
		{"Hello ${name}", []string{"name"}},
		{"{{host}}:{{port}}", []string{"host", "port"}},
		{"${var1} and {{var2}}", []string{"var1", "var2"}},
		{"No variables here", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
			got := ExtractVariables(tt.template)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractVariables(%q) returned %d variables, want %d", tt.template, len(got), len(tt.want))
				return
			}
			for i, v := range tt.want {
				if got[i] != v {
					t.Errorf("ExtractVariables(%q)[%d] = %q, want %q", tt.template, i, got[i], v)
				}
			}
		})
	}
}

func TestRemoveEmptyLines(t *testing.T) {
	input := []string{
		"line1",
		"",
		"line2",
		"   ",
		"line3",
	}

	result := RemoveEmptyLines(input)

	if len(result) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result))
	}

	expected := []string{"line1", "line2", "line3"}
	for i, want := range expected {
		if result[i] != want {
			t.Errorf("Line %d = %q, want %q", i, result[i], want)
		}
	}
}
