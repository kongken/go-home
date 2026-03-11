package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/testutil"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUserRepository_Integration(t *testing.T) {
	testutil.SkipIfShort(t)
	testutil.SkipIfNoMongoDB(t)

	ctx := context.Background()
	testDB := testutil.NewTestDB(t)
	repo := NewUserRepositoryWithClient(testDB.Client, testDB.DBName)

	t.Run("Create", func(t *testing.T) {
		user := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashedpassword",
			Nickname: "Test User",
			Status:   model.UserStatusNormal,
		}

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if user.ID == "" {
			t.Fatal("Expected ID to be set after insert")
		}

		if user.CreatedAt.IsZero() {
			t.Fatal("Expected CreatedAt to be set after insert")
		}

		t.Logf("Created user with ID: %s", user.ID)
	})

	t.Run("GetByID", func(t *testing.T) {
		// First create a user
		user := &model.User{
			Username: "testuser_getbyid",
			Email:    "getbyid@example.com",
			Password: "hashedpassword",
			Nickname: "Test User GetByID",
			Status:   model.UserStatusNormal,
		}

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Then get by ID
		found, err := repo.GetByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to get user by ID: %v", err)
		}

		if found.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, found.ID)
		}

		if found.Username != user.Username {
			t.Errorf("Expected Username %s, got %s", user.Username, found.Username)
		}

		if found.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, found.Email)
		}
	})

	t.Run("GetByUsername", func(t *testing.T) {
		// First create a user
		user := &model.User{
			Username: "testuser_getbyusername",
			Email:    "getbyusername@example.com",
			Password: "hashedpassword",
			Nickname: "Test User GetByUsername",
			Status:   model.UserStatusNormal,
		}

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Then get by username
		found, err := repo.GetByUsername(ctx, user.Username)
		if err != nil {
			t.Fatalf("Failed to get user by username: %v", err)
		}

		if found.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, found.ID)
		}

		if found.Username != user.Username {
			t.Errorf("Expected Username %s, got %s", user.Username, found.Username)
		}
	})

	t.Run("GetByEmail", func(t *testing.T) {
		// First create a user
		user := &model.User{
			Username: "testuser_getbyemail",
			Email:    "getbyemail@example.com",
			Password: "hashedpassword",
			Nickname: "Test User GetByEmail",
			Status:   model.UserStatusNormal,
		}

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Then get by email
		found, err := repo.GetByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("Failed to get user by email: %v", err)
		}

		if found.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, found.ID)
		}

		if found.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, found.Email)
		}
	})

	t.Run("Update", func(t *testing.T) {
		// First create a user
		user := &model.User{
			Username: "testuser_update",
			Email:    "update@example.com",
			Password: "hashedpassword",
			Nickname: "Test User Update",
			Status:   model.UserStatusNormal,
		}

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Update user
		user.Nickname = "Updated Nickname"
		user.Bio = "This is my bio"
		user.Location = "Beijing, China"
		err = repo.Update(ctx, user)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify update
		found, err := repo.GetByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to get user after update: %v", err)
		}

		if found.Nickname != "Updated Nickname" {
			t.Errorf("Expected Nickname 'Updated Nickname', got %s", found.Nickname)
		}

		if found.Bio != "This is my bio" {
			t.Errorf("Expected Bio 'This is my bio', got %s", found.Bio)
		}

		if found.Location != "Beijing, China" {
			t.Errorf("Expected Location 'Beijing, China', got %s", found.Location)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// First create a user
		user := &model.User{
			Username: "testuser_delete",
			Email:    "delete@example.com",
			Password: "hashedpassword",
			Nickname: "Test User Delete",
			Status:   model.UserStatusNormal,
		}

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Delete user
		err = repo.Delete(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify deletion
		_, err = repo.GetByID(ctx, user.ID)
		if err == nil {
			t.Fatal("Expected error when getting deleted user")
		}
	})

	t.Run("List", func(t *testing.T) {
		// Clean up first
		DeleteAllDocuments(ctx, repo)

		// Create multiple users
		for i := 0; i < 15; i++ {
			user := &model.User{
				Username: "testuser_list_" + string(rune('a'+i)),
				Email:    string(rune('a'+i)) + "list@example.com",
				Password: "hashedpassword",
				Nickname: "Test User List " + string(rune('a'+i)),
				Status:   model.UserStatusNormal,
			}
			err := repo.Create(ctx, user)
			if err != nil {
				t.Fatalf("Failed to create user: %v", err)
			}
		}

		// List first page
		users, total, err := repo.List(ctx, 0, 10)
		if err != nil {
			t.Fatalf("Failed to list users: %v", err)
		}

		if total < 10 {
			t.Errorf("Expected total >= 10, got %d", total)
		}

		if len(users) != 10 {
			t.Errorf("Expected 10 users, got %d", len(users))
		}

		// List second page
		users, _, err = repo.List(ctx, 10, 10)
		if err != nil {
			t.Fatalf("Failed to list users second page: %v", err)
		}

		if len(users) < 5 {
			t.Errorf("Expected at least 5 users on second page, got %d", len(users))
		}
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "nonexistent_id")
		if err == nil {
			t.Fatal("Expected error for non-existent user")
		}
	})

	t.Run("GetByUsername_NotFound", func(t *testing.T) {
		_, err := repo.GetByUsername(ctx, "nonexistent_username")
		if err == nil {
			t.Fatal("Expected error for non-existent username")
		}
	})

	t.Run("GetByEmail_NotFound", func(t *testing.T) {
		_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		if err == nil {
			t.Fatal("Expected error for non-existent email")
		}
	})
}

func TestUser_BeforeInsert(t *testing.T) {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	user.BeforeInsert(nil)

	if user.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	// Test that existing ID is not overwritten
	existingID := "existing_id_123"
	user2 := &model.User{
		ID:       existingID,
		Username: "testuser2",
		Email:    "test2@example.com",
	}

	user2.BeforeInsert(nil)

	if user2.ID != existingID {
		t.Errorf("Expected ID to remain %s, got %s", existingID, user2.ID)
	}
}

func TestUser_Genders(t *testing.T) {
	tests := []struct {
		name   string
		gender model.Gender
	}{
		{"Unknown", model.GenderUnknown},
		{"Male", model.GenderMale},
		{"Female", model.GenderFemale},
		{"Other", model.GenderOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &model.User{
				Username: "testuser_gender_" + tt.name,
				Email:    "gender_" + tt.name + "@example.com",
				Gender:   tt.gender,
			}

			if user.Gender != tt.gender {
				t.Errorf("Expected Gender %d, got %d", tt.gender, user.Gender)
			}
		})
	}
}

func TestUser_Statuses(t *testing.T) {
	tests := []struct {
		name   string
		status model.UserStatus
	}{
		{"Normal", model.UserStatusNormal},
		{"Frozen", model.UserStatusFrozen},
		{"Banned", model.UserStatusBanned},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &model.User{
				Username: "testuser_status_" + tt.name,
				Email:    "status_" + tt.name + "@example.com",
				Status:   tt.status,
			}

			if user.Status != tt.status {
				t.Errorf("Expected Status %d, got %d", tt.status, user.Status)
			}
		})
	}
}

func TestUser_WithBirthday(t *testing.T) {
	testutil.SkipIfShort(t)
	testutil.SkipIfNoMongoDB(t)

	ctx := context.Background()
	testDB := testutil.NewTestDB(t)
	repo := NewUserRepositoryWithClient(testDB.Client, testDB.DBName)

	birthday := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	user := &model.User{
		Username: "testuser_birthday",
		Email:    "birthday@example.com",
		Password: "hashedpassword",
		Nickname: "Test User Birthday",
		Birthday: &birthday,
		Status:   model.UserStatusNormal,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify birthday was saved
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if found.Birthday == nil {
		t.Fatal("Expected Birthday to be set")
	}

	// Compare dates without time component
	if !found.Birthday.Equal(birthday) {
		t.Errorf("Expected Birthday %v, got %v", birthday, found.Birthday)
	}
}

func TestUser_WithStats(t *testing.T) {
	testutil.SkipIfShort(t)
	testutil.SkipIfNoMongoDB(t)

	ctx := context.Background()
	testDB := testutil.NewTestDB(t)
	repo := NewUserRepositoryWithClient(testDB.Client, testDB.DBName)

	user := &model.User{
		Username:      "testuser_stats",
		Email:         "stats@example.com",
		Password:      "hashedpassword",
		Nickname:      "Test User Stats",
		FriendsCount:  10,
		FollowersCount: 20,
		FollowingCount: 15,
		BlogsCount:    5,
		AlbumsCount:   3,
		SharesCount:   2,
		Status:        model.UserStatusNormal,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify stats
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if found.FriendsCount != 10 {
		t.Errorf("Expected FriendsCount 10, got %d", found.FriendsCount)
	}

	if found.FollowersCount != 20 {
		t.Errorf("Expected FollowersCount 20, got %d", found.FollowersCount)
	}

	if found.FollowingCount != 15 {
		t.Errorf("Expected FollowingCount 15, got %d", found.FollowingCount)
	}

	if found.BlogsCount != 5 {
		t.Errorf("Expected BlogsCount 5, got %d", found.BlogsCount)
	}

	if found.AlbumsCount != 3 {
		t.Errorf("Expected AlbumsCount 3, got %d", found.AlbumsCount)
	}

	if found.SharesCount != 2 {
		t.Errorf("Expected SharesCount 2, got %d", found.SharesCount)
	}
}

// Benchmark tests
func BenchmarkUserRepository_Create(b *testing.B) {
	testutil.SkipIfShort(&testing.T{})
	testutil.SkipIfNoMongoDB(&testing.T{})

	ctx := context.Background()
	testDB := testutil.NewTestDB(&testing.T{})
	repo := NewUserRepositoryWithClient(testDB.Client, testDB.DBName)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &model.User{
			Username: "benchuser_" + string(rune(i)),
			Email:    "bench" + string(rune(i)) + "@example.com",
			Password: "hashedpassword",
			Nickname: "Benchmark User",
			Status:   model.UserStatusNormal,
		}
		repo.Create(ctx, user)
	}
}

func BenchmarkUserRepository_GetByID(b *testing.B) {
	testutil.SkipIfShort(&testing.T{})
	testutil.SkipIfNoMongoDB(&testing.T{})

	ctx := context.Background()
	testDB := testutil.NewTestDB(&testing.T{})
	repo := NewUserRepositoryWithClient(testDB.Client, testDB.DBName)

	// Create a test user
	user := &model.User{
		Username: "benchuser_get",
		Email:    "benchget@example.com",
		Password: "hashedpassword",
		Nickname: "Benchmark User",
		Status:   model.UserStatusNormal,
	}
	repo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByID(ctx, user.ID)
	}
}

// Helper to satisfy interface check
var _ UserRepository = (*UserRepositoryButterfly)(nil)

// Additional test for BSON serialization
func TestUser_BSONSerialization(t *testing.T) {
	user := &model.User{
		ID:       bson.NewObjectID().Hex(),
		Username: "bson_test",
		Email:    "bson@example.com",
		Password: "hashedpassword",
		Nickname: "BSON Test",
		Gender:   model.GenderMale,
		Status:   model.UserStatusNormal,
	}

	data, err := bson.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	var unmarshaledUser model.User
	err = bson.Unmarshal(data, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	if unmarshaledUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, unmarshaledUser.ID)
	}

	if unmarshaledUser.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, unmarshaledUser.Username)
	}

	if unmarshaledUser.Gender != user.Gender {
		t.Errorf("Expected Gender %d, got %d", user.Gender, unmarshaledUser.Gender)
	}
}