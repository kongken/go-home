package wire

import (
	"github.com/google/wire"
	"github.com/kongken/go-home/internal/cache"
	"github.com/kongken/go-home/internal/handler"
	"github.com/kongken/go-home/internal/repository"
	"github.com/kongken/go-home/internal/service"
)

// RepositorySet 仓库提供者集合
var RepositorySet = wire.NewSet(
	repository.NewUserRepositoryButterfly,
	repository.NewBlogRepositoryButterfly,
	repository.NewFeedRepositoryButterfly,
	repository.NewFriendRepositoryButterfly,
	repository.NewGroupRepositoryButterfly,
	repository.NewMessageRepositoryButterfly,
	repository.NewNotificationRepositoryButterfly,
	repository.NewCommentRepositoryButterfly,
	repository.NewAlbumRepositoryButterfly,
	repository.NewSettingsRepositoryButterfly,
)

// CacheSet 缓存提供者集合
var CacheSet = wire.NewSet(
	cache.NewRedisCache,
)

// ServiceSet 服务提供者集合
var ServiceSet = wire.NewSet(
	service.NewUserService,
	service.NewBlogService,
	service.NewFeedService,
	service.NewFriendService,
	service.NewGroupService,
	service.NewMessageService,
	service.NewNotificationService,
	service.NewCommentService,
	service.NewAlbumService,
	service.NewSettingsService,
)

// HandlerSet Handler 提供者集合
var HandlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewBlogHandler,
	handler.NewFeedHandler,
	handler.NewFriendHandler,
	handler.NewGroupHandler,
	handler.NewMessageHandler,
	handler.NewNotificationHandler,
	handler.NewCommentHandler,
	handler.NewAlbumHandler,
	handler.NewSettingsHandler,
)

// AppSet 应用提供者集合
var AppSet = wire.NewSet(
	RepositorySet,
	CacheSet,
	ServiceSet,
	HandlerSet,
)

// Handlers 所有 Handler 的集合
type Handlers struct {
	UserHandler    *handler.UserHandler
	BlogHandler    *handler.BlogHandler
	FeedHandler    *handler.FeedHandler
	FriendHandler  *handler.FriendHandler
	GroupHandler   *handler.GroupHandler
	MessageHandler *handler.MessageHandler
	NotifHandler   *handler.NotificationHandler
	CommentHandler *handler.CommentHandler
	AlbumHandler   *handler.AlbumHandler
	SettingsHandler *handler.SettingsHandler
}

// NewHandlers 创建所有 Handler
func NewHandlers(
	userHandler *handler.UserHandler,
	blogHandler *handler.BlogHandler,
	feedHandler *handler.FeedHandler,
	friendHandler *handler.FriendHandler,
	groupHandler *handler.GroupHandler,
	messageHandler *handler.MessageHandler,
	notifHandler *handler.NotificationHandler,
	commentHandler *handler.CommentHandler,
	albumHandler *handler.AlbumHandler,
	settingsHandler *handler.SettingsHandler,
) *Handlers {
	return &Handlers{
		UserHandler:     userHandler,
		BlogHandler:     blogHandler,
		FeedHandler:     feedHandler,
		FriendHandler:   friendHandler,
		GroupHandler:    groupHandler,
		MessageHandler:  messageHandler,
		NotifHandler:    notifHandler,
		CommentHandler:  commentHandler,
		AlbumHandler:    albumHandler,
		SettingsHandler: settingsHandler,
	}
}