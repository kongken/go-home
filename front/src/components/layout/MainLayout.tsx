import { Outlet } from 'react-router-dom'
import { Header } from './Header'
import { Sidebar } from './Sidebar'

export function MainLayout() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="container py-6">
        <div className="grid gap-6 md:grid-cols-12">
          {/* 左侧边栏 - 只在 md 及以上屏幕显示 */}
          <div className="hidden md:block md:col-span-3 lg:col-span-3">
            <div className="sticky top-20">
              <Sidebar />
            </div>
          </div>
          
          {/* 中间内容区 */}
          <div className="col-span-12 md:col-span-9 lg:col-span-6">
            <Outlet />
          </div>
          
          {/* 右侧边栏 - 只在 lg 及以上屏幕显示 */}
          <div className="hidden lg:block lg:col-span-3">
            <div className="sticky top-20">
              {/* 右侧内容可以在这里添加 */}
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
