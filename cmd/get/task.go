package get

import (
	"context"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/adminutils"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"

	"github.com/golang/protobuf/proto"
)

const (
	taskShort = "Gets task resources"
	taskLong  = `
Retrieves all the task within project and domain.(task,tasks can be used interchangeably in these commands)
::

 bin/flytectl get task -p flytesnacks -d development

Retrieves task by name within project and domain.

::

 bin/flytectl task -p flytesnacks -d development core.basic.lp.greet

Retrieves project by filters.
::

 Not yet implemented

Retrieves all the tasks within project and domain in yaml format.

::

 bin/flytectl get task -p flytesnacks -d development -o yaml

Retrieves all the tasks within project and domain in json format.

::

 bin/flytectl get task -p flytesnacks -d development -o json

Usage
`
)

//go:generate pflags TaskConfig --default-var taskConfig
var (
	taskConfig = &TaskConfig{}
)

// FilesConfig
type TaskConfig struct {
	ExecFile string `json:"execFile" pflag:",execution file name to be used for generating execution spec of a single task."`
	Version  string `json:"version" pflag:",version of the task to be fetched."`
}

var taskColumns = []printer.Column{
	{Header: "Version", JSONPath: "$.id.version"},
	{Header: "Name", JSONPath: "$.id.name"},
	{Header: "Type", JSONPath: "$.closure.compiledTask.template.type"},
	{Header: "Discoverable", JSONPath: "$.closure.compiledTask.template.metadata.discoverable"},
	{Header: "Discovery Version", JSONPath: "$.closure.compiledTask.template.metadata.discoveryVersion"},
	{Header: "Created At", JSONPath: "$.closure.createdAt"},
}

func TaskToProtoMessages(l []*admin.Task) []proto.Message {
	messages := make([]proto.Message, 0, len(l))
	for _, m := range l {
		messages = append(messages, m)
	}
	return messages
}

func getTaskFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	taskPrinter := printer.Printer{}
	if len(args) == 1 {
		var tasks []*admin.Task
		var err error
		// Right only support writing execution file for single task version.
		if taskConfig.ExecFile != "" {
			var task *admin.Task
			if task, err = FetchTaskVersionOrLatest(ctx, args[0], taskConfig.Version, cmdCtx); err != nil {
				return err
			}
			tasks = append(tasks, task)
			if err = createAndWriteExecConfigForTask(task, taskConfig.ExecFile); err != nil {
				return err
			}
		} else {
			tasks, err = getAllVerOfTask(ctx, args[0], cmdCtx)
			if err != nil {
				return err
			}
		}
		logger.Debugf(ctx, "Retrieved Task", tasks)
		return taskPrinter.Print(config.GetConfig().MustOutputFormat(), taskColumns, TaskToProtoMessages(tasks)...)
	}
	tasks, err := adminutils.GetAllNamedEntities(ctx, cmdCtx.AdminClient().ListTaskIds, adminutils.ListRequest{Project: config.GetConfig().Project, Domain: config.GetConfig().Domain})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v Task", len(tasks))
	return taskPrinter.Print(config.GetConfig().MustOutputFormat(), entityColumns, adminutils.NamedEntityToProtoMessage(tasks)...)
}
