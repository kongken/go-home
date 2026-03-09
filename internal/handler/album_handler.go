package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// AlbumHandler 相册处理器
type AlbumHandler struct {
	albumService service.AlbumService
}

// NewAlbumHandler 创建相册处理器
func NewAlbumHandler(albumService service.AlbumService) *AlbumHandler {
	return &AlbumHandler{albumService: albumService}
}

// CreateAlbumRequest 创建相册请求
type CreateAlbumRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	CoverPhoto  string             `json:"cover_photo"`
	Privacy     model.PrivacyLevel `json:"privacy"`
}

// AlbumResponse 相册响应
type AlbumResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	CoverPhoto    string `json:"cover_photo"`
	Privacy       int32  `json:"privacy"`
	PhotosCount   int32  `json:"photos_count"`
	ViewsCount    int32  `json:"views_count"`
	CommentsCount int32  `json:"comments_count"`
	CreatedAt     string `json:"created_at"`
}

// Create 创建相册
func (h *AlbumHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	album, err := h.albumService.Create(c.Request.Context(), userID, req.Name, req.Description, req.CoverPhoto, req.Privacy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"album": h.toAlbumResponse(album)})
}

// Get 获取相册
func (h *AlbumHandler) Get(c *gin.Context) {
	id := c.Param("id")

	album, err := h.albumService.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrAlbumNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"album": h.toAlbumResponse(album)})
}

// ListByUser 获取用户相册
func (h *AlbumHandler) ListByUser(c *gin.Context) {
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

	albums, total, err := h.albumService.ListByUser(c.Request.Context(), userID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*AlbumResponse
	for _, album := range albums {
		responses = append(responses, h.toAlbumResponse(album))
	}

	c.JSON(http.StatusOK, gin.H{
		"albums": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Delete 删除相册
func (h *AlbumHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.albumService.Delete(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrAlbumNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
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

// PhotoResponse 照片响应
type PhotoResponse struct {
	ID          string `json:"id"`
	AlbumID     string `json:"album_id"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
	CreatedAt   string `json:"created_at"`
}

// AddPhotoRequest 添加照片请求
type AddPhotoRequest struct {
	URL         string `json:"url" binding:"required"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
}

// AddPhoto 添加照片
func (h *AlbumHandler) AddPhoto(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	albumID := c.Param("id")

	var req AddPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo, err := h.albumService.AddPhoto(c.Request.Context(), albumID, userID, req.URL, req.Thumbnail, req.Description, req.Width, req.Height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"photo": h.toPhotoResponse(photo)})
}

// GetPhotos 获取照片
func (h *AlbumHandler) GetPhotos(c *gin.Context) {
	albumID := c.Param("id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := int64((page - 1) * pageSize)

	photos, total, err := h.albumService.GetPhotos(c.Request.Context(), albumID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*PhotoResponse
	for _, photo := range photos {
		responses = append(responses, h.toPhotoResponse(photo))
	}

	c.JSON(http.StatusOK, gin.H{
		"photos": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// toAlbumResponse 转换为相册响应
func (h *AlbumHandler) toAlbumResponse(album *model.Album) *AlbumResponse {
	return &AlbumResponse{
		ID:            album.ID,
		UserID:        album.UserID,
		Name:          album.Name,
		Description:   album.Description,
		CoverPhoto:    album.CoverPhoto,
		Privacy:       int32(album.Privacy),
		PhotosCount:   album.PhotosCount,
		ViewsCount:    album.ViewsCount,
		CommentsCount: album.CommentsCount,
		CreatedAt:     album.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// toPhotoResponse 转换为照片响应
func (h *AlbumHandler) toPhotoResponse(photo *model.Photo) *PhotoResponse {
	return &PhotoResponse{
		ID:          photo.ID,
		AlbumID:     photo.AlbumID,
		URL:         photo.URL,
		Thumbnail:   photo.Thumbnail,
		Description: photo.Description,
		Width:       photo.Width,
		Height:      photo.Height,
		CreatedAt:   photo.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}