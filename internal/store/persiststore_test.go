//go:build integration
// +build integration

package store_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/mysql"

	"phone-directory/internal/repository"
	"phone-directory/internal/repository/persistent"
	testUtil "phone-directory/internal/repository/persistent/testing"
	"phone-directory/internal/store"
)

type StoreTestSuite struct {
	suite.Suite
	container *testcontainers.MySQLContainer
	s         *persistent.Storage
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (ts *StoreTestSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = testUtil.SetupContainer(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = persistent.NewStorage()
	require.NoError(ts.T(), err)
}

func (ts *StoreTestSuite) TearDownTest() {
	ctx := context.Background()
	err := testUtil.TeardownContainer(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *StoreTestSuite) TestTransaction() {
	ctx := context.Background()
	mainStore := store.NewPersistentStore(ts.s.DB())

	usersRepo := mainStore.Users()

	ts.Run("Find all users outside of transaction before", func() {
		user, errOut := usersRepo.Get(ctx, 1)
		ts.Require().ErrorContains(errOut, "record not found")
		ts.Require().Nil(user)
	})

	ts.Run("create a user, phone and address inside a transaction", func() {
		if errOut := mainStore.ExecTx(ctx, func(s store.Store) error {
			usersRepoLocal := s.Users()
			phonesRepoLocal := s.Phones()
			addressesRepoLocal := s.Addresses()

			// create a user
			user := repository.User{
				Name: "John Doe",
			}
			err := usersRepoLocal.Create(ctx, &user)
			ts.Require().NoError(err)
			// create a phone
			phone := repository.Phone{
				Number: "123456789",
				UserID: user.ID,
			}
			err = phonesRepoLocal.Create(ctx, &phone)
			ts.Require().NoError(err)
			// create an address
			address := repository.Address{
				Address: "123 Main St",
				UserID:  user.ID,
			}
			err = addressesRepoLocal.Create(ctx, &address)
			ts.Require().NoError(err)

			// Find eager user by ID
			user2, err := usersRepoLocal.Get(ctx, user.ID)
			ts.Require().NoError(err)
			ts.Require().Equal("John Doe", user2.Name)
			ts.Require().Len(user2.Phones, 1)
			ts.Require().Len(user2.Addresses, 1)
			return nil
		}); errOut != nil {
			ts.T().Errorf("ExecTx: %v", errOut)
		}
	})

	ts.Run("Find all users outside of transaction after", func() {
		user, errOut := usersRepo.Get(ctx, 1)
		ts.Require().NoError(errOut)
		ts.Require().NotNil(user)
		ts.Require().Equal("John Doe", user.Name)
		ts.Require().Len(user.Phones, 1)
		ts.Require().Len(user.Addresses, 1)
	})
}
