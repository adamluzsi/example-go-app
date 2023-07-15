package contracts

import (
	"context"
	"example-go-app/domains/projectmanagement"
	"fmt"
	"github.com/adamluzsi/frameless/ports/crud/crudcontracts"
	"github.com/adamluzsi/testcase/assert"
	"github.com/adamluzsi/testcase/random"
	"testing"
)

func TestProjectRepository(t *testing.T, repo projectmanagement.ProjectRepository) {
	t.Run("when project created in the resource, then it can be found and deleted", func(t *testing.T) {
		var (
			project = makeRandomProject()
			ctx     = context.Background()
		)

		assert.NoError(t, repo.Create(ctx, &project))
		assert.NotEmpty(t, project.ID, "id was expected to be set by the ProjectRepository.Create Method call")
		defer repo.DeleteByID(ctx, project.ID)

		t.Log("when created, then it can be found")
		gotProject, found, err := repo.FindByID(ctx, project.ID)
		assert.NoError(t, err)
		assert.True(t, found, "expected to find the project in the repository")
		assert.Equal(t, project, gotProject)

		t.Log("it can be deleted")
		assert.NoError(t, repo.DeleteByID(ctx, project.ID))

		t.Log("when deleted, then it can't be found")
		_, found, err = repo.FindByID(ctx, project.ID)
		assert.NoError(t, err)
		assert.False(t, found, "should be deleted")
	})

	t.Run("when project is already deleted, it cannot be deleted again", func(t *testing.T) {
		var (
			project = makeRandomProject()
			ctx     = context.Background()
		)
		assert.NoError(t, repo.Create(ctx, &project))
		assert.NotEmpty(t, project.ID, "id was expected to be set by the ProjectRepository.Create Method call")
		assert.NoError(t, repo.DeleteByID(ctx, project.ID))
		assert.Error(t, repo.DeleteByID(ctx, project.ID))
	})

	t.Run("when context has an error, error is returned from the entities", func(t *testing.T) {
		var (
			project     = makeRandomProject()
			ctx, cancel = context.WithCancel(context.Background())
		)
		cancel()

		assert.ErrorIs(t, ctx.Err(), repo.Create(ctx, &project))
		assert.NoError(t, repo.Create(context.Background(), &project))
		defer repo.DeleteByID(context.Background(), project.ID)

		_, _, err := repo.FindByID(ctx, project.ID)
		assert.ErrorIs(t, ctx.Err(), err)
		assert.ErrorIs(t, ctx.Err(), repo.DeleteByID(ctx, project.ID))
	})
}

var rnd = random.New(random.CryptoSeed{})

func makeRandomProject() projectmanagement.Project {
	return projectmanagement.Project{
		Name: fmt.Sprintf("%s's (#%d) project", rnd.Contact().LastName, rnd.IntB(0, 128)),
	}
}

type ProjectRepositorySubject struct {
	Resource projectmanagement.ProjectRepository
}

func TestTaskRepository(t *testing.T, repo projectmanagement.TaskRepository) {
	// This is akin to creating our own contract,
	// except in this case, I've simply imported an existing test suite from the frameless project.
	//
	// I did this to save time writing this example. (laziness)
	crudcontracts.SuiteFor(func(tb testing.TB) crudcontracts.SuiteSubject[projectmanagement.Task, projectmanagement.TaskID, projectmanagement.TaskRepository] {
		return crudcontracts.SuiteSubject[projectmanagement.Task, projectmanagement.TaskID, projectmanagement.TaskRepository]{
			Resource:    repo,
			MakeContext: context.Background,
			MakeEntity: func() projectmanagement.Task {
				return projectmanagement.Task{
					ProjectID:   projectmanagement.ProjectID(rnd.UUID()),
					Subject:     rnd.String(),
					Description: rnd.String(),
				}
			},
		}
	}).Test(t)
}
