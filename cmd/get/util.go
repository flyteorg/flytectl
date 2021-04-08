package get

import (
	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

func buildResourceListRequestWithName(c *config.Config, name string) *admin.ResourceListRequest {
	request := &admin.ResourceListRequest{
		Limit:   uint32(c.Limit),
		Filters: c.Filters,
		Id: &admin.NamedEntityIdentifier{
			Project: c.Project,
			Domain:  c.Domain,
			Name:    name,
		},
	}
	if c.SortBy != "" {
		request.SortBy = &admin.Sort{
			Key:       c.SortBy,
			Direction: admin.Sort_DESCENDING,
		}
	}
	return request
}

func buildResourceListRequestWithoutName(c *config.Config) *admin.ResourceListRequest {
	request := &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: c.Project,
			Domain:  c.Domain,
		},
		Limit:   uint32(c.Limit),
		Filters: c.Filters,
	}
	if c.SortBy != "" {
		request.SortBy = &admin.Sort{
			Key:       c.SortBy,
			Direction: admin.Sort_DESCENDING,
		}
	}
	return request
}

func buildProjectListRequest(c *config.Config) *admin.ProjectListRequest {
	request := &admin.ProjectListRequest{
		Limit:   uint32(c.Limit),
		Filters: c.Filters,
	}
	if c.SortBy != "" {
		request.SortBy = &admin.Sort{
			Key:       c.SortBy,
			Direction: admin.Sort_DESCENDING,
		}
	}
	return request
}
