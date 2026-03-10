import { Blog } from '@/api/blog'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Heart, MessageCircle, Eye } from 'lucide-react'
import { formatDate, formatNumber } from '@/utils'
import { Link } from 'react-router-dom'

interface BlogCardProps {
  blog: Blog
  showAuthor?: boolean
}

export function BlogCard({ blog, showAuthor = true }: BlogCardProps) {
  return (
    <Card className="mb-4 hover:shadow-md transition-shadow">
      {showAuthor && (
        <CardHeader className="flex flex-row items-center gap-3 pb-3">
          <Avatar className="h-8 w-8">
            <AvatarImage src={blog.author?.avatar} />
            <AvatarFallback>{blog.author?.nickname?.[0]}</AvatarFallback>
          </Avatar>
          <div className="flex-1">
            <p className="font-medium text-sm">{blog.author?.nickname}</p>
            <p className="text-xs text-muted-foreground">
              {formatDate(blog.created_at)}
            </p>
          </div>
        </CardHeader>
      )}
      
      <CardContent className={showAuthor ? 'pb-3' : 'pt-6 pb-3'}>
        <Link to={`/blogs/${blog.id}`}>
          <h3 className="text-lg font-semibold mb-2 hover:text-primary transition-colors">
            {blog.title}
          </h3>
        </Link>
        <p className="text-sm text-muted-foreground line-clamp-3 mb-3">
          {blog.summary || blog.content.slice(0, 200)}
        </p>
        
        {blog.cover_image && (
          <div className="mb-3 rounded-lg overflow-hidden aspect-video">
            <img 
              src={blog.cover_image} 
              alt={blog.title}
              className="w-full h-full object-cover"
            />
          </div>
        )}
        
        {blog.tags && blog.tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {blog.tags.map((tag, index) => (
              <Badge key={index} variant="secondary" className="text-xs">
                {tag}
              </Badge>
            ))}
          </div>
        )}
      </CardContent>
      
      <CardFooter className="flex justify-between pt-0">
        <div className="flex items-center gap-4">
          <span className="flex items-center text-xs text-muted-foreground">
            <Eye className="h-3 w-3 mr-1" />
            {formatNumber(blog.views_count)}
          </span>
          <span className="flex items-center text-xs text-muted-foreground">
            <Heart className="h-3 w-3 mr-1" />
            {formatNumber(blog.likes_count)}
          </span>
          <span className="flex items-center text-xs text-muted-foreground">
            <MessageCircle className="h-3 w-3 mr-1" />
            {formatNumber(blog.comments_count)}
          </span>
        </div>
        <Button variant="ghost" size="sm" asChild>
          <Link to={`/blogs/${blog.id}`}>阅读更多</Link>
        </Button>
      </CardFooter>
    </Card>
  )
}
