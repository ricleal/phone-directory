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

type AddressSuite struct {
	suite.Suite
	container *testcontainers.MySQLContainer
	s         *persistent.Storage
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressSuite))
}

func (ts *AddressSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = testUtil.SetupContainer(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = persistent.NewStorage()
	require.NoError(ts.T(), err)
}

func (ts *AddressSuite) TearDownTest() {
	ctx := context.Background()
	err := testUtil.TeardownContainer(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *AddressSuite) TestData() {
	ctx := context.Background()
	a := persistent.NewAddressStorage(ts.s.DB())

	// create a user so it can be referenced by the address
	user := &repository.User{
		Name: "Jane Doe",
	}
	u := persistent.NewUserStorage(ts.s.DB())
	err := u.Create(ctx, user)
	ts.Require().NoError(err)

	var address *repository.Address
	ts.Run("create address", func() {
		address = &repository.Address{
			Address: "Some address",
			UserID:  user.ID,
		}
		err := a.Create(ctx, address)
		require.NoError(ts.T(), err)
	})

	ts.Run("get address", func() {
		address2, err := a.Get(ctx, address.ID)
		require.NoError(ts.T(), err)
		require.Equal(ts.T(), "Some address", address2.Address)
		require.Equal(ts.T(), address.UserID, address2.UserID)
	})

	ts.Run("update address", func() {
		address2 := &repository.Address{
			ID:      address.ID,
			Address: "Some new address",
			UserID:  user.ID,
		}
		err := a.Update(ctx, address2)
		require.NoError(ts.T(), err)

		address3, err := a.Get(ctx, address.ID)
		require.NoError(ts.T(), err)
		require.Equal(ts.T(), "Some new address", address3.Address)
	})

	ts.Run("delete address", func() {
		err := a.Delete(ctx, 1)
		require.NoError(ts.T(), err)
	})

}
