package register

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flytestdlib/logger"
)

//go:generate pflags FilesConfig

var (
	filesConfig = &FilesConfig{
		Version:     "v1",
		SkipOnError: false,
	}
)

const registrationProjectPattern = "{{ registration.project }}"
const registrationDomainPattern = "{{ registration.domain }}"
const registrationVersionPattern = "{{ registration.version }}"

// FilesConfig
type FilesConfig struct {
	Version     string `json:"version" pflag:",version of the entity to be registered with flyte."`
	SkipOnError bool   `json:"skipOnError" pflag:",fail fast when registering files."`
}

type Result struct {
	Name   string
	Status string
	Info   string
}

var projectColumns = []printer.Column{
	{Header: "Name", JSONPath: "$.Name"},
	{Header: "Status", JSONPath: "$.Status"},
	{Header: "Additional Info", JSONPath: "$.Info"},
}

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
	launchPlan := &admin.LaunchPlan{}
	if err := proto.Unmarshal(fileContents, launchPlan); err == nil {
		return launchPlan, nil
	}
	logger.Debugf(ctx, "Failed to unmarshal file %v for launch plan type", fname)
	return nil, fmt.Errorf("failed unmarshalling file %v", fname)

}

func register(ctx context.Context, message proto.Message, cmdCtx cmdCore.CommandContext) error {
	switch v := message.(type) {
	case *admin.LaunchPlan:
		launchPlan := message.(*admin.LaunchPlan)
		_, err := cmdCtx.AdminClient().CreateLaunchPlan(ctx, &admin.LaunchPlanCreateRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_LAUNCH_PLAN,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         launchPlan.Id.Name,
				Version:      filesConfig.Version,
			},
			Spec: launchPlan.Spec,
		})
		return err
	case *admin.WorkflowSpec:
		workflowSpec := message.(*admin.WorkflowSpec)
		_, err := cmdCtx.AdminClient().CreateWorkflow(ctx, &admin.WorkflowCreateRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_WORKFLOW,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         workflowSpec.Template.Id.Name,
				Version:      filesConfig.Version,
			},
			Spec: workflowSpec,
		})
		return err
	case *admin.TaskSpec:
		taskSpec := message.(*admin.TaskSpec)
		_, err := cmdCtx.AdminClient().CreateTask(ctx, &admin.TaskCreateRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_TASK,
				Project:      config.GetConfig().Project,
				Domain:       config.GetConfig().Domain,
				Name:         taskSpec.Template.Id.Name,
				Version:      filesConfig.Version,
			},
			Spec: taskSpec,
		})
		return err
	default:
		return fmt.Errorf("Failed registering unknown entity  %v", v)
	}
}

func hydrateNode(node *core.Node) error {
	targetNode := node.Target
	switch v := targetNode.(type) {
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
		default:
			return fmt.Errorf("unknown type %T", workflowNodeWrapper.WorkflowNode.Reference)
		}
	case *core.Node_BranchNode:
		branchNodeWrapper := targetNode.(*core.Node_BranchNode)
		if err := hydrateNode(branchNodeWrapper.BranchNode.IfElse.Case.ThenNode); err != nil {
			return fmt.Errorf("failed to hydrateNode")
		}
		if len(branchNodeWrapper.BranchNode.IfElse.Other) > 0 {
			for _, ifBlock := range branchNodeWrapper.BranchNode.IfElse.Other {
				if err := hydrateNode(ifBlock.ThenNode); err != nil {
					return fmt.Errorf("failed to hydrateNode")
				}
			}
		}
		switch branchNodeWrapper.BranchNode.IfElse.Default.(type) {
		case *core.IfElseBlock_ElseNode:
			elseNodeReference := branchNodeWrapper.BranchNode.IfElse.Default.(*core.IfElseBlock_ElseNode)
			if err := hydrateNode(elseNodeReference.ElseNode); err != nil {
				return fmt.Errorf("failed to hydrateNode")
			}

		case *core.IfElseBlock_Error:
			// Do nothing.
		default:
			return fmt.Errorf("unknown type %T", branchNodeWrapper.BranchNode.IfElse.Default)
		}
	default:
		return fmt.Errorf("unknown type %T", v)
	}
	return nil
}

func hydrateIdentifier(identifier *core.Identifier) {
	if identifier.Project == "" || identifier.Project == registrationProjectPattern {
		identifier.Project = config.GetConfig().Project
	}
	if identifier.Domain == "" || identifier.Domain == registrationDomainPattern {
		identifier.Domain = config.GetConfig().Domain
	}
	if identifier.Version == "" || identifier.Version == registrationVersionPattern {
		identifier.Version = filesConfig.Version
	}
}

func hydrateSpec(message proto.Message) error {
	switch v := message.(type) {
	case *admin.LaunchPlan:
		launchPlan := message.(*admin.LaunchPlan)
		hydrateIdentifier(launchPlan.Spec.WorkflowId)
	case *admin.WorkflowSpec:
		workflowSpec := message.(*admin.WorkflowSpec)
		for _, Noderef := range workflowSpec.Template.Nodes {
			if err := hydrateNode(Noderef); err != nil {
				return err
			}
		}
		hydrateIdentifier(workflowSpec.Template.Id)
		for _, subWorkflow := range workflowSpec.SubWorkflows {
			for _, Noderef := range subWorkflow.Nodes {
				if err := hydrateNode(Noderef); err != nil {
					return err
				}
			}
			hydrateIdentifier(subWorkflow.Id)
		}
	case *admin.TaskSpec:
		taskSpec := message.(*admin.TaskSpec)
		hydrateIdentifier(taskSpec.Template.Id)
	default:
		return fmt.Errorf("Unknown type %T", v)
	}
	return nil
}

func getJSONSpec(message proto.Message) string {
	marshaller := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	jsonSpec, _ := marshaller.MarshalToString(message)
	return jsonSpec
}
