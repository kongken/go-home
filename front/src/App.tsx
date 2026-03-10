import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MainLayout } from '@/components/layout/MainLayout'
import { Login } from '@/pages/Login'
import { Register } from '@/pages/Register'
import { Home } from '@/pages/Home'
import { BlogList } from '@/pages/Blog'

const queryClient = new QueryClient()

// 空页面组件
const Placeholder = ({ title }: { title: string }) => (
  <div className="flex items-center justify-center min-h-[400px]">
    <div className="text-center">
      <h2 className="text-2xl font-bold mb-2">{title}</h2>
      <p className="text-muted-foreground">开发中...</p>
    </div>
  </div>
)

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      
      <Route path="/" element={<MainLayout />}>
        <Route index element={<Home />} />
        <Route path="blogs" element={<BlogList />} />
        <Route path="blogs/:id" element={<Placeholder title="博客详情" />} />
        <Route path="blogs/new" element={<Placeholder title="写博客" />} />
        <Route path="profile" element={<Placeholder title="个人主页" />} />
        <Route path="friends" element={<Placeholder title="好友管理" />} />
        <Route path="groups" element={<Placeholder title="群组" />} />
        <Route path="messages" element={<Placeholder title="消息" />} />
        <Route path="settings" element={<Placeholder title="设置" />} />
      </Route>
    </Routes>
  )
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App
