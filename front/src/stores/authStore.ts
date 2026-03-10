import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { UserInfo } from '@/api/auth'

interface AuthState {
  isAuthenticated: boolean
  user: UserInfo | null
  accessToken: string | null
  refreshToken: string | null
  
  login: (accessToken: string, refreshToken: string, user: UserInfo) => void
  logout: () => void
  updateUser: (user: Partial<UserInfo>) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      isAuthenticated: false,
      user: null,
      accessToken: null,
      refreshToken: null,
      
      login: (accessToken, refreshToken, user) => {
        localStorage.setItem('access_token', accessToken)
        localStorage.setItem('refresh_token', refreshToken)
        set({
          isAuthenticated: true,
          user,
          accessToken,
          refreshToken,
        })
      },
      
      logout: () => {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        set({
          isAuthenticated: false,
          user: null,
          accessToken: null,
          refreshToken: null,
        })
      },
      
      updateUser: (updates) => {
        set((state) => ({
          user: state.user ? { ...state.user, ...updates } : null,
        }))
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        user: state.user,
      }),
    }
  )
)
