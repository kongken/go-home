# Task 8: 好友与群组

## 目标
实现好友管理和群组功能

## 好友功能

### 1. 好友列表 (/friends)
- 按分组显示
- 在线状态
- 点击进入私聊

### 2. 好友请求 (/friends/requests)
- 收到的请求列表
- 发送的请求列表
- 接受/拒绝操作

### 3. 搜索用户
- 搜索框
- 用户列表
- 发送好友请求

## 群组功能

### 1. 群组列表 (/groups)
- 我加入的群组
- 推荐群组

### 2. 群组详情 (/groups/:id)
- 群组信息
- 成员列表
- 加入/退出按钮

## API 对接
```typescript
// GET /api/v1/friends
// POST /api/v1/friends/requests
// GET /api/v1/friends/requests/received
// GET /api/v1/friends/requests/sent
// POST /api/v1/friends/requests/handle

// GET /api/v1/groups
// GET /api/v1/groups/:id
// POST /api/v1/groups
// POST /api/v1/groups/:id/join
// POST /api/v1/groups/:id/leave
```

## 输出
- Friends 页面
- Groups 页面
- GroupDetail 页面
- friends/groups API 封装
