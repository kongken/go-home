import { useState } from 'react'
import { Link } from 'react-router-dom'
// import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Search, MessageSquare, Bell } from 'lucide-react'
import { formatDate } from '@/utils'

// 模拟数据
const mockConversations = [
  { 
    id: '1', 
    user: { id: '1', nickname: '张三', username: 'zhangsan', avatar: '' },
    lastMessage: '你好，最近怎么样？',
    lastTime: '2024-01-15T10:30:00Z',
    unread: 2,
  },
  { 
    id: '2', 
    user: { id: '2', nickname: '李四', username: 'lisi', avatar: '' },
    lastMessage: '明天有空吗？',
    lastTime: '2024-01-14T15:20:00Z',
    unread: 0,
  },
]

const mockNotifications = [
  { id: '1', type: 'like', content: '张三赞了你的博客', time: '2024-01-15T10:30:00Z', read: false },
  { id: '2', type: 'comment', content: '李四评论了你的动态', time: '2024-01-14T15:20:00Z', read: true },
  { id: '3', type: 'follow', content: '王五关注了你', time: '2024-01-13T09:00:00Z', read: true },
]

export function Messages() {
  const [searchQuery, setSearchQuery] = useState('')
  const [activeTab, setActiveTab] = useState('conversations')

  const unreadCount = mockNotifications.filter(n => !n.read).length

  return (
    <div className="space-y-4">
      {/* 头部 */}
      <h1 className="text-2xl font-bold">消息</h1>

      {/* 搜索 */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="搜索消息..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="conversations">
            <MessageSquare className="h-4 w-4 mr-2" />
            私信
          </TabsTrigger>
          <TabsTrigger value="notifications">
            <Bell className="h-4 w-4 mr-2" />
            通知
            {unreadCount > 0 && (
              <Badge variant="destructive" className="ml-2 h-5 w-5 p-0 text-xs">
                {unreadCount}
              </Badge>
            )}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="conversations" className="space-y-2 mt-4">
          {mockConversations.map((conversation) => (
            <ConversationCard key={conversation.id} conversation={conversation} />
          ))}
          {mockConversations.length === 0 && (
            <EmptyState message="暂无私信" />
          )}
        </TabsContent>

        <TabsContent value="notifications" className="space-y-2 mt-4">
          {mockNotifications.map((notification) => (
            <NotificationCard key={notification.id} notification={notification} />
          ))}
          {mockNotifications.length === 0 && (
            <EmptyState message="暂无通知" />
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}

function ConversationCard({ conversation }: { conversation: typeof mockConversations[0] }) {
  return (
    <Link to={`/messages/${conversation.user.id}`}>
      <Card className="hover:bg-accent/50 transition-colors">
        <CardContent className="flex items-center gap-3 p-4">
          <div className="relative">
            <Avatar className="h-12 w-12">
              <AvatarImage src={conversation.user.avatar} />
              <AvatarFallback>{conversation.user.nickname[0]}</AvatarFallback>
            </Avatar>
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between">
              <p className="font-medium truncate">{conversation.user.nickname}</p>
              <span className="text-xs text-muted-foreground">
                {formatDate(conversation.lastTime)}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <p className="text-sm text-muted-foreground truncate pr-4">
                {conversation.lastMessage}
              </p>
              {conversation.unread > 0 && (
                <Badge variant="default" className="h-5 w-5 p-0 text-xs flex items-center justify-center">
                  {conversation.unread}
                </Badge>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  )
}

function NotificationCard({ notification }: { notification: typeof mockNotifications[0] }) {
  return (
    <Card className={notification.read ? 'opacity-60' : ''}>
      <CardContent className="flex items-center gap-3 p-4">
        <div className="flex-1">
          <p className="text-sm">{notification.content}</p>
          <p className="text-xs text-muted-foreground mt-1">
            {formatDate(notification.time)}
          </p>
        </div>
        {!notification.read && (
          <Badge variant="default" className="h-2 w-2 p-0 rounded-full" />
        )}
      </CardContent>
    </Card>
  )
}

function EmptyState({ message }: { message: string }) {
  return (
    <Card>
      <CardContent className="py-12 text-center">
        <p className="text-muted-foreground">{message}</p>
      </CardContent>
    </Card>
  )
}
