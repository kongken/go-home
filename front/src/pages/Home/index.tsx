import { useState, useEffect } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { feedApi, FeedItem, FeedType } from '@/api/feed'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import { FeedCard } from '@/components/common/FeedCard'
import { Image as ImageIcon } from 'lucide-react'

export function Home() {
  const { user } = useAuth()
  const [feeds, setFeeds] = useState<FeedItem[]>([])
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  
  const loadFeeds = async (pageNum: number = 1) => {
    try {
      const response = await feedApi.listHome(pageNum, 20)
      if (pageNum === 1) {
        setFeeds(response.data.feeds)
      } else {
        setFeeds(prev => [...prev, ...response.data.feeds])
      }
      setHasMore(response.data.feeds.length === 20)
    } catch (error) {
      console.error('Failed to load feeds:', error)
    }
  }
  
  useEffect(() => {
    loadFeeds()
  }, [])
  
  const handlePublish = async () => {
    if (!content.trim()) return
    setLoading(true)
    try {
      await feedApi.create({
        type: FeedType.Text,
        content: content.trim(),
      })
      setContent('')
      loadFeeds(1)
    } catch (error) {
      console.error('Failed to publish feed:', error)
    } finally {
      setLoading(false)
    }
  }
  
  const handleLike = async (id: string) => {
    try {
      await feedApi.like(id, 1)
      loadFeeds(1)
    } catch (error) {
      console.error('Failed to like:', error)
    }
  }
  
  const handleLoadMore = () => {
    const nextPage = page + 1
    setPage(nextPage)
    loadFeeds(nextPage)
  }
  
  return (
    <div className="space-y-4">
      {/* 发布框 */}
      {user && (
        <Card>
          <CardContent className="pt-4">
            <Textarea
              placeholder="分享你的想法..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              className="min-h-[100px] resize-none"
            />
            <div className="flex justify-between items-center mt-3">
              <Button variant="ghost" size="sm" className="text-muted-foreground">
                <ImageIcon className="h-4 w-4 mr-2" />
                添加图片
              </Button>
              <Button 
                onClick={handlePublish} 
                disabled={loading || !content.trim()}
                size="sm"
              >
                {loading ? '发布中...' : '发布'}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
      
      {/* 动态列表 */}
      <div className="space-y-4">
        {feeds.map((feed) => (
          <FeedCard 
            key={feed.id} 
            feed={feed} 
            onLike={handleLike}
          />
        ))}
      </div>
      
      {/* 加载更多 */}
      {hasMore && feeds.length > 0 && (
        <div className="text-center py-4">
          <Button variant="outline" onClick={handleLoadMore}>
            加载更多
          </Button>
        </div>
      )}
      
      {feeds.length === 0 && (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">暂无动态，来发布第一条吧！</p>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
