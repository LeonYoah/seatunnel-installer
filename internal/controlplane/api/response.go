package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一的API响应格式
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
	Error   *ErrorInfo  `json:"error"`   // 错误详情
}

// ErrorInfo 错误详情
type ErrorInfo struct {
	Type    string `json:"type"`    // 错误类型
	Details string `json:"details"` // 错误详情
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Error:   nil,
	})
}

// 成功响应（无数据）
func SuccessNoData(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    nil,
		Error:   nil,
	})
}

// 错误响应
func Error(c *gin.Context, code int, message string, err error) {
	var errorInfo *ErrorInfo
	if err != nil {
		errorInfo = &ErrorInfo{
			Type:    "internal_error",
			Details: err.Error(),
		}
	}

	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Error:   errorInfo,
	})
}

// 参数验证错误
func ValidationError(c *gin.Context, message string, details string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Data:    nil,
		Error: &ErrorInfo{
			Type:    "validation_error",
			Details: details,
		},
	})
}

// 认证失败
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: message,
		Data:    nil,
		Error: &ErrorInfo{
			Type:    "authentication_error",
			Details: "请检查认证信息",
		},
	})
}

// 权限不足
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: message,
		Data:    nil,
		Error: &ErrorInfo{
			Type:    "authorization_error",
			Details: "权限不足",
		},
	})
}

// 资源不存在
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
		Data:    nil,
		Error: &ErrorInfo{
			Type:    "resource_not_found",
			Details: "请求的资源不存在",
		},
	})
}

// 内部服务器错误
func InternalError(c *gin.Context, message string, err error) {
	Error(c, http.StatusInternalServerError, message, err)
}
