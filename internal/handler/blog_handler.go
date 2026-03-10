package handler

import (
	"net/http"
	"strconv"

	"butterfly.orx.me/core/log"
	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// BlogHandler 博客处理器
type BlogHandler struct {
	blogService service.BlogService
}

// NewBlogHandler 创建博客处理器
func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}

// CreateBlogRequest 创建博客请求
type CreateBlogRequest struct {
	Title      string             `json:"title" binding:"required"`
	Content    string             `json:"content" binding:"required"`
	Summary    string             `json:"summary"`
	CoverImage string             `json:"cover_image"`
	Tags       []string           `json:"tags"`
	Category   string             `json:"category"`
	Privacy    model.PrivacyLevel `json:"privacy"`
	Status     model.BlogStatus   `json:"status"`
}

// BlogResponse 博客响应
type BlogResponse struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Summary     string            `json:"summary"`
	CoverImage  string            `json:"cover_image"`
	Tags        []string          `json:"tags"`
	Category    string            `json:"category"`
	Privacy     int32             `json:"privacy"`
	Status      int32             `json:"status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	ViewsCount  int32             `json:"views_count"`
	LikesCount  int32             `json:"likes_count"`
	CommentsCount int32           `json:"comments_count"`
	Author      *UserInfo         `json:"author,omitempty"`
}

// Create 创建博客
func (h *BlogHandler) Create(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	userID := c.GetString("user_id")
	if userID == "" {
		logger.Warn("unauthorized blog create attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("invalid blog create request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog, err := h.blogService.Create(c.Request.Context(), userID, req.Title, req.Content, req.Summary, req.CoverImage, req.Tags, req.Category, req.Privacy, req.Status)
	if err != nil {
		logger.Error("failed to create blog", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("blog created via handler", "blog_id", blog.ID, "user_id", userID)
	c.JSON(http.StatusCreated, gin.H{"blog": h.toBlogResponse(blog)})
}

// Get 获取博客
func (h *BlogHandler) Get(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	id := c.Param("id")

	blog, err := h.blogService.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrBlogNotFound {
			logger.Warn("blog not found", "blog_id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
			return
		}
		logger.Error("failed to get blog", "blog_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog": h.toBlogResponse(blog)})
}

// UpdateBlogRequest 更新博客请求
type UpdateBlogRequest struct {
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Summary    string             `json:"summary"`
	CoverImage string             `json:"cover_image"`
	Tags       []string           `json:"tags"`
	Category   string             `json:"category"`
	Privacy    model.PrivacyLevel `json:"privacy"`
	Status     model.BlogStatus   `json:"status"`
}

// Update 更新博客
func (h *BlogHandler) Update(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	userID := c.GetString("user_id")
	if userID == "" {
		logger.Warn("unauthorized blog update attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	var req UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("invalid blog update request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	updates["privacy"] = req.Privacy
	updates["status"] = req.Status

	blog, err := h.blogService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		if err == service.ErrBlogNotFound {
			logger.Warn("blog not found for update", "blog_id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
			return
		}
		if err == service.ErrUnauthorized {
			logger.Warn("forbidden blog update", "blog_id", id, "user_id", userID)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		logger.Error("failed to update blog", "blog_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("blog updated via handler", "blog_id", id, "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"blog": h.toBlogResponse(blog)})
}

// Delete 删除博客
func (h *BlogHandler) Delete(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	userID := c.GetString("user_id")
	if userID == "" {
		logger.Warn("unauthorized blog delete attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.blogService.Delete(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrBlogNotFound {
			logger.Warn("blog not found for delete", "blog_id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
			return
		}
		if err == service.ErrUnauthorized {
			logger.Warn("forbidden blog delete", "blog_id", id, "user_id", userID)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		logger.Error("failed to delete blog", "blog_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("blog deleted via handler", "blog_id", id, "user_id", userID)
	c.JSON(http.StatusNoContent, nil)
}

// List 获取博客列表
func (h *BlogHandler) List(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	userID := c.Query("user_id")
	category := c.Query("category")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	blogs, total, err := h.blogService.List(c.Request.Context(), userID, category, nil, offset, pageSize)
	if err != nil {
		logger.Error("failed to list blogs", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*BlogResponse
	for _, blog := range blogs {
		responses = append(responses, h.toBlogResponse(blog))
	}

	logger.Debug("blogs listed", "count", len(blogs), "total", total, "page", page)
	c.JSON(http.StatusOK, gin.H{
		"blogs": responses,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ListByUser 获取指定用户的博客列表
func (h *BlogHandler) ListByUser(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	
	userID := c.Param("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	blogs, total, err := h.blogService.ListByUser(c.Request.Context(), userID, offset, pageSize)
	if err != nil {
		logger.Error("failed to list user blogs", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*BlogResponse
	for _, blog := range blogs {
		responses = append(responses, h.toBlogResponse(blog))
	}

	logger.Debug("user blogs listed", "user_id", userID, "count", len(blogs))
	c.JSON(http.StatusOK, gin.H{
		"blogs": responses,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// toBlogResponse 转换为博客响应
func (h *BlogHandler) toBlogResponse(blog *model.Blog) *BlogResponse {
	return &BlogResponse{
		ID:          blog.ID,
		UserID:      blog.UserID,
		Title:       blog.Title,
		Content:     blog.Content,
		Summary:     blog.Summary,
		CoverImage:  blog.CoverImage,
		Category:    blog.Category,
		Privacy:     int32(blog.Privacy),
		Status:      int32(blog.Status),
		CreatedAt:   blog.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   blog.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		ViewsCount:  blog.ViewsCount,
		LikesCount:  blog.LikesCount,
		CommentsCount: blog.CommentsCount,
	}
}
