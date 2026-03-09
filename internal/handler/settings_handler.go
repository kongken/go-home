package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// SettingsHandler 设置处理器
type SettingsHandler struct {
	settingsService service.SettingsService
}

// NewSettingsHandler 创建设置处理器
func NewSettingsHandler(settingsService service.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: settingsService}
}

// SettingsResponse 设置响应
type SettingsResponse struct {
	Privacy      PrivacySettingsResponse      `json:"privacy"`
	Notification NotificationSettingsResponse `json:"notification"`
	Blacklist    []string                     `json:"blacklist"`
}

// PrivacySettingsResponse 隐私设置响应
type PrivacySettingsResponse struct {
	DefaultBlogPrivacy  int32 `json:"default_blog_privacy"`
	DefaultAlbumPrivacy int32 `json:"default_album_privacy"`
	DefaultSharePrivacy int32 `json:"default_share_privacy"`
	DefaultFeedPrivacy  int32 `json:"default_feed_privacy"`
}

// NotificationSettingsResponse 通知设置响应
type NotificationSettingsResponse struct {
	NotifyFriendRequest bool `json:"notify_friend_request"`
	NotifyComment       bool `json:"notify_comment"`
	NotifyLike          bool `json:"notify_like"`
	NotifyMention       bool `json:"notify_mention"`
	NotifyGroup         bool `json:"notify_group"`
	NotifySystem        bool `json:"notify_system"`
}

// Get 获取设置
func (h *SettingsHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	settings, err := h.settingsService.Get(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": h.toSettingsResponse(settings)})
}

// UpdatePrivacyRequest 更新隐私设置请求
type UpdatePrivacyRequest struct {
	Privacy model.PrivacySettings `json:"privacy" binding:"required"`
}

// UpdatePrivacy 更新隐私设置
func (h *SettingsHandler) UpdatePrivacy(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdatePrivacyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsService.UpdatePrivacy(c.Request.Context(), userID, req.Privacy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// UpdateNotificationRequest 更新通知设置请求
type UpdateNotificationRequest struct {
	Notification model.NotificationSettings `json:"notification" binding:"required"`
}

// UpdateNotification 更新通知设置
func (h *SettingsHandler) UpdateNotification(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsService.UpdateNotification(c.Request.Context(), userID, req.Notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// AddToBlacklistRequest 添加到黑名单请求
type AddToBlacklistRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// AddToBlacklist 添加到黑名单
func (h *SettingsHandler) AddToBlacklist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AddToBlacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsService.AddToBlacklist(c.Request.Context(), userID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// RemoveFromBlacklist 从黑名单移除
func (h *SettingsHandler) RemoveFromBlacklist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	blockedUserID := c.Param("user_id")

	if err := h.settingsService.RemoveFromBlacklist(c.Request.Context(), userID, blockedUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// toSettingsResponse 转换为设置响应
func (h *SettingsHandler) toSettingsResponse(s *model.Settings) *SettingsResponse {
	return &SettingsResponse{
		Privacy: PrivacySettingsResponse{
			DefaultBlogPrivacy:  int32(s.Privacy.DefaultBlogPrivacy),
			DefaultAlbumPrivacy: int32(s.Privacy.DefaultAlbumPrivacy),
			DefaultSharePrivacy: int32(s.Privacy.DefaultSharePrivacy),
			DefaultFeedPrivacy:  int32(s.Privacy.DefaultFeedPrivacy),
		},
		Notification: NotificationSettingsResponse{
			NotifyFriendRequest: s.Notification.NotifyFriendRequest,
			NotifyComment:       s.Notification.NotifyComment,
			NotifyLike:          s.Notification.NotifyLike,
			NotifyMention:       s.Notification.NotifyMention,
			NotifyGroup:         s.Notification.NotifyGroup,
			NotifySystem:        s.Notification.NotifySystem,
		},
		Blacklist: s.Blacklist,
	}
}