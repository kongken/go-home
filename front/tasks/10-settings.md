# Task 10: 设置页面

## 目标
实现用户设置页面

## 功能

### 1. 个人资料设置
- 头像上传
- 昵称修改
- 简介修改

### 2. 隐私设置
- 博客可见性
- 动态可见性
- 好友验证设置

### 3. 通知设置
- 邮件通知开关
- 站内通知开关

### 4. 安全设置
- 修改密码
- 登录设备管理

### 5. 黑名单
- 拉黑用户列表
- 解除拉黑

## API 对接
```typescript
// GET /api/v1/settings
// PUT /api/v1/settings/privacy
// PUT /api/v1/settings/notification
// POST /api/v1/settings/blacklist
// DELETE /api/v1/settings/blacklist/:user_id
```

## 输出
- Settings 页面
- settings API 封装
