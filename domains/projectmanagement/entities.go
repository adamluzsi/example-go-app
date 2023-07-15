package projectmanagement

import (
	"context"
	"github.com/adamluzsi/frameless/ports/crud"
)

type Project struct {
	ID   ProjectID `ext:"id"`
	Name string
}

type ProjectID string

// ProjectRepository is a role interface, serving as a way to persist Project entities
// without revealing any implementation details about the backing resource.
//
// Remember that an interface is merely a collection of method signatures; it doesn't define behavior.
// Therefore, we have a dedicated interface testing suite for this role interface in the ./contracts directory.
// This ensures that any implementation adheres to the expected behavior of the domain package,
// inverting the dependency chain between consumer and its supplier.
type ProjectRepository interface {
	// Create is a function that takes a pointer to an entity and stores it in an external resource.
	// And external resource could be a backing service like PostgreSQL.
	// The use of a pointer type allows the function to update the entity's ID value,
	// which is significant in both the external resource and the domain layer.
	// The ID is essential because entities in the backing service are referenced using their IDs,
	// which is why the ID value is included as part of the entity structure fieldset.
	//
	// The pointer is also employed for other fields managed by the external resource, such as UpdatedAt, CreatedAt,
	// and any other fields present in the domain entity but controlled by the external resource.
	Create(ctx context.Context, ptr *Project) error
	// FindByID is a function that tries to find an Entity using its ID.
	// It will inform you if it successfully located the entity or if there was an unexpected issue during the process.
	// Instead of using an error to represent a "not found" situation,
	// a return boolean value is used to provide this information explicitly.
	//
	//
	// Why the return signature includes a found bool value?
	//
	// This approach serves two key purposes.
	// First, it ensures that the go-vet tool checks if the 'found' boolean variable is reviewed before using the entity.
	// Second, it enhances readability and demonstrates the function's cyclomatic complexity.
	//   total: 2^(n+1+1)
	//     -> found/bool 2^(n+1)  | An entity might be found or not.
	//     -> error 2^(n+1)       | An error might occur or not.
	//
	// Additionally, this method prevents returning an initialized pointer type with no value,
	// which could lead to a runtime error if a valid but nil pointer is given to an interface variable type.
	//   (MyInterface)((*Entity)(nil)) != nil
	//
	// Similar approaches can be found in the standard library,
	// such as SQL null value types and environment lookup in the os package.
	FindByID(ctx context.Context, id ProjectID) (ent Project, found bool, err error)
	// DeleteByID will remove a <V> type entity from the repository by a given ID
	DeleteByID(ctx context.Context, id ProjectID) error
}

type Task struct {
	ID          TaskID `ext:"id"`
	ProjectID   ProjectID
	Subject     string
	Description string
}

type TaskID string

type TaskRepository interface {
	crud.Creator[Task]
	crud.Finder[Task, TaskID]
	crud.Updater[Task]
	crud.Deleter[TaskID]
}
