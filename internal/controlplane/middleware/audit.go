package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"go.uber.org/zap"
)

// AuditMiddleware 审计日志中间件
func AuditMiddleware(repoManager repository.RepositoryManager, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录写操作（POST、PUT、PATCH、DELETE）
		if !isWriteOperation(c.Request.Method) {
			c.Next()
			return
		}

		// 获取用户信息
		userID, _ := c.Get("user_id")
		tenantID, _ := c.Get("tenant_id")

		// 如果没有用户信息，跳过审计（可能是公开接口）
		if userID == nil || tenantID == nil {
			c.Next()
			return
		}

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，以便后续处理器可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录开始时间
		startTime := time.Now()

		// 创建响应写入器包装器来捕获响应
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// 处理请求
		c.Next()

		// 创建审计日志
		auditLog := &models.AuditLog{
			TenantID:  tenantID.(string),
			UserID:    userID.(string),
			Action:    getActionFromMethod(c.Request.Method),
			Resource:  getResourceFromPath(c.Request.URL.Path),
			Details:   buildAuditDetails(c, requestBody, responseWriter.body.Bytes()),
			CreatedAt: startTime,
		}

		// 设置资源ID（如果可以从路径中提取）
		if resourceID := extractResourceID(c.Request.URL.Path); resourceID != "" {
			auditLog.ResourceID = resourceID
		}

		// 根据响应状态码设置结果
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			auditLog.Result = "success"
		} else {
			auditLog.Result = "failure"
			// 尝试从响应中提取错误信息
			if errorMsg := extractErrorMessage(responseWriter.body.Bytes()); errorMsg != "" {
				auditLog.ErrorMsg = errorMsg
			}
		}

		// 异步保存审计日志，避免影响主请求性能
		go func() {
			if err := repoManager.AuditLog().Create(c.Request.Context(), auditLog); err != nil {
				logger.Error("保存审计日志失败",
					zap.Error(err),
					zap.String("user_id", auditLog.UserID),
					zap.String("action", auditLog.Action),
					zap.String("resource", auditLog.Resource),
				)
			}
		}()
	}
}

// responseWriter 响应写入器包装器，用于捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// isWriteOperation 判断是否为写操作
func isWriteOperation(method string) bool {
	writeMethods := []string{"POST", "PUT", "PATCH", "DELETE"}
	for _, writeMethod := range writeMethods {
		if method == writeMethod {
			return true
		}
	}
	return false
}

// getActionFromMethod 根据HTTP方法获取操作类型
func getActionFromMethod(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

// getResourceFromPath 从路径中提取资源类型
func getResourceFromPath(path string) string {
	// 移除API版本前缀
	path = strings.TrimPrefix(path, "/api/v1/")

	// 移除开头的斜杠
	path = strings.TrimPrefix(path, "/")

	// 分割路径
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "unknown"
}

// extractResourceID 从路径中提取资源ID
func extractResourceID(path string) string {
	// 移除API版本前缀
	path = strings.TrimPrefix(path, "/api/v1/")

	// 分割路径，查找UUID格式的ID
	parts := strings.Split(path, "/")
	for _, part := range parts {
		// 简单的UUID格式检查（36个字符，包含连字符）
		if len(part) == 36 && strings.Count(part, "-") == 4 {
			return part
		}
	}
	return ""
}

// buildAuditDetails 构建审计详情
func buildAuditDetails(c *gin.Context, requestBody, responseBody []byte) string {
	details := map[string]interface{}{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.RawQuery,
		"user_agent": c.Request.UserAgent(),
		"client_ip":  c.ClientIP(),
		"status":     c.Writer.Status(),
	}

	// 添加请求体（如果不为空且不是敏感信息）
	if len(requestBody) > 0 && !containsSensitiveData(requestBody) {
		var requestJSON interface{}
		if err := json.Unmarshal(requestBody, &requestJSON); err == nil {
			details["request"] = requestJSON
		}
	}

	// 添加响应体（仅在失败时或特定操作时）
	if c.Writer.Status() >= 400 && len(responseBody) > 0 {
		var responseJSON interface{}
		if err := json.Unmarshal(responseBody, &responseJSON); err == nil {
			details["response"] = responseJSON
		}
	}

	// 转换为JSON字符串
	detailsJSON, _ := json.Marshal(details)
	return string(detailsJSON)
}

// containsSensitiveData 检查请求体是否包含敏感数据
func containsSensitiveData(body []byte) bool {
	bodyStr := strings.ToLower(string(body))
	sensitiveFields := []string{"password", "secret", "token", "key", "credential"}

	for _, field := range sensitiveFields {
		if strings.Contains(bodyStr, field) {
			return true
		}
	}
	return false
}

// extractErrorMessage 从响应体中提取错误信息
func extractErrorMessage(responseBody []byte) string {
	if len(responseBody) == 0 {
		return ""
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return ""
	}

	// 尝试提取错误信息
	if message, ok := response["message"].(string); ok {
		return message
	}

	if errorInfo, ok := response["error"].(map[string]interface{}); ok {
		if details, ok := errorInfo["details"].(string); ok {
			return details
		}
	}

	return ""
}
