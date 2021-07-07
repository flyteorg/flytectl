package get

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/disiqueira/gotree"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flyteidl/clients/go/coreutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/event"
	_struct "github.com/golang/protobuf/ptypes/struct"
)

var nodeExecutionColumns = []printer.Column{
	{Header: "Name", JSONPath: "$.id.nodeID"},
	{Header: "Exec", JSONPath: "$.id.executionId.name"},
	{Header: "EndedAt", JSONPath: "$.endedAt"},
	{Header: "StartedAt", JSONPath: "$.startedAt"},
	{Header: "Phase", JSONPath: "$.phase"},
}

const (
	taskAttemptPrefix          = "Attempt :"
	taskExecPrefix             = "Task - "
	taskTypePrefix             = "Task Type - "
	taskReasonPrefix           = "Reason - "
	taskMetadataPrefix         = "Metadata"
	taskGeneratedNamePrefix    = "Generated Name : "
	taskPluginIDPrefix         = "Plugin Identifier : "
	taskExtResourcesPrefix     = "External Resources"
	taskExtResourcePrefix      = "Ext Resource : "
	taskExtResourceTokenPrefix = "Ext Resource Token : " //nolint
	taskResourcePrefix         = "Resource Pool Info"
	taskLogsPrefix             = "Logs :"
	taskLogsNamePrefix         = "Name :"
	taskLogURIPrefix           = "URI :"
	hyphenPrefix               = " - "
)

func getExecutionDetails(ctx context.Context, project, domain, execName, nodeName string, cmdCtx cmdCore.CommandContext) ([]*NodeExecutionWrapper, error) {
	// Fetching Node execution details
	nodeExecDetailsMap := map[string]*NodeExecutionWrapper{}
	nExecDetails, err := getNodeExecDetailsInt(ctx, project, domain, execName, nodeName, "", nodeExecDetailsMap, cmdCtx)
	if err != nil {
		return nil, err
	}

	var nExecDetailsForView []*NodeExecutionWrapper
	// Get the execution details only for the nodeId passed
	if len(nodeName) > 0 {
		// Fetch the last one which contains the nodeId details as previous ones are used to reach the nodeId
		if nodeExecDetailsMap[nodeName] != nil {
			nExecDetailsForView = append(nExecDetailsForView, nodeExecDetailsMap[nodeName])
		}
	} else {
		nExecDetailsForView = nExecDetails
	}
	return nExecDetailsForView, nil
}

type StatusTrackerWrapper struct {
	Phase string `json:"phase,omitempty"`
	// Time at which the node execution began running.
	StartedAt time.Time `json:"started_at,omitempty"`
	// The amount of time the node execution spent running.
	EndedAt time.Time `json:"ended_at,omitempty"`
	// Time at which the node execution was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Time at which the node execution was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type NodeExecutionWrapper struct {
	ID *core.NodeExecutionIdentifier `json:"id,omitempty"`
	// Inputs for the node
	Inputs map[string]interface{} `json:"inputs,omitempty"`
	// Outputs for the node
	Outputs    map[string]interface{} `json:"outputs,omitempty"`
	RetryGroup string                 `json:"retry_group,omitempty"`

	ChildNodes []*NodeExecutionWrapper `json:"child_nodes,omitempty"`
	SpecNodeID string                  `json:"spec_node_id,omitempty"`

	TaskExecutions []*TaskExecutionWrapper `json:"task_execs,omitempty"`

	StatusTrackerWrapper
}

type TaskExecutionWrapper struct {
	ID *core.TaskExecutionIdentifier `json:"id,omitempty"`
	// Detailed log information output by the task execution.
	Logs []*core.TaskLog `json:"logs,omitempty"`
	// Custom data specific to the task plugin.
	CustomInfo *_struct.Struct `json:"custom_info,omitempty"`
	Reason     string          `json:"reason,omitempty"`
	// A predefined yet extensible Task type identifier.
	TaskType string `json:"task_type,omitempty"`
	// Metadata around how a task was executed.
	Metadata *event.TaskExecutionMetadata `json:"metadata,omitempty"`

	StatusTrackerWrapper
}

func transformToWrapperNode(entity *admin.NodeExecution) NodeExecutionWrapper {
	statusTracker := StatusTrackerWrapper{
		Phase:     entity.Closure.Phase.String(),
		StartedAt: entity.Closure.StartedAt.AsTime(),
		EndedAt:   entity.Closure.StartedAt.AsTime().Add(entity.Closure.Duration.AsDuration()),
		CreatedAt: entity.Closure.CreatedAt.AsTime(),
		UpdatedAt: entity.Closure.UpdatedAt.AsTime(),
	}
	return NodeExecutionWrapper{
		ID:                   entity.Id,
		RetryGroup:           entity.Metadata.RetryGroup,
		SpecNodeID:           entity.Metadata.SpecNodeId,
		StatusTrackerWrapper: statusTracker,
	}
}

func transformToWrapperTask(entity *admin.TaskExecution) TaskExecutionWrapper {
	statusTracker := StatusTrackerWrapper{
		Phase:     entity.Closure.Phase.String(),
		StartedAt: entity.Closure.StartedAt.AsTime(),
		EndedAt:   entity.Closure.StartedAt.AsTime().Add(entity.Closure.Duration.AsDuration()),
		CreatedAt: entity.Closure.CreatedAt.AsTime(),
		UpdatedAt: entity.Closure.UpdatedAt.AsTime(),
	}
	return TaskExecutionWrapper{
		ID:                   entity.Id,
		TaskType:             entity.Closure.TaskType,
		Reason:               entity.Closure.Reason,
		Logs:                 entity.Closure.Logs,
		Metadata:             entity.Closure.Metadata,
		StatusTrackerWrapper: statusTracker,
	}
}

func getNodeExecDetailsInt(ctx context.Context, project, domain, execName, nodeName, uniqueParentID string,
	nodeExecDetailsMap map[string]*NodeExecutionWrapper, cmdCtx cmdCore.CommandContext) ([]*NodeExecutionWrapper, error) {

	nExecDetails, err := cmdCtx.AdminFetcherExt().FetchNodeExecutionDetails(ctx, execName, project, domain, uniqueParentID)
	if err != nil {
		return nil, err
	}

	var nodeExecWrappers []*NodeExecutionWrapper
	for _, nodeExec := range nExecDetails.NodeExecutions {
		nodeExecWrapper := transformToWrapperNode(nodeExec)
		nodeExecWrappers = append(nodeExecWrappers, &nodeExecWrapper)

		// Check if this is parent node. If yes do recursive call to get child nodes.
		if nodeExec.Metadata != nil && nodeExec.Metadata.IsParentNode {
			nodeExecWrapper.ChildNodes, err = getNodeExecDetailsInt(ctx, project, domain, execName, nodeName, nodeExec.Id.NodeId, nodeExecDetailsMap, cmdCtx)
			if err != nil {
				return nil, err
			}
		} else {
			// Bug in admin https://github.com/flyteorg/flyte/issues/1221
			if strings.HasSuffix(nodeExec.Id.NodeId, "start-node") {
				continue
			}
			taskExecList, err := cmdCtx.AdminFetcherExt().FetchTaskExecutionsOnNode(ctx,
				nodeExec.Id.NodeId, execName, project, domain)
			if err != nil {
				return nil, err
			}
			for _, taskExec := range taskExecList.TaskExecutions {
				t := transformToWrapperTask(taskExec)
				nodeExecWrapper.TaskExecutions = append(nodeExecWrapper.TaskExecutions, &t)
			}
			// Fetch the node inputs and outputs
			nExecDataResp, err := cmdCtx.AdminFetcherExt().FetchNodeExecutionData(ctx, nodeExec.Id.NodeId, execName, project, domain)
			if err != nil {
				return nil, err
			}
			// Extract the inputs from the literal map
			nodeExecWrapper.Inputs, err = extractLiteralMap(nExecDataResp.FullInputs)
			if err != nil {
				return nil, err
			}
			// Extract the outputs from the literal map
			nodeExecWrapper.Outputs, err = extractLiteralMap(nExecDataResp.FullOutputs)
			if err != nil {
				return nil, err
			}
		}
		nodeExecDetailsMap[nodeExec.Id.NodeId] = &nodeExecWrapper
		// Found the node
		if len(nodeName) > 0 && nodeName == nodeExec.Id.NodeId {
			return nodeExecWrappers, err
		}
	}
	return nodeExecWrappers, nil
}

func createNodeTaskExecTreeView(rootView gotree.Tree, taskExecWrappers []*TaskExecutionWrapper) {
	if len(taskExecWrappers) == 0 {
		return
	}
	if rootView == nil {
		rootView = gotree.New("")
	}
	// TODO: Replace this by filter to sort in the admin
	sort.Slice(taskExecWrappers[:], func(i, j int) bool {
		return taskExecWrappers[i].ID.RetryAttempt < taskExecWrappers[j].ID.RetryAttempt
	})
	for _, taskExecWrapper := range taskExecWrappers {
		attemptView := rootView.Add(taskAttemptPrefix + strconv.Itoa(int(taskExecWrapper.ID.RetryAttempt)))
		attemptView.Add(taskExecPrefix + taskExecWrapper.Phase +
			hyphenPrefix + taskExecWrapper.StartedAt.String() +
			hyphenPrefix + taskExecWrapper.EndedAt.String())
		attemptView.Add(taskTypePrefix + taskExecWrapper.TaskType)
		attemptView.Add(taskReasonPrefix + taskExecWrapper.Reason)
		if taskExecWrapper.Metadata != nil {
			metadata := attemptView.Add(taskMetadataPrefix)
			metadata.Add(taskGeneratedNamePrefix + taskExecWrapper.Metadata.GeneratedName)
			metadata.Add(taskPluginIDPrefix + taskExecWrapper.Metadata.PluginIdentifier)
			extResourcesView := metadata.Add(taskExtResourcesPrefix)
			for _, extResource := range taskExecWrapper.Metadata.ExternalResources {
				extResourcesView.Add(taskExtResourcePrefix + extResource.ExternalId)
			}
			resourcePoolInfoView := metadata.Add(taskResourcePrefix)
			for _, rsPool := range taskExecWrapper.Metadata.ResourcePoolInfo {
				resourcePoolInfoView.Add(taskExtResourcePrefix + rsPool.Namespace)
				resourcePoolInfoView.Add(taskExtResourceTokenPrefix + rsPool.AllocationToken)
			}
		}

		sort.Slice(taskExecWrapper.Logs[:], func(i, j int) bool {
			return taskExecWrapper.Logs[i].Name < taskExecWrapper.Logs[j].Name
		})

		logsView := attemptView.Add(taskLogsPrefix)
		for _, logData := range taskExecWrapper.Logs {
			logsView.Add(taskLogsNamePrefix + logData.Name)
			logsView.Add(taskLogURIPrefix + logData.Uri)
		}
	}
}

func createNodeDetailsTreeView(rootView gotree.Tree, nodeExecutionWrappers []*NodeExecutionWrapper) gotree.Tree {
	if rootView == nil {
		rootView = gotree.New("")
	}
	if len(nodeExecutionWrappers) == 0 {
		return rootView
	}
	// TODO : Move to sorting using filters.
	sort.Slice(nodeExecutionWrappers[:], func(i, j int) bool {
		return nodeExecutionWrappers[i].StartedAt.Before(nodeExecutionWrappers[j].StartedAt)
	})

	for _, nodeExecWrapper := range nodeExecutionWrappers {
		nExecView := rootView.Add(nodeExecWrapper.ID.NodeId + hyphenPrefix + nodeExecWrapper.Phase +
			hyphenPrefix + nodeExecWrapper.StartedAt.String() +
			hyphenPrefix + nodeExecWrapper.EndedAt.String())
		if len(nodeExecWrapper.ChildNodes) > 0 {
			createNodeDetailsTreeView(nExecView, nodeExecWrapper.ChildNodes)
		}
		createNodeTaskExecTreeView(nExecView, nodeExecWrapper.TaskExecutions)
	}
	return rootView
}

func extractLiteralMap(literalMap *core.LiteralMap) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for key, literalVal := range literalMap.Literals {
		extractedLiteralVal, err := coreutils.ExtractFromLiteral(literalVal)
		if err != nil {
			return nil, err
		}
		m[key] = extractedLiteralVal
	}
	return m, nil
}
