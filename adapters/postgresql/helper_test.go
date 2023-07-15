package postgresql_test

import (
	"context"
	"github.com/adamluzsi/testcase/assert"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func databaseURL(tb testing.TB) string {
	dsn, ok := os.LookupEnv("TEST_DATABASE_URL")
	assert.True(tb, ok, "TEST_DATABASE_URL is required for this test")
	return dsn
}

func MigrateDB(tb testing.TB) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL(tb))
	assert.NoError(tb, err)
	_, err = pool.Exec(ctx, testMigrateDOWN)
	assert.Nil(tb, err)
	_, err = pool.Exec(ctx, testMigrateUP)
	assert.Nil(tb, err)
	tb.Cleanup(func() {
		_, err := pool.Exec(ctx, testMigrateDOWN)
		assert.Nil(tb, err)
		pool.Close()
	})
}

const testMigrateUP = `
CREATE TABLE "projects" (
    id  TEXT  NOT  NULL  PRIMARY KEY,
  name  TEXT  NOT  NULL
);

CREATE TABLE "tasks" (
    id  TEXT  NOT  NULL  PRIMARY KEY,
  project_id  TEXT  NOT  NULL,
  subject  TEXT  NOT  NULL,
  description  TEXT  NOT  NULL
);
`

const testMigrateDOWN = `
DROP TABLE IF EXISTS "projects";
DROP TABLE IF EXISTS "tasks";
`
