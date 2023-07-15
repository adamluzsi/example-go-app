package memory

import (
	"example-go-app/domains/projectmanagement"
	"github.com/adamluzsi/frameless/adapters/memory"
)

// NewProjectRepository returns a in-memory repository variant, for testing purposes.
//
// Out of laziness, I'm just using an existing crud compliant in-memory implementation from the frameless project,
// implementing the in-memory variant should be not difficult using a map as the table.
func NewProjectRepository() *memory.Repository[projectmanagement.Project, projectmanagement.ProjectID] {
	m := memory.NewMemory()
	return memory.NewRepository[projectmanagement.Project, projectmanagement.ProjectID](m)
}
