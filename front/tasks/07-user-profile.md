# Task 7: 用户个人主页

## 目标
实现用户个人主页

## 页面结构
- 顶部: 封面图 + 用户信息
- Tab 导航: 动态、博客、相册、好友
- 内容区: 根据 Tab 显示不同内容

## 功能

### 1. 用户信息区
- 大头像
- 昵称 + 用户名
- 简介
- 统计: 博客数、好友数
- 操作按钮: 加好友 (如果是他人)、编辑资料 (如果是自己)

### 2. 动态 Tab
- 该用户的动态列表

### 3. 博客 Tab
- 该用户的博客列表

### 4. 相册 Tab (可选)
- 相册列表

### 5. 好友 Tab (可选)
- 好友列表

## API 对接
```typescript
// GET /api/v1/users/:id
// PUT /api/v1/users/me
{ nickname, avatar, bio }
```

## 需要组件
- shadcn: Tabs, Avatar, Button, Card, Separator
- 自定义: ProfileHeader, StatsCard

## 输出
- Profile 页面
- user API 封装
