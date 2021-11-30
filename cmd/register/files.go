package register

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	rconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/register"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flytestdlib/logger"
)

const (
	registerFilesShort = "Registers file resources"
	registerFilesLong  = `
This registers all the serialized protobuf files including tasks, workflows and launchplans with default v1 version.
If the registered entities with v1 version are already present, the command will fail on first such encounter.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks
	
There is no difference between registration and fast registration. In fast registration, the input provided by the user is fast serialized proto which is generated by pyflyte. 
If Flytectl finds any source code in the user's input, such registration is considered as fast registration. Flytectl finds input file by searching an archive file whose name begins with fast and has .tar.gz extension. When the user runs pyflyte with --fast flag, pyflyte creates serialize proto. It archives the create source code archive file in the same directory. 
SourceUploadPath is an optional flag. By default, flytectl will create SourceUploadPath from your storage config. In case of s3, flytectl will upload code base in s3://{{DEFINE_BUCKET_IN_STORAGE_CONFIG}}/fast/{{VERSION}}-fast{{MD5_CREATED_BY_PYFLYTE}.tar.gz}. 
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks  --version v2 
	
In case of fast registration, if the SourceUploadPath flag is defined, flytectl won't use the default directory to upload the source code. It will override the destination path on the registration. 
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks  --version v2 --SourceUploadPath="s3://dummy/fast" 
	
The archive file command format is shown below. Currently it supports .tgz and .tar extension files and can be local or remote file served through http/https.
The --archive flag should be used.

::

 bin/flytectl register files  http://localhost:8080/_pb_output.tar -d development  -p flytesnacks --archive

Usage of local tgz file with archive flag:

::

 bin/flytectl register files  _pb_output.tgz -d development  -p flytesnacks --archive

If you wish to continue executing registration on other files by ignoring the errors with inclusion of version conflicts, then the continueOnError flag should be passed.
It has been shown below:
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError

The short format of continueOnError flag usage has been shown below:
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError

Overriding the default version v1 using version string is shown below:
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --version v2

Changing the output format has no effect on the registration. Currently, the output is available only in tabular format.
It has been shown below:
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError -o yaml

The IamRole can be overridden during registration. This has been shown below:

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError --version v2 -i "arn:aws:iam::123456789:role/dummy"

The Kubernetes service account can be overridden during registration. This has been shown below:

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError --version v2 -k "kubernetes-service-account"

The Output location prefix can be overridden during registration. This has been shown below:

::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks --continueOnError --version v2 -l "s3://dummy/prefix"
	
Usage
`
	sourceCodeExtension = ".tar.gz"
)

func registerFromFilesFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	return Register(ctx, args, cmdCtx)
}

func Register(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	var regErr error
	var dataRefs []string

	// Deprecated checks for --k8Service
	deprecatedCheck(ctx, &rconfig.DefaultFilesConfig.K8sServiceAccount, rconfig.DefaultFilesConfig.K8ServiceAccount)

	// getSerializeOutputFiles will return you all proto and  source code compress file in sorted order
	dataRefs, tmpDir, err := getSerializeOutputFiles(ctx, args, rconfig.DefaultFilesConfig.Archive)
	if err != nil {
		logger.Errorf(ctx, "error while un-archiving files in tmp dir due to %v", err)
		return err
	}
	logger.Infof(ctx, "Parsing file... Total(%v)", len(dataRefs))

	// It will segregate serialize output files in valid proto,Invalid files if have any and source code(In case of fast serialize input files)
	sourceCode, validProto, InvalidFiles := segregateSourceAndProtos(dataRefs)

	// If any invalid files provide in input then through an error
	if len(InvalidFiles) > 0 {
		return fmt.Errorf("input package have some invalid files. try to run pyflyte package again %v", InvalidFiles)
	}

	// In case of fast serialize input, upload source code to destination bucket
	var sourceCodeName string
	if len(sourceCode) > 0 {
		logger.Infof(ctx, "Fast Registration detected")
		_, sourceCodeName = filepath.Split(sourceCode)
		if err = uploadFastRegisterArtifact(ctx, sourceCode, sourceCodeName, rconfig.DefaultFilesConfig.Version, &rconfig.DefaultFilesConfig.SourceUploadPath); err != nil {
			return fmt.Errorf("please check your Storage Config. It failed while uploading the source code. %v", err)
		}
		logger.Infof(ctx, "Source code successfully uploaded %v/%v ", rconfig.DefaultFilesConfig.SourceUploadPath, sourceCodeName)
	}

	var registerResults []Result
	fastFail := rconfig.DefaultFilesConfig.ContinueOnError
	for i := 0; i < len(validProto) && !(fastFail && regErr != nil); i++ {
		registerResults, regErr = registerFile(ctx, validProto[i], sourceCodeName, registerResults, cmdCtx, *rconfig.DefaultFilesConfig)
	}

	payload, _ := json.Marshal(registerResults)
	registerPrinter := printer.Printer{}
	_ = registerPrinter.JSONToTable(payload, projectColumns)
	if tmpDir != "" {
		if _err := os.RemoveAll(tmpDir); _err != nil {
			logger.Errorf(ctx, "unable to delete temp dir %v due to %v", tmpDir, _err)
			return _err
		}
	}
	return regErr
}
