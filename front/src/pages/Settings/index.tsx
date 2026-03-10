import { useState } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'
import { Camera, Lock, Bell, Shield, User } from 'lucide-react'

export function Settings() {
  const { user, updateProfile } = useAuth()
  const [profileForm, setProfileForm] = useState({
    nickname: user?.nickname || '',
    bio: user?.bio || '',
    avatar: user?.avatar || '',
  })
  const [saving, setSaving] = useState(false)

  const handleSaveProfile = async () => {
    try {
      setSaving(true)
      await updateProfile(profileForm)
      alert('保存成功')
    } catch (err) {
      alert('保存失败')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">设置</h1>

      <Tabs defaultValue="profile" className="space-y-4">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="profile">
            <User className="h-4 w-4 mr-2" />
            个人资料
          </TabsTrigger>
          <TabsTrigger value="account">
            <Lock className="h-4 w-4 mr-2" />
            账号
          </TabsTrigger>
          <TabsTrigger value="notifications">
            <Bell className="h-4 w-4 mr-2" />
            通知
          </TabsTrigger>
          <TabsTrigger value="privacy">
            <Shield className="h-4 w-4 mr-2" />
            隐私
          </TabsTrigger>
        </TabsList>

        {/* 个人资料 */}
        <TabsContent value="profile" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>个人资料</CardTitle>
              <CardDescription>更新您的个人信息和头像</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* 头像 */}
              <div className="flex items-center gap-4">
                <Avatar className="h-20 w-20">
                  <AvatarImage src={profileForm.avatar} />
                  <AvatarFallback className="text-2xl">
                    {profileForm.nickname[0]}
                  </AvatarFallback>
                </Avatar>
                <Button variant="outline">
                  <Camera className="h-4 w-4 mr-2" />
                  更换头像
                </Button>
              </div>

              <Separator />

              {/* 昵称 */}
              <div className="space-y-2">
                <Label htmlFor="nickname">昵称</Label>
                <Input
                  id="nickname"
                  value={profileForm.nickname}
                  onChange={(e) =>
                    setProfileForm({ ...profileForm, nickname: e.target.value })
                  }
                />
              </div>

              {/* 简介 */}
              <div className="space-y-2">
                <Label htmlFor="bio">简介</Label>
                <Textarea
                  id="bio"
                  value={profileForm.bio}
                  onChange={(e) =>
                    setProfileForm({ ...profileForm, bio: e.target.value })
                  }
                  placeholder="介绍一下自己..."
                  rows={4}
                />
              </div>

              <Button onClick={handleSaveProfile} disabled={saving}>
                {saving ? '保存中...' : '保存更改'}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        {/* 账号 */}
        <TabsContent value="account" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>修改密码</CardTitle>
              <CardDescription>更新您的登录密码</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="current">当前密码</Label>
                <Input id="current" type="password" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="new">新密码</Label>
                <Input id="new" type="password" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="confirm">确认新密码</Label>
                <Input id="confirm" type="password" />
              </div>
              <Button>修改密码</Button>
            </CardContent>
          </Card>
        </TabsContent>

        {/* 通知 */}
        <TabsContent value="notifications" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>通知设置</CardTitle>
              <CardDescription>选择您希望接收的通知类型</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">邮件通知</p>
                  <p className="text-sm text-muted-foreground">
                    接收重要更新的邮件通知
                  </p>
                </div>
                <Switch defaultChecked />
              </div>
              <Separator />
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">评论通知</p>
                  <p className="text-sm text-muted-foreground">
                    有人评论您的内容时通知
                  </p>
                </div>
                <Switch defaultChecked />
              </div>
              <Separator />
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">点赞通知</p>
                  <p className="text-sm text-muted-foreground">
                    有人点赞您的内容时通知
                  </p>
                </div>
                <Switch defaultChecked />
              </div>
              <Separator />
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">关注通知</p>
                  <p className="text-sm text-muted-foreground">
                    有人关注您时通知
                  </p>
                </div>
                <Switch />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* 隐私 */}
        <TabsContent value="privacy" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>隐私设置</CardTitle>
              <CardDescription>控制您的内容可见性</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">公开动态</p>
                  <p className="text-sm text-muted-foreground">
                    允许所有人查看您的动态
                  </p>
                </div>
                <Switch defaultChecked />
              </div>
              <Separator />
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">公开博客</p>
                  <p className="text-sm text-muted-foreground">
                    允许所有人查看您的博客
                  </p>
                </div>
                <Switch defaultChecked />
              </div>
              <Separator />
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">好友验证</p>
                  <p className="text-sm text-muted-foreground">
                    添加好友需要验证
                  </p>
                </div>
                <Switch defaultChecked />
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
