package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/api"
	"go.uber.org/zap"
)

// ErrorHandler 错误处理中间件
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic信息
				logger.Error("HTTP请求发生panic",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// 返回500错误
				api.InternalError(c, "服务器内部错误", fmt.Errorf("%v", err))
				c.Abort()
			}
		}()

		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 记录错误日志
			logger.Error("HTTP请求处理错误",
				zap.Error(err.Err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			// 如果还没有响应，返回错误响应
			if !c.Writer.Written() {
				switch err.Type {
				case gin.ErrorTypeBind:
					api.ValidationError(c, "请求参数错误", err.Error())
				case gin.ErrorTypePublic:
					api.Error(c, http.StatusBadRequest, "请求错误", err.Err)
				default:
					api.InternalError(c, "服务器内部错误", err.Err)
				}
			}
		}
	}
}
