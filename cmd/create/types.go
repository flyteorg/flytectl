package create

type CreateProject struct {
	Name string `json:"name" yaml:"name"`
	Id string `json:"id" yaml:"id"`
	Labels map[string]string `json:"labels" yaml:"labels"`
	Description string `json:"description" yaml:"description"`
}
