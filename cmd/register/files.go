package register

import (
	"context"
	"encoding/json"
	"io"
	"sort"

	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flytestdlib/logger"
)

const (
	registerFilesShort = "Registers file resources"
	registerFilesLong  = `
Registers all the serialized protobuf files including tasks, workflows and launchplans with default v1 version.
If there are already registered entities with v1 version then the command will fail immediately on the first such encounter.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks

If you want to continue executing registration on other files ignoring the errors including version conflicts then pass in
the skipOnError flag.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --skipOnError

Using short format of skipOnError flag
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks -s

Overriding the default version v1 using version string.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks -v v2

Change the o/p format has not effect on registration. The O/p is currently available only in table format.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks -s -o yaml

Usage
`
)

func registerFromFilesFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	dataRefs := args
	registerPrinter := printer.Printer{}
	sort.Strings(dataRefs)
	logger.Infof(ctx, "Parsing files... Total(%v)", len(dataRefs))
	logger.Infof(ctx, "Registering with version %v", filesConfig.Version)
	var _err error
	fastFail := !filesConfig.SkipOnError
	isArchive := filesConfig.Archive
	var registerResults [] Result
	logger.Infof(ctx, "Archive flag state =  %v ", isArchive)
	for i := 0; i < len(dataRefs) && !(fastFail && _err != nil) ; i++ {
		dataRefReader, err := getReader(ctx, dataRefs[i])
		if err != nil {
			_err = err
			continue
		}
		for {
			name, fileContents, err := readContents(ctx, dataRefReader, isArchive)
			if err == nil {
				if fileContents == nil {
					// Mostly directory
					continue
				}
				registerResults, _err = registerContent(ctx, fileContents, name, registerResults, cmdCtx)
			} else {
				_err = err
				if err != io.EOF {
					logger.Errorf(ctx,"failed to readContents from reader due to %v", err)
				}
			}
			if !isArchive || (fastFail && _err != nil) || _err == io.EOF {
				break
			}
		}
	}
	payload, _ := json.Marshal(registerResults)
	registerPrinter.JSONToTable(payload, projectColumns)
	return nil
}
