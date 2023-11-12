package persistent

import (
	"context"
	"log"

	"phone-directory/internal/repository"

	"gorm.io/gorm"
)

type AddressStorage struct {
	db     repository.DBTx
	gormDB *gorm.DB
}

func NewAddressStorage(db repository.DBTx) *AddressStorage {
	gormDB, err := standardToGormConnection(db)
	if err != nil {
		log.Fatal("failed to convert sql.DB to gorm.DB")
	}
	return &AddressStorage{
		db:     db,
		gormDB: gormDB,
	}
}

func (s *AddressStorage) Create(ctx context.Context, address *repository.Address) error {

	u := Address{
		Address: address.Address,
		UserID:  address.UserID,
	}

	if err := s.gormDB.Create(&u).Error; err != nil {
		return err
	}

	address.ID = u.ID

	return nil

}

func (s *AddressStorage) Get(ctx context.Context, id uint) (*repository.Address, error) {

	var u Address

	if err := s.gormDB.Model(&Address{}).First(&u, id).Error; err != nil {
		return nil, err
	}

	address := &repository.Address{
		ID:      u.ID,
		Address: u.Address,
		UserID:  u.UserID,
	}

	return address, nil

}

func (s *AddressStorage) Update(ctx context.Context, address *repository.Address) error {

	u := Address{
		Address: address.Address,
	}

	if err := s.gormDB.Model(&Address{}).Where("id = ?", address.ID).Updates(&u).Error; err != nil {
		return err
	}

	return nil

}

func (s *AddressStorage) Delete(ctx context.Context, id uint) error {

	if err := s.gormDB.Where("id = ?", id).Delete(&Address{}).Error; err != nil {
		return err
	}

	return nil

}
