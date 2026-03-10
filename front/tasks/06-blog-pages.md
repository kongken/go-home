# Task 6: 博客功能

## 目标
实现博客列表、详情、编辑功能

## 页面

### 1. 博客列表 (/blogs)
- 分类筛选
- 标签筛选
- 分页/无限滚动
- 博客卡片列表

### 2. 博客详情 (/blogs/:id)
- 标题、作者、发布时间
- 正文内容 (Markdown 渲染)
- 标签
- 操作: 点赞、评论、收藏
- 作者信息卡片

### 3. 博客编辑 (/blogs/new, /blogs/:id/edit)
- 标题输入
- Markdown 编辑器
- 摘要输入
- 封面图片上传
- 标签输入
- 隐私设置 (公开/私密)
- 保存/发布按钮

### 4. 用户博客列表 (/users/:id/blogs)
- 指定用户的所有博客
- 按时间排序

## API 对接
```typescript
// GET /api/v1/blogs?page=1&page_size=10&category=&user_id=
// GET /api/v1/blogs/:id
// POST /api/v1/blogs
{ title, content, summary, cover_image, tags, category, privacy, status }
// PUT /api/v1/blogs/:id
// DELETE /api/v1/blogs/:id
// GET /api/v1/users/:user_id/blogs
```

## 需要组件
- shadcn: Card, Input, Textarea, Select, Badge, Tabs, Dialog
- 自定义: MarkdownEditor, BlogCard, TagInput

## 输出
- BlogList 页面
- BlogDetail 页面
- BlogEdit 页面
- blog API 封装
