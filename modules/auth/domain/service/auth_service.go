package service

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/jwt"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/users/domain/entity"
	"ping-uptime/modules/users/domain/repository"
	"sync"
)

// Errors
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyUsed      = errors.New("email already in use")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrRegistrationDisabled  = errors.New("registration is disabled")
)

// AuthService handles user authentication
type AuthService struct {
	userRepo repository.UserRepository
	jwt      jwt.JWT
	mu       sync.Mutex
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}
	return &AuthService{
		userRepo: userRepo,
	}
}

// IsSetupNeeded checks if there are no registered users in the database
func (s *AuthService) IsSetupNeeded(ctx context.Context) (bool, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return false, err
	}
	return len(users) == 0, nil
}

// CreateUser creates a new user (only permitted during initial setup)
func (s *AuthService) CreateUser(ctx context.Context, user *entity.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user.Email == "" || user.Password == "" {
		return errors.New("email and password cannot be empty")
	}

	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return ErrRegistrationDisabled
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && err != repository.ERR_RECORD_NOT_FOUND {
		return err
	}
	if existingUser != nil {
		return ErrEmailAlreadyUsed
	}

	// Hash the password before saving the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.Role = "admin" // First user is automatically Admin

	return s.userRepo.Create(ctx, user)
}

// ProcessLogin handles user login and password verification
func (s *AuthService) ProcessLogin(ctx context.Context, email, password string) (*entity.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password cannot be empty")
	}

	// Find user by email
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Compare the provided password with the hashed password in the database
	if !utils.CompareHashAndPassword(existingUser.Password, password) {
		return nil, ErrInvalidPassword
	}

	// Return the authenticated user
	return existingUser, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint, password string) (*entity.User, error) {
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Password = hashedPassword

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.New("failed to update password")
	}

	return user, nil
}
