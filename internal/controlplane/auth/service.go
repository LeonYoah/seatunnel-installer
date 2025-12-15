package auth

import (
	"context"
	"errors"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserInactive       = errors.New("用户已被禁用")
)

// AuthService 认证服务
type AuthService struct {
	repoManager repository.RepositoryManager
	jwtService  *JWTService
}

// NewAuthService 创建认证服务
func NewAuthService(repoManager repository.RepositoryManager, jwtService *JWTService) *AuthService {
	return &AuthService{
		repoManager: repoManager,
		jwtService:  jwtService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	TenantID    string     `json:"tenant_id"`
	WorkspaceID string     `json:"workspace_id"`
	Roles       []string   `json:"roles"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 根据用户名查找用户
	user, err := s.repoManager.User().GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	// 获取用户角色信息
	userWithRoles, err := s.repoManager.User().GetWithRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// 生成访问令牌
	accessToken, err := s.jwtService.GenerateAccessToken(userWithRoles)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	err = s.repoManager.User().UpdateLastLogin(ctx, user.ID)
	if err != nil {
		// 记录错误但不影响登录流程
		// TODO: 添加日志记录
	}

	// 构建响应
	roles := make([]string, len(userWithRoles.Roles))
	for i, role := range userWithRoles.Roles {
		roles[i] = role.Name
	}

	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.jwtService.accessTokenTTL),
		User: UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			TenantID:    user.TenantID,
			WorkspaceID: user.WorkspaceID,
			Roles:       roles,
			LastLoginAt: user.LastLoginAt,
		},
	}

	return response, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshRequest) (*LoginResponse, error) {
	// 验证刷新令牌并获取用户ID
	userID, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	user, err := s.repoManager.User().GetWithRoles(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	// 生成新的访问令牌
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// 构建响应
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // 保持原刷新令牌
		ExpiresAt:    time.Now().Add(s.jwtService.accessTokenTTL),
		User: UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			TenantID:    user.TenantID,
			WorkspaceID: user.WorkspaceID,
			Roles:       roles,
			LastLoginAt: user.LastLoginAt,
		},
	}

	return response, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*UserInfo, error) {
	user, err := s.repoManager.User().GetWithRoles(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	userInfo := &UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		TenantID:    user.TenantID,
		WorkspaceID: user.WorkspaceID,
		Roles:       roles,
		LastLoginAt: user.LastLoginAt,
	}

	return userInfo, nil
}
