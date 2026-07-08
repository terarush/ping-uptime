package service

import (
	"context"
	"errors"
	"testing"

	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/users/domain/entity"
	"ping-uptime/modules/users/domain/repository"
)

type fakeUserRepository struct {
	allUsers       []*entity.User
	usersByID      map[uint]*entity.User
	usersByEmail   map[string]*entity.User
	findAllErr     error
	findByIDErr    error
	findByEmailErr error
	createErr      error
	updateErr      error
	createdUser    *entity.User
	updatedUser    *entity.User
}

func (f *fakeUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	return f.allUsers, f.findAllErr
}

func (f *fakeUserRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.usersByID[id], nil
}

func (f *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if f.findByEmailErr != nil {
		return nil, f.findByEmailErr
	}
	if user, ok := f.usersByEmail[email]; ok {
		return user, nil
	}
	return nil, repository.ERR_RECORD_NOT_FOUND
}

func (f *fakeUserRepository) Create(ctx context.Context, user *entity.User) error {
	f.createdUser = user
	return f.createErr
}

func (f *fakeUserRepository) Update(ctx context.Context, user *entity.User) error {
	f.updatedUser = user
	return f.updateErr
}

func (f *fakeUserRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

func TestCreateUserHashesPasswordAndAssignsAdmin(t *testing.T) {
	repo := &fakeUserRepository{
		usersByEmail: map[string]*entity.User{},
	}
	service := NewAuthService(repo)
	user := entity.NewUser("Admin", "admin@example.com", "secret123")

	err := service.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser returned an error: %v", err)
	}

	if repo.createdUser == nil {
		t.Fatal("expected repository Create to be called")
	}
	if repo.createdUser.Role != "admin" {
		t.Fatalf("expected first user to receive admin role, got %q", repo.createdUser.Role)
	}
	if repo.createdUser.Password == "secret123" {
		t.Fatal("expected password to be hashed before create")
	}
	if !utils.CompareHashAndPassword(repo.createdUser.Password, "secret123") {
		t.Fatal("expected hashed password to match original input")
	}
}

func TestCreateUserRejectsWhenSetupAlreadyCompleted(t *testing.T) {
	repo := &fakeUserRepository{
		allUsers: []*entity.User{{ID: 1}},
	}
	service := NewAuthService(repo)

	err := service.CreateUser(context.Background(), entity.NewUser("Admin", "admin@example.com", "secret123"))
	if !errors.Is(err, ErrRegistrationDisabled) {
		t.Fatalf("expected ErrRegistrationDisabled, got %v", err)
	}
}

func TestProcessLoginRejectsInvalidPassword(t *testing.T) {
	hashedPassword, err := utils.HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	repo := &fakeUserRepository{
		usersByEmail: map[string]*entity.User{
			"user@example.com": {
				ID:       1,
				Email:    "user@example.com",
				Password: hashedPassword,
			},
		},
	}
	service := NewAuthService(repo)

	_, err = service.ProcessLogin(context.Background(), "user@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestValidateUserExistsRejectsBlockedUsers(t *testing.T) {
	repo := &fakeUserRepository{
		usersByID: map[uint]*entity.User{
			7: {ID: 7, IsBlocked: true},
		},
	}
	service := NewAuthService(repo)

	err := service.ValidateUserExists(context.Background(), 7)
	if err == nil || err.Error() != "user is blocked" {
		t.Fatalf("expected blocked user error, got %v", err)
	}
}

func TestRequestPasswordResetStoresTokenAndExpiry(t *testing.T) {
	user := &entity.User{
		ID:    1,
		Email: "user@example.com",
	}
	repo := &fakeUserRepository{
		usersByEmail: map[string]*entity.User{
			user.Email: user,
		},
	}
	service := NewAuthService(repo)

	returnedUser, token, err := service.RequestPasswordReset(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("RequestPasswordReset returned an error: %v", err)
	}

	if returnedUser != user {
		t.Fatal("expected RequestPasswordReset to return the existing user")
	}
	if token == "" {
		t.Fatal("expected reset token to be generated")
	}
	if repo.updatedUser == nil {
		t.Fatal("expected repository Update to be called")
	}
	if repo.updatedUser.ResetToken == nil || *repo.updatedUser.ResetToken != token {
		t.Fatal("expected reset token to be stored on the user")
	}
	if repo.updatedUser.ResetTokenExpiry == nil {
		t.Fatal("expected reset token expiry to be stored on the user")
	}
}
