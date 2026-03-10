import apiClient from './client'

export enum FeedType {
  Text = 1,
  Blog = 2,
  Image = 3,
}

export interface MediaAttachment {
  type: 'image' | 'video'
  url: string
  thumbnail?: string
}

export interface FeedItem {
  id: string
  user_id: string
  type: FeedType
  content: string
  target_id?: string
  target_type?: string
  attachments: MediaAttachment[]
  privacy: number
  created_at: string
  updated_at: string
  likes_count: number
  comments_count: number
  user?: {
    id: string
    username: string
    nickname: string
    avatar: string
  }
}

export interface CreateFeedRequest {
  type: FeedType
  content: string
  target_id?: string
  target_type?: string
  attachments?: MediaAttachment[]
  privacy?: number
}

export interface FeedListResponse {
  feeds: FeedItem[]
  pagination: {
    page: number
    page_size: number
    total: number
  }
}

export const feedApi = {
  listHome: (page?: number, page_size?: number) =>
    apiClient.get<FeedListResponse>('/feeds', { params: { page, page_size } }),
  
  get: (id: string) =>
    apiClient.get<{ feed: FeedItem }>(`/feeds/${id}`),
  
  create: (data: CreateFeedRequest) =>
    apiClient.post<{ feed: FeedItem }>('/feeds', data),
  
  delete: (id: string) =>
    apiClient.delete(`/feeds/${id}`),
  
  like: (id: string, delta: 1 | -1 = 1) =>
    apiClient.post(`/feeds/${id}/like`, { delta }),
  
  listByUser: (userId: string, page?: number, page_size?: number) =>
    apiClient.get<FeedListResponse>(`/users/${userId}/feeds`, { params: { page, page_size } }),
}
