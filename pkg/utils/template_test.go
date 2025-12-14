package utils

import (
	"strings"
	"testing"
	"text/template"
)

func TestRenderTemplate(t *testing.T) {
	t.Run("SimpleTemplate", func(t *testing.T) {
		tmpl := "Hello, {{.Name}}!"
		data := map[string]string{"Name": "World"}

		result, err := RenderTemplate(tmpl, data)
		if err != nil {
			t.Fatalf("RenderTemplate failed: %v", err)
		}

		expected := "Hello, World!"
		if result != expected {
			t.Errorf("Result = %q, want %q", result, expected)
		}
	})

	t.Run("ComplexTemplate", func(t *testing.T) {
		tmpl := `
Host: {{.Host}}
Port: {{.Port}}
User: {{.User}}
`
		data := map[string]interface{}{
			"Host": "localhost",
			"Port": 8080,
			"User": "admin",
		}

		result, err := RenderTemplate(tmpl, data)
		if err != nil {
			t.Fatalf("RenderTemplate failed: %v", err)
		}

		if !strings.Contains(result, "Host: localhost") {
			t.Error("Result should contain 'Host: localhost'")
		}
		if !strings.Contains(result, "Port: 8080") {
			t.Error("Result should contain 'Port: 8080'")
		}
	})

	t.Run("InvalidTemplate", func(t *testing.T) {
		tmpl := "{{.Invalid"
		data := map[string]string{}

		_, err := RenderTemplate(tmpl, data)
		if err == nil {
			t.Error("Expected error for invalid template")
		}
	})
}

func TestRenderTemplateWithFuncs(t *testing.T) {
	tmpl := "Hello, {{upper .Name}}!"
	data := map[string]string{"Name": "world"}
	funcs := template.FuncMap{
		"upper": strings.ToUpper,
	}

	result, err := RenderTemplateWithFuncs(tmpl, data, funcs)
	if err != nil {
		t.Fatalf("RenderTemplateWithFuncs failed: %v", err)
	}

	expected := "Hello, WORLD!"
	if result != expected {
		t.Errorf("Result = %q, want %q", result, expected)
	}
}

func TestDefaultTemplateFuncs(t *testing.T) {
	funcs := DefaultTemplateFuncs()

	tests := []struct {
		name     string
		template string
		data     interface{}
		want     string
	}{
		{
			name:     "upper",
			template: "{{upper .}}",
			data:     "hello",
			want:     "HELLO",
		},
		{
			name:     "lower",
			template: "{{lower .}}",
			data:     "HELLO",
			want:     "hello",
		},
		{
			name:     "trim",
			template: "{{trim .}}",
			data:     "  hello  ",
			want:     "hello",
		},
		{
			name:     "default",
			template: "{{default \"fallback\" .}}",
			data:     "",
			want:     "fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderTemplateWithFuncs(tt.template, tt.data, funcs)
			if err != nil {
				t.Fatalf("RenderTemplateWithFuncs failed: %v", err)
			}
			if result != tt.want {
				t.Errorf("Result = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestSimpleReplace(t *testing.T) {
	tests := []struct {
		name string
		text string
		vars map[string]string
		want string
	}{
		{
			name: "DollarBrace",
			text: "Hello, ${NAME}!",
			vars: map[string]string{"NAME": "World"},
			want: "Hello, World!",
		},
		{
			name: "DoubleBrace",
			text: "Host: {{HOST}}, Port: {{PORT}}",
			vars: map[string]string{"HOST": "localhost", "PORT": "8080"},
			want: "Host: localhost, Port: 8080",
		},
		{
			name: "Mixed",
			text: "${VAR1} and {{VAR2}}",
			vars: map[string]string{"VAR1": "first", "VAR2": "second"},
			want: "first and second",
		},
		{
			name: "NoVariables",
			text: "Plain text",
			vars: map[string]string{},
			want: "Plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SimpleReplace(tt.text, tt.vars)
			if result != tt.want {
				t.Errorf("SimpleReplace() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestRenderSimpleTemplate(t *testing.T) {
	tmpl := "Server: ${HOST}:${PORT}"
	vars := map[string]string{
		"HOST": "192.168.1.1",
		"PORT": "8080",
	}

	result := RenderSimpleTemplate(tmpl, vars)
	expected := "Server: 192.168.1.1:8080"

	if result != expected {
		t.Errorf("Result = %q, want %q", result, expected)
	}
}

func TestTemplateToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := tmpDir + "/output.txt"

	tmpl := "Hello, {{.Name}}!"
	data := map[string]string{"Name": "World"}

	err := TemplateToFile(tmpl, data, outputPath)
	if err != nil {
		t.Fatalf("TemplateToFile failed: %v", err)
	}

	// Verify file content
	content, err := ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "Hello, World!"
	if string(content) != expected {
		t.Errorf("File content = %q, want %q", content, expected)
	}
}

func TestRenderTemplateFile(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := tmpDir + "/template.txt"

	// Create template file
	tmplContent := "Config: {{.Key}}={{.Value}}"
	if err := WriteFile(templatePath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Render template
	data := map[string]string{"Key": "setting", "Value": "enabled"}
	result, err := RenderTemplateFile(templatePath, data)
	if err != nil {
		t.Fatalf("RenderTemplateFile failed: %v", err)
	}

	expected := "Config: setting=enabled"
	if result != expected {
		t.Errorf("Result = %q, want %q", result, expected)
	}
}

func TestRenderSimpleTemplateFile(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := tmpDir + "/simple_template.txt"

	// Create template file
	tmplContent := "Server: ${HOST}:${PORT}"
	if err := WriteFile(templatePath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Render template
	vars := map[string]string{"HOST": "localhost", "PORT": "3000"}
	result, err := RenderSimpleTemplateFile(templatePath, vars)
	if err != nil {
		t.Fatalf("RenderSimpleTemplateFile failed: %v", err)
	}

	expected := "Server: localhost:3000"
	if result != expected {
		t.Errorf("Result = %q, want %q", result, expected)
	}
}
