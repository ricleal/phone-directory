package service

import (
	"context"
	"fmt"

	"phone-directory/internal/entities"
	"phone-directory/internal/repository"
	"phone-directory/internal/store"
)

func validateName(name string) error {
	if name == "" {
		return entities.ErrInvalidName
	}
	return nil
}

// UserService is a domain service for users.
type UserService interface {
	Create(ctx context.Context, u *entities.User) error
	Get(ctx context.Context, id uint) (*entities.User, error)
	Update(ctx context.Context, u *entities.User) error
	Delete(ctx context.Context, id uint) error
}

// userService is an implementation of the UserService interface.
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(s store.Store) *userService {
	repo := s.Users()
	return &userService{repo}
}

// Create creates a new user.
func (s *userService) Create(ctx context.Context, u *entities.User) error {
	if err := validateName(u.Name); err != nil {
		return entities.ErrInvalidName
	}

	user := &repository.User{
		Name: u.Name,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	u.ID = user.ID
	return nil
}

// Get returns a user by id.
func (s *userService) Get(ctx context.Context, id uint) (*entities.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	entUser := &entities.User{
		ID:   user.ID,
		Name: user.Name,
	}

	for _, p := range user.Phones {
		entUser.Phones = append(entUser.Phones, entities.Phone{
			ID:     p.ID,
			Number: p.Number,
		})
	}

	for _, a := range user.Addresses {
		entUser.Addresses = append(entUser.Addresses, entities.Address{
			ID:      a.ID,
			Address: a.Address,
		})
	}

	return entUser, nil
}

// Update updates a user.
func (s *userService) Update(ctx context.Context, u *entities.User) error {
	if err := validateName(u.Name); err != nil {
		return entities.ErrInvalidName
	}

	user := &repository.User{
		Name: u.Name,
	}

	return s.repo.Update(ctx, user)
}

// Delete deletes a user.
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
