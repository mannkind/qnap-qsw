package cmd

type loginCmdOptions struct {
	Host     string
	Password string
}

type poeModeCmdOptions struct {
	Host     string
	Password string
	Token    string
	Ports    []string
	Mode     string
}
