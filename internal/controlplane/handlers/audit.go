package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/api"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
)

// AuditHandler 审计日志处理器
type AuditHandler struct {
	repoManager repository.RepositoryManager
}

// NewAuditHandler 创建审计日志处理器
func NewAuditHandler(repoManager repository.RepositoryManager) *AuditHandler {
	return &AuditHandler{
		repoManager: repoManager,
	}
}

// GetAuditLogs 获取审计日志列表
// @Summary 获取审计日志列表
// @Description 获取当前租户的审计日志列表，支持分页和筛选
// @Tags 审计日志
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页大小" default(20)
// @Param user_id query string false "用户ID筛选"
// @Param resource query string false "资源类型筛选"
// @Param resource_id query string false "资源ID筛选"
// @Param action query string false "操作类型筛选"
// @Param result query string false "结果筛选(success/failure)"
// @Success 200 {object} api.Response{data=AuditLogListResponse} "获取成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "未认证"
// @Failure 403 {object} api.Response "权限不足"
// @Router /api/v1/audit/logs [get]
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	// 获取租户ID
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		api.Unauthorized(c, "缺少租户信息")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size

	// 获取筛选参数
	userID := c.Query("user_id")
	resource := c.Query("resource")
	resourceID := c.Query("resource_id")

	var logs interface{}
	var total int64
	var err error

	// 根据筛选条件选择查询方法
	if userID != "" {
		// 按用户ID查询
		logs, total, err = h.repoManager.AuditLog().GetByUserID(c.Request.Context(), userID, offset, size)
	} else if resource != "" {
		// 按资源查询
		logs, total, err = h.repoManager.AuditLog().GetByResource(c.Request.Context(), tenantID.(string), resource, resourceID, offset, size)
	} else {
		// 按租户ID查询所有
		logs, total, err = h.repoManager.AuditLog().GetByTenantID(c.Request.Context(), tenantID.(string), offset, size)
	}

	if err != nil {
		api.InternalError(c, "获取审计日志失败", err)
		return
	}

	response := AuditLogListResponse{
		Logs:  logs,
		Total: total,
		Page:  page,
		Size:  size,
	}

	api.Success(c, response)
}

// GetAuditLogsByUser 获取指定用户的审计日志
// @Summary 获取指定用户的审计日志
// @Description 获取指定用户的审计日志列表
// @Tags 审计日志
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页大小" default(20)
// @Success 200 {object} api.Response{data=AuditLogListResponse} "获取成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "未认证"
// @Failure 403 {object} api.Response "权限不足"
// @Router /api/v1/audit/users/{user_id}/logs [get]
func (h *AuditHandler) GetAuditLogsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		api.ValidationError(c, "用户ID不能为空", "user_id参数是必需的")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size

	logs, total, err := h.repoManager.AuditLog().GetByUserID(c.Request.Context(), userID, offset, size)
	if err != nil {
		api.InternalError(c, "获取用户审计日志失败", err)
		return
	}

	response := AuditLogListResponse{
		Logs:  logs,
		Total: total,
		Page:  page,
		Size:  size,
	}

	api.Success(c, response)
}

// GetAuditLogsByResource 获取指定资源的审计日志
// @Summary 获取指定资源的审计日志
// @Description 获取指定资源的审计日志列表
// @Tags 审计日志
// @Produce json
// @Security BearerAuth
// @Param resource path string true "资源类型"
// @Param resource_id query string false "资源ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页大小" default(20)
// @Success 200 {object} api.Response{data=AuditLogListResponse} "获取成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "未认证"
// @Failure 403 {object} api.Response "权限不足"
// @Router /api/v1/audit/resources/{resource}/logs [get]
func (h *AuditHandler) GetAuditLogsByResource(c *gin.Context) {
	resource := c.Param("resource")
	if resource == "" {
		api.ValidationError(c, "资源类型不能为空", "resource参数是必需的")
		return
	}

	// 获取租户ID
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		api.Unauthorized(c, "缺少租户信息")
		return
	}

	resourceID := c.Query("resource_id")

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size

	logs, total, err := h.repoManager.AuditLog().GetByResource(c.Request.Context(), tenantID.(string), resource, resourceID, offset, size)
	if err != nil {
		api.InternalError(c, "获取资源审计日志失败", err)
		return
	}

	response := AuditLogListResponse{
		Logs:  logs,
		Total: total,
		Page:  page,
		Size:  size,
	}

	api.Success(c, response)
}

// AuditLogListResponse 审计日志列表响应
type AuditLogListResponse struct {
	Logs  interface{} `json:"logs"`  // 审计日志列表
	Total int64       `json:"total"` // 总数
	Page  int         `json:"page"`  // 当前页码
	Size  int         `json:"size"`  // 每页大小
}
