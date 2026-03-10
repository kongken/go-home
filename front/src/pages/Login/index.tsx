import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'

export function Login() {
  const navigate = useNavigate()
  const { signIn } = useAuth()
  const [account, setAccount] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    
    try {
      await signIn(account, password)
      navigate('/')
    } catch (err) {
      setError('登录失败，请检查账号和密码')
    } finally {
      setLoading(false)
    }
  }
  
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>登录</CardTitle>
          <CardDescription>欢迎回来，请登录您的账号</CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && (
              <div className="text-sm text-destructive">{error}</div>
            )}
            <div className="space-y-2">
              <label className="text-sm font-medium">账号/邮箱</label>
              <Input
                type="text"
                value={account}
                onChange={(e) => setAccount(e.target.value)}
                placeholder="请输入账号或邮箱"
                required
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">密码</label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="请输入密码"
                required
              />
            </div>
          </CardContent>
          <CardFooter className="flex flex-col space-y-4">
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? '登录中...' : '登录'}
            </Button>
            <p className="text-sm text-muted-foreground">
              还没有账号？{' '}
              <Link to="/register" className="text-primary hover:underline">
                立即注册
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  )
}
