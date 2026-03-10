import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { blogApi, Blog } from '@/api/blog'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Heart, MessageCircle, Eye, Edit, Trash2, ArrowLeft } from 'lucide-react'
import { formatDate, formatNumber } from '@/utils'

export function BlogDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const [blog, setBlog] = useState<Blog | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (id) {
      loadBlog(id)
    }
  }, [id])

  const loadBlog = async (blogId: string) => {
    try {
      setLoading(true)
      const response = await blogApi.get(blogId)
      setBlog(response.data.blog)
    } catch (err) {
      setError('加载博客失败')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!id || !blog) return
    if (!confirm('确定要删除这篇博客吗？')) return

    try {
      await blogApi.delete(id)
      navigate('/blogs')
    } catch (err) {
      setError('删除失败')
    }
  }

  const isOwner = user?.id === blog?.user_id

  if (loading) {
    return (
      <Card>
        <CardContent className="py-12 text-center">
          <p className="text-muted-foreground">加载中...</p>
        </CardContent>
      </Card>
    )
  }

  if (error || !blog) {
    return (
      <Card>
        <CardContent className="py-12 text-center">
          <p className="text-destructive mb-4">{error || '博客不存在'}</p>
          <Button asChild>
            <Link to="/blogs">返回博客列表</Link>
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-4">
      {/* 返回按钮 */}
      <Button variant="ghost" size="sm" asChild>
        <Link to="/blogs">
          <ArrowLeft className="h-4 w-4 mr-2" />
          返回列表
        </Link>
      </Button>

      <Card>
        {/* 头部 */}
        <CardHeader className="pb-4">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <Avatar className="h-10 w-10">
                <AvatarImage src={blog.author?.avatar} />
                <AvatarFallback>{blog.author?.nickname?.[0]}</AvatarFallback>
              </Avatar>
              <div>
                <p className="font-medium">{blog.author?.nickname}</p>
                <p className="text-sm text-muted-foreground">
                  {formatDate(blog.created_at)}
                </p>
              </div>
            </div>
            {isOwner && (
              <div className="flex gap-2">
                <Button variant="ghost" size="sm" asChild>
                  <Link to={`/blogs/${blog.id}/edit`}>
                    <Edit className="h-4 w-4 mr-1" />
                    编辑
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" onClick={handleDelete}>
                  <Trash2 className="h-4 w-4 mr-1" />
                  删除
                </Button>
              </div>
            )}
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          {/* 标题 */}
          <h1 className="text-2xl font-bold">{blog.title}</h1>

          {/* 标签 */}
          {blog.tags && blog.tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {blog.tags.map((tag, index) => (
                <Badge key={index} variant="secondary">
                  {tag}
                </Badge>
              ))}
            </div>
          )}

          {/* 封面图 */}
          {blog.cover_image && (
            <div className="rounded-lg overflow-hidden">
              <img
                src={blog.cover_image}
                alt={blog.title}
                className="w-full max-h-[400px] object-cover"
              />
            </div>
          )}

          {/* 内容 */}
          <div className="prose prose-sm max-w-none">
            <div className="whitespace-pre-wrap">{blog.content}</div>
          </div>

          <Separator />

          {/* 统计 */}
          <div className="flex items-center gap-6">
            <span className="flex items-center text-sm text-muted-foreground">
              <Eye className="h-4 w-4 mr-1" />
              {formatNumber(blog.views_count)} 阅读
            </span>
            <span className="flex items-center text-sm text-muted-foreground">
              <Heart className="h-4 w-4 mr-1" />
              {formatNumber(blog.likes_count)} 点赞
            </span>
            <span className="flex items-center text-sm text-muted-foreground">
              <MessageCircle className="h-4 w-4 mr-1" />
              {formatNumber(blog.comments_count)} 评论
            </span>
          </div>
        </CardContent>
      </Card>

      {/* 评论区域 */}
      <Card>
        <CardHeader>
          <h3 className="font-semibold">评论</h3>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-center py-8">
            评论功能开发中...
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
