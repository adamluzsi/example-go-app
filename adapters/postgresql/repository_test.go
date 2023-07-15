package postgresql_test

import (
	"context"
	"example-go-app/adapters/postgresql"
	"example-go-app/domains/projectmanagement/contracts"
	"github.com/adamluzsi/testcase/assert"
	"github.com/jackc/pgx/v5/pgxpool"
	"testing"
)

func TestProjectRepository_implementsProjectRepositoryContract(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), databaseURL(t))
	assert.NoError(t, err)
	defer pool.Close()
	MigrateDB(t)
	contracts.TestProjectRepository(t, postgresql.ProjectRepository{PGXPool: pool})
}

// Test other Database driver related things such as retry mechanism

func TestProjectRepository_retry(t *testing.T) {
	// ...
}

func TestNewTaskRepository_implementsTaskRepositoryContract(t *testing.T) {
	taskRepo, err := postgresql.NewTaskRepository(databaseURL(t))
	assert.NoError(t, err)
	defer taskRepo.Connection.Close()
	MigrateDB(t)
	contracts.TestTaskRepository(t, taskRepo)
}
