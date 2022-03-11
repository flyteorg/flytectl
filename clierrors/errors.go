package clierrors

var (
	ErrInvalidStateUpdate = "Invalid state passed. Specify either activate or archive\n"

	ErrProjectNotPassed     = "Project id wasn't passed\n" // #nosec
	ErrProjectNameNotPassed = "Project name is a required flag"
	ErrFailedProjectUpdate  = "Project %v failed to update due to %v\n"

	ErrLPNotPassed        = "Launch plan name wasn't passed\n"
	ErrLPVersionNotPassed = "Launch plan version wasn't passed\n" //nolint
	ErrFailedLPUpdate     = "Launch plan %v failed to update due to %v\n"

	ErrExecutionNotPassed    = "Execution name wasn't passed\n"
	ErrFailedExecutionUpdate = "Execution %v failed to update due to %v\n"

	ErrWorkflowNotPassed    = "Workflow name wasn't passed\n"
	ErrFailedWorkflowUpdate = "Workflow %v failed to update to due to %v\n"

	ErrTaskNotPassed    = "Task name wasn't passed\n" // #nosec
	ErrFailedTaskUpdate = "Task %v failed to update to due to %v\n"

	ErrSandboxExists = "Sandbox already exists!\n"
)
