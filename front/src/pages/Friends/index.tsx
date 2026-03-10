import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Search, UserPlus, MessageCircle, UserMinus } from 'lucide-react'

// 模拟数据
const mockFriends = [
  { id: '1', nickname: '张三', username: 'zhangsan', avatar: '', status: 'online' },
  { id: '2', nickname: '李四', username: 'lisi', avatar: '', status: 'offline' },
  { id: '3', nickname: '王五', username: 'wangwu', avatar: '', status: 'online' },
]

const mockRequests = [
  { id: '1', nickname: '赵六', username: 'zhaoliu', avatar: '', type: 'received' },
]

export function Friends() {
  const [searchQuery, setSearchQuery] = useState('')
  const [activeTab, setActiveTab] = useState('all')

  return (
    <div className="space-y-4">
      {/* 头部 */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">好友</h1>
        <Button>
          <UserPlus className="h-4 w-4 mr-2" />
          添加好友
        </Button>
      </div>

      {/* 搜索 */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="搜索好友..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="all">全部</TabsTrigger>
          <TabsTrigger value="online">在线</TabsTrigger>
          <TabsTrigger value="requests">
            请求
            {mockRequests.length > 0 && (
              <Badge variant="destructive" className="ml-2 h-5 w-5 p-0 text-xs">
                {mockRequests.length}
              </Badge>
            )}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="space-y-2 mt-4">
          {mockFriends.map((friend) => (
            <FriendCard key={friend.id} friend={friend} />
          ))}
          {mockFriends.length === 0 && (
            <EmptyState message="暂无好友" />
          )}
        </TabsContent>

        <TabsContent value="online" className="space-y-2 mt-4">
          {mockFriends
            .filter((f) => f.status === 'online')
            .map((friend) => (
              <FriendCard key={friend.id} friend={friend} />
            ))}
        </TabsContent>

        <TabsContent value="requests" className="space-y-2 mt-4">
          {mockRequests.map((request) => (
            <RequestCard key={request.id} request={request} />
          ))}
          {mockRequests.length === 0 && (
            <EmptyState message="暂无好友请求" />
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}

function FriendCard({ friend }: { friend: typeof mockFriends[0] }) {
  return (
    <Card>
      <CardContent className="flex items-center justify-between p-4">
        <Link to={`/users/${friend.id}`} className="flex items-center gap-3 flex-1">
          <div className="relative">
            <Avatar className="h-10 w-10">
              <AvatarImage src={friend.avatar} />
              <AvatarFallback>{friend.nickname[0]}</AvatarFallback>
            </Avatar>
            {friend.status === 'online' && (
              <span className="absolute bottom-0 right-0 h-3 w-3 rounded-full bg-green-500 border-2 border-background" />
            )}
          </div>
          <div>
            <p className="font-medium">{friend.nickname}</p>
            <p className="text-sm text-muted-foreground">@{friend.username}</p>
          </div>
        </Link>
        <div className="flex gap-2">
          <Button variant="ghost" size="icon">
            <MessageCircle className="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="icon">
            <UserMinus className="h-4 w-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

function RequestCard({ request }: { request: typeof mockRequests[0] }) {
  return (
    <Card>
      <CardContent className="flex items-center justify-between p-4">
        <Link to={`/users/${request.id}`} className="flex items-center gap-3 flex-1">
          <Avatar className="h-10 w-10">
            <AvatarImage src={request.avatar} />
            <AvatarFallback>{request.nickname[0]}</AvatarFallback>
          </Avatar>
          <div>
            <p className="font-medium">{request.nickname}</p>
            <p className="text-sm text-muted-foreground">@{request.username}</p>
          </div>
        </Link>
        <div className="flex gap-2">
          <Button size="sm">接受</Button>
          <Button variant="outline" size="sm">拒绝</Button>
        </div>
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
