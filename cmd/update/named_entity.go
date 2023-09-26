package update

import (
	"context"
	"fmt"
	"os"

	"github.com/flyteorg/flytectl/clierrors"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
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
	Force       bool   `json:"force" pflag:",do not ask for an acknowledgement during updates."`
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

	id := &admin.NamedEntityIdentifier{
		Project: project,
		Domain:  domain,
		Name:    name,
	}

	namedEntity, err := cmdCtx.AdminClient().GetNamedEntity(ctx, &admin.NamedEntityGetRequest{
		ResourceType: rsType,
		Id:           id,
	})
	if err != nil {
		return fmt.Errorf("update metadata for %s: could not fetch metadata: %w", name, err)
	}

	oldMetadata := namedEntity.Metadata
	newMetadata := &admin.NamedEntityMetadata{
		Description: cfg.Description,
		State:       state,
	}

	patch, err := diffAsYaml(oldMetadata, newMetadata)
	if err != nil {
		panic(err)
	}

	if patch == "" {
		fmt.Printf("No changes detected. Skipping the update.\n")
		return nil
	}

	fmt.Printf("The following changes are to be applied.\n%s\n", patch)

	if cfg.DryRun {
		fmt.Printf("skipping UpdateNamedEntity request (dryRun)\n")
		return nil
	}

	if !cfg.Force && !cmdUtil.AskForConfirmation("Continue?", os.Stdin) {
		return fmt.Errorf("update aborted by user")
	}

	_, err = cmdCtx.AdminClient().UpdateNamedEntity(ctx, &admin.NamedEntityUpdateRequest{
		ResourceType: rsType,
		Id:           id,
		Metadata:     newMetadata,
	})
	if err != nil {
		return fmt.Errorf("update metadata for %s: update failed: %w", name, err)
	}

	return nil
}
