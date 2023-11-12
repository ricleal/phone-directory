package persistent_test

import (
	"context"

	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/mysql"

	"phone-directory/internal/repository"
	"phone-directory/internal/repository/persistent"
	testUtil "phone-directory/internal/repository/persistent/testing"
)

type UserTestSuite struct {
	suite.Suite
	container *testcontainers.MySQLContainer
	s         *persistent.Storage
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (ts *UserTestSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = testUtil.SetupContainer(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = persistent.NewStorage()
	require.NoError(ts.T(), err)
}

func (ts *UserTestSuite) TearDownTest() {
	ctx := context.Background()
	err := testUtil.TeardownContainer(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *UserTestSuite) TestData() {
	ctx := context.Background()
	u := persistent.NewUserStorage(ts.s.DB())

	ts.Run("create user with phone and address", func() {
		user := &repository.User{
			Name: "John Doe",
			Phones: []repository.Phone{
				{
					Number: "123456789",
				},
			},
			Addresses: []repository.Address{
				{
					Address: "123 Main St",
				},
			},
		}
		err := u.Create(ctx, user)
		ts.Require().NoError(err)

		// get the user back
		user, err = u.Get(ctx, user.ID)
		ts.Require().NoError(err)
		ts.Require().Equal("John Doe", user.Name)
		ts.Require().Len(user.Phones, 1)
		ts.Require().Len(user.Addresses, 1)
		ts.Require().Equal("123456789", user.Phones[0].Number)
		ts.Require().Equal("123 Main St", user.Addresses[0].Address)

	})

	ts.Run("create user without phone and address", func() {
		user := &repository.User{
			Name: "Jane Doe",
		}
		err := u.Create(ctx, user)
		ts.Require().NoError(err)
		ts.Require().Equal(uint(2), user.ID)
	})

	ts.Run("create user with phone and without address", func() {
		user := &repository.User{
			Name: "John Smith",
			Phones: []repository.Phone{
				{
					Number: "987654321",
				},
			},
		}
		err := u.Create(ctx, user)
		ts.Require().NoError(err)
		ts.Require().Equal(uint(3), user.ID)
	})

	ts.Run("create user without phone and with address", func() {
		user := &repository.User{
			Name: "Jane Smith",
			Addresses: []repository.Address{
				{
					Address: "456 Main St",
				},
			},
		}
		err := u.Create(ctx, user)
		ts.Require().NoError(err)
		ts.Require().Equal(uint(4), user.ID)
	})

	ts.Run("create user with multiple phones and addresses", func() {
		user := &repository.User{
			Name: "John Doe",
			Phones: []repository.Phone{
				{
					Number: "123456789",
				},
				{
					Number: "987654321",
				},
			},
			Addresses: []repository.Address{
				{
					Address: "123 Main St",
				},
				{
					Address: "456 Main St",
				},
			},
		}
		err := u.Create(ctx, user)
		ts.Require().NoError(err)
		ts.Require().Equal(uint(5), user.ID)
	})

	ts.Run("get user with phone and address", func() {
		user, err := u.Get(ctx, 1)
		ts.Require().NoError(err)
		ts.Require().Equal("John Doe", user.Name)
		ts.Require().Len(user.Phones, 1)
		ts.Require().Len(user.Addresses, 1)
		ts.Require().Equal("123456789", user.Phones[0].Number)
		ts.Require().Equal("123 Main St", user.Addresses[0].Address)
	})

	ts.Run("get user without phone and address", func() {
		user, err := u.Get(ctx, 2)
		ts.Require().NoError(err)
		ts.Require().Equal("Jane Doe", user.Name)
		ts.Require().Len(user.Phones, 0)
		ts.Require().Len(user.Addresses, 0)
	})

	ts.Run("get user with phone and without address", func() {
		user, err := u.Get(ctx, 3)
		ts.Require().NoError(err)
		ts.Require().Equal("John Smith", user.Name)
		ts.Require().Len(user.Phones, 1)
		ts.Require().Len(user.Addresses, 0)
		ts.Require().Equal("987654321", user.Phones[0].Number)
	})

	ts.Run("get user without phone and with address", func() {
		user, err := u.Get(ctx, 4)
		ts.Require().NoError(err)
		ts.Require().Equal("Jane Smith", user.Name)
		ts.Require().Len(user.Phones, 0)
		ts.Require().Len(user.Addresses, 1)
		ts.Require().Equal("456 Main St", user.Addresses[0].Address)
	})

	ts.Run("get user with multiple phones and addresses", func() {
		user, err := u.Get(ctx, 5)
		ts.Require().NoError(err)
		ts.Require().Equal("John Doe", user.Name)
		ts.Require().Len(user.Phones, 2)
		ts.Require().Len(user.Addresses, 2)
		ts.Require().Equal("123456789", user.Phones[0].Number)
		ts.Require().Equal("987654321", user.Phones[1].Number)
		ts.Require().Equal("123 Main St", user.Addresses[0].Address)
		ts.Require().Equal("456 Main St", user.Addresses[1].Address)
	})

	ts.Run("get user that does not exist", func() {
		user, err := u.Get(ctx, 6)
		ts.Require().Error(err)
		ts.Require().Nil(user)
	})

}
