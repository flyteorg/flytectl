package impl

import (
	"github.com/flyteorg/flytectl/pkg/commandutils/interfaces"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

// FetcherImpl implementing the fetcher interface
type FetcherImpl struct {
	service.AdminServiceClient
}

// NewFetcherImpl Singleton object
func NewFetcherImpl(adminClient service.AdminServiceClient) interfaces.Fetcher {
	return FetcherImpl{adminClient}
}
