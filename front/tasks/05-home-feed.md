# Task 5: 首页动态流

## 目标
实现 uchome 风格的首页动态流

## 页面结构
- 左侧: 用户信息卡片
- 中间: 动态发布框 + 动态列表
- 右侧: 推荐用户/热门博客 (可选)

## 功能

### 1. 动态发布框
- 文本输入区
- 发布按钮
- 可选: 添加图片、@用户

### 2. 动态列表
- 无限滚动加载
- 动态卡片:
  - 用户头像+昵称
  - 发布时间
  - 内容文本
  - 图片预览 (如果有)
  - 操作栏: 点赞、评论、分享

### 3. 动态类型
- 纯文本动态
- 博客分享动态
- 图片动态

## API 对接
```typescript
// GET /api/v1/feeds?page=1&page_size=20
// POST /api/v1/feeds
{ content: string, type: number, attachments?: [] }
// POST /api/v1/feeds/:id/like
{ delta: 1 | -1 }
```

## 需要组件
- shadcn: Card, Avatar, Button, Textarea, Skeleton
- 自定义: FeedCard, FeedComposer, InfiniteScroll

## 输出
- Home 页面
- FeedCard 组件
- FeedComposer 组件
- feed API 封装
