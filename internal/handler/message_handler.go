package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	messageService service.MessageService
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ToUserID    string                 `json:"to_user_id" binding:"required"`
	Content     string                 `json:"content" binding:"required"`
	Attachments []model.MediaAttachment `json:"attachments"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID          string                  `json:"id"`
	FromUserID  string                  `json:"from_user_id"`
	ToUserID    string                  `json:"to_user_id"`
	Content     string                  `json:"content"`
	Attachments []model.MediaAttachment `json:"attachments,omitempty"`
	IsRead      bool                    `json:"is_read"`
	CreatedAt   string                  `json:"created_at"`
}

// Send 发送消息
func (h *MessageHandler) Send(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.messageService.Send(c.Request.Context(), userID, req.ToUserID, req.Content, req.Attachments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": h.toMessageResponse(message)})
}

// ListMessages 获取消息列表
func (h *MessageHandler) ListMessages(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	otherUserID := c.Param("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := int64((page - 1) * pageSize)

	messages, total, err := h.messageService.ListByConversation(c.Request.Context(), userID, otherUserID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*MessageResponse
	for _, m := range messages {
		responses = append(responses, h.toMessageResponse(m))
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// MarkAsRead 标记已读
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	otherUserID := c.Param("user_id")

	if err := h.messageService.MarkAsRead(c.Request.Context(), userID, otherUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ConversationResponse 会话响应
type ConversationResponse struct {
	ID           string           `json:"id"`
	UserID       string           `json:"user_id"`
	OtherUserID  string           `json:"other_user_id"`
	UnreadCount  int32            `json:"unread_count"`
	UpdatedAt    string           `json:"updated_at"`
}

// ListConversations 获取会话列表
func (h *MessageHandler) ListConversations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := int64((page - 1) * pageSize)

	conversations, total, err := h.messageService.ListConversations(c.Request.Context(), userID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*ConversationResponse
	for _, conv := range conversations {
		responses = append(responses, &ConversationResponse{
			ID:          conv.ID,
			UserID:      conv.UserID,
			OtherUserID: conv.OtherUserID,
			UnreadCount: conv.UnreadCount,
			UpdatedAt:   conv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// UnreadCount 获取未读消息数
func (h *MessageHandler) UnreadCount(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, err := h.messageService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// toMessageResponse 转换为消息响应
func (h *MessageHandler) toMessageResponse(m *model.Message) *MessageResponse {
	return &MessageResponse{
		ID:          m.ID,
		FromUserID:  m.FromUserID,
		ToUserID:    m.ToUserID,
		Content:     m.Content,
		Attachments: m.Attachments,
		IsRead:      m.IsRead,
		CreatedAt:   m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}