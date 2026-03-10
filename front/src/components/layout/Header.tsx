import { Link, useNavigate } from 'react-router-dom'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/hooks/useAuth'
import { Search, Bell, MessageSquare, Home, Book, Users, Settings } from 'lucide-react'

export function Header() {
  const navigate = useNavigate()
  const { user, signOut } = useAuth()
  
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center">
        <div className="mr-4 flex">
          <Link to="/" className="mr-6 flex items-center space-x-2">
            <span className="font-bold">Go Home</span>
          </Link>
          <nav className="flex items-center space-x-6 text-sm font-medium">
            <Link to="/" className="transition-colors hover:text-foreground/80 flex items-center gap-1">
              <Home className="h-4 w-4" />
              首页
            </Link>
            <Link to="/blogs" className="transition-colors hover:text-foreground/80 flex items-center gap-1">
              <Book className="h-4 w-4" />
              博客
            </Link>
            <Link to="/groups" className="transition-colors hover:text-foreground/80 flex items-center gap-1">
              <Users className="h-4 w-4" />
              群组
            </Link>
          </nav>
        </div>
        
        <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
          <div className="w-full flex-1 md:w-auto md:flex-none">
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                type="search"
                placeholder="搜索..."
                className="pl-8 md:w-[200px] lg:w-[300px]"
              />
            </div>
          </div>
          
          {user ? (
            <div className="flex items-center gap-4">
              <Button variant="ghost" size="icon" className="relative">
                <MessageSquare className="h-5 w-5" />
              </Button>
              <Button variant="ghost" size="icon" className="relative">
                <Bell className="h-5 w-5" />
              </Button>
              <Link to="/profile">
                <Avatar className="h-8 w-8 cursor-pointer">
                  <AvatarImage src={user.avatar} alt={user.nickname} />
                  <AvatarFallback>{user.nickname[0]?.toUpperCase()}</AvatarFallback>
                </Avatar>
              </Link>
              <Button variant="ghost" size="icon" onClick={() => navigate('/settings')}>
                <Settings className="h-5 w-5" />
              </Button>
              <Button variant="ghost" size="sm" onClick={signOut}>
                退出
              </Button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Button variant="ghost" size="sm" onClick={() => navigate('/login')}>
                登录
              </Button>
              <Button size="sm" onClick={() => navigate('/register')}>
                注册
              </Button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}
