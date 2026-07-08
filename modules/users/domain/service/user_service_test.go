package service

import (
	"context"
	"errors"
	"testing"

	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/users/domain/entity"
)

type fakeUserRepo struct {
	usersByID      map[uint]*entity.User
	usersByEmail   map[string]*entity.User
	findByIDErr    error
	findByEmailErr error
	createErr      error
	updateErr      error
	deleteErr      error
	createdUser    *entity.User
	updatedUser    *entity.User
	deletedID      uint
}

func (f *fakeUserRepo) FindAll(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (f *fakeUserRepo) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.usersByID[id], nil
}

func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if f.findByEmailErr != nil {
		return nil, f.findByEmailErr
	}
	return f.usersByEmail[email], nil
}

func (f *fakeUserRepo) Create(ctx context.Context, user *entity.User) error {
	f.createdUser = user
	return f.createErr
}

func (f *fakeUserRepo) Update(ctx context.Context, user *entity.User) error {
	f.updatedUser = user
	return f.updateErr
}

func (f *fakeUserRepo) Delete(ctx context.Context, id uint) error {
	f.deletedID = id
	return f.deleteErr
}

func TestCreateUserHashesPasswordBeforePersisting(t *testing.T) {
	repo := &fakeUserRepo{
		usersByEmail: map[string]*entity.User{},
	}
	service := NewUserService(repo)
	user := entity.NewUser("User", "user@example.com", "secret123")

	err := service.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	if repo.createdUser == nil {
		t.Fatal("expected repository Create to be called")
	}
	if repo.createdUser.Password == "secret123" {
		t.Fatal("expected password to be hashed before create")
	}
	if !utils.CompareHashAndPassword(repo.createdUser.Password, "secret123") {
		t.Fatal("expected hashed password to match original password")
	}
}

func TestCreateUserRejectsDuplicateEmail(t *testing.T) {
	repo := &fakeUserRepo{
		usersByEmail: map[string]*entity.User{
			"user@example.com": {ID: 1, Email: "user@example.com"},
		},
	}
	service := NewUserService(repo)

	err := service.CreateUser(context.Background(), entity.NewUser("User", "user@example.com", "secret123"))
	if !errors.Is(err, ErrEmailAlreadyUsed) {
		t.Fatalf("expected ErrEmailAlreadyUsed, got %v", err)
	}
}

func TestGetUserByIDReturnsErrUserNotFoundForMissingUser(t *testing.T) {
	repo := &fakeUserRepo{
		usersByID: map[uint]*entity.User{},
	}
	service := NewUserService(repo)

	_, err := service.GetUserByID(context.Background(), 99)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdateUserHashesChangedPassword(t *testing.T) {
	existingPassword, err := utils.HashPassword("old-password")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	repo := &fakeUserRepo{
		usersByID: map[uint]*entity.User{
			7: {ID: 7, Password: existingPassword},
		},
	}
	service := NewUserService(repo)
	user := &entity.User{
		ID:       7,
		Email:    "user@example.com",
		Password: "new-password",
	}

	err = service.UpdateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("UpdateUser returned an error: %v", err)
	}

	if repo.updatedUser == nil {
		t.Fatal("expected repository Update to be called")
	}
	if repo.updatedUser.Password == "new-password" {
		t.Fatal("expected changed password to be hashed before update")
	}
	if !utils.CompareHashAndPassword(repo.updatedUser.Password, "new-password") {
		t.Fatal("expected hashed password to match updated password")
	}
}

func TestDeleteUserRejectsMissingUser(t *testing.T) {
	repo := &fakeUserRepo{
		usersByID: map[uint]*entity.User{},
	}
	service := NewUserService(repo)

	err := service.DeleteUser(context.Background(), 42)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
