package get

import (
"context"
"github.com/lyft/flytectl/cmd/config"
"encoding/json"
cmdCore "github.com/lyft/flytectl/cmd/core"
"github.com/lyft/flytectl/pkg/printer"
"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
"github.com/lyft/flytestdlib/logger"
)

// PrintableExcution is the structure for printing Excution
type PrintableExcution struct {
	Version          string `header:"Version"`
	Name             string `header:"Name"`
	Type             string `header:"Type"`
	Discoverable     bool   `header:"Discoverable"`
	DiscoveryVersion string `header:"DiscoveryVersion"`
}

var excutionStructure = map[string]string{
	"Version" : "$.id.version",
	"Name" : "$.id.name",
	"Type" : "$.closure.compiledTask.template.type",
	"Discoverable" : "$.closure.compiledTask.template.metadata.discoverable",
	"DiscoveryVersion" : "$.closure.compiledTask.template.metadata.discovery_version",
}

func transformExcution(jsonbody [] byte)(interface{},error){
	results := PrintableExcution{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results,err
	}
	return results,nil
}

func getExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	executionPrinter := printer.Printer{}
	excution, err := cmdCtx.AdminClient().ListExecutions(ctx, &admin.ResourceListRequest{
		Limit: 10,
		Id : &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name:    args[0],
		},
	})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v excution", len(excution.Executions))
	executionPrinter.Print(config.GetConfig().MustOutputFormat(),excution.Executions,excutionStructure,transformExcution)
	return nil
}
