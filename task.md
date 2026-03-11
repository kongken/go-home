技术栈:  go/gin/protobuf
当前项目使用 protobuf 管理 api, protobuf 文件目录为: `proto`

已经生成了 go 代码到 pkg/proto 目录.

任务为,根据当前的 api protobuf文档,生成 go 后端项目

## 进度

### 已完成 ✅
- [x] 项目基础结构搭建
- [x] 用户服务 (User) - 注册/登录/JWT认证/CRUD
- [x] 博客服务 (Blog) - CRUD/列表/分页
- [x] 动态服务 (Feed) - 创建/删除/列表/点赞
- [x] 好友服务 (Friend) - 请求/接受/删除/分组
- [x] 群组服务 (Group) - CRUD/成员管理/加入/离开
- [x] 消息服务 (Message) - 发送/会话/已读状态
- [x] 通知服务 (Notification) - 列表/已读/未读数
- [x] 相册服务 (Album) - 创建/删除/列表/照片管理
- [x] 评论服务 (Comment) - 创建/删除/列表
- [x] 设置服务 (Settings) - 隐私设置/通知设置/黑名单
- [x] 搜索服务 (Search) - 全局搜索/用户搜索/博客搜索/群组搜索 (基础实现)
- [x] 中间件 (Auth/CORS/Logger)
- [x] MongoDB 存储
- [x] Redis 缓存
- [x] Docker 支持
- [x] 项目成功构建

### 进行中 🔄
- [ ] Content 服务 (Activity/Poll/Share) - 待实现
- [ ] 搜索服务优化 - MongoDB 文本搜索集成

### 存储
- 数据库: MongoDB
- 缓存: Redis

## 已实现的 API 端点 (58)

### 认证 (3)
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/refresh

### 用户 (2)
- GET /api/v1/users/:id
- PUT /api/v1/users/me

### 博客 (6)
- GET /api/v1/blogs
- GET /api/v1/blogs/:id
- GET /api/v1/users/:user_id/blogs
- POST /api/v1/blogs
- PUT /api/v1/blogs/:id
- DELETE /api/v1/blogs/:id

### 动态 (6)
- GET /api/v1/feeds
- GET /api/v1/feeds/:id
- GET /api/v1/users/:user_id/feeds
- POST /api/v1/feeds
- DELETE /api/v1/feeds/:id
- POST /api/v1/feeds/:id/like

### 好友 (7)
- POST /api/v1/friends/requests
- POST /api/v1/friends/requests/handle
- GET /api/v1/friends/requests/received
- GET /api/v1/friends/requests/sent
- GET /api/v1/friends
- DELETE /api/v1/friends/:id
- PUT /api/v1/friends/:id/group

### 群组 (9)
- GET /api/v1/groups
- GET /api/v1/groups/search
- GET /api/v1/groups/:id
- GET /api/v1/groups/:id/members
- POST /api/v1/groups
- PUT /api/v1/groups/:id
- DELETE /api/v1/groups/:id
- POST /api/v1/groups/:id/join
- POST /api/v1/groups/:id/leave
- DELETE /api/v1/groups/:id/members/:user_id

### 消息 (6)
- GET /api/v1/messages/conversations
- GET /api/v1/messages/unread
- GET /api/v1/messages/:user_id
- POST /api/v1/messages
- POST /api/v1/messages/:user_id/read

### 通知 (5)
- GET /api/v1/notifications
- GET /api/v1/notifications/unread
- PUT /api/v1/notifications/:id/read
- PUT /api/v1/notifications/read-all
- DELETE /api/v1/notifications/:id

### 相册 (6)
- GET /api/v1/users/:user_id/albums
- GET /api/v1/albums/:id
- GET /api/v1/albums/:id/photos
- POST /api/v1/albums
- DELETE /api/v1/albums/:id
- POST /api/v1/albums/:id/photos

### 评论 (3)
- GET /api/v1/comments
- POST /api/v1/comments
- DELETE /api/v1/comments/:id

### 设置 (5)
- GET /api/v1/settings
- PUT /api/v1/settings/privacy
- PUT /api/v1/settings/notification
- POST /api/v1/settings/blacklist
- DELETE /api/v1/settings/blacklist/:user_id