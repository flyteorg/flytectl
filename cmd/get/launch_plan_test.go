package get

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLaunchPlanFunc(t *testing.T) {
	setup()
	err = getLaunchPlanFunc(ctx, argsLp, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListLaunchPlans", ctx, resourceListRequest)
	teardownAndVerify(t, `[
	{
		"id": {
			"name": "launchplan1",
			"version": "v2"
		},
		"closure": {
			"createdAt": "1970-01-01T00:00:01Z"
		}
	},
	{
		"id": {
			"name": "launchplan1",
			"version": "v1"
		},
		"closure": {
			"createdAt": "1970-01-01T00:00:00Z"
		}
	}
]`)
}

func TestGetLaunchPlanFuncLatest(t *testing.T) {
	setup()
	launchPlanConfig.Latest = true
	err = getLaunchPlanFunc(ctx, argsLp, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListLaunchPlans", ctx, resourceListRequest)
	teardownAndVerify(t, `{
	"id": {
		"name": "launchplan1",
		"version": "v2"
	},
	"closure": {
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}

func TestGetLaunchPlanWithVersion(t *testing.T) {
	setup()
	launchPlanConfig.Version = "v2"
	err = getLaunchPlanFunc(ctx, argsLp, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "GetLaunchPlan", ctx, objectGetRequest)
	teardownAndVerify(t, `{
	"id": {
		"name": "launchplan1",
		"version": "v2"
	},
	"closure": {
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}

func TestGetLaunchPlans(t *testing.T) {
	setup()
	argsLp = []string{}
	err = getLaunchPlanFunc(ctx, argsLp, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListLaunchPlanIds", ctx, namedIDRequest)
	teardownAndVerify(t, `[
	{
		"project": "dummyProject",
		"domain": "dummyDomain",
		"name": "launchplan1"
	},
	{
		"project": "dummyProject",
		"domain": "dummyDomain",
		"name": "launchplan2"
	}
]`)
}

func TestGetLaunchPlansWithExecFile(t *testing.T) {
	setup()
	launchPlanConfig.Version = "v2"
	launchPlanConfig.ExecFile = "exec_file"
	err = getLaunchPlanFunc(ctx, argsLp, cmdCtx)
	os.Remove(launchPlanConfig.ExecFile)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "GetLaunchPlan", ctx, objectGetRequest)
	teardownAndVerify(t, `{
	"id": {
		"name": "launchplan1",
		"version": "v2"
	},
	"spec": {
		"defaultInputs": {
			"parameters": {
				"numbers": {
					"var": {
						"type": {
							"collectionType": {
								"simple": "INTEGER"
							}
						}
					}
				},
				"numbers_count": {
					"var": {
						"type": {
							"simple": "INTEGER"
						}
					}
				},
				"run_local_at_count": {
					"var": {
						"type": {
							"simple": "INTEGER"
						}
					},
					"default": {
						"scalar": {
							"primitive": {
								"integer": "10"
							}
						}
					}
				}
			}
		}
	},
	"closure": {
		"expectedInputs": {
			"parameters": {
				"numbers": {
					"var": {
						"type": {
							"collectionType": {
								"simple": "INTEGER"
							}
						}
					}
				},
				"numbers_count": {
					"var": {
						"type": {
							"simple": "INTEGER"
						}
					}
				},
				"run_local_at_count": {
					"var": {
						"type": {
							"simple": "INTEGER"
						}
					},
					"default": {
						"scalar": {
							"primitive": {
								"integer": "10"
							}
						}
					}
				}
			}
		},
		"createdAt": "1970-01-01T00:00:01Z"
	}
}`)
}
