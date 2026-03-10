import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'

export function Register() {
  const navigate = useNavigate()
  const { signUp } = useAuth()
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    
    if (formData.password !== formData.confirmPassword) {
      setError('两次输入的密码不一致')
      return
    }
    
    if (formData.password.length < 6) {
      setError('密码长度至少为 6 位')
      return
    }
    
    setLoading(true)
    
    try {
      await signUp(formData.username, formData.email, formData.password)
      navigate('/')
    } catch (err) {
      setError('注册失败，请检查输入信息')
    } finally {
      setLoading(false)
    }
  }
  
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>注册</CardTitle>
          <CardDescription>创建您的账号，开始使用 Go Home</CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && (
              <div className="text-sm text-destructive">{error}</div>
            )}
            <div className="space-y-2">
              <label className="text-sm font-medium">用户名</label>
              <Input
                name="username"
                value={formData.username}
                onChange={handleChange}
                placeholder="请输入用户名"
                required
                minLength={3}
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">邮箱</label>
              <Input
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                placeholder="请输入邮箱"
                required
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">密码</label>
              <Input
                name="password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                placeholder="请输入密码"
                required
                minLength={6}
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">确认密码</label>
              <Input
                name="confirmPassword"
                type="password"
                value={formData.confirmPassword}
                onChange={handleChange}
                placeholder="请再次输入密码"
                required
              />
            </div>
          </CardContent>
          <CardFooter className="flex flex-col space-y-4">
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? '注册中...' : '注册'}
            </Button>
            <p className="text-sm text-muted-foreground">
              已有账号？{' '}
              <Link to="/login" className="text-primary hover:underline">
                立即登录
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  )
}
