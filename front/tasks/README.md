# Frontend Tasks - uchome React + shadcn/ui

## 项目概述
使用 React + Vite + TypeScript + shadcn/ui 构建仿 uchome 风格的前端项目，对接 go-home 后端服务。

## Task 列表

| 序号 | Task | 描述 | 状态 |
|------|------|------|------|
| 01 | [初始化项目](./01-init-project.md) | React + Vite + TypeScript + shadcn/ui 初始化 | ✅ 已完成 |
| 02 | [项目结构](./02-project-structure.md) | 目录结构设计和基础文件创建 | ✅ 已完成 |
| 03 | [认证页面](./03-auth-pages.md) | 登录/注册页面 | ✅ 已完成 |
| 04 | [布局组件](./04-layout-components.md) | Header, Sidebar, MainLayout | ✅ 已完成 |
| 05 | [首页动态流](./05-home-feed.md) | Feed 流、发布框 | ✅ 已完成 |
| 06 | [博客功能](./06-blog-pages.md) | 博客列表、详情、编辑 | ✅ 已完成 |
| 07 | [个人主页](./07-user-profile.md) | 用户 Profile 页面 | ✅ 已完成 |
| 08 | [好友与群组](./08-friends-groups.md) | 好友管理、群组功能 | ⏳ 待开始 |
| 09 | [消息系统](./09-messages.md) | 私信、通知 | ⏳ 待开始 |
| 10 | [设置页面](./10-settings.md) | 用户设置 | ⏳ 待开始 |

## 技术栈
- **框架**: React 18 + Vite
- **语言**: TypeScript
- **UI**: shadcn/ui + Tailwind CSS
- **路由**: React Router DOM
- **状态**: Zustand
- **请求**: Axios
- **查询**: TanStack Query (React Query)

## 后端 API Base
```
http://localhost:2222/api/v1
```

## 开发顺序建议
1. Task 01-02: 项目初始化
2. Task 03: 认证 (先做登录才能测试其他)
3. Task 04: 布局
4. Task 05: 首页 (核心功能)
5. Task 06: 博客 (核心功能)
6. Task 07: 个人主页
7. Task 08-10: 其他功能
