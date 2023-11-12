package persistent

import (
	"context"
	"log"

	"gorm.io/gorm"

	"phone-directory/internal/repository"
)

type UserStorage struct {
	db     repository.DBTx
	gormDB *gorm.DB
}

func NewUserStorage(db repository.DBTx) *UserStorage {
	gormDB, err := standardToGormConnection(db)
	if err != nil {
		log.Fatal("failed to convert sql.DB to gorm.DB")
	}
	return &UserStorage{
		db:     db,
		gormDB: gormDB,
	}
}

func (s *UserStorage) Create(ctx context.Context, user *repository.User) error {

	u := User{
		Name: user.Name,
	}

	for _, p := range user.Phones {
		u.Phones = append(u.Phones, Phone{
			Number: p.Number,
		})
	}

	for _, a := range user.Addresses {
		u.Addresses = append(u.Addresses, Address{
			Address: a.Address,
		})
	}

	if err := s.gormDB.Create(&u).Error; err != nil {
		return err
	}

	user.ID = u.ID

	return nil

}

func (s *UserStorage) Get(ctx context.Context, id uint) (*repository.User, error) {

	var u User

	if err := s.gormDB.Model(&User{}).Preload("Phones").Preload("Addresses").First(&u, id).Error; err != nil {
		return nil, err
	}

	user := &repository.User{
		ID:   u.ID,
		Name: u.Name,
	}

	for _, p := range u.Phones {
		user.Phones = append(user.Phones, repository.Phone{
			ID:     p.ID,
			Number: p.Number,
		})
	}

	for _, a := range u.Addresses {
		user.Addresses = append(user.Addresses, repository.Address{
			ID:      a.ID,
			Address: a.Address,
		})
	}

	return user, nil

}

func (s *UserStorage) Update(ctx context.Context, user *repository.User) error {

	u := User{
		Model: gorm.Model{
			ID: user.ID,
		},
		Name: user.Name,
	}

	for _, p := range user.Phones {
		u.Phones = append(u.Phones, Phone{
			Model: gorm.Model{
				ID: p.ID,
			},
			Number: p.Number,
			UserID: user.ID,
		})
	}

	for _, a := range user.Addresses {
		u.Addresses = append(u.Addresses, Address{
			Model: gorm.Model{
				ID: a.ID,
			},
			Address: a.Address,
			UserID:  user.ID,
		})
	}

	if err := s.gormDB.Save(&u).Error; err != nil {
		return err
	}

	return nil

}

func (s *UserStorage) Delete(ctx context.Context, id uint) error {

	if err := s.gormDB.Delete(&User{}, id).Error; err != nil {
		return err
	}

	return nil

}
