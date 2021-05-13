package delete

import (
	cmdcore "github.com/flyteorg/flytectl/cmd/core"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config/subcommand/clusterresourceattribute"
	u "github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SetupFunc func()

func Test1ClusterResourceattribute(t *testing.T) {
	clusterresourceattribute.DefaultDelConfig.AttrFile = "testdata/valid_project_domain_cluster_attribute.yaml"

	defaultConfigData := clusterresourceattribute.DefaultDelConfig
	MatchableAttributeTest1(t, deleteClusterResourceAttributeSetup, deleteClusterResourceAttributes,
		&(defaultConfigData.AttrFile))
}

func MatchableAttributeTest1(t *testing.T, setupFunc SetupFunc,
	commandFunc cmdcore.CommandFunc,
	attrFileRef *string) {
	t.Run("successful project domain attribute deletion file", func(t *testing.T) {
		setup()
		setupFunc()
		// Empty attribute file
//		clusterresourceattribute.DefaultDelConfig.AttrFile = "testdata/valid_project_domain_cluster_attribute.yaml"
		*attrFileRef = "testdata/valid_project_domain_cluster_attribute.yaml"

		// No args implying project domain attribute deletion
		u.DeleterExt.OnDeleteProjectDomainAttributesMatch(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(nil)
		err = commandFunc(ctx, args, cmdCtx)
		assert.Nil(t, err)
		u.DeleterExt.AssertCalled(t, "DeleteProjectDomainAttributes",
			ctx, "flytectldemo", "development", admin.MatchableResource_CLUSTER_RESOURCE)
	})
}
