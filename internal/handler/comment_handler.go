package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService service.CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	TargetID    string                 `json:"target_id" binding:"required"`
	TargetType  string                 `json:"target_type" binding:"required"`
	Content     string                 `json:"content" binding:"required"`
	ParentID    string                 `json:"parent_id"`
	Attachments []model.MediaAttachment `json:"attachments"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	TargetID    string                 `json:"target_id"`
	TargetType  string                 `json:"target_type"`
	Content     string                 `json:"content"`
	ParentID    string                 `json:"parent_id,omitempty"`
	Attachments []model.MediaAttachment `json:"attachments,omitempty"`
	LikesCount  int32                  `json:"likes_count"`
	RepliesCount int32                 `json:"replies_count"`
	CreatedAt   string                 `json:"created_at"`
}

// Create 创建评论
func (h *CommentHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.commentService.Create(c.Request.Context(), userID, req.TargetID, req.TargetType, req.Content, req.ParentID, req.Attachments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"comment": h.toCommentResponse(comment)})
}

// List 获取评论列表
func (h *CommentHandler) List(c *gin.Context) {
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")
	if targetID == "" || targetType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_id and target_type required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := int64((page - 1) * pageSize)

	comments, total, err := h.commentService.ListByTarget(c.Request.Context(), targetID, targetType, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*CommentResponse
	for _, c := range comments {
		responses = append(responses, h.toCommentResponse(c))
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.commentService.Delete(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrCommentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
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

// toCommentResponse 转换为评论响应
func (h *CommentHandler) toCommentResponse(c *model.Comment) *CommentResponse {
	return &CommentResponse{
		ID:           c.ID,
		UserID:       c.UserID,
		TargetID:     c.TargetID,
		TargetType:   c.TargetType,
		Content:      c.Content,
		ParentID:     c.ParentID,
		Attachments:  c.Attachments,
		LikesCount:   c.LikesCount,
		RepliesCount: c.RepliesCount,
		CreatedAt:    c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}