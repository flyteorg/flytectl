package get

var executionStructure = map[string]string{
	"Version": "$.id.version",
	"Name":    "$.id.name",
	"Type":    "$.closure.compiledTask.template.type",
}

var executionSingleStructure = map[string]string{
	"Version":    "$.id.version",
	"Name":       "$.id.name",
	"LaunchPlan": "$.spec.launchplan.name",
	"Phase":      "$.spec.phase",
	"Duration":   "$.closure.duration",
	"StartedAt":  "$.closure.started_at",
	"Workflow":   "$.closure.workflow_id.name",
	"Metadata":   "$.spec.metadata",
}

var launchPlanStructure = map[string]string{
	"Version": "$.id.version",
	"Name":    "$.id.name",
	"Type":    "$.closure.compiledTask.template.type",
}

var tableStructure = map[string]string{
	"ID":          "$.id",
	"Name":        "$.name",
	"Description": "$.description",
}

var taskStructure = map[string]string{
	"Version":          "$.id.version",
	"Name":             "$.id.name",
	"Type":             "$.closure.compiledTask.template.type",
	"Discoverable":     "$.closure.compiledTask.template.metadata.discoverable",
	"DiscoveryVersion": "$.closure.compiledTask.template.metadata.discovery_version",
}

var workflowStructure = map[string]string{
	"Version": "$.id.version",
	"Name":    "$.id.name",
}
