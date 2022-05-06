package compile

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/flyteorg/flytectl/cmd/register"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flytepropeller/pkg/compiler"
	"github.com/flyteorg/flytepropeller/pkg/compiler/common"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
)

// Utility function for compiling a list of Tasks
func compileTasks(tasks []*core.TaskTemplate) ([]*core.CompiledTask, error) {
	res := make([]*core.CompiledTask, 0, len(tasks))
	for _, task := range tasks {
		compiledTask, err := compiler.CompileTask(task)
		if err != nil {
			return nil, err
		}
		res = append(res, compiledTask)
	}
	return res, nil
}

/*
 Utility to compile a packaged workflow locally.
 compilation is done locally so no flyte cluster is required.
*/
func compileFromPackage(packagePath string) error {
	args := []string{packagePath}
	fileList, tmpDir, err := register.GetSerializeOutputFiles(context.Background(), args, true)
	defer os.RemoveAll(tmpDir)
	if err != nil {
		fmt.Println("Error found while extracting package..")
		return err
	}
	fmt.Println("Successfuly extracted package...")
	fmt.Println("Processing Protobuff files...")
	workflows := make(map[string]admin.WorkflowSpec)
	plans := make(map[string]admin.LaunchPlan)
	tasks := []admin.TaskSpec{}

	for _, pbFilePath := range fileList {
		// deserializing packaged protocolbuffers
		// based on filename as defined in flytekit
		// https://github.com/flyteorg/flytekit/blob/master/flytekit/tools/serialize_helpers.py#L102
		// filename suffix tell us how to deserialize file
		// _1.pb : are Tasks protobufs
		// _2.pb : are workflows protobufs
		// _3.pb are Launchplans protobufs

		if strings.HasSuffix(pbFilePath, "_1.pb") {
			fmt.Println("Task:", pbFilePath)
			task := admin.TaskSpec{}
			rawTsk, err := ioutil.ReadFile(pbFilePath)
			if err != nil {
				fmt.Printf("error unmarshalling task..")
				return err
			}
			err = proto.Unmarshal(rawTsk, &task)
			tasks = append(tasks, task)
		}

		if strings.HasSuffix(pbFilePath, "_2.pb") {
			fmt.Println("Workflow:", pbFilePath)
			wfSpec := admin.WorkflowSpec{}
			rawWf, err := ioutil.ReadFile(pbFilePath)
			if err != nil {
				return err
			}
			err = proto.Unmarshal(rawWf, &wfSpec)
			if err != nil {
				fmt.Println("error unmarshalling workflow: ", pbFilePath)
				return err
			}
			workflows[wfSpec.Template.Id.Name] = wfSpec
		}
		if strings.HasSuffix(pbFilePath, "_3.pb") {
			fmt.Println("Launch Plan:", pbFilePath)
			var launchPlan admin.LaunchPlan
			rawLp, err := ioutil.ReadFile(pbFilePath)
			if err != nil {
				return err
			}
			err = proto.Unmarshal(rawLp, &launchPlan)
			if err != nil {
				fmt.Println("error unmarshalling Launch Plan: ", pbFilePath)
				return err
			}
			plans[launchPlan.Id.Name] = launchPlan
		}
	}

	// compile tasks
	taskTemplates := []*core.TaskTemplate{}
	for _, task := range tasks {
		taskTemplates = append(taskTemplates, task.Template)
	}

	fmt.Println("Compiling tasks...")
	compiledTasks, err := compileTasks(taskTemplates)
	if err != nil {
		fmt.Println("Error while compiling tasks...")
		return err
	}

	// compile workflows
	for wfName, workflow := range workflows {

		fmt.Println("Compiling workflow:", wfName)
		plan := plans[wfName]

		_, err := compiler.CompileWorkflow(workflow.Template,
			workflow.SubWorkflows,
			compiledTasks,
			[]common.InterfaceProvider{compiler.NewLaunchPlanInterfaceProvider(plan)})
		if err != nil {
			fmt.Println(":( Error Compiling workflow:", wfName)
			return err
		}

	}

	fmt.Println("All Workflows compiled successfully!")
	fmt.Println("Summary:")
	fmt.Println("X workflows found in package")
	fmt.Println("X Tasks found in package")
	fmt.Println("X Launch plans found in package")

	fmt.Println("Workflows:")
	fmt.Println("Tasks:")
	return nil
}

const (
	compileShort = `Validates your flyte packages without registration needed..`
	compileLong  = `Validates your flyte packages without registration needed. Validation is done by compiling your workflow's protocolbuffers files. This ensures typesafeity and composition without needing a running flyte cluster.`
)

func CreateCompileCommand() *cobra.Command {
	var file string
	compile := &cobra.Command{
		Use:   "compile",
		Short: compileShort,
		Long:  compileLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			packageFilePath, _ := cmd.Flags().GetString("file")
			return compileFromPackage(packageFilePath)
		},
	}
	compile.Flags().StringVarP(&file, "file", "f", "", "path to file with packaged workflow..")
	compile.MarkFlagRequired("file")
	return compile
}
