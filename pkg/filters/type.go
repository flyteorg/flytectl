package filters

import "github.com/flyteorg/flytectl/cmd/config"

var DefaultFilter = Filters{
	Limit: config.DefaultLimit,
	Asc:   false,
}

type Filters struct {
	FieldSelector string `json:"field-selector" pflag:",Specifies the Field selector"`
	SortBy        string `json:"sort-by" pflag:",Specifies which field to sort results "`
	// TODO: Support paginated queries
	Limit int32 `json:"limit" pflag:",Specifies the limit"`
	Asc   bool  `json:"asc"  pflag:",Specifies the sorting order. By default flytectl sort result in descending order"`
}
