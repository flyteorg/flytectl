package get

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flyteidl/clients/go/coreutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

// ExecutionConfig is duplicated struct from create with the same structure. This is to avoid the circular dependency. Only works with go-yaml.
// TODO : replace this with a cleaner design
type ExecutionConfig struct {
	IamRoleARN      string      `yaml:"iamRoleARN"`
	Inputs          []yaml.Node `yaml:"inputs"`
	KubeServiceAcct string      `yaml:"kubeServiceAcct"`
	TargetDomain    string      `yaml:"targetDomain"`
	TargetProject   string      `yaml:"targetProject"`
	Task            string      `yaml:"task,omitempty"`
	Version         string      `yaml:"version"`
	Workflow        string      `yaml:"workflow,omitempty"`
}

func WriteExecConfigToFile(executionConfig ExecutionConfig, fileName string) error {
	d, err := yaml.Marshal(executionConfig)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	if _, err = os.Stat(fileName); err == nil {
		if !cmdUtil.AskForConfirmation(fmt.Sprintf("warning file %v will be overwritten", fileName), os.Stdin) {
			return errors.New("backup the file before continuing")
		}
	}
	return ioutil.WriteFile(fileName, d, 0600)
}

func CreateAndWriteExecConfigForTask(task *admin.Task, fileName string) error {
	var err error
	executionConfig := ExecutionConfig{Task: task.Id.Name, Version: task.Id.Version}
	if executionConfig.Inputs, err = ParamMapForTask(task); err != nil {
		return err
	}
	return WriteExecConfigToFile(executionConfig, fileName)
}

func CreateAndWriteExecConfigForWorkflow(wlp *admin.LaunchPlan, fileName string) error {
	var err error
	executionConfig := ExecutionConfig{Workflow: wlp.Id.Name, Version: wlp.Id.Version}
	if executionConfig.Inputs, err = ParamMapForWorkflow(wlp); err != nil {
		return err
	}
	return WriteExecConfigToFile(executionConfig, fileName)
}

func TaskInputs(task *admin.Task) []*core.VariableMapEntry {
	taskInputs := []*core.VariableMapEntry{}
	if task == nil || task.Closure == nil {
		return taskInputs
	}
	if task.Closure.CompiledTask == nil {
		return taskInputs
	}
	if task.Closure.CompiledTask.Template == nil {
		return taskInputs
	}
	if task.Closure.CompiledTask.Template.Interface == nil {
		return taskInputs
	}
	if task.Closure.CompiledTask.Template.Interface.Inputs == nil {
		return taskInputs
	}
	return task.Closure.CompiledTask.Template.Interface.Inputs.Variables
}

func ParamMapForTask(task *admin.Task) ([]yaml.Node, error) {
	taskInputs := TaskInputs(task)
	paramMap := make([]yaml.Node, 0, len(taskInputs))
	for _, e := range taskInputs {
		varTypeValue, err := coreutils.MakeDefaultLiteralForType(e.Var.Type)
		if err != nil {
			fmt.Println("error creating default value for literal type ", e.Var.Type)
			return nil, err
		}
		var nativeLiteral interface{}
		if nativeLiteral, err = coreutils.ExtractFromLiteral(varTypeValue); err != nil {
			return nil, err
		}

		if e.Name == e.Var.Description {
			// a: # a isn't very helpful
			var n yaml.Node
			n, err = getCommentedYamlNode(nativeLiteral, "")
			paramMap = append(paramMap, n)
		} else {
			var n yaml.Node
			n, err = getCommentedYamlNode(nativeLiteral, e.Var.Description)
			paramMap = append(paramMap, n)
		}
		if err != nil {
			return nil, err
		}
	}
	return paramMap, nil
}

func WorkflowParams(lp *admin.LaunchPlan) []*core.ParameterMapEntry {
	workflowParams := []*core.ParameterMapEntry{}
	if lp == nil || lp.Spec == nil {
		return workflowParams
	}
	if lp.Spec.DefaultInputs == nil {
		return workflowParams
	}
	return lp.Spec.DefaultInputs.Parameters
}

func ParamMapForWorkflow(lp *admin.LaunchPlan) ([]yaml.Node, error) {
	workflowParams := WorkflowParams(lp)
	paramMap := make([]yaml.Node, len(workflowParams))
	for _, e := range workflowParams {
		varTypeValue, err := coreutils.MakeDefaultLiteralForType(e.Var.Var.Type)
		if err != nil {
			fmt.Println("error creating default value for literal type ", e.Var.Var.Type)
			return nil, err
		}
		var nativeLiteral interface{}
		if nativeLiteral, err = coreutils.ExtractFromLiteral(varTypeValue); err != nil {
			return nil, err
		}
		// Override if there is a default value
		if paramsDefault, ok := e.Var.Behavior.(*core.Parameter_Default); ok {
			if nativeLiteral, err = coreutils.ExtractFromLiteral(paramsDefault.Default); err != nil {
				return nil, err
			}
		}
		if e.Name == e.Var.Var.Description {
			// a: # a isn't very helpful
			var n yaml.Node
			n, err = getCommentedYamlNode(nativeLiteral, "")
			paramMap = append(paramMap, n)
		} else {
			var n yaml.Node
			n, err = getCommentedYamlNode(nativeLiteral, e.Var.Var.Description)
			paramMap = append(paramMap, n)
		}

		if err != nil {
			return nil, err
		}
	}
	return paramMap, nil
}

func getCommentedYamlNode(input interface{}, comment string) (yaml.Node, error) {
	var node yaml.Node
	err := node.Encode(input)
	node.LineComment = comment
	return node, err
}
