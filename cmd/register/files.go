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
	//logger.Infof(ctx, "Failed to unmarshal file %v for workflow type", fname)
	taskSpec := &admin.TaskSpec{}
	if err := proto.Unmarshal(fileContents, taskSpec); err == nil {
		return taskSpec, nil
	}
	//logger.Infof(ctx, "Failed to unmarshal  file %v for task type", fname)
	launchPlanSpec := &admin.LaunchPlanSpec{}
	if err := proto.Unmarshal(fileContents, launchPlanSpec); err == nil {
		return launchPlanSpec, nil
	}
	//logger.Infof(ctx, "Failed to unmarshal file %v for launch plan type", fname)
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

func hydrateWorkflowNode_LaunchplanRef(launchPlanNodeReference *core.WorkflowNode_LaunchplanRef) *core.WorkflowNode_LaunchplanRef {
	return &core.WorkflowNode_LaunchplanRef{
		LaunchplanRef: &core.Identifier{
			ResourceType: launchPlanNodeReference.LaunchplanRef.ResourceType,
			Project:      config.GetConfig().Project,
			Domain:       config.GetConfig().Domain,
			Name:         launchPlanNodeReference.LaunchplanRef.Name,
			Version:      GetConfig().version,
		},
	}
}

func hydrateWorkflowNode_SubWorkflowRef(subWorkflowNodeReference *core.WorkflowNode_SubWorkflowRef) *core.WorkflowNode_SubWorkflowRef {
	return &core.WorkflowNode_SubWorkflowRef{
		SubWorkflowRef: &core.Identifier{
			ResourceType: subWorkflowNodeReference.SubWorkflowRef.ResourceType,
			Project:      config.GetConfig().Project,
			Domain:       config.GetConfig().Domain,
			Name:         subWorkflowNodeReference.SubWorkflowRef.Name,
			Version:      GetConfig().version,
		},
	}
}

func hydrateNode(node *core.Node) *core.Node {
	targetNode := node.Target
	switch targetNode.(type) {
	case *core.Node_TaskNode:
		taskNodeWrapper := targetNode.(*core.Node_TaskNode)
		taskNodeReference := taskNodeWrapper.TaskNode.Reference.(*core.TaskNode_ReferenceId)
		hydratedTaskNodeReferenceId := &core.TaskNode_ReferenceId{
			ReferenceId: &core.Identifier{
				ResourceType: taskNodeReference.ReferenceId.ResourceType,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         taskNodeReference.ReferenceId.Name,
				Version:      GetConfig().version,
			},
		}
		hydratedTaskNode := &core.TaskNode{
			Reference: hydratedTaskNodeReferenceId,
		}
		hydratedWrapperTaskNode := &core.Node_TaskNode{
			TaskNode: hydratedTaskNode,
		}
		hydratedNode := &core.Node{
			Id:              node.Id,
			Metadata:        node.Metadata,
			Inputs:          node.Inputs,
			UpstreamNodeIds: node.UpstreamNodeIds,
			OutputAliases:   node.OutputAliases,
			Target:          hydratedWrapperTaskNode,
		}
		return hydratedNode
	case *core.Node_WorkflowNode:
		workflowNodeWrapper := targetNode.(*core.Node_WorkflowNode)
		var hydratedWorkflowNode *core.WorkflowNode
		switch workflowNodeWrapper.WorkflowNode.Reference.(type) {
		case *core.WorkflowNode_SubWorkflowRef:
			subWorkflowNodeReference := workflowNodeWrapper.WorkflowNode.Reference.(*core.WorkflowNode_SubWorkflowRef)
			hydratedSubWorkflowNodeReference := hydrateWorkflowNode_SubWorkflowRef(subWorkflowNodeReference)
			hydratedWorkflowNode = &core.WorkflowNode{
				Reference: hydratedSubWorkflowNodeReference,
			}
		case *core.WorkflowNode_LaunchplanRef:
			launchPlanNodeReference := workflowNodeWrapper.WorkflowNode.Reference.(*core.WorkflowNode_LaunchplanRef)
			hydratedlaunchPlanNodeReference := hydrateWorkflowNode_LaunchplanRef(launchPlanNodeReference)
			hydratedWorkflowNode = &core.WorkflowNode{
				Reference: hydratedlaunchPlanNodeReference,
			}
		}
		hydratedWrapperWorkflowNode := &core.Node_WorkflowNode{
			WorkflowNode: hydratedWorkflowNode,
		}
		hydratedNode := &core.Node{
			Id:              node.Id,
			Metadata:        node.Metadata,
			Inputs:          node.Inputs,
			UpstreamNodeIds: node.UpstreamNodeIds,
			OutputAliases:   node.OutputAliases,
			Target:          hydratedWrapperWorkflowNode,
		}
		return hydratedNode
	}
	return nil
}

func hydrateSpec(ctx context.Context, message proto.Message) proto.Message {
	switch message.(type) {
	case *admin.LaunchPlanSpec:
		launchSpec := message.(*admin.LaunchPlanSpec)
		hydratedLaunchSpec := &admin.LaunchPlanSpec{
			WorkflowId: &core.Identifier{
				ResourceType: launchSpec.WorkflowId.ResourceType,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         launchSpec.WorkflowId.Name,
				Version:      GetConfig().version,
			},
			Auth:                launchSpec.Auth,
			AuthRole:            launchSpec.AuthRole,
			RawOutputDataConfig: launchSpec.RawOutputDataConfig,
			EntityMetadata:      launchSpec.EntityMetadata,
			DefaultInputs:       launchSpec.DefaultInputs,
			FixedInputs:         launchSpec.FixedInputs,
			Labels:              launchSpec.Labels,
			Annotations:         launchSpec.Annotations,
			QualityOfService:    launchSpec.QualityOfService,
			Role:                launchSpec.Role,
		}
		return hydratedLaunchSpec
	case *admin.WorkflowSpec:
		workflowSpec := message.(*admin.WorkflowSpec)
		hydrateNodes := make([]*core.Node, len(workflowSpec.Template.Nodes))
		for index, Noderef := range workflowSpec.Template.Nodes {
			hydrateNodes[index] = hydrateNode(Noderef)
		}
		hydratedWorkflowTemplate := &core.WorkflowTemplate{
			Id: &core.Identifier{
				ResourceType: workflowSpec.Template.Id.ResourceType,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         workflowSpec.Template.Id.Name,
				Version:      GetConfig().version,
			},
			Metadata:         workflowSpec.Template.Metadata,
			Interface:        workflowSpec.Template.Interface,
			Outputs:          workflowSpec.Template.Outputs,
			FailureNode:      workflowSpec.Template.FailureNode,
			MetadataDefaults: workflowSpec.Template.MetadataDefaults,
			Nodes: hydrateNodes,
		}

		hydratedWorkflowSpec := &admin.WorkflowSpec{
			Template: hydratedWorkflowTemplate,
		}
		return hydratedWorkflowSpec
	case *admin.TaskSpec:
		taskSpec := message.(*admin.TaskSpec)
		hydratedTaskTemplate := &core.TaskTemplate{
			Id: &core.Identifier{
				ResourceType: taskSpec.Template.Id.ResourceType,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         taskSpec.Template.Id.Name,
				Version:      GetConfig().version,
			},
			Type: taskSpec.Template.Type,
			Metadata: taskSpec.Template.Metadata,
			Interface: taskSpec.Template.Interface,
			Custom: taskSpec.Template.Custom,
			Target: taskSpec.Template.Target,
		}
		hydratedTaskSpec := &admin.TaskSpec{
			Template: hydratedTaskTemplate,
		}
		return hydratedTaskSpec
	}
	return nil
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
		hydratedSpec := hydrateSpec(ctx, spec)
		err = register(ctx, hydratedSpec, name, cmdCtx)
		if err != nil {
			//logger.Errorf(ctx, "Error registering entity %v with spec %v due to : %v", name, getJsonSpec(hydratedSpec), err)
			logger.Errorf(ctx, "Error registering entity %v due to : %v", name, err)
			continue
		}
		//logger.Infof(ctx, "Registered successfully entity %v with spec %v", name, getJsonSpec(hydratedSpec))
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

func getEntityNameFromPath(absFilePath string) string {
	pathComponents := strings.Split(absFilePath, "/")
	fileName := strings.SplitAfterN(pathComponents[len(pathComponents)-1], "_", 2)
	fileNameWithoutSuffix := strings.TrimSuffix(fileName[len(fileName)-1], ".pb")
	return fileNameWithoutSuffix
}
