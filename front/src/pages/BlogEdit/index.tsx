import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { blogApi, CreateBlogRequest } from '@/api/blog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { X, Save, Loader2 } from 'lucide-react'

export function BlogEdit() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEditing = !!id

  const [formData, setFormData] = useState<CreateBlogRequest>({
    title: '',
    content: '',
    summary: '',
    cover_image: '',
    tags: [],
    category: '',
    privacy: 1,
    status: 1,
  })
  const [tagInput, setTagInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (id) {
      loadBlog(id)
    }
  }, [id])

  const loadBlog = async (blogId: string) => {
    try {
      setLoading(true)
      const response = await blogApi.get(blogId)
      const blog = response.data.blog
      setFormData({
        title: blog.title,
        content: blog.content,
        summary: blog.summary,
        cover_image: blog.cover_image,
        tags: blog.tags || [],
        category: blog.category,
        privacy: blog.privacy,
        status: blog.status,
      })
    } catch (err) {
      console.error('Failed to load blog:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleAddTag = () => {
    if (tagInput.trim() && !formData.tags?.includes(tagInput.trim())) {
      setFormData((prev) => ({
        ...prev,
        tags: [...(prev.tags || []), tagInput.trim()],
      }))
      setTagInput('')
    }
  }

  const handleRemoveTag = (tagToRemove: string) => {
    setFormData((prev) => ({
      ...prev,
      tags: prev.tags?.filter((tag) => tag !== tagToRemove) || [],
    }))
  }

  const handleSubmit = async () => {
    if (!formData.title.trim() || !formData.content.trim()) {
      alert('标题和内容不能为空')
      return
    }

    try {
      setSaving(true)
      if (isEditing && id) {
        await blogApi.update(id, formData)
        navigate(`/blogs/${id}`)
      } else {
        const response = await blogApi.create(formData)
        navigate(`/blogs/${response.data.blog.id}`)
      }
    } catch (err) {
      console.error('Failed to save blog:', err)
      alert('保存失败')
    } finally {
      setSaving(false)
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

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">
          {isEditing ? '编辑博客' : '写博客'}
        </h1>
        <Button onClick={handleSubmit} disabled={saving}>
          {saving ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              保存中...
            </>
          ) : (
            <>
              <Save className="h-4 w-4 mr-2" />
              发布
            </>
          )}
        </Button>
      </div>

      <Card>
        <CardContent className="space-y-6 pt-6">
          {/* 标题 */}
          <div className="space-y-2">
            <Label htmlFor="title">标题</Label>
            <Input
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              placeholder="请输入博客标题"
            />
          </div>

          {/* 摘要 */}
          <div className="space-y-2">
            <Label htmlFor="summary">摘要</Label>
            <Textarea
              id="summary"
              name="summary"
              value={formData.summary}
              onChange={handleChange}
              placeholder="请输入博客摘要（可选）"
              rows={2}
            />
          </div>

          {/* 封面图 */}
          <div className="space-y-2">
            <Label htmlFor="cover_image">封面图片 URL</Label>
            <Input
              id="cover_image"
              name="cover_image"
              value={formData.cover_image}
              onChange={handleChange}
              placeholder="https://example.com/image.jpg"
            />
          </div>

          {/* 内容 */}
          <div className="space-y-2">
            <Label htmlFor="content">内容</Label>
            <Textarea
              id="content"
              name="content"
              value={formData.content}
              onChange={handleChange}
              placeholder="请输入博客内容"
              rows={15}
              className="font-mono"
            />
          </div>

          {/* 标签 */}
          <div className="space-y-2">
            <Label>标签</Label>
            <div className="flex gap-2">
              <Input
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                placeholder="添加标签，按回车确认"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    handleAddTag()
                  }
                }}
              />
              <Button type="button" onClick={handleAddTag} variant="secondary">
                添加
              </Button>
            </div>
            <div className="flex flex-wrap gap-2 mt-2">
              {formData.tags?.map((tag) => (
                <Badge key={tag} variant="secondary" className="gap-1">
                  {tag}
                  <button
                    onClick={() => handleRemoveTag(tag)}
                    className="ml-1 hover:text-destructive"
                  >
                    <X className="h-3 w-3" />
                  </button>
                </Badge>
              ))}
            </div>
          </div>

          {/* 分类 */}
          <div className="space-y-2">
            <Label htmlFor="category">分类</Label>
            <Input
              id="category"
              name="category"
              value={formData.category}
              onChange={handleChange}
              placeholder="请输入分类"
            />
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
