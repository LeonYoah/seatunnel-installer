package utils

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// RenderTemplate renders a template string with the given data
func RenderTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderTemplateWithFuncs renders a template with custom functions
func RenderTemplateWithFuncs(templateStr string, data interface{}, funcMap template.FuncMap) (string, error) {
	tmpl, err := template.New("template").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderTemplateFile renders a template file with the given data
func RenderTemplateFile(templatePath string, data interface{}) (string, error) {
	content, err := ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return RenderTemplate(string(content), data)
}

// RenderTemplateFileWithFuncs renders a template file with custom functions
func RenderTemplateFileWithFuncs(templatePath string, data interface{}, funcMap template.FuncMap) (string, error) {
	content, err := ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return RenderTemplateWithFuncs(string(content), data, funcMap)
}

// DefaultTemplateFuncs returns a map of commonly used template functions
func DefaultTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"upper":     strings.ToUpper,
		"lower":     strings.ToLower,
		"title":     strings.Title,
		"trim":      strings.TrimSpace,
		"replace":   strings.ReplaceAll,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"join":      strings.Join,
		"split":     strings.Split,
		"default": func(defaultValue, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
	}
}

// SimpleReplace performs simple variable replacement in a string
// Variables are in the format ${VAR_NAME} or {{VAR_NAME}}
func SimpleReplace(text string, vars map[string]string) string {
	result := text
	for key, value := range vars {
		// Replace ${VAR_NAME} format
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", key), value)
		// Replace {{VAR_NAME}} format
		result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", key), value)
	}
	return result
}

// RenderSimpleTemplate renders a simple template with variable replacement
func RenderSimpleTemplate(templateStr string, vars map[string]string) string {
	return SimpleReplace(templateStr, vars)
}

// RenderSimpleTemplateFile renders a simple template file with variable replacement
func RenderSimpleTemplateFile(templatePath string, vars map[string]string) (string, error) {
	content, err := ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return RenderSimpleTemplate(string(content), vars), nil
}

// TemplateToFile renders a template and writes it to a file
func TemplateToFile(templateStr string, data interface{}, outputPath string) error {
	rendered, err := RenderTemplate(templateStr, data)
	if err != nil {
		return err
	}

	return WriteFile(outputPath, []byte(rendered), 0644)
}

// TemplateFileToFile renders a template file and writes it to another file
func TemplateFileToFile(templatePath string, data interface{}, outputPath string) error {
	rendered, err := RenderTemplateFile(templatePath, data)
	if err != nil {
		return err
	}

	return WriteFile(outputPath, []byte(rendered), 0644)
}
