package repository

import (
	"context"
	"testing"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/testutil"
)

func TestBlogRepository_Integration(t *testing.T) {
	testutil.SkipIfShort(t)
	testutil.SkipIfNoMongoDB(t)

	ctx := context.Background()
	testDB := testutil.NewTestDB(t)
	repo := NewBlogRepositoryWithClient(testDB.Client, testDB.DBName)

	t.Run("Create", func(t *testing.T) {
		blog := &model.Blog{
			UserID:   testutil.GenerateTestID(),
			Title:    "Test Blog Title",
			Content:  "This is the content of the test blog.",
			Summary:  "Test summary",
			Category: "Technology",
			Tags:     "go,test,mongodb",
			Privacy:  model.PrivacyPublic,
			Status:   model.BlogStatusPublished,
		}

		err := repo.Create(ctx, blog)
		if err != nil {
			t.Fatalf("Failed to create blog: %v", err)
		}

		if blog.ID == "" {
			t.Fatal("Expected ID to be set after insert")
		}

		t.Logf("Created blog with ID: %s", blog.ID)
	})

	t.Run("GetByID", func(t *testing.T) {
		// First create a blog
		blog := &model.Blog{
			UserID:   testutil.GenerateTestID(),
			Title:    "Test Blog GetByID",
			Content:  "Content for GetByID test",
			Category: "Technology",
			Privacy:  model.PrivacyPublic,
			Status:   model.BlogStatusPublished,
		}

		err := repo.Create(ctx, blog)
		if err != nil {
			t.Fatalf("Failed to create blog: %v", err)
		}

		// Then get by ID
		found, err := repo.GetByID(ctx, blog.ID)
		if err != nil {
			t.Fatalf("Failed to get blog by ID: %v", err)
		}

		if found.ID != blog.ID {
			t.Errorf("Expected ID %s, got %s", blog.ID, found.ID)
		}

		if found.Title != blog.Title {
			t.Errorf("Expected Title %s, got %s", blog.Title, found.Title)
		}

		if found.Content != blog.Content {
			t.Errorf("Expected Content %s, got %s", blog.Content, found.Content)
		}
	})

	t.Run("Update", func(t *testing.T) {
		// First create a blog
		blog := &model.Blog{
			UserID:   testutil.GenerateTestID(),
			Title:    "Test Blog Update",
			Content:  "Original content",
			Category: "Technology",
			Privacy:  model.PrivacyPublic,
			Status:   model.BlogStatusPublished,
		}

		err := repo.Create(ctx, blog)
		if err != nil {
			t.Fatalf("Failed to create blog: %v", err)
		}

		// Update blog
		blog.Title = "Updated Title"
		blog.Content = "Updated content"
		blog.Category = "Programming"
		err = repo.Update(ctx, blog)
		if err != nil {
			t.Fatalf("Failed to update blog: %v", err)
		}

		// Verify update
		found, err := repo.GetByID(ctx, blog.ID)
		if err != nil {
			t.Fatalf("Failed to get blog after update: %v", err)
		}

		if found.Title != "Updated Title" {
			t.Errorf("Expected Title 'Updated Title', got %s", found.Title)
		}

		if found.Content != "Updated content" {
			t.Errorf("Expected Content 'Updated content', got %s", found.Content)
		}

		if found.Category != "Programming" {
			t.Errorf("Expected Category 'Programming', got %s", found.Category)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// First create a blog
		blog := &model.Blog{
			UserID:   testutil.GenerateTestID(),
			Title:    "Test Blog Delete",
			Content:  "Content for delete test",
			Category: "Technology",
			Privacy:  model.PrivacyPublic,
			Status:   model.BlogStatusPublished,
		}

		err := repo.Create(ctx, blog)
		if err != nil {
			t.Fatalf("Failed to create blog: %v", err)
		}

		// Delete blog
		err = repo.Delete(ctx, blog.ID)
		if err != nil {
			t.Fatalf("Failed to delete blog: %v", err)
		}

		// Verify deletion
		_, err = repo.GetByID(ctx, blog.ID)
		if err == nil {
			t.Fatal("Expected error when getting deleted blog")
		}
	})

	t.Run("List", func(t *testing.T) {
		// Clean up first
		DeleteAllDocuments(ctx, repo)

		userID := testutil.GenerateTestID()

		// Create multiple blogs
		for i := 0; i < 15; i++ {
			blog := &model.Blog{
				UserID:   userID,
				Title:    "Test Blog List",
				Content:  "Content for list test",
				Category: "Technology",
				Privacy:  model.PrivacyPublic,
				Status:   model.BlogStatusPublished,
			}
			err := repo.Create(ctx, blog)
			if err != nil {
				t.Fatalf("Failed to create blog: %v", err)
			}
		}

		// List first page
		blogs, total, err := repo.List(ctx, "", "", 0, 10)
		if err != nil {
			t.Fatalf("Failed to list blogs: %v", err)
		}

		if total < 10 {
			t.Errorf("Expected total >= 10, got %d", total)
		}

		if len(blogs) != 10 {
			t.Errorf("Expected 10 blogs, got %d", len(blogs))
		}
	})

	t.Run("ListByUser", func(t *testing.T) {
		// Clean up first
		DeleteAllDocuments(ctx, repo)

		userID := testutil.GenerateTestID()
		otherUserID := testutil.GenerateTestID()

		// Create blogs for user
		for i := 0; i < 5; i++ {
			blog := &model.Blog{
				UserID:   userID,
				Title:    "User Blog",
				Content:  "Content",
				Category: "Technology",
				Privacy:  model.PrivacyPublic,
				Status:   model.BlogStatusPublished,
			}
			err := repo.Create(ctx, blog)
			if err != nil {
				t.Fatalf("Failed to create blog: %v", err)
			}
		}

		// Create blogs for other user
		for i := 0; i < 3; i++ {
			blog := &model.Blog{
				UserID:   otherUserID,
				Title:    "Other User Blog",
				Content:  "Content",
				Category: "Technology",
				Privacy:  model.PrivacyPublic,
				Status:   model.BlogStatusPublished,
			}
			err := repo.Create(ctx, blog)
			if err != nil {
				t.Fatalf("Failed to create blog: %v", err)
			}
		}

		// List by user
		blogs, total, err := repo.ListByUser(ctx, userID, 0, 10)
		if err != nil {
			t.Fatalf("Failed to list blogs by user: %v", err)
		}

		if total != 5 {
			t.Errorf("Expected total 5, got %d", total)
		}

		if len(blogs) != 5 {
			t.Errorf("Expected 5 blogs, got %d", len(blogs))
		}

		// Verify all blogs belong to the user
		for _, blog := range blogs {
			if blog.UserID != userID {
				t.Errorf("Expected UserID %s, got %s", userID, blog.UserID)
			}
		}
	})

	t.Run("ListWithCategory", func(t *testing.T) {
		// Clean up first
		DeleteAllDocuments(ctx, repo)

		userID := testutil.GenerateTestID()

		// Create blogs with different categories
		categories := []string{"Technology", "Technology", "Programming", "Life"}
		for _, cat := range categories {
			blog := &model.Blog{
				UserID:   userID,
				Title:    "Blog with category",
				Content:  "Content",
				Category: cat,
				Privacy:  model.PrivacyPublic,
				Status:   model.BlogStatusPublished,
			}
			err := repo.Create(ctx, blog)
			if err != nil {
				t.Fatalf("Failed to create blog: %v", err)
			}
		}

		// List by category
		blogs, total, err := repo.List(ctx, "", "Technology", 0, 10)
		if err != nil {
			t.Fatalf("Failed to list blogs by category: %v", err)
		}

		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}

		if len(blogs) != 2 {
			t.Errorf("Expected 2 blogs, got %d", len(blogs))
		}

		// Verify all blogs have the correct category
		for _, blog := range blogs {
			if blog.Category != "Technology" {
				t.Errorf("Expected Category 'Technology', got %s", blog.Category)
			}
		}
	})

	t.Run("DraftBlogsNotListed", func(t *testing.T) {
		// Clean up first
		DeleteAllDocuments(ctx, repo)

		userID := testutil.GenerateTestID()

		// Create published blog
		publishedBlog := &model.Blog{
			UserID:   userID,
			Title:    "Published Blog",
			Content:  "Content",
			Category: "Technology",
			Privacy:  model.PrivacyPublic,
			Status:   model.BlogStatusPublished,
		}
		err := repo.Create(ctx, publishedBlog)
		if err != nil {
			t.Fatalf("Failed to create blog: %v", err)
		}

		// Create draft blog
		draftBlog := &model.Blog{
			UserID:   userID,
			Title:    "Draft Blog",
			Content:  "Content",
			Category: "Technology",
			Privacy:  model.PrivacyPublic,
			Status:   model.BlogStatusDraft,
		}
		err = repo.Create(ctx, draftBlog)
		if err != nil {
			t.Fatalf("Failed to create blog: %v", err)
		}

		// List blogs
		blogs, total, err := repo.List(ctx, "", "", 0, 10)
		if err != nil {
			t.Fatalf("Failed to list blogs: %v", err)
		}

		// Should only get published blog
		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}

		if len(blogs) != 1 {
			t.Errorf("Expected 1 blog, got %d", len(blogs))
		}

		if blogs[0].Status != model.BlogStatusPublished {
			t.Errorf("Expected published blog, got status %d", blogs[0].Status)
		}
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "nonexistent_id")
		if err == nil {
			t.Fatal("Expected error for non-existent blog")
		}
	})
}

func TestBlog_BeforeInsert(t *testing.T) {
	blog := &model.Blog{
		UserID:  "user123",
		Title:   "Test Blog",
		Content: "Content",
	}

	blog.BeforeInsert(nil)

	if blog.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if blog.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if blog.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestBlog_Statuses(t *testing.T) {
	tests := []struct {
		name   string
		status model.BlogStatus
	}{
		{"Draft", model.BlogStatusDraft},
		{"Published", model.BlogStatusPublished},
		{"Deleted", model.BlogStatusDeleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blog := &model.Blog{
				UserID: "user123",
				Title:  "Test Blog",
				Status: tt.status,
			}

			if blog.Status != tt.status {
				t.Errorf("Expected Status %d, got %d", tt.status, blog.Status)
			}
		})
	}
}

func TestBlog_PrivacyLevels(t *testing.T) {
	tests := []struct {
		name   string
		privacy model.PrivacyLevel
	}{
		{"Public", model.PrivacyPublic},
		{"Friends", model.PrivacyFriends},
		{"Private", model.PrivacyPrivate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blog := &model.Blog{
				UserID:  "user123",
				Title:   "Test Blog",
				Privacy: tt.privacy,
			}

			if blog.Privacy != tt.privacy {
				t.Errorf("Expected Privacy %d, got %d", tt.privacy, blog.Privacy)
			}
		})
	}
}

func TestBlog_WithStats(t *testing.T) {
	testutil.SkipIfShort(t)
	testutil.SkipIfNoMongoDB(t)

	ctx := context.Background()
	testDB := testutil.NewTestDB(t)
	repo := NewBlogRepositoryWithClient(testDB.Client, testDB.DBName)

	blog := &model.Blog{
		UserID:         testutil.GenerateTestID(),
		Title:          "Test Blog With Stats",
		Content:        "Content",
		Category:       "Technology",
		Privacy:        model.PrivacyPublic,
		Status:         model.BlogStatusPublished,
		ViewsCount:     100,
		LikesCount:     50,
		CommentsCount:  25,
		FavoritesCount: 10,
	}

	err := repo.Create(ctx, blog)
	if err != nil {
		t.Fatalf("Failed to create blog: %v", err)
	}

	// Verify stats
	found, err := repo.GetByID(ctx, blog.ID)
	if err != nil {
		t.Fatalf("Failed to get blog: %v", err)
	}

	if found.ViewsCount != 100 {
		t.Errorf("Expected ViewsCount 100, got %d", found.ViewsCount)
	}

	if found.LikesCount != 50 {
		t.Errorf("Expected LikesCount 50, got %d", found.LikesCount)
	}

	if found.CommentsCount != 25 {
		t.Errorf("Expected CommentsCount 25, got %d", found.CommentsCount)
	}

	if found.FavoritesCount != 10 {
		t.Errorf("Expected FavoritesCount 10, got %d", found.FavoritesCount)
	}
}

// Helper to satisfy interface check
var _ BlogRepository = (*BlogRepositoryButterfly)(nil)