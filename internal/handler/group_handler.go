package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/service"
)

// GroupHandler 群组处理器
type GroupHandler struct {
	groupService service.GroupService
}

// NewGroupHandler 创建群组处理器
func NewGroupHandler(groupService service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// CreateGroupRequest 创建群组请求
type CreateGroupRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	Avatar      string          `json:"avatar"`
	Category    string          `json:"category"`
	Type        model.GroupType `json:"type"`
	JoinMode    model.JoinMode  `json:"join_mode"`
	MemberLimit int32           `json:"member_limit"`
}

// GroupResponse 群组响应
type GroupResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Avatar       string `json:"avatar"`
	Category     string `json:"category"`
	OwnerID      string `json:"owner_id"`
	Type         int32  `json:"type"`
	JoinMode     int32  `json:"join_mode"`
	MemberLimit  int32  `json:"member_limit"`
	MembersCount int32  `json:"members_count"`
	CreatedAt    string `json:"created_at"`
}

// Create 创建群组
func (h *GroupHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.groupService.Create(c.Request.Context(), userID, req.Name, req.Description, req.Avatar, req.Category, req.Type, req.JoinMode, req.MemberLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"group": h.toGroupResponse(group)})
}

// Get 获取群组
func (h *GroupHandler) Get(c *gin.Context) {
	id := c.Param("id")

	group, err := h.groupService.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": h.toGroupResponse(group)})
}

// UpdateGroupRequest 更新群组请求
type UpdateGroupRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Avatar      string         `json:"avatar"`
	Category    string         `json:"category"`
	JoinMode    model.JoinMode `json:"join_mode"`
	MemberLimit int32          `json:"member_limit"`
}

// Update 更新群组
func (h *GroupHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	updates["join_mode"] = req.JoinMode
	updates["member_limit"] = req.MemberLimit

	group, err := h.groupService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		if err == service.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err == service.ErrNotGroupAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": h.toGroupResponse(group)})
}

// Delete 删除群组
func (h *GroupHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.groupService.Delete(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err == service.ErrNotGroupOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// List 获取群组列表
func (h *GroupHandler) List(c *gin.Context) {
	category := c.Query("category")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := int64((page - 1) * pageSize)

	groups, total, err := h.groupService.List(c.Request.Context(), category, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*GroupResponse
	for _, group := range groups {
		responses = append(responses, h.toGroupResponse(group))
	}

	c.JSON(http.StatusOK, gin.H{
		"groups": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Search 搜索群组
func (h *GroupHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword required"})
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

	groups, total, err := h.groupService.Search(c.Request.Context(), keyword, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*GroupResponse
	for _, group := range groups {
		responses = append(responses, h.toGroupResponse(group))
	}

	c.JSON(http.StatusOK, gin.H{
		"groups": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Join 加入群组
func (h *GroupHandler) Join(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.groupService.Join(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err == service.ErrGroupFull {
			c.JSON(http.StatusBadRequest, gin.H{"error": "group is full"})
			return
		}
		if err == service.ErrGroupMemberExists {
			c.JSON(http.StatusConflict, gin.H{"error": "already a member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined successfully"})
}

// Leave 离开群组
func (h *GroupHandler) Leave(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	if err := h.groupService.Leave(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left successfully"})
}

// KickMember 踢出成员
func (h *GroupHandler) KickMember(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	groupID := c.Param("id")
	targetUserID := c.Param("user_id")

	if err := h.groupService.KickMember(c.Request.Context(), groupID, userID, targetUserID); err != nil {
		if err == service.ErrNotGroupAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member kicked"})
}

// MemberResponse 成员响应
type MemberResponse struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	UserID   string `json:"user_id"`
	Role     int32  `json:"role"`
	JoinedAt string `json:"joined_at"`
}

// ListMembers 获取成员列表
func (h *GroupHandler) ListMembers(c *gin.Context) {
	groupID := c.Param("id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := int64((page - 1) * pageSize)

	members, total, err := h.groupService.ListMembers(c.Request.Context(), groupID, offset, int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*MemberResponse
	for _, m := range members {
		responses = append(responses, &MemberResponse{
			ID:       m.ID,
			GroupID:  m.GroupID,
			UserID:   m.UserID,
			Role:     int32(m.Role),
			JoinedAt: m.JoinedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"members": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// toGroupResponse 转换为群组响应
func (h *GroupHandler) toGroupResponse(group *model.Group) *GroupResponse {
	return &GroupResponse{
		ID:           group.ID,
		Name:         group.Name,
		Description:  group.Description,
		Avatar:       group.Avatar,
		Category:     group.Category,
		OwnerID:      group.OwnerID,
		Type:         int32(group.Type),
		JoinMode:     int32(group.JoinMode),
		MemberLimit:  group.MemberLimit,
		MembersCount: group.MembersCount,
		CreatedAt:    group.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}