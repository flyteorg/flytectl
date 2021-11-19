// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package sandbox

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (Config) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func (Config) mustJsonMarshal(v interface{}) string {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(raw)
}

func (Config) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in Config and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg Config) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("Config", pflag.ExitOnError)
	cmdFlags.StringVar(&DefaultConfig.Source, fmt.Sprintf("%v%v", prefix, "source"), DefaultConfig.Source, "Path of your source code")
	cmdFlags.StringVar(&DefaultConfig.Version, fmt.Sprintf("%v%v", prefix, "version"), DefaultConfig.Version, "Version of flyte. Only supports flyte releases greater than v0.10.0")
	cmdFlags.StringVar(&DefaultConfig.Image, fmt.Sprintf("%v%v", prefix, "image"), DefaultConfig.Image, "Optional. Provide a fully qualified path to a Flyte compliant docker image.")
	cmdFlags.IntVar((*int)(&DefaultConfig.ImagePullPolicy), fmt.Sprintf("%v%v", prefix, "imagePullPolicy"), int(DefaultConfig.ImagePullPolicy), "Optional. Defines the image pull behavior. (0/1/2) Always/IfNotPresent/Never")
	return cmdFlags
}
