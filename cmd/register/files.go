package register

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
	
Fast Register will register all the fast serialized protobuf files including tasks, workflows and launchplans with default v1 version. Learn more about fast registration  https://docs.flyte.org/projects/cookbook/en/stable/auto/deployment/workflow/fast_registration.html
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks  -v v2 -l "s3://dummy/prefix"  --additionalDistributionPath="s3://dummy/fast" 
	
Fast Register will fail if you didn't pass additional flags like --additionalDistributionPath flags. Fast register o/p only support work with additional flags
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks  -v v2 -l "s3://dummy/prefix" 

	
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
	sourceCodeExtension = ".tar.gz"
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

	// Validate Input files
	sourceCode, validProto, InvalidFiles, err := validateRegisterFiles(dataRefs)
	if err != nil {
		return err
	}

	if len(InvalidFiles) > 0 {
		return fmt.Errorf("input package have some invalid files. try to run pyflyte package again")
	}

	if len(sourceCode) > 0 {
		if err = uploadFastRegisterArtifact(ctx, sourceCode, rconfig.DefaultFilesConfig.AdditionalDistributionPath, rconfig.DefaultFilesConfig.Version); err != nil {
			return fmt.Errorf("please check your Storage Config. It failed while uploading the source code. %v", err)
		}
	}

	fastFail := rconfig.DefaultFilesConfig.ContinueOnError
	for i := 0; i < len(validProto) && !(fastFail && _err != nil); i++ {
		registerResults, _err = registerFile(ctx, validProto[i], registerResults, cmdCtx)
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
