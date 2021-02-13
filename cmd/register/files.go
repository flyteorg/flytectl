package register

import (
	"context"
	"encoding/json"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flytestdlib/logger"
	"io/ioutil"
	"os"
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
	dataRefs, tmpDir, err := getSortedFileList(ctx, args)
	if err != nil {
		logger.Errorf(ctx, "error while un-archiving files in tmp dir due to %v", err)
		return nil
	}
	logger.Infof(ctx, "Parsing files... Total(%v)", len(dataRefs))
	fastFail := !filesConfig.SkipOnError
	var _err error
	var registerResults [] Result
	for i := 0; i < len(dataRefs) && !(fastFail && _err != nil) ; i++ {
		fileContents, err := ioutil.ReadFile(dataRefs[i])
		if err != nil {
			_err = err
			continue
		}
		registerResults, _err = registerContent(ctx, fileContents, dataRefs[i], registerResults, cmdCtx)
	}
	payload, _ := json.Marshal(registerResults)
	registerPrinter := printer.Printer{}
	registerPrinter.JSONToTable(payload, projectColumns)
	if _err = os.RemoveAll(tmpDir); _err != nil {
		logger.Errorf(ctx, "unable to delete temp dir %v due to %v", tmpDir, _err)
	}
	return nil
}
