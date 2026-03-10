import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { blogApi, Blog } from '@/api/blog'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { BlogCard } from '@/components/common/BlogCard'
import { Plus, Loader2 } from 'lucide-react'

export function BlogList() {
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  
  const loadBlogs = async (pageNum: number = 1) => {
    try {
      setLoading(true)
      const response = await blogApi.list({ page: pageNum, page_size: 10 })
      if (pageNum === 1) {
        setBlogs(response.data.blogs)
      } else {
        setBlogs(prev => [...prev, ...response.data.blogs])
      }
      setHasMore(response.data.blogs.length === 10)
    } catch (error) {
      console.error('Failed to load blogs:', error)
    } finally {
      setLoading(false)
    }
  }
  
  useEffect(() => {
    loadBlogs()
  }, [])
  
  const handleLoadMore = () => {
    const nextPage = page + 1
    setPage(nextPage)
    loadBlogs(nextPage)
  }
  
  return (
    <div className="space-y-4">
      {/* 头部 */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">博客</h1>
        <Button asChild>
          <Link to="/blogs/new">
            <Plus className="h-4 w-4 mr-2" />
            写博客
          </Link>
        </Button>
      </div>
      
      {/* 博客列表 */}
      {loading && blogs.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-muted-foreground" />
            <p className="text-muted-foreground">加载中...</p>
          </CardContent>
        </Card>
      ) : (
        <>
          <div className="space-y-4">
            {blogs.map((blog) => (
              <BlogCard key={blog.id} blog={blog} />
            ))}
          </div>
          
          {hasMore && (
            <div className="text-center py-4">
              <Button 
                variant="outline" 
                onClick={handleLoadMore}
                disabled={loading}
              >
                {loading ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    加载中...
                  </>
                ) : (
                  '加载更多'
                )}
              </Button>
            </div>
          )}
          
          {blogs.length === 0 && (
            <Card>
              <CardContent className="py-12 text-center">
                <p className="text-muted-foreground mb-4">暂无博客</p>
                <Button asChild>
                  <Link to="/blogs/new">
                    <Plus className="h-4 w-4 mr-2" />
                    写第一篇博客
                  </Link>
                </Button>
              </CardContent>
            </Card>
          )}
        </>
      )}
    </div>
  )
}
