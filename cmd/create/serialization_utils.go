package create

import (
	"github.com/flyteorg/flyteidl/clients/go/coreutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

// TODO: Move all functions to flyteidl
func MakeLiteralForVariables(serialize map[string]interface{}, variables map[string]*core.Variable) (map[string]*core.Literal, error) {
	result := make(map[string]*core.Literal)
	var err error
	for k, v := range variables {
		// Only serialize inputs that are provided
		if input, provided := serialize[k]; provided {
			if result[k], err = coreutils.MakeLiteralForType(v.Type, input); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func MakeLiteralForParams(serialize map[string]interface{}, parameters map[string]*core.Parameter) (map[string]*core.Literal, error) {
	result := make(map[string]*core.Literal)
	var err error
	for k, v := range parameters {
		// Only serialize inputs that are provided
		if input, provided := serialize[k]; provided {
			if result[k], err = coreutils.MakeLiteralForType(v.GetVar().Type, input); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
