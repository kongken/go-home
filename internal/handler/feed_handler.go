package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// FeedHandler 动态处理器
type FeedHandler struct {
	feedService service.FeedService
}

// NewFeedHandler 创建动态处理器
func NewFeedHandler(feedService service.FeedService) *FeedHandler {
	return &FeedHandler{feedService: feedService}
}

// CreateFeedRequest 创建动态请求
type CreateFeedRequest struct {
	Type        model.FeedType         `json:"type" binding:"required"`
	Content     string                 `json:"content"`
	TargetID    string                 `json:"target_id"`
	TargetType  string                 `json:"target_type"`
	Attachments []model.MediaAttachment `json:"attachments"`
	Privacy     model.PrivacyLevel     `json:"privacy"`
}

// FeedResponse 动态响应
type FeedResponse struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Type        int32                  `json:"type"`
	Content     string                 `json:"content"`
	TargetID    string                 `json:"target_id,omitempty"`
	TargetType  string                 `json:"target_type,omitempty"`
	Attachments []model.MediaAttachment `json:"attachments,omitempty"`
	Privacy     int32                  `json:"privacy"`
	CreatedAt   string                 `json:"created_at"`
	LikesCount  int32                  `json:"likes_count"`
	CommentsCount int32                `json:"comments_count"`
	SharesCount int32                  `json:"shares_count"`
}

// Create 创建动态
func (h *FeedHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateFeedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed, err := h.feedService.Create(c.Request.Context(), userID, req.Type, req.Content, req.TargetID, req.TargetType, req.Attachments, req.Privacy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"feed": h.toFeedResponse(feed)})
}

// Get 获取动态
func (h *FeedHandler) Get(c *gin.Context) {
	id := c.Param("id")

	feed, err := h.feedService.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrFeedNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feed": h.toFeedResponse(feed)})
}

// Delete 删除动态
func (h *FeedHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.feedService.Delete(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrFeedNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "feed not found"})
			return
		}
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListHome 获取首页动态流
func (h *FeedHandler) ListHome(c *gin.Context) {
	userID := c.GetString("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := int64((page - 1) * pageSize)

	// 这里简化处理，实际应该获取好友列表
	var userIDs []string
	if userID != "" {
		userIDs = append(userIDs, userID)
	}

	feeds, total, err := h.feedService.ListHome(c.Request.Context(), userIDs, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*FeedResponse
	for _, feed := range feeds {
		responses = append(responses, h.toFeedResponse(feed))
	}

	c.JSON(http.StatusOK, gin.H{
		"feeds": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ListByUser 获取用户动态
func (h *FeedHandler) ListByUser(c *gin.Context) {
	userID := c.Param("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := int64((page - 1) * pageSize)

	feeds, total, err := h.feedService.ListByUser(c.Request.Context(), userID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*FeedResponse
	for _, feed := range feeds {
		responses = append(responses, h.toFeedResponse(feed))
	}

	c.JSON(http.StatusOK, gin.H{
		"feeds": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// LikeRequest 点赞请求
type LikeRequest struct {
	Action string `json:"action" binding:"required,oneof=like unlike"`
}

// Like 点赞/取消点赞
func (h *FeedHandler) Like(c *gin.Context) {
	id := c.Param("id")

	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delta := int32(1)
	if req.Action == "unlike" {
		delta = -1
	}

	if err := h.feedService.Like(c.Request.Context(), id, delta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// toFeedResponse 转换为动态响应
func (h *FeedHandler) toFeedResponse(feed *model.FeedItem) *FeedResponse {
	return &FeedResponse{
		ID:            feed.ID,
		UserID:        feed.UserID,
		Type:          int32(feed.Type),
		Content:       feed.Content,
		TargetID:      feed.TargetID,
		TargetType:    feed.TargetType,
		Attachments:   feed.Attachments,
		Privacy:       int32(feed.Privacy),
		CreatedAt:     feed.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LikesCount:    feed.LikesCount,
		CommentsCount: feed.CommentsCount,
		SharesCount:   feed.SharesCount,
	}
}
