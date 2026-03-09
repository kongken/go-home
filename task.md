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
- [x] 中间件 (Auth/CORS/Logger)
- [x] MongoDB 存储 (users/blogs/feeds/friendships/friend_requests)
- [x] Redis 缓存
- [x] Docker 支持
- [x] 项目成功构建 (21MB 二进制)

### 进行中 🔄
- [ ] Group 服务
- [ ] Message 服务
- [ ] Notification 服务
- [ ] Content 服务 (Activity/Album/Comment/Poll/Share)
- [ ] Search 服务
- [ ] Settings 服务

### 存储
- 数据库: MongoDB
- 缓存: Redis

### 已实现的 API 端点

#### 认证
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/refresh

#### 用户
- GET /api/v1/users/:id
- PUT /api/v1/users/me (需认证)

#### 博客
- GET /api/v1/blogs
- GET /api/v1/blogs/:id
- GET /api/v1/users/:user_id/blogs
- POST /api/v1/blogs (需认证)
- PUT /api/v1/blogs/:id (需认证)
- DELETE /api/v1/blogs/:id (需认证)

#### 动态
- GET /api/v1/feeds
- GET /api/v1/feeds/:id
- GET /api/v1/users/:user_id/feeds
- POST /api/v1/feeds (需认证)
- DELETE /api/v1/feeds/:id (需认证)
- POST /api/v1/feeds/:id/like (需认证)

#### 好友
- POST /api/v1/friends/requests (需认证)
- POST /api/v1/friends/requests/handle (需认证)
- GET /api/v1/friends/requests/received (需认证)
- GET /api/v1/friends/requests/sent (需认证)
- GET /api/v1/friends (需认证)
- DELETE /api/v1/friends/:id (需认证)
- PUT /api/v1/friends/:id/group (需认证)
