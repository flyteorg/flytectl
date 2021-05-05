package interfaces

import (
	"context"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

//go:generate mockery -all -case=underscore

// Fetcher Interface for exposing the fetch capabilities from the admin and also allow this to be injectable into other
// modules. eg : create execution which requires to fetch launchplan details to construct the execution spec.
type Fetcher interface {
	// FetchExecution fetches the execution based on name, project, domain
	FetchExecution(ctx context.Context, name, project, domain string) (*admin.Execution, error)

	// FetchAllVerOfLP fetches all versions of launch plan in a  project, domain
	FetchAllVerOfLP(ctx context.Context, lpName, project, domain string) ([]*admin.LaunchPlan, error)

	// FetchLPLatestVersion fetches latest version of launch plan in a  project, domain
	FetchLPLatestVersion(ctx context.Context, name, project, domain string) (*admin.LaunchPlan, error)

	// FetchLPVersion fetches particular version of launch plan in a  project, domain
	FetchLPVersion(ctx context.Context, name, version, project, domain string) (*admin.LaunchPlan, error)

	// FetchAllVerOfTask fetches all versions of task in a  project, domain
	FetchAllVerOfTask(ctx context.Context, name, project, domain string) ([]*admin.Task, error)

	// FetchTaskLatestVersion fetches latest version of task in a  project, domain
	FetchTaskLatestVersion(ctx context.Context, name, project, domain string) (*admin.Task, error)

	// FetchTaskVersion fetches particular version of task in a  project, domain
	FetchTaskVersion(ctx context.Context, name, version, project, domain string) (*admin.Task, error)
}
