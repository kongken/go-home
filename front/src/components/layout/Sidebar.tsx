import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
// import { Separator } from '@/components/ui/separator'
import { 
  Home, 
  BookOpen, 
  Users, 
  MessageSquare, 
  Settings,
  Image,
  UserPlus
} from 'lucide-react'

const menuItems = [
  { icon: Home, label: '首页', path: '/' },
  { icon: BookOpen, label: '我的博客', path: '/blogs' },
  { icon: Users, label: '好友', path: '/friends' },
  { icon: MessageSquare, label: '消息', path: '/messages' },
  { icon: Image, label: '相册', path: '/albums' },
  { icon: Settings, label: '设置', path: '/settings' },
]

export function Sidebar() {
  const { user } = useAuth()
  
  if (!user) return null
  
  return (
    <div className="space-y-4">
      {/* 用户信息卡片 */}
      <Card>
        <CardHeader className="text-center pb-2">
          <Avatar className="h-20 w-20 mx-auto">
            <AvatarImage src={user.avatar} />
            <AvatarFallback className="text-2xl">{user.nickname?.[0]}</AvatarFallback>
          </Avatar>
        </CardHeader>
        <CardContent className="text-center space-y-2">
          <h3 className="font-semibold text-lg">{user.nickname}</h3>
          <p className="text-sm text-muted-foreground">@{user.username}</p>
          {user.bio && (
            <p className="text-sm text-muted-foreground line-clamp-2">{user.bio}</p>
          )}
          <div className="flex justify-center gap-4 pt-2">
            <div className="text-center">
              <p className="font-semibold">0</p>
              <p className="text-xs text-muted-foreground">博客</p>
            </div>
            <div className="text-center">
              <p className="font-semibold">0</p>
              <p className="text-xs text-muted-foreground">好友</p>
            </div>
            <div className="text-center">
              <p className="font-semibold">0</p>
              <p className="text-xs text-muted-foreground">关注</p>
            </div>
          </div>
        </CardContent>
      </Card>
      
      {/* 快捷导航 */}
      <Card>
        <CardContent className="p-2">
          <nav className="space-y-1">
            {menuItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </Link>
            ))}
          </nav>
        </CardContent>
      </Card>
      
      {/* 推荐用户 */}
      <Card>
        <CardHeader className="pb-3">
          <h3 className="font-semibold text-sm">推荐关注</h3>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Avatar className="h-8 w-8">
                <AvatarFallback>U1</AvatarFallback>
              </Avatar>
              <div>
                <p className="text-sm font-medium">用户1</p>
                <p className="text-xs text-muted-foreground">@user1</p>
              </div>
            </div>
            <Button size="icon" variant="ghost" className="h-8 w-8">
              <UserPlus className="h-4 w-4" />
            </Button>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Avatar className="h-8 w-8">
                <AvatarFallback>U2</AvatarFallback>
              </Avatar>
              <div>
                <p className="text-sm font-medium">用户2</p>
                <p className="text-xs text-muted-foreground">@user2</p>
              </div>
            </div>
            <Button size="icon" variant="ghost" className="h-8 w-8">
              <UserPlus className="h-4 w-4" />
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
