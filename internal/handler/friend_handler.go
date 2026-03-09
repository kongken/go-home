package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// FriendHandler 好友处理器
type FriendHandler struct {
	friendService service.FriendService
}

// NewFriendHandler 创建好友处理器
func NewFriendHandler(friendService service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

// SendFriendRequestRequest 发送好友请求
type SendFriendRequestRequest struct {
	ToUserID string `json:"to_user_id" binding:"required"`
	Message  string `json:"message"`
}

// FriendRequestResponse 好友请求响应
type FriendRequestResponse struct {
	ID         string `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Message    string `json:"message"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// SendRequest 发送好友请求
func (h *FriendHandler) SendRequest(c *gin.Context) {
	fromUserID := c.GetString("user_id")
	if fromUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.friendService.SendRequest(c.Request.Context(), fromUserID, req.ToUserID, req.Message)
	if err != nil {
		if err == service.ErrAlreadyFriends {
			c.JSON(http.StatusConflict, gin.H{"error": "already friends"})
			return
		}
		if err == service.ErrRequestPending {
			c.JSON(http.StatusConflict, gin.H{"error": "friend request already pending"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"request": h.toRequestResponse(request)})
}

// HandleRequestRequest 处理好友请求
type HandleRequestRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Accept    bool   `json:"accept"`
	GroupName string `json:"group_name"`
}

// HandleRequest 处理好友请求
func (h *FriendHandler) HandleRequest(c *gin.Context) {
	if c.GetString("user_id") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req HandleRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.friendService.HandleRequest(c.Request.Context(), req.RequestID, req.Accept, req.GroupName); err != nil {
		if err == service.ErrRequestNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ListReceivedRequests 获取收到的好友请求
func (h *FriendHandler) ListReceivedRequests(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

	requests, total, err := h.friendService.ListReceivedRequests(c.Request.Context(), userID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*FriendRequestResponse
	for _, req := range requests {
		responses = append(responses, h.toRequestResponse(req))
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ListSentRequests 获取发送的好友请求
func (h *FriendHandler) ListSentRequests(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

	requests, total, err := h.friendService.ListSentRequests(c.Request.Context(), userID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*FriendRequestResponse
	for _, req := range requests {
		responses = append(responses, h.toRequestResponse(req))
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// FriendshipResponse 好友关系响应
type FriendshipResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	FriendID  string `json:"friend_id"`
	GroupName string `json:"group_name"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
}

// ListFriends 获取好友列表
func (h *FriendHandler) ListFriends(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

	friendships, total, err := h.friendService.ListFriends(c.Request.Context(), userID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*FriendshipResponse
	for _, fs := range friendships {
		responses = append(responses, h.toFriendshipResponse(fs))
	}

	c.JSON(http.StatusOK, gin.H{
		"friends": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// DeleteFriend 删除好友
func (h *FriendHandler) DeleteFriend(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	friendID := c.Param("id")

	if err := h.friendService.DeleteFriend(c.Request.Context(), userID, friendID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateFriendGroupRequest 更新好友分组请求
type UpdateFriendGroupRequest struct {
	GroupName string `json:"group_name" binding:"required"`
}

// UpdateFriendGroup 更新好友分组
func (h *FriendHandler) UpdateFriendGroup(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	friendID := c.Param("id")

	var req UpdateFriendGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.friendService.UpdateFriendGroup(c.Request.Context(), userID, friendID, req.GroupName); err != nil {
		if err == service.ErrFriendshipNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "friendship not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// toRequestResponse 转换为好友请求响应
func (h *FriendHandler) toRequestResponse(req *model.FriendRequest) *FriendRequestResponse {
	return &FriendRequestResponse{
		ID:         req.ID,
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Message:    req.Message,
		Status:     int32(req.Status),
		CreatedAt:  req.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// toFriendshipResponse 转换为好友关系响应
func (h *FriendHandler) toFriendshipResponse(fs *model.Friendship) *FriendshipResponse {
	return &FriendshipResponse{
		ID:        fs.ID,
		UserID:    fs.UserID,
		FriendID:  fs.FriendID,
		GroupName: fs.GroupName,
		Remark:    fs.Remark,
		CreatedAt: fs.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
