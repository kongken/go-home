package service

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

// SearchResult 搜索结果
type SearchResult struct {
	Users  []*model.User
	Blogs  []*model.Blog
	Albums []*model.Album
	Groups []*model.Group
}

// SearchService 搜索服务接口
type SearchService interface {
	GlobalSearch(ctx context.Context, keyword string, offset, limit int64) (*SearchResult, error)
	SearchUsers(ctx context.Context, keyword string, offset, limit int64) ([]*model.User, int64, error)
	SearchBlogs(ctx context.Context, keyword string, offset, limit int64) ([]*model.Blog, int64, error)
	SearchGroups(ctx context.Context, keyword string, offset, limit int64) ([]*model.Group, int64, error)
}

// searchService 搜索服务实现
type searchService struct {
	userRepo  repository.UserRepository
	blogRepo  repository.BlogRepository
	groupRepo repository.GroupRepository
}

// NewSearchService 创建搜索服务
func NewSearchService(userRepo repository.UserRepository, blogRepo repository.BlogRepository, groupRepo repository.GroupRepository) SearchService {
	return &searchService{
		userRepo:  userRepo,
		blogRepo:  blogRepo,
		groupRepo: groupRepo,
	}
}

// GlobalSearch 全局搜索
func (s *searchService) GlobalSearch(ctx context.Context, keyword string, offset, limit int64) (*SearchResult, error) {
	result := &SearchResult{}

	// 搜索用户
	users, _, _ := s.SearchUsers(ctx, keyword, 0, 5)
	result.Users = users

	// 搜索博客
	blogs, _, _ := s.SearchBlogs(ctx, keyword, 0, 5)
	result.Blogs = blogs

	// 搜索群组
	groups, _, _ := s.SearchGroups(ctx, keyword, 0, 5)
	result.Groups = groups

	return result, nil
}

// SearchUsers 搜索用户
func (s *searchService) SearchUsers(ctx context.Context, keyword string, offset, limit int64) ([]*model.User, int64, error) {
	// 这里简化处理，实际应该使用 MongoDB 的文本搜索
	// 暂时返回空结果
	return []*model.User{}, 0, nil
}

// SearchBlogs 搜索博客
func (s *searchService) SearchBlogs(ctx context.Context, keyword string, offset, limit int64) ([]*model.Blog, int64, error) {
	// 这里简化处理，实际应该使用 MongoDB 的文本搜索
	// 暂时返回空结果
	return []*model.Blog{}, 0, nil
}

// SearchGroups 搜索群组
func (s *searchService) SearchGroups(ctx context.Context, keyword string, offset, limit int64) ([]*model.Group, int64, error) {
	return s.groupRepo.Search(ctx, keyword, offset, limit)
}