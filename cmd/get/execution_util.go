package get

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flyteidl/clients/go/coreutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"sigs.k8s.io/yaml"
)

// ExecutionConfig is duplicated struct from create with the same structure. This is to avoid the circular dependency.
// TODO : replace this with a cleaner design
type ExecutionConfig struct {
	File            string                 `json:"file,omitempty"`
	TargetDomain    string                 `json:"targetDomain"`
	TargetProject   string                 `json:"targetProject"`
	KubeServiceAcct string                 `json:"kubeServiceAcct"`
	IamRoleARN      string                 `json:"iamRoleARN"`
	Workflow        string                 `json:"workflow,omitempty"`
	Task            string                 `json:"task,omitempty"`
	Version         string                 `json:"version"`
	Inputs          map[string]interface{} `json:"inputs"`
}

func writeExecConfigToFile(executionConfig ExecutionConfig, fileName string) error {
	d, err := yaml.Marshal(executionConfig)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	if _, err = os.Stat(fileName); err == nil {
		if !cmdUtil.AskForConfirmation(fmt.Sprintf("warning file %v will be overwritten", fileName)) {
			return errors.New("backup the file before continuing")
		}
	}
	err = ioutil.WriteFile(fileName, d, 0600)
	if err != nil {
		return fmt.Errorf("unable to write in %v yaml file due to %v", executionConfig.File, err)
	}
	return nil
}

func createAndWriteExecConfigForTask(task *admin.Task, fileName string) error {
	var err error
	executionConfig := ExecutionConfig{Task: task.Id.Name}
	if executionConfig.Inputs, err = getParamMapForTask(task); err != nil {
		return err
	}
	if err = writeExecConfigToFile(executionConfig, fileName); err != nil {
		return err
	}
	return nil
}

func createAndWriteExecConfigForWorkflow(wlp *admin.LaunchPlan, fileName string) error {
	var err error
	executionConfig := ExecutionConfig{Workflow: wlp.Id.Name}
	if executionConfig.Inputs, err = getParamMapForWorkflow(wlp); err != nil {
		return err
	}
	if err = writeExecConfigToFile(executionConfig, fileName); err != nil {
		return err
	}
	return nil
}

func getParamMapForTask(task *admin.Task) (map[string]interface{}, error) {
	paramMap := make(map[string]interface{})
	for k, v := range task.Closure.CompiledTask.Template.Interface.Inputs.Variables {
		varTypeValue, err := coreutils.MakeDefaultLiteralForType(v.Type)
		if err != nil {
			fmt.Println("error creating default value for literal type ", v.Type)
			return nil, err
		}
		if paramMap[k], err = coreutils.ExtractFromLiteral(varTypeValue); err != nil {
			return nil, err
		}
	}
	return paramMap, nil
}

func getParamMapForWorkflow(lp *admin.LaunchPlan) (map[string]interface{}, error) {
	paramMap := make(map[string]interface{})
	for k, v := range lp.Spec.DefaultInputs.Parameters {
		varTypeValue, err := coreutils.MakeDefaultLiteralForType(v.Var.Type)
		if err != nil {
			fmt.Println("error creating default value for literal type ", v.Var.Type)
			return nil, err
		}
		if paramMap[k], err = coreutils.ExtractFromLiteral(varTypeValue); err != nil {
			return nil, err
		}
		// Override if there is a default value
		if paramsDefault, ok := v.Behavior.(*core.Parameter_Default); ok {
			if paramMap[k], err = coreutils.ExtractFromLiteral(paramsDefault.Default); err != nil {
				return nil, err
			}
		}
	}
	return paramMap, nil
}
