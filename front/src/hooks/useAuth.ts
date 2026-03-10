import { useAuthStore } from '@/stores/authStore'
import { authApi, UserInfo } from '@/api'

export function useAuth() {
  const { isAuthenticated, user, login, logout, updateUser } = useAuthStore()
  
  const signIn = async (account: string, password: string) => {
    const response = await authApi.login({ account, password })
    const { access_token, refresh_token, user } = response.data
    login(access_token, refresh_token, user)
    return user
  }
  
  const signUp = async (username: string, email: string, password: string) => {
    const response = await authApi.register({ username, email, password })
    const { access_token, refresh_token, user } = response.data
    login(access_token, refresh_token, user)
    return user
  }
  
  const signOut = () => {
    logout()
  }
  
  const updateProfile = async (updates: Partial<UserInfo>) => {
    const response = await authApi.updateUser(updates)
    updateUser(response.data.user)
    return response.data.user
  }
  
  return {
    isAuthenticated,
    user,
    signIn,
    signUp,
    signOut,
    updateProfile,
  }
}
