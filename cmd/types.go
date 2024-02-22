package cmd

type rootCommandOptions struct {
	Host     string
	Password string

	Verbosity int
}

type poeModeCmdOptions struct {
	Token            string
	DisablePorts     []string
	PoePorts         []string
	PoePlusPorts     []string
	PoePlusPlusPorts []string
	Mode             string
}
