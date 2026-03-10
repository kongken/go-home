import apiClient from './client'
import { UserInfo } from './auth'

export interface Blog {
  id: string
  user_id: string
  title: string
  content: string
  summary: string
  cover_image: string
  tags: string[]
  category: string
  privacy: number
  status: number
  created_at: string
  updated_at: string
  views_count: number
  likes_count: number
  comments_count: number
  author?: UserInfo
}

export interface CreateBlogRequest {
  title: string
  content: string
  summary?: string
  cover_image?: string
  tags?: string[]
  category?: string
  privacy?: number
  status?: number
}

export interface BlogListResponse {
  blogs: Blog[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export const blogApi = {
  list: (params?: { page?: number; page_size?: number; category?: string; user_id?: string }) =>
    apiClient.get<BlogListResponse>('/blogs', { params }),
  
  get: (id: string) =>
    apiClient.get<{ blog: Blog }>(`/blogs/${id}`),
  
  create: (data: CreateBlogRequest) =>
    apiClient.post<{ blog: Blog }>('/blogs', data),
  
  update: (id: string, data: Partial<CreateBlogRequest>) =>
    apiClient.put<{ blog: Blog }>(`/blogs/${id}`, data),
  
  delete: (id: string) =>
    apiClient.delete(`/blogs/${id}`),
  
  listByUser: (userId: string, page?: number, page_size?: number) =>
    apiClient.get<BlogListResponse>(`/users/${userId}/blogs`, { params: { page, page_size } }),
}
