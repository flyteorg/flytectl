package register

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	rconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/register"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flytestdlib/logger"
)

const (
	registerFilesShort = "Registers file resources"
	registerFilesLong  = `
Registers all the serialized protobuf files including tasks, workflows and launchplans with default v1 version.
If there are already registered entities with v1 version then the command will fail immediately on the first such encounter.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks
	
Fast Register will register all the fast serialized protobuf files including tasks, workflows and launchplans with default v1 version. Learn more about registration  https://docs.flyte.org/projects/cookbook/en/stable/auto/deployment/workflow/fast_registration.html
Fast Register required --additionalDistributionDir and --destinationDir flags 	
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks  -v v2 -l "s3://dummy/prefix" --destinationDir="" --additionalDistributionDir="s3://dummy/fast" 

	
Using archive file.Currently supported are .tgz and .tar extension files and can be local or remote file served through http/https.
Use --archive flag.

::

 bin/flytectl register files  http://localhost:8080/_pb_output.tar -d development  -p flytesnacks --archive

Using  local tgz file.

::

 bin/flytectl register files  _pb_output.tgz -d development  -p flytesnacks --archive

If you want to continue executing registration on other files ignoring the errors including version conflicts then pass in
the continueOnError flag.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError

Using short format of continueOnError flag
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError

Overriding the default version v1 using version string.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks -v v2

Change the o/p format has not effect on registration. The O/p is currently available only in table format.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError -o yaml

Override IamRole during registration.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError -v v2 -i "arn:aws:iam::123456789:role/dummy"

Override Kubernetes service account during registration.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError -v v2 -k "kubernetes-service-account"

Override Output location prefix during registration.

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError -v v2 -l "s3://dummy/prefix"
	
Usage
`
)

func registerFromFilesFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	return Register(ctx, args, cmdCtx)
}

func Register(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	var _err error
	var dataRefs []string
	var tmpDir string
	var registerResults []Result

	// getSerializeOutputFiles will return you all proto and  source code compress file in sorted order
	dataRefs, tmpDir, err := getSerializeOutputFiles(ctx, args)
	if err != nil {
		logger.Errorf(ctx, "error while un-archiving files in tmp dir due to %v", err)
		return err
	}
	logger.Infof(ctx, "Parsing file... Total(%v)", len(dataRefs))

	fastFail := rconfig.DefaultFilesConfig.ContinueOnError
	for i := 0; i < len(dataRefs) && !(fastFail && _err != nil); i++ {
		if len(rconfig.DefaultFilesConfig.AdditionalDistributionDir) > 0 && len(rconfig.DefaultFilesConfig.DestinationDir) > 0 && strings.HasSuffix(dataRefs[i], ".tar.gz") {
			logger.Infof(ctx, "Fast register started for file %v", dataRefs[i])
			if _err = uploadFastRegisterArtifact(ctx, dataRefs[i], rconfig.DefaultFilesConfig.AdditionalDistributionDir, rconfig.DefaultFilesConfig.Version); _err != nil {
				registerResults = append(registerResults, Result{Name: dataRefs[i], Status: "Failed", Info: "Failed while uploading the source code"})
			}
		} else if strings.HasSuffix(dataRefs[i], ".pb") {
			registerResults, _err = registerFile(ctx, dataRefs[i], registerResults, cmdCtx)
		} else {
			registerResults = append(registerResults, Result{Name: dataRefs[i], Status: "Failed", Info: "Invalid files"})
		}
	}

	payload, _ := json.Marshal(registerResults)
	registerPrinter := printer.Printer{}
	_ = registerPrinter.JSONToTable(payload, projectColumns)
	if tmpDir != "" {
		if _err = os.RemoveAll(tmpDir); _err != nil {
			logger.Errorf(ctx, "unable to delete temp dir %v due to %v", tmpDir, _err)
		}
	}
	return _err
}
