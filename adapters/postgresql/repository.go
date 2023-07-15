package postgresql

import (
	"context"
	"errors"
	"example-go-app/domains/projectmanagement"
	"fmt"
	"github.com/adamluzsi/frameless/adapters/postgresql"
	"github.com/adamluzsi/frameless/ports/iterators"
	"github.com/adamluzsi/testcase/random"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProjectRepository is an example projectmanagement.ProjectRepository implementation
// which supplies the behaviour defined by the ProjectRepository contract.
type ProjectRepository struct {
	PGXPool *pgxpool.Pool
}

var rnd = random.New(random.CryptoSeed{})

func (r ProjectRepository) Create(ctx context.Context, ptr *projectmanagement.Project) error {
	ptr.ID = projectmanagement.ProjectID(rnd.UUID())
	_, err := r.PGXPool.Exec(ctx, "INSERT INTO projects (id, name) VALUES ($1, $2)", ptr.ID, ptr.Name)
	return err
}

func (r ProjectRepository) FindByID(ctx context.Context, id projectmanagement.ProjectID) (projectmanagement.Project, bool, error) {
	var ent projectmanagement.Project

	err := r.PGXPool.QueryRow(ctx, "SELECT id, name FROM projects WHERE id = $1 LIMIT 1", id).
		Scan(&ent.ID, &ent.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		// we communicate not found not as an sql error but as a business logic of value not exist in the repository.
		return projectmanagement.Project{}, false, nil
	}
	if err != nil {
		return projectmanagement.Project{}, false, err
	}

	return ent, true, nil
}

func (r ProjectRepository) DeleteByID(ctx context.Context, id projectmanagement.ProjectID) error {
	res, err := r.PGXPool.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("project with ID of %v is not found", id)
	}
	return nil
}

func NewTaskRepository(dsn string) (postgresql.Repository[projectmanagement.Task, projectmanagement.TaskID], error) {
	conn, err := postgresql.Connect(dsn)
	if err != nil {
		return postgresql.Repository[projectmanagement.Task, projectmanagement.TaskID]{}, err
	}
	return postgresql.Repository[projectmanagement.Task, projectmanagement.TaskID]{
		Mapping: postgresql.Mapper[projectmanagement.Task, projectmanagement.TaskID]{
			Table:   "tasks",
			ID:      "id",
			Columns: []string{"id", "project_id", "subject", "description"},
			ToArgsFn: func(ptr *projectmanagement.Task) ([]interface{}, error) {
				return []any{ptr.ID, ptr.ProjectID, ptr.Subject, ptr.Description}, nil
			},
			MapFn: func(scanner iterators.SQLRowScanner) (projectmanagement.Task, error) {
				var task projectmanagement.Task
				err := scanner.Scan(&task.ID, &task.ProjectID, &task.Subject, &task.Description)
				return task, err
			},
			NewIDFn: func(ctx context.Context) (projectmanagement.TaskID, error) {
				return projectmanagement.TaskID(rnd.UUID()), nil
			},
		},
		Connection: conn,
	}, nil
}
