package handler

import (
	"net/http"

	"butterfly.orx.me/core/log"
	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	User         *UserInfo   `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio,omitempty"`
	Email    string `json:"email,omitempty"`
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("invalid register request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		if err == service.ErrUserExists {
			logger.Warn("user already exists", "username", req.Username)
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		logger.Error("failed to register user", "error", err, "username", req.Username)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("user registered", "user_id", user.ID, "username", user.Username)
	c.JSON(http.StatusCreated, gin.H{
		"user": &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Email:    user.Email,
		},
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("invalid login request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.userService.Login(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrInvalidPassword {
			logger.Warn("login failed - invalid credentials", "account", req.Account)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		logger.Error("login failed", "error", err, "account", req.Account)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("user logged in", "user_id", result.User.ID, "account", req.Account)
	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: &UserInfo{
			ID:       result.User.ID,
			Username: result.User.Username,
			Nickname: result.User.Nickname,
			Avatar:   result.User.Avatar,
			Email:    result.User.Email,
		},
	})
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Bio:      user.Bio,
			Email:    user.Email,
		},
	})
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, updates)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Email:    user.Email,
		},
	})
}

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == service.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: &UserInfo{
			ID:       result.User.ID,
			Username: result.User.Username,
			Nickname: result.User.Nickname,
			Avatar:   result.User.Avatar,
			Email:    result.User.Email,
		},
	})
}
