# Task 2: 项目目录结构设计

## 目标
设计仿照 uchome 的前端项目结构

## 目录结构
```
front/
├── src/
│   ├── api/              # API 接口封装
│   │   ├── client.ts     # axios 实例
│   │   ├── auth.ts       # 认证相关 API
│   │   ├── user.ts       # 用户相关 API
│   │   ├── blog.ts       # 博客相关 API
│   │   ├── feed.ts       # 动态相关 API
│   │   └── index.ts      # API 导出
│   ├── components/       # 公共组件
│   │   ├── ui/           # shadcn 组件
│   │   ├── layout/       # 布局组件
│   │   │   ├── Header.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── MainLayout.tsx
│   │   └── common/       # 通用组件
│   │       ├── UserCard.tsx
│   │       ├── BlogCard.tsx
│   │       └── FeedCard.tsx
│   ├── hooks/            # 自定义 hooks
│   │   ├── useAuth.ts
│   │   ├── useUser.ts
│   │   └── useApi.ts
│   ├── pages/            # 页面组件
│   │   ├── Home/         # 首页/动态流
│   │   ├── Login/        # 登录
│   │   ├── Register/     # 注册
│   │   ├── Profile/      # 个人主页
│   │   ├── Blog/         # 博客列表/详情
│   │   ├── BlogEdit/     # 博客编辑
│   │   ├── Friends/      # 好友管理
│   │   ├── Groups/       # 群组
│   │   ├── Messages/     # 消息
│   │   └── Settings/     # 设置
│   ├── stores/           # 状态管理 (zustand)
│   │   ├── authStore.ts
│   │   ├── userStore.ts
│   │   └── themeStore.ts
│   ├── types/            # TypeScript 类型定义
│   │   ├── user.ts
│   │   ├── blog.ts
│   │   ├── feed.ts
│   │   └── index.ts
│   ├── utils/            # 工具函数
│   │   ├── storage.ts
│   │   ├── format.ts
│   │   └── constants.ts
│   ├── App.tsx
│   └── main.tsx
├── public/
├── index.html
├── package.json
├── tsconfig.json
├── tailwind.config.js
└── vite.config.ts
```

## 输出
- 创建所有目录
- 创建基础文件模板
