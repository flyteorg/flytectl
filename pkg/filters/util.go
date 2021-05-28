package filters

import (
	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

func BuildResourceListRequestWithName(c Filters, name string) *admin.ResourceListRequest {
	fieldSelector, err := Transform(SplitTerms(c.FieldSelector))
	if err != nil {
		fieldSelector = ""
	}
	request := &admin.ResourceListRequest{
		Limit:   uint32(c.Limit),
		Filters: fieldSelector,
		Id: &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		},
	}
	if len(name) > 0 {
		request.Id.Name = name
	}
	if sort := buildSortingRequest(c); sort != nil {
		request.SortBy = sort
	}
	return request
}

func BuildProjectListRequest(c Filters) *admin.ProjectListRequest {
	fieldSelector, err := Transform(SplitTerms(c.FieldSelector))
	if err != nil {
		fieldSelector = ""
	}
	request := &admin.ProjectListRequest{
		Limit:   uint32(c.Limit),
		Filters: fieldSelector,
	}
	if sort := buildSortingRequest(c); sort != nil {
		request.SortBy = sort
	}
	return request
}

func buildSortingRequest(c Filters) *admin.Sort {
	if len(c.SortBy) > 0 {
		sortingOrder := admin.Sort_DESCENDING
		if c.Asc {
			sortingOrder = admin.Sort_ASCENDING
		}
		return &admin.Sort{
			Key:       c.SortBy,
			Direction: sortingOrder,
		}
	}
	return nil
}
