# Task 9: 消息系统

## 目标
实现私信和通知功能

## 私信功能

### 1. 会话列表 (/messages)
- 最近联系人列表
- 未读消息数
- 最后消息预览

### 2. 聊天界面 (/messages/:user_id)
- 消息气泡列表
- 输入框
- 发送按钮
- 自动滚动到底部

## 通知功能

### 1. 通知列表
- 评论回复
- 点赞通知
- 好友请求
- 系统通知
- 标记已读

## API 对接
```typescript
// GET /api/v1/messages/conversations
// GET /api/v1/messages/:user_id
// POST /api/v1/messages
// POST /api/v1/messages/:user_id/read

// GET /api/v1/notifications
// PUT /api/v1/notifications/:id/read
// PUT /api/v1/notifications/read-all
```

## 输出
- Messages 页面
- Conversation 组件
- Notification 组件
- messages/notifications API 封装
