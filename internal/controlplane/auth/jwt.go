package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
)

var (
	ErrInvalidToken = errors.New("无效的令牌")
	ErrExpiredToken = errors.New("令牌已过期")
)

// JWTClaims JWT声明
type JWTClaims struct {
	UserID      string   `json:"user_id"`
	TenantID    string   `json:"tenant_id"`
	WorkspaceID string   `json:"workspace_id"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTService JWT服务
type JWTService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTService 创建JWT服务
func NewJWTService(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTService {
	return &JWTService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// GenerateAccessToken 生成访问令牌
func (s *JWTService) GenerateAccessToken(user *models.User) (string, error) {
	now := time.Now()
	
	// 提取角色名称
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	claims := JWTClaims{
		UserID:      user.ID,
		TenantID:    user.TenantID,
		WorkspaceID: user.WorkspaceID,
		Username:    user.Username,
		Roles:       roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "seatunnel-enterprise-platform",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// GenerateRefreshToken 生成刷新令牌
func (s *JWTService) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "seatunnel-enterprise-platform",
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateAccessToken 验证访问令牌
func (s *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateRefreshToken 验证刷新令牌
func (s *JWTService) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", ErrInvalidToken
}

// RefreshAccessToken 刷新访问令牌
func (s *JWTService) RefreshAccessToken(refreshToken string, user *models.User) (string, error) {
	// 验证刷新令牌
	userID, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 检查用户ID是否匹配
	if userID != user.ID {
		return "", ErrInvalidToken
	}

	// 生成新的访问令牌
	return s.GenerateAccessToken(user)
}