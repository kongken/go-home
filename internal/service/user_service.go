package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kongken/go-home/internal/cache"
	"github.com/kongken/go-home/internal/config"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidToken     = errors.New("invalid token")
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, username, password, email string) (*model.User, error)
	Login(ctx context.Context, account, password string) (*AuthResult, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error)
}

// AuthResult 认证结果
type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         *model.User
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	userCache *cache.UserCache
	jwtSecret []byte
	accessExpiry time.Duration
	refreshExpiry time.Duration
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, redisCache *cache.RedisCache) UserService {
	cfg := config.Get().JWT
	return &userService{
		userRepo:      userRepo,
		userCache:     redisCache.UserCache(),
		jwtSecret:     []byte(cfg.Secret),
		accessExpiry:  time.Duration(cfg.AccessExpiry) * time.Second,
		refreshExpiry: time.Duration(cfg.RefreshExpiry) * time.Second,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, username, password, email string) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(ctx, username); err == nil {
		return nil, ErrUserExists
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(ctx, email); err == nil {
		return nil, ErrUserExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Status:   model.UserStatusNormal,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, account, password string) (*AuthResult, error) {
	// 尝试通过用户名或邮箱查找用户
	user, err := s.userRepo.GetByUsername(ctx, account)
	if err != nil {
		user, err = s.userRepo.GetByEmail(ctx, account)
		if err != nil {
			return nil, ErrUserNotFound
		}
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// 生成 token
	return s.generateTokens(user)
}

// GetUser 获取用户信息 (带缓存)
func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	// 先查缓存
	if s.userCache != nil {
		if user, err := s.userCache.Get(ctx, id); err == nil && user != nil {
			return user, nil
		}
	}
	
	// 查数据库
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// 写入缓存
	if s.userCache != nil {
		s.userCache.Set(ctx, user, time.Hour)
	}
	
	return user, nil
}

// UpdateUser 更新用户信息 (更新后删除缓存)
func (s *userService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 更新字段
	if nickname, ok := updates["nickname"].(string); ok {
		user.Nickname = nickname
	}
	if avatar, ok := updates["avatar"].(string); ok {
		user.Avatar = avatar
	}
	if bio, ok := updates["bio"].(string); ok {
		user.Bio = bio
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	
	// 删除缓存
	if s.userCache != nil {
		s.userCache.Delete(ctx, id)
	}

	return user, nil
}

// RefreshToken 刷新 token
func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	// 解析 refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokens(user)
}

// generateTokens 生成 token
func (s *userService) generateTokens(user *model.User) (*AuthResult, error) {
	// 生成 access token
	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(s.accessExpiry).Unix(),
		"type":     "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// 生成 refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.refreshExpiry).Unix(),
		"type":    "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
		User:         user,
	}, nil
}
