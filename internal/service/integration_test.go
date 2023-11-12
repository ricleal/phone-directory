//go:build integration
// +build integration

package service_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/mysql"

	"phone-directory/internal/entities"
	"phone-directory/internal/repository/persistent"
	testUtil "phone-directory/internal/repository/persistent/testing"
	"phone-directory/internal/service"
	"phone-directory/internal/store"
)

type PhoneListTestSuite struct {
	suite.Suite
	container *testcontainers.MySQLContainer
	s         *persistent.Storage
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestPhoneListTestSuite(t *testing.T) {
	suite.Run(t, new(PhoneListTestSuite))
}

func (ts *PhoneListTestSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = testUtil.SetupContainer(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = persistent.NewStorage()
	require.NoError(ts.T(), err)
}

func (ts *PhoneListTestSuite) TearDownTest() {
	ctx := context.Background()
	err := testUtil.TeardownContainer(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *PhoneListTestSuite) TestValid() {
	s := store.NewPersistentStore(ts.s.DB())
	sa := service.NewAddressService(s)
	sp := service.NewPhoneService(s)
	su := service.NewUserService(s)
	ctx := context.Background()

	ts.Run("get user 1 empty DB", func() {
		users, err := su.Get(ctx, 0)
		ts.Require().Contains(err.Error(), "record not found")
		ts.Require().Nil(users)
	})

	ts.Run("create user empty DB", func() {
		phoneNumber := "123-234-1234"
		address := "123 Main St"
		user := entities.User{
			Name: "John Doe",
			Phones: []entities.Phone{
				{
					Number: phoneNumber,
				},
			},
			Addresses: []entities.Address{
				{
					Address: address,
				},
			},
		}

		err := su.Create(ctx, &user)
		ts.Require().NoError(err)
		user2, err := su.Get(ctx, user.ID)
		ts.Require().Equal("John Doe", user2.Name)
		ts.Require().Equal(phoneNumber, user2.Phones[0].Number)
		ts.Require().Equal(address, user2.Addresses[0].Address)
	})

	ts.Run("create a phone and address and assign to user", func() {
		phoneNumber := "123-234-1234"
		address := "123 Main St"

		user := entities.User{
			Name: "John Doe",
		}
		err := su.Create(ctx, &user)
		ts.Require().NoError(err)

		phone := entities.Phone{
			Number: phoneNumber,
			UserID: user.ID,
		}
		err = sp.Create(ctx, &phone)
		ts.Require().NoError(err)

		addr := entities.Address{
			Address: address,
			UserID:  user.ID,
		}
		err = sa.Create(ctx, &addr)
		ts.Require().NoError(err)

		user2, err := su.Get(ctx, user.ID)
		ts.Require().NoError(err)
		ts.Require().Equal("John Doe", user2.Name)
		ts.Require().Len(user2.Phones, 1)
		ts.Require().Len(user2.Addresses, 1)
		ts.Require().Equal(phoneNumber, user2.Phones[0].Number)
		ts.Require().Equal(address, user2.Addresses[0].Address)
	})
}

func (ts *PhoneListTestSuite) TestInvalid() {
	s := store.NewPersistentStore(ts.s.DB())
	sp := service.NewPhoneService(s)
	su := service.NewUserService(s)
	ctx := context.Background()

	user := entities.User{
		Name: "John Doe",
	}
	err := su.Create(ctx, &user)
	ts.Require().NoError(err)

	// create a phone with invalid number
	err = sp.Create(ctx, &entities.Phone{
		Number: "123-invalid",
		UserID: user.ID,
	})
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, entities.ErrInvalidPhoneNumber)
}
