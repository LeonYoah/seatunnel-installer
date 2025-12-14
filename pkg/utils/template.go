package utils

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// RenderTemplate 使用给定的数据渲染模板字符串
func RenderTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}

// RenderTemplateWithFuncs 使用自定义函数渲染模板
func RenderTemplateWithFuncs(templateStr string, data interface{}, funcMap template.FuncMap) (string, error) {
	tmpl, err := template.New("template").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}

// RenderTemplateFile 使用给定的数据渲染模板文件
func RenderTemplateFile(templatePath string, data interface{}) (string, error) {
	content, err := ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %w", err)
	}

	return RenderTemplate(string(content), data)
}

// RenderTemplateFileWithFuncs 使用自定义函数渲染模板文件
func RenderTemplateFileWithFuncs(templatePath string, data interface{}, funcMap template.FuncMap) (string, error) {
	content, err := ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %w", err)
	}

	return RenderTemplateWithFuncs(string(content), data, funcMap)
}

// DefaultTemplateFuncs 返回常用模板函数的映射
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

// SimpleReplace 在字符串中执行简单的变量替换
// 变量格式为 ${VAR_NAME} 或 {{VAR_NAME}}
func SimpleReplace(text string, vars map[string]string) string {
	result := text
	for key, value := range vars {
		// 替换 ${VAR_NAME} 格式
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", key), value)
		// 替换 {{VAR_NAME}} 格式
		result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", key), value)
	}
	return result
}

// RenderSimpleTemplate 使用变量替换渲染简单模板
func RenderSimpleTemplate(templateStr string, vars map[string]string) string {
	return SimpleReplace(templateStr, vars)
}

// RenderSimpleTemplateFile 使用变量替换渲染简单模板文件
func RenderSimpleTemplateFile(templatePath string, vars map[string]string) (string, error) {
	content, err := ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %w", err)
	}

	return RenderSimpleTemplate(string(content), vars), nil
}

// TemplateToFile 渲染模板并写入文件
func TemplateToFile(templateStr string, data interface{}, outputPath string) error {
	rendered, err := RenderTemplate(templateStr, data)
	if err != nil {
		return err
	}

	return WriteFile(outputPath, []byte(rendered), 0644)
}

// TemplateFileToFile 渲染模板文件并写入另一个文件
func TemplateFileToFile(templatePath string, data interface{}, outputPath string) error {
	rendered, err := RenderTemplateFile(templatePath, data)
	if err != nil {
		return err
	}

	return WriteFile(outputPath, []byte(rendered), 0644)
}
