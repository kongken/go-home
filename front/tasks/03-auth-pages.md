# Task 3: 认证页面

## 目标
实现登录和注册页面

## 页面

### 1. 登录页 (/login)
- 账号/邮箱输入框
- 密码输入框
- 登录按钮
- 跳转到注册链接

### 2. 注册页 (/register)
- 用户名输入框
- 邮箱输入框
- 密码输入框
- 确认密码输入框
- 注册按钮
- 跳转到登录链接

## 需要组件
- shadcn: Card, Input, Button, Label, Form
- 自定义: AuthLayout (左右分栏布局)

## API 对接
```typescript
// POST /api/v1/auth/login
{ account: string, password: string }

// POST /api/v1/auth/register
{ username: string, email: string, password: string }

// POST /api/v1/auth/refresh
{ refresh_token: string }
```

## 状态管理
- authStore: token, user, isAuthenticated, login(), logout()

## 输出
- Login 页面组件
- Register 页面组件
- authStore 状态管理
- API 接口封装
