package persistent

import (
	"context"
	"log"

	"phone-directory/internal/repository"

	"gorm.io/gorm"
)

type PhoneStorage struct {
	db     repository.DBTx
	gormDB *gorm.DB
}

func NewPhoneStorage(db repository.DBTx) *PhoneStorage {
	gormDB, err := standardToGormConnection(db)
	if err != nil {
		log.Fatal("failed to convert sql.DB to gorm.DB")
	}
	return &PhoneStorage{
		db:     db,
		gormDB: gormDB,
	}
}

func (s *PhoneStorage) Create(ctx context.Context, phone *repository.Phone) error {

	u := Phone{
		Number: phone.Number,
		UserID: phone.UserID,
	}

	if err := s.gormDB.Create(&u).Error; err != nil {
		return err
	}

	phone.ID = u.ID

	return nil

}

func (s *PhoneStorage) Get(ctx context.Context, id uint) (*repository.Phone, error) {

	var u Phone

	if err := s.gormDB.Model(&Phone{}).First(&u, id).Error; err != nil {
		return nil, err
	}

	phone := &repository.Phone{
		ID:     u.ID,
		Number: u.Number,
		UserID: u.UserID,
	}

	return phone, nil

}

func (s *PhoneStorage) Update(ctx context.Context, phone *repository.Phone) error {

	u := Phone{
		Number: phone.Number,
	}

	if err := s.gormDB.Model(&Phone{}).Where("id = ?", phone.ID).Updates(&u).Error; err != nil {
		return err
	}

	return nil

}

func (s *PhoneStorage) Delete(ctx context.Context, id uint) error {

	if err := s.gormDB.Where("id = ?", id).Delete(&Phone{}).Error; err != nil {
		return err
	}

	return nil

}
