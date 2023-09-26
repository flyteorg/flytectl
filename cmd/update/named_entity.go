package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/flyteorg/flytectl/clierrors"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flytestdlib/logger"
)

//go:generate pflags NamedEntityConfig --default-var namedEntityConfig --bind-default-var

var (
	namedEntityConfig = &NamedEntityConfig{}
)

type NamedEntityConfig struct {
	Archive     bool   `json:"archive" pflag:",archive named entity."`
	Activate    bool   `json:"activate" pflag:",activate the named entity."`
	Description string `json:"description" pflag:",description of the named entity."`
	DryRun      bool   `json:"dryRun" pflag:",execute command without making any modifications."`
}

func (cfg NamedEntityConfig) UpdateNamedEntity(ctx context.Context, name string, project string, domain string, rsType core.ResourceType, cmdCtx cmdCore.CommandContext) error {
	archive := cfg.Archive
	activate := cfg.Activate
	if activate == archive && activate {
		return fmt.Errorf(clierrors.ErrInvalidStateUpdate)
	}

	var state admin.NamedEntityState
	if activate {
		state = admin.NamedEntityState_NAMED_ENTITY_ACTIVE
	} else if archive {
		state = admin.NamedEntityState_NAMED_ENTITY_ARCHIVED
	}

	namedEntity, err := cmdCtx.AdminClient().GetNamedEntity(ctx, &admin.NamedEntityGetRequest{
		ResourceType: rsType,
		Id: &admin.NamedEntityIdentifier{
			Project: project,
			Domain:  domain,
			Name:    name,
		},
	})
	if err != nil {
		return err
	}

	v, _ := json.MarshalIndent(namedEntity, "", "    ")
	fmt.Println(string(v))

	// TODO: kamal - ack/force

	if cfg.DryRun {
		logger.Infof(ctx, "skipping UpdateNamedEntity request (dryRun)")
		return nil
	}

	_, err = cmdCtx.AdminClient().UpdateNamedEntity(ctx, &admin.NamedEntityUpdateRequest{
		ResourceType: rsType,
		Id: &admin.NamedEntityIdentifier{
			Project: project,
			Domain:  domain,
			Name:    name,
		},
		Metadata: &admin.NamedEntityMetadata{
			Description: cfg.Description,
			State:       state,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
