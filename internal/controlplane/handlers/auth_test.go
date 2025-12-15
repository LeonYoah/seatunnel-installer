package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService 模拟认证服务
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, req *auth.RefreshRequest) (*auth.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.LoginResponse), args.Error(1)
}

func (m *MockAuthService) GetCurrentUser(ctx context.Context, userID string) (*auth.UserInfo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.UserInfo), args.Error(1)
}

func TestAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("登录成功", func(t *testing.T) {
		// 创建模拟服务
		mockService := new(MockAuthService)
		handler := &AuthHandler{authService: mockService}

		// 设置期望
		loginReq := &auth.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		loginResp := &auth.LoginResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresAt:    time.Now().Add(time.Hour),
			User: auth.UserInfo{
				ID:       "user-id",
				Username: "testuser",
				Email:    "test@example.com",
			},
		}
		mockService.On("Login", mock.Anything, loginReq).Return(loginResp, nil)

		// 创建请求
		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 执行处理器
		handler.Login(c)

		// 验证结果
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusOK), response["code"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("登录失败-无效凭证", func(t *testing.T) {
		// 创建模拟服务
		mockService := new(MockAuthService)
		handler := &AuthHandler{authService: mockService}

		// 设置期望
		loginReq := &auth.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		mockService.On("Login", mock.Anything, loginReq).Return(nil, auth.ErrInvalidCredentials)

		// 创建请求
		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 执行处理器
		handler.Login(c)

		// 验证结果
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusUnauthorized), response["code"])

		mockService.AssertExpectations(t)
	})

	t.Run("刷新令牌成功", func(t *testing.T) {
		// 创建模拟服务
		mockService := new(MockAuthService)
		handler := &AuthHandler{authService: mockService}

		// 设置期望
		refreshReq := &auth.RefreshRequest{
			RefreshToken: "valid-refresh-token",
		}
		loginResp := &auth.LoginResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "valid-refresh-token",
			ExpiresAt:    time.Now().Add(time.Hour),
			User: auth.UserInfo{
				ID:       "user-id",
				Username: "testuser",
			},
		}
		mockService.On("RefreshToken", mock.Anything, refreshReq).Return(loginResp, nil)

		// 创建请求
		reqBody, _ := json.Marshal(refreshReq)
		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 执行处理器
		handler.RefreshToken(c)

		// 验证结果
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusOK), response["code"])

		mockService.AssertExpectations(t)
	})

	t.Run("获取当前用户信息", func(t *testing.T) {
		// 创建模拟服务
		mockService := new(MockAuthService)
		handler := &AuthHandler{authService: mockService}

		// 设置期望
		userInfo := &auth.UserInfo{
			ID:       "user-id",
			Username: "testuser",
			Email:    "test@example.com",
		}
		mockService.On("GetCurrentUser", mock.Anything, "user-id").Return(userInfo, nil)

		// 创建请求
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		w := httptest.NewRecorder()

		// 创建Gin上下文并设置用户ID
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", "user-id")

		// 执行处理器
		handler.GetCurrentUser(c)

		// 验证结果
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusOK), response["code"])

		mockService.AssertExpectations(t)
	})

	t.Run("登出", func(t *testing.T) {
		handler := &AuthHandler{}

		// 创建请求
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 执行处理器
		handler.Logout(c)

		// 验证结果
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusOK), response["code"])
	})
}
