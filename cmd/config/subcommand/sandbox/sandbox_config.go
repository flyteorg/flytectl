package sandbox

//go:generate pflags Config --default-var DefaultConfig --bind-default-var
var (
	DefaultConfig = &Config{
		Local: false,
	}
)

//Config
type Config struct {
	Source string `json:"source" pflag:",Path of your source code"`

	// Flytectl sandbox only supports flyte version available in Github release https://github.com/flyteorg/flyte/tags
	// Flytectl sandbox will only work for v0.10.0+
	// Default value dind represents the latest release
	Version string `json:"version" pflag:",Version of flyte. Only supports flyte releases greater than v0.10.0"`

	// Optionally it is possible to specify a specific fqn for the docker image with the tag. This should be
	// Flyte compliant sandbox image. Usually useful, if you want to push the image to your own registry and relaunch
	// from there.
	Image string `json:"image" pflag:",Optional. Provide a fully qualified path to a Flyte compliant docker image."`

	// Optionally it is possible to use local sandbox image
	// If local flag pass then flytectl will not pull image from registry. Usually useful, if you want to test your local images without pushing them to a registry
	Local bool `json:"local" pflag:",Optional. Enable if you want to use local available image."`
}
