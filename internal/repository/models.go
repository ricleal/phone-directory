package repository

type User struct {
	ID        uint
	Name      string
	Phones    []Phone
	Addresses []Address
}

type Phone struct {
	ID     uint
	Number string
	UserID uint
}

type Address struct {
	ID      uint
	Address string
	UserID  uint
}
