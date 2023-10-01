package update

import (
	"context"
	"fmt"
	"os"

	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/launchplan"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

const (
	updateLPShort = "Updates launch plan status"
	updateLPLong  = `
Activates a ` + "`launch plan <https://docs.flyte.org/projects/cookbook/en/latest/auto/core/scheduled_workflows/lp_schedules.html#activating-a-schedule>`__" + ` which activates the scheduled job associated with it:
::

 flytectl update launchplan -p flytesnacks -d development core.control_flow.merge_sort.merge_sort --version v1 --activate

Archives ` + "`(deactivates) <https://docs.flyte.org/projects/cookbook/en/latest/auto/core/scheduled_workflows/lp_schedules.html#deactivating-a-schedule>`__" + ` a launch plan which deschedules any scheduled job associated with it:
::

 flytectl update launchplan -p flytesnacks -d development core.control_flow.merge_sort.merge_sort --version v1 --archive

Usage
`
)

func updateLPFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	project := config.GetConfig().Project
	domain := config.GetConfig().Domain
	if len(args) != 1 {
		return fmt.Errorf(clierrors.ErrLPNotPassed)
	}
	name := args[0]
	version := launchplan.UConfig.Version
	if len(version) == 0 {
		return fmt.Errorf(clierrors.ErrLPVersionNotPassed)
	}

	activate := launchplan.UConfig.Activate
	archive := launchplan.UConfig.Archive
	if activate == archive && archive {
		return fmt.Errorf(clierrors.ErrInvalidStateUpdate)
	}

	var newState admin.LaunchPlanState
	if activate {
		newState = admin.LaunchPlanState_ACTIVE
	} else if archive {
		newState = admin.LaunchPlanState_INACTIVE
	}

	id := &core.Identifier{
		Project: project,
		Domain:  domain,
		Name:    name,
		Version: version,
	}

	launchPlan, err := cmdCtx.AdminClient().GetLaunchPlan(ctx, &admin.ObjectGetRequest{Id: id})
	if err != nil {
		return fmt.Errorf("update launch plan %s: could not fetch launch plan: %w", name, err)
	}
	oldState := launchPlan.Closure.GetState()

	type LaunchPlan struct {
		State admin.LaunchPlanState `json:"state"`
	}
	patch, err := diffAsYaml(LaunchPlan{oldState}, LaunchPlan{newState})
	if err != nil {
		panic(err)
	}

	if patch == "" {
		fmt.Printf("No changes detected. Skipping the update.\n")
		return nil
	}

	fmt.Printf("The following changes are to be applied.\n%s\n", patch)

	if launchplan.UConfig.DryRun {
		fmt.Printf("skipping LaunchPlanUpdate request (DryRun)")
		return nil
	}

	if !launchplan.UConfig.Force && !cmdUtil.AskForConfirmation("Continue?", os.Stdin) {
		return fmt.Errorf("update aborted by user")
	}

	_, err = cmdCtx.AdminClient().UpdateLaunchPlan(ctx, &admin.LaunchPlanUpdateRequest{
		Id:    id,
		State: newState,
	})
	if err != nil {
		return fmt.Errorf(clierrors.ErrFailedLPUpdate, name, err)
	}

	fmt.Printf("updated launch plan successfully on %s", name)

	return nil
}
