package memory_test

import (
	"example-go-app/adapters/memory"
	"example-go-app/domains/projectmanagement/contracts"
	"testing"
)

func TestNewProjectRepository(t *testing.T) {
	contracts.TestProjectRepository(t, memory.NewProjectRepository())
}
