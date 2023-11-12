package persistent

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string
	Phones    []Phone
	Addresses []Address
}

type Phone struct {
	gorm.Model
	Number string `gorm:"column:phone"`
	UserID uint
}

type Address struct {
	gorm.Model
	Address string
	UserID  uint
}
