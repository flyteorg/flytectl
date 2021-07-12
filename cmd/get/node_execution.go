package get

import (
	"bytes"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"sort"
	"strconv"
	"strings"

	"github.com/disiqueira/gotree"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flyteidl/clients/go/coreutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
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

func getExecutionDetails(ctx context.Context, project, domain, execName, nodeName string, cmdCtx cmdCore.CommandContext) ([]*NodeExecutionClosure, error) {
	// Fetching Node execution details
	nodeExecDetailsMap := map[string]*NodeExecutionClosure{}
	nExecDetails, err := getNodeExecDetailsInt(ctx, project, domain, execName, nodeName, "", nodeExecDetailsMap, cmdCtx)
	if err != nil {
		return nil, err
	}

	var nExecDetailsForView []*NodeExecutionClosure
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

type TaskExecution struct {
	*admin.TaskExecution
}
//
//func (in *TaskExecution) MarshalJSON() ([]byte, error) {
//	var buf bytes.Buffer
//	marshaller := jsonpb.Marshaler{}
//	if err := marshaller.Marshal(&buf, in.TaskExecution); err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}
//
//func (in *TaskExecution) UnmarshalJSON(b []byte) error {
//	in.TaskExecution = &admin.TaskExecution{}
//	return jsonpb.Unmarshal(bytes.NewReader(b), in.TaskExecution)
//}

type NodeExecution struct {
	*admin.NodeExecution
}

func (in *NodeExecutionClosure) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	marshaller := jsonpb.Marshaler{}
	if err := marshaller.Marshal(&buf, in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (in *NodeExecutionClosure) UnmarshalJSON(b []byte) error {
	*in = NodeExecutionClosure{}
	return jsonpb.Unmarshal(bytes.NewReader(b), in)
}

type NodeExecutionClosure struct {
	*NodeExecution
	ChildNodes     []*NodeExecutionClosure `json:"child_nodes,omitempty"`
	TaskExecutions []*TaskExecutionClosure `json:"task_execs,omitempty"`
	// Inputs for the node
	Inputs map[string]interface{} `json:"inputs,omitempty"`
	// Outputs for the node
	Outputs map[string]interface{} `json:"outputs,omitempty"`
}

type TaskExecutionClosure struct {
	*TaskExecution
}

/*
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
*/

func getNodeExecDetailsInt(ctx context.Context, project, domain, execName, nodeName, uniqueParentID string,
	nodeExecDetailsMap map[string]*NodeExecutionClosure, cmdCtx cmdCore.CommandContext) ([]*NodeExecutionClosure, error) {

	nExecDetails, err := cmdCtx.AdminFetcherExt().FetchNodeExecutionDetails(ctx, execName, project, domain, uniqueParentID)
	if err != nil {
		return nil, err
	}

	var nodeExecClosures []*NodeExecutionClosure
	for _, nodeExec := range nExecDetails.NodeExecutions {
		nodeExecClosure := &NodeExecutionClosure{
			NodeExecution: &NodeExecution{nodeExec},
		}
		nodeExecClosures = append(nodeExecClosures, nodeExecClosure)

		// Check if this is parent node. If yes do recursive call to get child nodes.
		if nodeExec.Metadata != nil && nodeExec.Metadata.IsParentNode {
			nodeExecClosure.ChildNodes, err = getNodeExecDetailsInt(ctx, project, domain, execName, nodeName, nodeExec.Id.NodeId, nodeExecDetailsMap, cmdCtx)
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
				taskExecClosure := &TaskExecutionClosure{
					TaskExecution: &TaskExecution{taskExec},
				}
				nodeExecClosure.TaskExecutions = append(nodeExecClosure.TaskExecutions, taskExecClosure)
			}
			// Fetch the node inputs and outputs
			nExecDataResp, err := cmdCtx.AdminFetcherExt().FetchNodeExecutionData(ctx, nodeExec.Id.NodeId, execName, project, domain)
			if err != nil {
				return nil, err
			}
			// Extract the inputs from the literal map
			nodeExecClosure.Inputs, err = extractLiteralMap(nExecDataResp.FullInputs)
			if err != nil {
				return nil, err
			}
			// Extract the outputs from the literal map
			nodeExecClosure.Outputs, err = extractLiteralMap(nExecDataResp.FullOutputs)
			if err != nil {
				return nil, err
			}
		}
		nodeExecDetailsMap[nodeExec.Id.NodeId] = nodeExecClosure
		// Found the node
		if len(nodeName) > 0 && nodeName == nodeExec.Id.NodeId {
			return nodeExecClosures, err
		}
	}
	return nodeExecClosures, nil
}

func createNodeTaskExecTreeView(rootView gotree.Tree, taskExecClosures []*TaskExecutionClosure) {
	if len(taskExecClosures) == 0 {
		return
	}
	if rootView == nil {
		rootView = gotree.New("")
	}
	// TODO: Replace this by filter to sort in the admin
	sort.Slice(taskExecClosures[:], func(i, j int) bool {
		return taskExecClosures[i].Id.RetryAttempt < taskExecClosures[j].Id.RetryAttempt
	})
	for _, taskExecClosure := range taskExecClosures {
		attemptView := rootView.Add(taskAttemptPrefix + strconv.Itoa(int(taskExecClosure.Id.RetryAttempt)))
		attemptView.Add(taskExecPrefix + taskExecClosure.Closure.Phase.String() +
			hyphenPrefix + taskExecClosure.Closure.StartedAt.String() +
			hyphenPrefix + taskExecClosure.Closure.Duration.String())
		attemptView.Add(taskTypePrefix + taskExecClosure.Closure.TaskType)
		attemptView.Add(taskReasonPrefix + taskExecClosure.Closure.Reason)
		if taskExecClosure.Closure.Metadata != nil {
			metadata := attemptView.Add(taskMetadataPrefix)
			metadata.Add(taskGeneratedNamePrefix + taskExecClosure.Closure.Metadata.GeneratedName)
			metadata.Add(taskPluginIDPrefix + taskExecClosure.Closure.Metadata.PluginIdentifier)
			extResourcesView := metadata.Add(taskExtResourcesPrefix)
			for _, extResource := range taskExecClosure.Closure.Metadata.ExternalResources {
				extResourcesView.Add(taskExtResourcePrefix + extResource.ExternalId)
			}
			resourcePoolInfoView := metadata.Add(taskResourcePrefix)
			for _, rsPool := range taskExecClosure.Closure.Metadata.ResourcePoolInfo {
				resourcePoolInfoView.Add(taskExtResourcePrefix + rsPool.Namespace)
				resourcePoolInfoView.Add(taskExtResourceTokenPrefix + rsPool.AllocationToken)
			}
		}

		sort.Slice(taskExecClosure.Closure.Logs[:], func(i, j int) bool {
			return taskExecClosure.Closure.Logs[i].Name < taskExecClosure.Closure.Logs[j].Name
		})

		logsView := attemptView.Add(taskLogsPrefix)
		for _, logData := range taskExecClosure.Closure.Logs {
			logsView.Add(taskLogsNamePrefix + logData.Name)
			logsView.Add(taskLogURIPrefix + logData.Uri)
		}
	}
}

func createNodeDetailsTreeView(rootView gotree.Tree, nodeExecutionClosures []*NodeExecutionClosure) gotree.Tree {
	if rootView == nil {
		rootView = gotree.New("")
	}
	if len(nodeExecutionClosures) == 0 {
		return rootView
	}
	// TODO : Move to sorting using filters.
	sort.Slice(nodeExecutionClosures[:], func(i, j int) bool {
		return nodeExecutionClosures[i].Closure.StartedAt.AsTime().Before(nodeExecutionClosures[j].Closure.StartedAt.AsTime())
	})

	for _, nodeExecWrapper := range nodeExecutionClosures {
		nExecView := rootView.Add(nodeExecWrapper.Id.NodeId + hyphenPrefix + nodeExecWrapper.Closure.Phase.String() +
			hyphenPrefix + nodeExecWrapper.Closure.StartedAt.String() +
			hyphenPrefix + nodeExecWrapper.Closure.Duration.String())
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
