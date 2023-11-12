package service

import (
	"context"
	"fmt"
	"regexp"

	"phone-directory/internal/entities"
	"phone-directory/internal/repository"
	"phone-directory/internal/store"
)

// Create a regular expression to match valid phone numbers.
var rePhoneNumber = regexp.MustCompile(`^(\+\d{1,2}\s?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)

// validatePhoneNumber validates a phone number.
func validatePhoneNumber(phoneNumber string) bool {
	// Return true if the email address matches the regular expression.
	return rePhoneNumber.MatchString(phoneNumber)
}

//go:generate mockgen -source=phones.go -package=mock -destination=mock/phones.go

// PhoneService is a domain service for phones.
type PhoneService interface {
	Create(ctx context.Context, u *entities.Phone) error
	Get(ctx context.Context, id uint) (*entities.Phone, error)
	Update(ctx context.Context, u *entities.Phone) error
	Delete(ctx context.Context, id uint) error
}

// phoneService is an implementation of the PhoneService interface.
type phoneService struct {
	repo repository.PhoneRepository
}

// NewPhoneService creates a new PhoneService.
func NewPhoneService(s store.Store) *phoneService {
	repo := s.Phones()
	return &phoneService{repo}
}

// Create creates a new phone.
func (s *phoneService) Create(ctx context.Context, u *entities.Phone) error {
	if !validatePhoneNumber(u.Number) {
		return entities.ErrInvalidPhoneNumber
	}

	phone := &repository.Phone{
		Number: u.Number,
		UserID: u.UserID,
	}

	if err := s.repo.Create(ctx, phone); err != nil {
		return fmt.Errorf("could not create phone: %w", err)
	}
	u.ID = phone.ID
	return nil
}

// Get returns a phone by id.
func (s *phoneService) Get(ctx context.Context, id uint) (*entities.Phone, error) {
	phone, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find phone: %w", err)
	}

	entPhone := &entities.Phone{
		ID:     phone.ID,
		Number: phone.Number,
		UserID: phone.UserID,
	}

	return entPhone, nil
}

// Update updates a phone.
func (s *phoneService) Update(ctx context.Context, u *entities.Phone) error {
	if !validatePhoneNumber(u.Number) {
		return entities.ErrInvalidPhoneNumber
	}

	phone := &repository.Phone{
		Number: u.Number,
		UserID: u.UserID,
	}

	return s.repo.Update(ctx, phone) //nolint:wrapcheck //no need to wrap here
}

// Delete deletes a phone.
func (s *phoneService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id) //nolint:wrapcheck //no need to wrap here
}
