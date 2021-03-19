package get

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/clients/go/admin/mocks"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const projectValue = "dummyProject"
const domainValue = "dummyDomain"
const output = "json"
const executionNameValue = "e124"
const launchPlanNameValue = "lp_name"
const launchPlanVersionValue = "lp_version"
const workflowNameValue = "wf_name"
const workflowVersionValue = "wf_version"

var (
	reader                  *os.File
	writer                  *os.File
	ctx                     context.Context
	argsLp                  []string
	argsTask                []string
	stdOut                  *os.File
	stderr                  *os.File
	err                     error
	cmdCtx                  cmdCore.CommandContext
	resourceListRequest     *admin.ResourceListRequest
	resourceListRequestTask *admin.ResourceListRequest
	objectGetRequest        *admin.ObjectGetRequest
	objectGetRequestTask    *admin.ObjectGetRequest
	namedIDRequest          *admin.NamedEntityIdentifierListRequest
	launchPlanListResponse  *admin.LaunchPlanList
	taskListResponse        *admin.TaskList
	mockClient              *mocks.AdminServiceClient
)

func setup() {
	ctx = context.Background()
	argsLp = []string{}
	argsTask = []string{}
	config.GetConfig().Project = projectValue
	config.GetConfig().Domain = domainValue
	config.GetConfig().Output = output
	reader, writer, err = os.Pipe()
	if err != nil {
		panic(err)
	}
	stdOut = os.Stdout
	stderr = os.Stderr
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	argsLp = append(argsLp, "launchplan1")
	argsTask = append(argsTask, "task1")
	mockClient = new(mocks.AdminServiceClient)
	mockOutStream := new(io.Writer)
	cmdCtx = cmdCore.NewCommandContext(mockClient, *mockOutStream)

	launchPlan1 := &admin.LaunchPlan{
		Id: &core.Identifier{
			Name:    "launchplan1",
			Version: "v1",
		},
		Closure: &admin.LaunchPlanClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 0, Nanos: 0},
		},
	}
	launchPlan2 := &admin.LaunchPlan{
		Id: &core.Identifier{
			Name:    "launchplan1",
			Version: "v2",
		},
		Closure: &admin.LaunchPlanClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 1, Nanos: 0},
		},
	}

	launchPlans := []*admin.LaunchPlan{launchPlan2, launchPlan1}

	task1 := &admin.Task{
		Id: &core.Identifier{
			Name:    "task1",
			Version: "v1",
		},
		Closure: &admin.TaskClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 0, Nanos: 0},
		},
	}
	task2 := &admin.Task{
		Id: &core.Identifier{
			Name:    "task1",
			Version: "v2",
		},
		Closure: &admin.TaskClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 1, Nanos: 0},
		},
	}

	tasks := []*admin.Task{task2, task1}
	resourceListRequest = &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: projectValue,
			Domain:  domainValue,
			Name:    argsLp[0],
		},
		SortBy: &admin.Sort{
			Key:       "created_at",
			Direction: admin.Sort_DESCENDING,
		},
		Limit: 100,
	}
	resourceListRequestTask = &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: projectValue,
			Domain:  domainValue,
			Name:    argsTask[0],
		},
		SortBy: &admin.Sort{
			Key:       "created_at",
			Direction: admin.Sort_DESCENDING,
		},
		Limit: 100,
	}

	launchPlanListResponse = &admin.LaunchPlanList{
		LaunchPlans: launchPlans,
	}
	taskListResponse = &admin.TaskList{
		Tasks: tasks,
	}
	objectGetRequest = &admin.ObjectGetRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_LAUNCH_PLAN,
			Project:      projectValue,
			Domain:       domainValue,
			Name:         argsLp[0],
			Version:      "v2",
		},
	}

	objectGetRequestTask = &admin.ObjectGetRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_TASK,
			Project:      projectValue,
			Domain:       domainValue,
			Name:         argsTask[0],
			Version:      "v2",
		},
	}
	namedIDRequest = &admin.NamedEntityIdentifierListRequest{
		Project: projectValue,
		Domain:  domainValue,
		SortBy: &admin.Sort{
			Key:       "name",
			Direction: admin.Sort_ASCENDING,
		},
		Limit: 100,
	}

	var entities []*admin.NamedEntityIdentifier
	id1 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "launchplan1",
	}
	id2 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "launchplan2",
	}
	entities = append(entities, id1, id2)
	namedIdentifierList := &admin.NamedEntityIdentifierList{
		Entities: entities,
	}

	var taskEntities []*admin.NamedEntityIdentifier
	idTask1 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "task1",
	}
	idTask2 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "task2",
	}
	taskEntities = append(taskEntities, idTask1, idTask2)
	namedIdentifierListTask := &admin.NamedEntityIdentifierList{
		Entities: taskEntities,
	}

	mockClient.OnListLaunchPlansMatch(ctx, resourceListRequest).Return(launchPlanListResponse, nil)
	mockClient.OnGetLaunchPlanMatch(ctx, objectGetRequest).Return(launchPlan2, nil)
	mockClient.OnListLaunchPlanIdsMatch(ctx, namedIDRequest).Return(namedIdentifierList, nil)

	mockClient.OnListTasksMatch(ctx, resourceListRequestTask).Return(taskListResponse, nil)
	mockClient.OnGetTaskMatch(ctx, objectGetRequestTask).Return(task2, nil)
	mockClient.OnListTaskIdsMatch(ctx, namedIDRequest).Return(namedIdentifierListTask, nil)

	taskConfig.Latest = false
	launchPlanConfig.Latest = false
	taskConfig.Version = ""
	launchPlanConfig.Version = ""
}

func teardownAndVerify(t *testing.T, expectedLog string) {
	writer.Close()
	os.Stdout = stdOut
	os.Stderr = stderr
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, reader); err == nil {
		assert.Equal(t, strings.Trim(expectedLog, "\n"), strings.Trim(buf.String(), "\n"))
	}
}
