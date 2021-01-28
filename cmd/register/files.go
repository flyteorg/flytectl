package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flytestdlib/logger"
	"io/ioutil"
	"sort"
	"strings"
)

const identifierFileSuffix = "identifier"

func unMarshalContents(ctx context.Context, fileContents []byte, fname string) (proto.Message, error) {
	workflowSpec := &admin.WorkflowSpec{}
	if err := proto.Unmarshal(fileContents, workflowSpec); err == nil {
		return workflowSpec, nil
	}
	logger.Debugf(ctx, "Failed to unmarshal file %v for workflow type", fname)
	taskSpec := &admin.TaskSpec{}
	if err := proto.Unmarshal(fileContents, taskSpec); err == nil {
		return taskSpec, nil
	}
	logger.Debugf(ctx, "Failed to unmarshal  file %v for task type", fname)
	launchPlanSpec := &admin.LaunchPlanSpec{}
	if err := proto.Unmarshal(fileContents, launchPlanSpec); err == nil {
		return launchPlanSpec, nil
	}
	logger.Debugf(ctx, "Failed to unmarshal file %v for launch plan type", fname)
	return nil, errors.New(fmt.Sprintf("Failed unmarshalling file %v", fname))

}

func register(ctx context.Context, message proto.Message, name string, cmdCtx cmdCore.CommandContext) error {
	switch message.(type) {
	case *admin.LaunchPlanSpec:
		_, err := cmdCtx.AdminClient().CreateLaunchPlan(ctx, &admin.LaunchPlanCreateRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_LAUNCH_PLAN,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         name,
				Version:      GetConfig().version,
			},
			Spec: message.(*admin.LaunchPlanSpec),
		})
		return err
	case *admin.WorkflowSpec:
		_, err := cmdCtx.AdminClient().CreateWorkflow(ctx, &admin.WorkflowCreateRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_WORKFLOW,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         name,
				Version:      GetConfig().version,
			},
			Spec: message.(*admin.WorkflowSpec),
		})
		return err
	case *admin.TaskSpec:
		_, err := cmdCtx.AdminClient().CreateTask(ctx, &admin.TaskCreateRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_TASK,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         name,
				Version:      GetConfig().version,
			},
			Spec: message.(*admin.TaskSpec),
		})
		return err
	default:
		return errors.New(fmt.Sprintf("Failed registering unknown entity  %v type as part of %v file", message, name))
	}
}

func hydrateNode(node *core.Node) {
	targetNode := node.Target
	switch targetNode.(type) {
	case *core.Node_TaskNode:
		taskNodeWrapper := targetNode.(*core.Node_TaskNode)
		taskNodeReference := taskNodeWrapper.TaskNode.Reference.(*core.TaskNode_ReferenceId)
		hydrateIdentifier(taskNodeReference.ReferenceId)
	case *core.Node_WorkflowNode:
		workflowNodeWrapper := targetNode.(*core.Node_WorkflowNode)
		switch workflowNodeWrapper.WorkflowNode.Reference.(type) {
		case *core.WorkflowNode_SubWorkflowRef:
			subWorkflowNodeReference := workflowNodeWrapper.WorkflowNode.Reference.(*core.WorkflowNode_SubWorkflowRef)
			hydrateIdentifier(subWorkflowNodeReference.SubWorkflowRef)
		case *core.WorkflowNode_LaunchplanRef:
			launchPlanNodeReference := workflowNodeWrapper.WorkflowNode.Reference.(*core.WorkflowNode_LaunchplanRef)
			hydrateIdentifier(launchPlanNodeReference.LaunchplanRef)
		}
	}
}

func hydrateIdentifier(identifier *core.Identifier) {
	identifier.Project = config.GetConfig().Project
	identifier.Domain = config.GetConfig().Domain
	identifier.Version = GetConfig().version
}

func hydrateSpec(message proto.Message) {
	switch message.(type) {
	case *admin.LaunchPlanSpec:
		launchSpec := message.(*admin.LaunchPlanSpec)
		hydrateIdentifier(launchSpec.WorkflowId)
	case *admin.WorkflowSpec:
		workflowSpec := message.(*admin.WorkflowSpec)
		for _, Noderef := range workflowSpec.Template.Nodes {
			hydrateNode(Noderef)
		}
		hydrateIdentifier(workflowSpec.Template.Id)
	case *admin.TaskSpec:
		taskSpec := message.(*admin.TaskSpec)
		hydrateIdentifier(taskSpec.Template.Id)
	}
}

func registerFromFilesFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	files := args
	sort.Strings(files)
	logger.Infof(ctx, "Parsing files... Total(%v)", len(files))
	logger.Infof(ctx, "Params version %v", GetConfig().version)
	for _, absFilePath := range files {
		if strings.Contains(absFilePath, identifierFileSuffix) {
			continue
		}
		logger.Infof(ctx, "Parsing  %v", absFilePath)
		fileContents, err := ioutil.ReadFile(absFilePath)
		if err != nil {
			logger.Errorf(ctx, "Error reading file %v  due to : %v. Skipping", absFilePath, err)
			continue
		}
		spec, err := unMarshalContents(ctx, fileContents, absFilePath)
		if err != nil {
			logger.Errorf(ctx, "Error unmarshalling Skipping file %v due to : %v", absFilePath, err)
			continue
		}
		name := getEntityNameFromPath(absFilePath)
		logger.Debugf(ctx, "Spec : %v ", getJsonSpec(spec))
		hydrateSpec(spec)
		logger.Debugf(ctx, "Hydrated Spec : %v ", getJsonSpec(spec))
		err = register(ctx, spec, name, cmdCtx)
		if err != nil {
			logger.Errorf(ctx, "Error registering entity %v due to : %v", name, err)
			continue
		}
		logger.Infof(ctx, "Registered successfully entity %v", name)
	}
	return nil
}

func getJsonSpec(message proto.Message) string {
	marshaller := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	jsonSpec, _ := marshaller.MarshalToString(message)
	return jsonSpec
}

// Temporarily using file name for generating the enity name. This would be changed to use it directly from serialized
// version of the entity protobuf files.
func getEntityNameFromPath(absFilePath string) string {
	pathComponents := strings.Split(absFilePath, "/")
	fileName := strings.SplitAfterN(pathComponents[len(pathComponents)-1], "_", 2)
	fileNameWithoutSuffix := strings.TrimSuffix(fileName[len(fileName)-1], ".pb")
	return fileNameWithoutSuffix
}
