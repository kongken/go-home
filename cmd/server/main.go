package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kongken/go-home/internal/config"
	"github.com/kongken/go-home/internal/handler"
	"github.com/kongken/go-home/internal/middleware"
	"github.com/kongken/go-home/internal/repository"
	"github.com/kongken/go-home/internal/service"
	"github.com/kongken/go-home/pkg/cache"
	"github.com/kongken/go-home/pkg/database"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		cfg = config.Get()
	}

	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 连接 MongoDB
	mongoDB, err := database.NewMongoDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoDB.Close()

	// 连接 Redis
	redisClient, err := cache.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Warn("Failed to connect to Redis", zap.Error(err))
		redisClient = nil
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	// 初始化仓库
	userRepo := repository.NewUserRepositoryMongo(mongoDB.Database)
	blogRepo := repository.NewBlogRepositoryMongo(mongoDB.Database)
	feedRepo := repository.NewFeedRepository(mongoDB.Database)
	friendRepo := repository.NewFriendRepository(mongoDB.Database)
	groupRepo := repository.NewGroupRepository(mongoDB.Database)
	messageRepo := repository.NewMessageRepository(mongoDB.Database)
	notifRepo := repository.NewNotificationRepository(mongoDB.Database)

	// 初始化服务
	userService := service.NewUserService(userRepo)
	blogService := service.NewBlogService(blogRepo)
	feedService := service.NewFeedService(feedRepo)
	friendService := service.NewFriendService(friendRepo)
	groupService := service.NewGroupService(groupRepo)
	messageService := service.NewMessageService(messageRepo)
	notifService := service.NewNotificationService(notifRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	blogHandler := handler.NewBlogHandler(blogService)
	feedHandler := handler.NewFeedHandler(feedService)
	friendHandler := handler.NewFriendHandler(friendService)
	groupHandler := handler.NewGroupHandler(groupService)
	messageHandler := handler.NewMessageHandler(messageService)
	notifHandler := handler.NewNotificationHandler(notifService)

	// 创建 Gin 引擎
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware(logger))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
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
		}
	}

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.String("addr", srv.Addr))

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}