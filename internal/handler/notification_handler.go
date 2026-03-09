package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	notifService service.NotificationService
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(notifService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Type       int32  `json:"type"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ActorID    string `json:"actor_id"`
	TargetID   string `json:"target_id"`
	TargetType string `json:"target_type"`
	IsRead     bool   `json:"is_read"`
	CreatedAt  string `json:"created_at"`
}

// List 获取通知列表
func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filter := c.Query("filter") // all, unread, read
	unreadOnly := filter == "unread"

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := int64((page - 1) * pageSize)

	notifications, total, err := h.notifService.List(c.Request.Context(), userID, unreadOnly, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*NotificationResponse
	for _, n := range notifications {
		responses = append(responses, h.toNotificationResponse(n))
	}

	// 获取未读数
	unreadCount, _ := h.notifService.GetUnreadCount(c.Request.Context(), userID)

	c.JSON(http.StatusOK, gin.H{
		"notifications": responses,
		"unread_count":  unreadCount,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// MarkAsRead 标记通知已读
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.notifService.MarkAsRead(c.Request.Context(), id); err != nil {
		if err == service.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// MarkAllAsRead 标记所有通知已读
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.notifService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// UnreadCount 获取未读通知数
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, err := h.notifService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// Delete 删除通知
func (h *NotificationHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.notifService.Delete(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// toNotificationResponse 转换为通知响应
func (h *NotificationHandler) toNotificationResponse(n *model.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:         n.ID,
		UserID:     n.UserID,
		Type:       int32(n.Type),
		Title:      n.Title,
		Content:    n.Content,
		ActorID:    n.ActorID,
		TargetID:   n.TargetID,
		TargetType: n.TargetType,
		IsRead:     n.IsRead,
		CreatedAt:  n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}