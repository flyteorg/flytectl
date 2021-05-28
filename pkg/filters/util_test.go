package filters

import (
	"testing"

	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/stretchr/testify/assert"
)

var (
	project = "flytesnack"
	domain  = "staging"
	name    = "test"
	output  = "json"
)

func TestListRequestWithoutNameFunc(t *testing.T) {
	config.GetConfig().Output = output
	config.GetConfig().Project = project
	config.GetConfig().Domain = domain
	filter := Filters{
		Limit: 100,
	}
	request := BuildResourceListRequestWithName(filter, "")
	expectedResponse := &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: project,
			Domain:  domain,
		},
		Limit:   100,
		Filters: "",
	}
	assert.Equal(t, expectedResponse, request)
}

func TestProjectListRequestFunc(t *testing.T) {
	config.GetConfig().Output = output
	config.GetConfig().Project = project
	config.GetConfig().Domain = domain
	filter := Filters{
		Limit: 100,
	}
	request := BuildProjectListRequest(filter)
	expectedResponse := &admin.ProjectListRequest{
		Limit:   100,
		Filters: "",
	}
	assert.Equal(t, expectedResponse, request)
}

func TestListRequestWithNameFunc(t *testing.T) {
	config.GetConfig().Output = output
	config.GetConfig().Project = project
	config.GetConfig().Domain = domain
	filter := Filters{
		Limit: 100,
	}
	request := BuildResourceListRequestWithName(filter, name)
	expectedResponse := &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: project,
			Domain:  domain,
			Name:    name,
		},
		Limit:   100,
		Filters: "",
	}
	assert.Equal(t, expectedResponse, request)
}
