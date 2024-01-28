package cmd

type rootCommandOptions struct {
	Host     string
	Password string

	Verbosity int
}

type poeModeCmdOptions struct {
	Token string
	Ports []string
	Mode  string
}
