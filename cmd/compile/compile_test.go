package compile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileCommand(t *testing.T) {
	compileCommand, err := CreateCompileCommand()
	assert.Nil(t, err)
	assert.Equal(t, compileCommand.Use, "compile")
	assert.Equal(t, compileCommand.Flags().Lookup("file").Shorthand, "f")
}

func TestCompilePackage(t *testing.T) {
	// valid package contains two workflows
	// with three tasks
	err := compileFromPackage("testdata/valid-package.tgz")
	assert.Nil(t, err, "unable to compile a valid package")

	// compiling via cobra command
	compileCommand, err := CreateCompileCommand()
	assert.Nil(t, err)
	err = compileCommand.Flags().Set("file", "testdata/valid-package.tgz")
	assert.Nil(t, err, "unable to set file flag")

	err = compileCommand.RunE(compileCommand, []string{})
	assert.Nil(t, err, "unable to compile a valid package")

	// invalid gzip header
	err = compileFromPackage("testdata/invalid.tgz")
	assert.NotNil(t, err, "compiling an invalid package returns no error")

	// invalid workflow, types do not match
	err = compileFromPackage("testdata/bad-workflow-package.tgz")
	assert.NotNil(t, err, "compilin an invalid workflow returns no error")

	// testing badly serialized task
	err = compileFromPackage("testdata/invalidtask.tgz")
	assert.NotNil(t, err, "unable to handle invalid task")

	// testing badly serialized launchplan
	err = compileFromPackage("testdata/invalidlaunchplan.tgz")
	assert.NotNil(t, err, "unable to handle invalid launchplan")

	// testing badly serialized workflow
	err = compileFromPackage("testdata/invalidworkflow.tgz")
	assert.NotNil(t, err, "unable to handle invalid workflow")
}
