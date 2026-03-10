import apiClient from './client'

export interface LoginRequest {
  account: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

export interface UserInfo {
  id: string
  username: string
  nickname: string
  avatar: string
  bio?: string
  email?: string
}

export const authApi = {
  login: (data: LoginRequest) => 
    apiClient.post<AuthResponse>('/auth/login', data),
  
  register: (data: RegisterRequest) => 
    apiClient.post<AuthResponse>('/auth/register', data),
  
  refresh: (refreshToken: string) => 
    apiClient.post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken }),
  
  getUser: (id: string) => 
    apiClient.get<{ user: UserInfo }>(`/users/${id}`),
  
  updateUser: (updates: Partial<UserInfo>) => 
    apiClient.put<{ user: UserInfo }>('/users/me', updates),
}
