package get

import (
	"context"
	"encoding/json"
	"github.com/lyft/flytectl/cmd/config"
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
	"Version":          "$.id.version",
	"Name":             "$.id.name",
	"Type":             "$.closure.compiledTask.template.type",
	"Discoverable":     "$.closure.compiledTask.template.metadata.discoverable",
	"DiscoveryVersion": "$.closure.compiledTask.template.metadata.discovery_version",
}

func transformExcution(jsonbody []byte) (interface{}, error) {
	results := PrintableExcution{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}

func getExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	executionPrinter := printer.Printer{}
	excutions, err := cmdCtx.AdminClient().ListExecutions(ctx, &admin.ResourceListRequest{
		Limit: 10,
		Id: &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		},
	})
	if err != nil {
		return err
	}
	if len(args) == 1 {
		name := args[0]
		if err != nil {
			return err
		}
		logger.Debugf(ctx, "Retrieved excutions")
		for _, v := range excutions.Executions {
			if v.Id.Name == name {
				err := executionPrinter.Print(config.GetConfig().MustOutputFormat(), v, excutionStructure, transformExcution)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}

	logger.Debugf(ctx, "Retrieved %v excutions", len(excutions.Executions))
	executionPrinter.Print(config.GetConfig().MustOutputFormat(), excutions.Executions, excutionStructure, transformExcution)
	return nil
}
