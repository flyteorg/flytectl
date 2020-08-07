package printer

type PrintableTask struct {
	Version          string `header:"Version"`
	Name             string `header:"Name"`
	Type             string `header:"Type"`
	Discoverable     bool   `header:"Discoverable"`
	DiscoveryVersion string `header:"DiscoveryVersion"`
}

type PrintableNamedEntityIdentifier struct {
	Name    string `header:"Name"`
	Project string `header:"Project"`
	Domain  string `header:"Domain"`
}

type PrintableProject struct {
	Id          string `header:"Id"`
	Name        string `header:"Name"`
	Description string `header:"Description"`
}

type PrintableDomain struct {
	Id   string `header:"Id"`
	Name string `header:"Name"`
}

type PrintableWorkflow struct {
	Name    string `header:"Name"`
	Version string `header:"Version"`
}
