import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MainLayout } from '@/components/layout/MainLayout'
import { Login } from '@/pages/Login'
import { Register } from '@/pages/Register'
import { Home } from '@/pages/Home'
import { BlogList } from '@/pages/Blog'
import { BlogDetail } from '@/pages/BlogDetail'
import { BlogEdit } from '@/pages/BlogEdit'
import { Profile } from '@/pages/Profile'
import { Friends } from '@/pages/Friends'
import { Groups } from '@/pages/Groups'
import { Messages } from '@/pages/Messages'
import { Settings } from '@/pages/Settings'

const queryClient = new QueryClient()

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      
      <Route path="/" element={<MainLayout />}>
        <Route index element={<Home />} />
        <Route path="blogs" element={<BlogList />} />
        <Route path="blogs/:id" element={<BlogDetail />} />
        <Route path="blogs/new" element={<BlogEdit />} />
        <Route path="blogs/:id/edit" element={<BlogEdit />} />
        <Route path="profile" element={<Profile />} />
        <Route path="users/:id" element={<Profile />} />
        <Route path="friends" element={<Friends />} />
        <Route path="groups" element={<Groups />} />
        <Route path="messages" element={<Messages />} />
        <Route path="settings" element={<Settings />} />
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