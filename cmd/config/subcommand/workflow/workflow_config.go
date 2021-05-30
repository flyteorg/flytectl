package workflow

//go:generate pflags Config --default-var DefaultConfig --bind-default-var

var (
	DefaultConfig = &Config{
		OutputFileName: "/tmp/workflow.svg",
	}
)

// Config commandline configuration
type Config struct {
	Version        string `json:"version" pflag:",version of the workflow to be fetched."`
	Latest         bool   `json:"latest" pflag:", flag to indicate to fetch the latest version, version flag will be ignored in this case"`
	OutputFileName string `json:"opFileName" pflag:",output file name to be used for the fetched workflow.Currently support dumping the SVG image.Use .svg extension for the filename. Defaults to /tmp/workflow.svg"`
}
