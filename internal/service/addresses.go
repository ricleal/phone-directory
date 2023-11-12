package service

import (
	"context"
	"fmt"

	"phone-directory/internal/entities"
	"phone-directory/internal/repository"
	"phone-directory/internal/store"
)

func validateAddress(address string) error {
	if address == "" {
		return entities.ErrInvalidAddress
	}
	return nil
}

//go:generate mockgen -source=addresses.go -package=mock -destination=mock/addresses.go

// AddressService is a domain service for addresses.
type AddressService interface {
	Create(ctx context.Context, u *entities.Address) error
	Get(ctx context.Context, id uint) (*entities.Address, error)
	Update(ctx context.Context, u *entities.Address) error
	Delete(ctx context.Context, id uint) error
}

type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(s store.Store) *addressService {
	repo := s.Addresses()
	return &addressService{repo}
}

func (s *addressService) Create(ctx context.Context, u *entities.Address) error {
	if err := validateAddress(u.Address); err != nil {
		return entities.ErrInvalidAddress
	}

	address := &repository.Address{
		Address: u.Address,
		UserID:  u.UserID,
	}

	if err := s.repo.Create(ctx, address); err != nil {
		return fmt.Errorf("could not create address: %w", err)
	}
	u.ID = address.ID
	return nil
}

func (s *addressService) Get(ctx context.Context, id uint) (*entities.Address, error) {
	address, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find address: %w", err)
	}

	entAddress := &entities.Address{
		ID:      address.ID,
		Address: address.Address,
		UserID:  address.UserID,
	}

	return entAddress, nil
}

func (s *addressService) Update(ctx context.Context, u *entities.Address) error {
	if err := validateAddress(u.Address); err != nil {
		return entities.ErrInvalidAddress
	}

	address := &repository.Address{
		Address: u.Address,
		UserID:  u.UserID,
	}

	return s.repo.Update(ctx, address)
}

func (s *addressService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
