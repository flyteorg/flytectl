package register

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	rconfig "github.com/flyteorg/flytectl/cmd/config/subcommand/register"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flytestdlib/contextutils"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/flyteorg/flytestdlib/promutils/labeled"
	"github.com/flyteorg/flytestdlib/storage"
)

const (
	registerFilesShort = "Registers file resources"
	registerFilesLong  = `
Registers all the serialized protobuf files including tasks, workflows and launchplans with default v1 version.
If there are already registered entities with v1 version then the command will fail immediately on the first such encounter.
::

 bin/flytectl register file  _pb_output/* -d development  -p flytesnacks

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
	dataRefs, tmpDir, _err := getSortedFileList(ctx, args)
	if _err != nil {
		logger.Errorf(ctx, "error while un-archiving files in tmp dir due to %v", _err)
		return _err
	}
	logger.Infof(ctx, "Parsing files... Total(%v)", len(dataRefs))
	fastFail := !rconfig.DefaultFilesConfig.ContinueOnError
	var registerResults []Result
	var s *storage.DataStore
	var dataRefReaderCloser io.ReadCloser
	pbFiles := []string{}
	testScope := promutils.NewTestScope()
	// Set Keys
	labeled.SetMetricKeys(contextutils.AppNameKey, contextutils.ProjectKey, contextutils.DomainKey)
	if rconfig.DefaultFilesConfig.FastRegister {

		s, _err = storage.NewDataStore(storage.GetConfig(), testScope.NewSubScope("flytectl"))
		if _err != nil {
			logger.Errorf(ctx, "error while creating storage client %v", _err)
			return _err
		}

		fastRegisterCheck := false
		fullRemotePath := getAdditionalDistributionLoc(rconfig.DefaultFilesConfig.AdditionalDistributionDir, rconfig.DefaultFilesConfig.Version)
		for i := 0; i < len(dataRefs) && !(fastFail && _err != nil); i++ {
			if strings.Contains(dataRefs[i], ".tar.gz") {
				raw, err := json.Marshal(dataRefs[i])
				if err != nil {
					return err
				}
				dataRefReaderCloser, err = os.Open(dataRefs[i])
				if err != nil {
					return err
				}
				dataRefReaderCloser, err = gzip.NewReader(dataRefReaderCloser)
				if err != nil {
					return err
				}
				fastRegisterCheck = true
				if err := s.WriteRaw(ctx, fullRemotePath, int64(len(raw)), storage.Options{}, dataRefReaderCloser); err != nil {
					return err
				}
				continue
			}
			pbFiles = append(pbFiles,dataRefs[1])
		}
		if !fastRegisterCheck {
			return fmt.Errorf("could not discover compressed source, did you remember to run 'pyflyte serialize fast'")
		}
	}
	for i := 0; i < len(pbFiles) && !(fastFail && _err != nil); i++ {
			registerResults, _err = registerFile(ctx, pbFiles[i], registerResults, cmdCtx)
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
