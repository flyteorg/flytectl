package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flytestdlib/logger"
	"io/ioutil"
	"sort"
)


func unMarshalProtoFile(ctx context.Context, fname string, cmdCtx cmdCore.CommandContext) (proto.Message, core.ResourceType, error) {
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		logger.Errorf(ctx,"Error reading file %v  due to :", fname, err)
		return nil, core.ResourceType_UNSPECIFIED, errors.New(fmt.Sprintf("Failed reading file %v",fname))
	}
	workflowSpec := &admin.WorkflowSpec{}
	if err := proto.Unmarshal(in, workflowSpec); err == nil {
		logger.Infof(ctx, "Registering workflow %v with flyte", fname)
		regsiterWorkflow(ctx, fname, workflowSpec, cmdCtx)
		return workflowSpec, core.ResourceType_WORKFLOW, nil
	}
	logger.Infof(ctx, "Failed to parse  file %v for workflow type:", fname, err)

	taskSpec := &admin.TaskSpec{}
	if err := proto.Unmarshal(in, taskSpec); err == nil {
		logger.Infof(ctx, "Registering Task  %v with flyte", fname)
		regsiterTask(ctx, fname, taskSpec, cmdCtx)
		return taskSpec, core.ResourceType_TASK, nil
	}
	logger.Infof(ctx, "Failed to parse  file %v for task type:", fname, err)
	launchPlanSpec := &admin.LaunchPlanSpec{}
	if err := proto.Unmarshal(in, launchPlanSpec); err == nil {
		logger.Infof(ctx, "Registering launch spec  %v with flyte", fname)
		regsiterLaunchPlan(ctx, fname, launchPlanSpec, cmdCtx)
		return launchPlanSpec, core.ResourceType_LAUNCH_PLAN,  nil
	}
	logger.Infof(ctx, "Failed to parse  file %v for launch plan type:", fname, err)
	return nil, core.ResourceType_UNSPECIFIED, errors.New(fmt.Sprintf("Failed parsing file %v",fname))

}

func registerFromFilesFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	files := args
	sort.Strings(files)
	logger.Infof(ctx, "Parsing files... Total(%v)", len(files))
	for _, fileName := range files {
		logger.Infof(ctx, "Parsing  %v", fileName)
		_, _, err := unMarshalProtoFile(ctx, fileName, cmdCtx);
		if err != nil {
			logger.Infof(ctx, "Skipping file %v", fileName)
		}
	}
	return nil
}

func regsiterTask(ctx context.Context, name string, flyteTaskSpec *admin.TaskSpec, cmdCtx cmdCore.CommandContext) {
	_, err := cmdCtx.AdminClient().CreateTask(ctx, &admin.TaskCreateRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_TASK,
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name: name,
			Version: "1",
		},
		Spec : flyteTaskSpec,
	})
	if err != nil {
		logger.Infof(ctx, "Failed to register task with identifier %v due to %v", flyteTaskSpec, err)
	}
}

func regsiterWorkflow(ctx context.Context, name string, flyteWorkflowSpec *admin.WorkflowSpec,
					  cmdCtx cmdCore.CommandContext) {

	_, err := cmdCtx.AdminClient().CreateWorkflow(ctx, &admin.WorkflowCreateRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_WORKFLOW,
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name: name,
			Version: "1",
		},
		Spec : flyteWorkflowSpec,
	})
	if err != nil {
		logger.Infof(ctx, "Failed to register workflow with identifier %v due to err", flyteWorkflowSpec, err)
	}
}

func regsiterLaunchPlan(ctx context.Context, name string, flyteLaunchPlanSpec *admin.LaunchPlanSpec,
					    cmdCtx cmdCore.CommandContext) {

	_, err := cmdCtx.AdminClient().CreateLaunchPlan(ctx, &admin.LaunchPlanCreateRequest{
		Id: &core.Identifier{
			ResourceType: core.ResourceType_LAUNCH_PLAN,
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name: name,
			Version: "1",
		},
		Spec : flyteLaunchPlanSpec,
	})
	if err != nil {
		logger.Infof(ctx, "Failed to register launch spec with identifier %v due to err", flyteLaunchPlanSpec, err)
	}
}