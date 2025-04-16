package migrations

import (
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestCRUDOperations verifies that basic CRUD operations work after migrations
func TestCRUDOperations(t *testing.T) {
	// Set up in-memory database and run migrations
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	logger := zerolog.Nop()

	// Use RunConsolidatedMigrations
	if err := RunConsolidatedMigrations(db, &logger); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Test CRUD operations on User entity
	testUserCRUD(t, db)
}

// Test CRUD operations for the User entity
func testUserCRUD(t *testing.T, db *gorm.DB) {
	// Create
	user := &entity.UserEntity{
		ID:    "test-user-1",
		Email: "test@example.com",
		Name:  "Test User",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Read
	var readUser entity.UserEntity
	if err := db.First(&readUser, "id = ?", user.ID).Error; err != nil {
		t.Fatalf("failed to read user: %v", err)
	}
	if readUser.ID != user.ID || readUser.Email != user.Email {
		t.Errorf("read user doesn't match created user: %+v vs %+v", readUser, user)
	}

	// Update
	if err := db.Model(&entity.UserEntity{}).Where("id = ?", user.ID).Update("name", "Updated Name").Error; err != nil {
		t.Fatalf("failed to update user: %v", err)
	}
	var updatedUser entity.UserEntity
	if err := db.First(&updatedUser, "id = ?", user.ID).Error; err != nil {
		t.Fatalf("failed to read updated user: %v", err)
	}
	if updatedUser.Name != "Updated Name" {
		t.Errorf("update didn't work, name is still: %s", updatedUser.Name)
	}

	// Delete
	if err := db.Delete(&entity.UserEntity{}, "id = ?", user.ID).Error; err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	var count int64
	db.Model(&entity.UserEntity{}).Where("id = ?", user.ID).Count(&count)
	if count != 0 {
		t.Errorf("delete didn't work, user still exists")
	}
}
