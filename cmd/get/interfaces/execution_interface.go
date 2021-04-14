package interfaces

import (
	"context"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

//go:generate mockery -all -case=underscore

type Fetcher interface {
	FetchExecution(ctx context.Context, name string, project string, domain string, cmdCtx cmdCore.CommandContext) (*admin.Execution, error)
	FetchLPVersion(ctx context.Context, name string, version string, project string, domain string, cmdCtx cmdCore.CommandContext) (*admin.LaunchPlan, error)
}
