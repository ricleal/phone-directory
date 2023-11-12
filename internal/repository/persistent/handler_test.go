//go:build integration
// +build integration

package persistent_test

import (
	"context"

	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/mysql"

	"phone-directory/internal/repository/persistent"
	testUtil "phone-directory/internal/repository/persistent/testing"
)

type MysqlTestSuite struct {
	suite.Suite
	container *testcontainers.MySQLContainer
	s         *persistent.Storage
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}

// func (ts *MysqlTestSuite) SetupTest() {
// 	var err error
// 	ctx := context.Background()
// 	ts.container, err = testUtil.SetupDB(ctx)
// 	ts.Require().NoError(err)
// 	ts.db, err = persistent.NewDatabase()
// 	ts.Require().NoError(err)
// 	err = ts.db.Ping()
// 	ts.Require().NoError(err)
// 	ts.db.Migrate()
// }

// func (ts *MysqlTestSuite) TearDownTest() {
// 	ctx := context.Background()
// 	err := ts.db.Close()
// 	ts.Require().NoError(err)
// 	err = testUtil.TeardownDB(ctx, ts.container)
// 	ts.Require().NoError(err)
// }

func (ts *MysqlTestSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = testUtil.SetupContainer(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = persistent.NewStorage()
	require.NoError(ts.T(), err)
}

func (ts *MysqlTestSuite) TearDownTest() {
	ctx := context.Background()
	err := testUtil.TeardownContainer(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

// This tests only that the connection to the DB is working and also the migrations.
func (ts *MysqlTestSuite) TestMysql() {
	s, err := persistent.NewStorage()
	ts.Require().NoError(err)
	ts.Require().NotNil(s.DB)
	err = s.Ping()
	ts.Require().NoError(err)

	// check the existing tables
	rows, err := s.DB().Query("SELECT table_name FROM information_schema.tables")
	ts.Require().NoError(err)
	ts.Require().NoError(rows.Err())
	defer rows.Close()
	var tables []string //nolint:prealloc // this is a test
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		require.NoError(ts.T(), err)
		tables = append(tables, table)
	}
	ts.Require().Contains(tables, "schema_migrations")
	ts.Require().Contains(tables, "users")
	ts.Require().Contains(tables, "phones")
	ts.Require().Contains(tables, "addresses")
	s.Close()
}
