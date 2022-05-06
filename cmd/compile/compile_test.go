package compile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompilePackage(t *testing.T) {
	// valid package contains two workflows
	// with three tasks
	err := compileFromPackage("testdata/valid-package.tgz")
	assert.Nil(t, err, "unable to compile a valid package")

	// invalid gzip header
	err = compileFromPackage("testdata/invalid.tgz")
	assert.NotNil(t, err, "compiling an invalid package returns no error")

	// invalid workflow, types do not match
	err = compileFromPackage("testdata/bad-workflow-package.tgz")
	assert.NotNil(t, err, "compilin an invalid workflow returns no error")
}
