import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Search, Plus, Users } from 'lucide-react'

// 模拟数据
const mockGroups = [
  { id: '1', name: '技术交流', description: '讨论前端后端技术', memberCount: 128, avatar: '' },
  { id: '2', name: '生活分享', description: '分享日常生活', memberCount: 256, avatar: '' },
  { id: '3', name: '读书俱乐部', description: '一起读书交流', memberCount: 64, avatar: '' },
]

const mockJoinedGroups = [
  { id: '1', name: '技术交流', description: '讨论前端后端技术', memberCount: 128, avatar: '', role: 'member' },
]

export function Groups() {
  const [searchQuery, setSearchQuery] = useState('')
  const [activeTab, setActiveTab] = useState('joined')

  return (
    <div className="space-y-4">
      {/* 头部 */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">群组</h1>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          创建群组
        </Button>
      </div>

      {/* 搜索 */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="搜索群组..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="joined">我的群组</TabsTrigger>
          <TabsTrigger value="discover">发现</TabsTrigger>
        </TabsList>

        <TabsContent value="joined" className="space-y-2 mt-4">
          {mockJoinedGroups.map((group) => (
            <JoinedGroupCard key={group.id} group={group} />
          ))}
          {mockJoinedGroups.length === 0 && (
            <EmptyState message="还没有加入任何群组" />
          )}
        </TabsContent>

        <TabsContent value="discover" className="space-y-2 mt-4">
          {mockGroups.map((group) => (
            <GroupCard key={group.id} group={group} />
          ))}
        </TabsContent>
      </Tabs>
    </div>
  )
}

function GroupCard({ group }: { group: typeof mockGroups[0] }) {
  return (
    <Card>
      <CardContent className="flex items-center justify-between p-4">
        <div className="flex items-center gap-3 flex-1">
          <Avatar className="h-12 w-12">
            <AvatarImage src={group.avatar} />
            <AvatarFallback>{group.name[0]}</AvatarFallback>
          </Avatar>
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <p className="font-medium">{group.name}</p>
            </div>
            <p className="text-sm text-muted-foreground line-clamp-1">
              {group.description}
            </p>
            <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
              <Users className="h-3 w-3" />
              {group.memberCount} 成员
            </div>
          </div>
        </div>
        <Button size="sm">加入</Button>
      </CardContent>
    </Card>
  )
}

function JoinedGroupCard({ group }: { group: typeof mockJoinedGroups[0] }) {
  return (
    <Card>
      <CardContent className="flex items-center justify-between p-4">
        <Link to={`/groups/${group.id}`} className="flex items-center gap-3 flex-1">
          <Avatar className="h-12 w-12">
            <AvatarImage src={group.avatar} />
            <AvatarFallback>{group.name[0]}</AvatarFallback>
          </Avatar>
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <p className="font-medium">{group.name}</p>
              {group.role === 'admin' && (
                <Badge variant="secondary" className="text-xs">管理员</Badge>
              )}
            </div>
            <p className="text-sm text-muted-foreground line-clamp-1">
              {group.description}
            </p>
            <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
              <Users className="h-3 w-3" />
              {group.memberCount} 成员
            </div>
          </div>
        </Link>
        <Button variant="outline" size="sm">进入</Button>
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
