package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/api"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/auth"
)

// AuthServiceInterface 认证服务接口
type AuthServiceInterface interface {
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	RefreshToken(ctx context.Context, req *auth.RefreshRequest) (*auth.LoginResponse, error)
	GetCurrentUser(ctx context.Context, userID string) (*auth.UserInfo, error)
}

// AuthHandler 认证处理器
type AuthHandler struct {
	authService AuthServiceInterface
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户通过用户名和密码登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录请求"
// @Success 200 {object} api.Response{data=auth.LoginResponse} "登录成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "用户名或密码错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.ValidationError(c, "请求参数错误", err.Error())
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			api.Unauthorized(c, "用户名或密码错误")
		case auth.ErrUserInactive:
			api.Unauthorized(c, "用户已被禁用")
		default:
			api.InternalError(c, "登录失败", err)
		}
		return
	}

	api.Success(c, response)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.RefreshRequest true "刷新令牌请求"
// @Success 200 {object} api.Response{data=auth.LoginResponse} "刷新成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "刷新令牌无效或已过期"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.ValidationError(c, "请求参数错误", err.Error())
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case auth.ErrExpiredToken:
			api.Unauthorized(c, "刷新令牌已过期")
		case auth.ErrInvalidToken:
			api.Unauthorized(c, "无效的刷新令牌")
		case auth.ErrUserNotFound:
			api.Unauthorized(c, "用户不存在")
		case auth.ErrUserInactive:
			api.Unauthorized(c, "用户已被禁用")
		default:
			api.InternalError(c, "刷新令牌失败", err)
		}
		return
	}

	api.Success(c, response)
}

// GetCurrentUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} api.Response{data=auth.UserInfo} "获取成功"
// @Failure 401 {object} api.Response "未认证"
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		api.Unauthorized(c, "未认证的用户")
		return
	}

	userInfo, err := h.authService.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		switch err {
		case auth.ErrUserNotFound:
			api.NotFound(c, "用户不存在")
		default:
			api.InternalError(c, "获取用户信息失败", err)
		}
		return
	}

	api.Success(c, userInfo)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出系统（客户端需要清除本地令牌）
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} api.Response "登出成功"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT是无状态的，服务端不需要做特殊处理
	// 客户端需要清除本地存储的令牌
	api.SuccessNoData(c)
}
