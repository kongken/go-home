package main

import (
	"log"
	
	"github.com/gin-gonic/gin"
	"butterfly.orx.me/core/app"

	"github.com/kongken/go-home/internal/config"
	"github.com/kongken/go-home/internal/handler"
	"github.com/kongken/go-home/internal/middleware"
)

func main() {
	// Butterfly 应用配置
	butterflyConfig := &app.Config{
		Service: "go-home",
		Config:  &config.ButterflyConfig{},
		
		// HTTP 路由注册
		Router: setupRoutes,
		
		// 初始化函数链
		InitFunc: []func() error{
			initServices,
		},
		
		// 清理函数链
		TeardownFunc: []func() error{
			cleanupServices,
		},
	}

	// 创建并运行应用
	application := app.New(butterflyConfig)
	application.Run()
}

// 全局服务实例 (butterfly 管理生命周期)
var (
	userHandler    *handler.UserHandler
	blogHandler    *handler.BlogHandler
	feedHandler    *handler.FeedHandler
	friendHandler  *handler.FriendHandler
	groupHandler   *handler.GroupHandler
	messageHandler *handler.MessageHandler
	notifHandler   *handler.NotificationHandler
	commentHandler *handler.CommentHandler
	albumHandler   *handler.AlbumHandler
	settingsHandler *handler.SettingsHandler
)

// initServices 初始化所有服务 (使用 Wire)
func initServices() error {
	log.Println("initializing services with Wire...")
	
	// 设置全局配置
	cfg := &config.ButterflyConfig{}
	config.SetGlobalConfig(cfg)
	
	// 使用 Wire 初始化所有 Handler
	handlers, err := InitializeHandlers()
	if err != nil {
		return err
	}
	
	// 赋值给全局变量
	userHandler = handlers.UserHandler
	blogHandler = handlers.BlogHandler
	feedHandler = handlers.FeedHandler
	friendHandler = handlers.FriendHandler
	groupHandler = handlers.GroupHandler
	messageHandler = handlers.MessageHandler
	notifHandler = handlers.NotifHandler
	commentHandler = handlers.CommentHandler
	albumHandler = handlers.AlbumHandler
	settingsHandler = handlers.SettingsHandler
	
	log.Println("services initialized successfully with Wire")
	return nil
}

// cleanupServices 清理服务
func cleanupServices() error {
	log.Println("cleaning up services...")
	// butterfly 会自动处理 store 连接的关闭
	return nil
}

// setupRoutes 设置 HTTP 路由
func setupRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "go-home",
		})
	})
	
	// 就绪检查
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})
	
	// API 路由组
	api := r.Group("/api/v1")
	{
		// 公开路由
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)
		api.POST("/auth/refresh", userHandler.RefreshToken)
		
		// 用户路由
		api.GET("/users/:id", userHandler.GetUser)
		
		// 博客路由（公开）
		api.GET("/blogs", blogHandler.List)
		api.GET("/blogs/:id", blogHandler.Get)
		api.GET("/users/:user_id/blogs", blogHandler.ListByUser)
		
		// 动态路由（公开）
		api.GET("/feeds", feedHandler.ListHome)
		api.GET("/feeds/:id", feedHandler.Get)
		api.GET("/users/:user_id/feeds", feedHandler.ListByUser)
		
		// 群组路由（公开）
		api.GET("/groups", groupHandler.List)
		api.GET("/groups/:id", groupHandler.Get)
		api.GET("/groups/:id/members", groupHandler.ListMembers)
		api.GET("/groups/search", groupHandler.Search)
		
		// 评论路由（公开）
		api.GET("/comments", commentHandler.List)
		
		// 相册路由（公开）
		api.GET("/users/:user_id/albums", albumHandler.ListByUser)
		api.GET("/albums/:id", albumHandler.Get)
		api.GET("/albums/:id/photos", albumHandler.GetPhotos)
		
		// 需要认证的路由
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户
			authorized.PUT("/users/me", userHandler.UpdateUser)
			
			// 博客
			authorized.POST("/blogs", blogHandler.Create)
			authorized.PUT("/blogs/:id", blogHandler.Update)
			authorized.DELETE("/blogs/:id", blogHandler.Delete)
			
			// 动态
			authorized.POST("/feeds", feedHandler.Create)
			authorized.DELETE("/feeds/:id", feedHandler.Delete)
			authorized.POST("/feeds/:id/like", feedHandler.Like)
			
			// 好友
			authorized.POST("/friends/requests", friendHandler.SendRequest)
			authorized.POST("/friends/requests/handle", friendHandler.HandleRequest)
			authorized.GET("/friends/requests/received", friendHandler.ListReceivedRequests)
			authorized.GET("/friends/requests/sent", friendHandler.ListSentRequests)
			authorized.GET("/friends", friendHandler.ListFriends)
			authorized.DELETE("/friends/:id", friendHandler.DeleteFriend)
			authorized.PUT("/friends/:id/group", friendHandler.UpdateFriendGroup)
			
			// 群组
			authorized.POST("/groups", groupHandler.Create)
			authorized.PUT("/groups/:id", groupHandler.Update)
			authorized.DELETE("/groups/:id", groupHandler.Delete)
			authorized.POST("/groups/:id/join", groupHandler.Join)
			authorized.POST("/groups/:id/leave", groupHandler.Leave)
			authorized.DELETE("/groups/:id/members/:user_id", groupHandler.KickMember)
			
			// 消息
			authorized.POST("/messages", messageHandler.Send)
			authorized.GET("/messages/conversations", messageHandler.ListConversations)
			authorized.GET("/messages/unread", messageHandler.UnreadCount)
			authorized.GET("/messages/:user_id", messageHandler.ListMessages)
			authorized.POST("/messages/:user_id/read", messageHandler.MarkAsRead)
			
			// 通知
			authorized.GET("/notifications", notifHandler.List)
			authorized.GET("/notifications/unread", notifHandler.UnreadCount)
			authorized.PUT("/notifications/:id/read", notifHandler.MarkAsRead)
			authorized.PUT("/notifications/read-all", notifHandler.MarkAllAsRead)
			authorized.DELETE("/notifications/:id", notifHandler.Delete)
			
			// 评论
			authorized.POST("/comments", commentHandler.Create)
			authorized.DELETE("/comments/:id", commentHandler.Delete)
			
			// 相册
			authorized.POST("/albums", albumHandler.Create)
			authorized.DELETE("/albums/:id", albumHandler.Delete)
			authorized.POST("/albums/:id/photos", albumHandler.AddPhoto)
			
			// 设置
			authorized.GET("/settings", settingsHandler.Get)
			authorized.PUT("/settings/privacy", settingsHandler.UpdatePrivacy)
			authorized.PUT("/settings/notification", settingsHandler.UpdateNotification)
			authorized.POST("/settings/blacklist", settingsHandler.AddToBlacklist)
			authorized.DELETE("/settings/blacklist/:user_id", settingsHandler.RemoveFromBlacklist)
		}
	}
}