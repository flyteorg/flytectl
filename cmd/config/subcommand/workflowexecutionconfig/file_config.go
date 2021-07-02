package workflowexecutionconfig

import (
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

// TaskResourceAttrFileConfig shadow Config for TaskResourceAttribute.
// The shadow Config is not using ProjectDomainAttribute/Workflowattribute directly inorder to simplify the inputs.
// As the same structure is being used for both ProjectDomainAttribute/Workflowattribute
type WorkflowExecutionConfigFileConfig struct {
	Project  string `json:"project"`
	Domain   string `json:"domain"`
	Workflow string `json:"workflow,omitempty"`
	*admin.WorkflowExecutionConfig
}

// Decorate decorator over TaskResourceAttributes.
func (t WorkflowExecutionConfigFileConfig) Decorate() *admin.MatchingAttributes {
	return &admin.MatchingAttributes{
		Target: &admin.MatchingAttributes_WorkflowExecutionConfig{
			WorkflowExecutionConfig: t.WorkflowExecutionConfig,
		},
	}
}

// UnDecorate to uncover TaskResourceAttributes.
func (t *WorkflowExecutionConfigFileConfig) UnDecorate(matchingAttribute *admin.MatchingAttributes) {
	if matchingAttribute == nil {
		return
	}
	t.WorkflowExecutionConfig = matchingAttribute.GetWorkflowExecutionConfig()
}

// GetProject from the TaskResourceAttrFileConfig
func (t WorkflowExecutionConfigFileConfig) GetProject() string {
	return t.Project
}

// GetDomain from the TaskResourceAttrFileConfig
func (t WorkflowExecutionConfigFileConfig) GetDomain() string {
	return t.Domain
}

// GetWorkflow from the TaskResourceAttrFileConfig
func (t WorkflowExecutionConfigFileConfig) GetWorkflow() string {
	return t.Workflow
}
