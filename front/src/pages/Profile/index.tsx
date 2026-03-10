import { useState, useEffect } from 'react'
import { Link, useParams } from 'react-router-dom'
import { authApi, UserInfo } from '@/api/auth'
import { blogApi, Blog } from '@/api/blog'
import { feedApi, FeedItem } from '@/api/feed'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { BlogCard } from '@/components/common/BlogCard'
import { FeedCard } from '@/components/common/FeedCard'
import { Loader2, Edit } from 'lucide-react'

export function Profile() {
  const { id } = useParams<{ id: string }>()
  const { user: currentUser } = useAuth()
  const [user, setUser] = useState<UserInfo | null>(null)
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [feeds, setFeeds] = useState<FeedItem[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('feeds')

  const userId = id || currentUser?.id
  const isOwnProfile = !id || id === currentUser?.id

  useEffect(() => {
    if (userId) {
      loadUserData(userId)
    }
  }, [userId])

  const loadUserData = async (uid: string) => {
    try {
      setLoading(true)
      // 加载用户信息
      const userRes = await authApi.getUser(uid)
      setUser(userRes.data.user)

      // 加载博客
      const blogsRes = await blogApi.listByUser(uid, 1, 5)
      setBlogs(blogsRes.data.blogs)

      // 加载动态
      const feedsRes = await feedApi.listByUser(uid, 1, 10)
      setFeeds(feedsRes.data.feeds)
    } catch (err) {
      console.error('Failed to load user data:', err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <Card>
        <CardContent className="py-12 text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4" />
          <p className="text-muted-foreground">加载中...</p>
        </CardContent>
      </Card>
    )
  }

  if (!user) {
    return (
      <Card>
        <CardContent className="py-12 text-center">
          <p className="text-muted-foreground">用户不存在</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-4">
      {/* 封面和个人信息 */}
      <Card>
        {/* 封面背景 */}
        <div className="h-32 bg-gradient-to-r from-blue-500 to-purple-600 rounded-t-lg" />

        <CardContent className="pt-0">
          <div className="flex flex-col md:flex-row md:items-end -mt-12 mb-4">
            <Avatar className="h-24 w-24 border-4 border-background">
              <AvatarImage src={user.avatar} />
              <AvatarFallback className="text-2xl">
                {user.nickname?.[0]}
              </AvatarFallback>
            </Avatar>
            <div className="mt-4 md:mt-0 md:ml-4 flex-1">
              <h1 className="text-2xl font-bold">{user.nickname}</h1>
              <p className="text-muted-foreground">@{user.username}</p>
            </div>
            <div className="mt-4 md:mt-0">
              {isOwnProfile ? (
                <Button asChild>
                  <Link to="/settings">
                    <Edit className="h-4 w-4 mr-2" />
                    编辑资料
                  </Link>
                </Button>
              ) : (
                <Button>关注</Button>
              )}
            </div>
          </div>

          {/* 简介 */}
          {user.bio && (
            <p className="text-sm mb-4">{user.bio}</p>
          )}

          {/* 统计 */}
          <div className="flex gap-6 py-4 border-t">
            <div className="text-center">
              <p className="font-semibold text-lg">{blogs.length}</p>
              <p className="text-sm text-muted-foreground">博客</p>
            </div>
            <div className="text-center">
              <p className="font-semibold text-lg">0</p>
              <p className="text-sm text-muted-foreground">好友</p>
            </div>
            <div className="text-center">
              <p className="font-semibold text-lg">0</p>
              <p className="text-sm text-muted-foreground">关注</p>
            </div>
            <div className="text-center">
              <p className="font-semibold text-lg">0</p>
              <p className="text-sm text-muted-foreground">粉丝</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="feeds">动态</TabsTrigger>
          <TabsTrigger value="blogs">博客</TabsTrigger>
        </TabsList>

        <TabsContent value="feeds" className="space-y-4 mt-4">
          {feeds.length > 0 ? (
            feeds.map((feed) => <FeedCard key={feed.id} feed={feed} />)
          ) : (
            <Card>
              <CardContent className="py-12 text-center">
                <p className="text-muted-foreground">暂无动态</p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="blogs" className="space-y-4 mt-4">
          {blogs.length > 0 ? (
            blogs.map((blog) => <BlogCard key={blog.id} blog={blog} showAuthor={false} />)
          ) : (
            <Card>
              <CardContent className="py-12 text-center">
                <p className="text-muted-foreground mb-4">暂无博客</p>
                {isOwnProfile && (
                  <Button asChild>
                    <Link to="/blogs/new">写第一篇博客</Link>
                  </Button>
                )}
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
