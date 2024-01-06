package cmd

type rootCommandOptions struct {
	Host     string
	Password string
}

type poeModeCmdOptions struct {
	Token string
	Ports []string
	Mode  string
}
