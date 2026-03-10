import { FeedItem } from '@/api/feed'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Heart, MessageCircle, Share2 } from 'lucide-react'
import { formatDate, formatNumber } from '@/utils'

interface FeedCardProps {
  feed: FeedItem
  onLike?: (id: string) => void
  onComment?: (id: string) => void
  onShare?: (id: string) => void
}

export function FeedCard({ feed, onLike, onComment, onShare }: FeedCardProps) {
  return (
    <Card className="mb-4">
      <CardHeader className="flex flex-row items-center gap-3 pb-3">
        <Avatar className="h-10 w-10">
          <AvatarImage src={feed.user?.avatar} />
          <AvatarFallback>{feed.user?.nickname?.[0]}</AvatarFallback>
        </Avatar>
        <div className="flex-1">
          <p className="font-medium text-sm">{feed.user?.nickname}</p>
          <p className="text-xs text-muted-foreground">
            @{feed.user?.username} · {formatDate(feed.created_at)}
          </p>
        </div>
      </CardHeader>
      
      <CardContent className="pb-3">
        <p className="whitespace-pre-wrap text-sm">{feed.content}</p>
        {feed.attachments && feed.attachments.length > 0 && (
          <div className="mt-3 grid grid-cols-3 gap-2">
            {feed.attachments.map((attachment, index) => (
              <div key={index} className="aspect-square rounded-lg overflow-hidden bg-muted">
                <img 
                  src={attachment.url} 
                  alt="attachment" 
                  className="w-full h-full object-cover"
                />
              </div>
            ))}
          </div>
        )}
      </CardContent>
      
      <CardFooter className="flex justify-between pt-0">
        <Button 
          variant="ghost" 
          size="sm" 
          onClick={() => onLike?.(feed.id)}
          className="text-muted-foreground hover:text-red-500"
        >
          <Heart className="h-4 w-4 mr-1" />
          {formatNumber(feed.likes_count)}
        </Button>
        <Button 
          variant="ghost" 
          size="sm" 
          onClick={() => onComment?.(feed.id)}
          className="text-muted-foreground"
        >
          <MessageCircle className="h-4 w-4 mr-1" />
          {formatNumber(feed.comments_count)}
        </Button>
        <Button 
          variant="ghost" 
          size="sm" 
          onClick={() => onShare?.(feed.id)}
          className="text-muted-foreground"
        >
          <Share2 className="h-4 w-4 mr-1" />
          分享
        </Button>
      </CardFooter>
    </Card>
  )
}
