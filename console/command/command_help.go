package command

type commandHelp struct{}

func init() {
	cmd := new(commandHelp)
	Register(cmd)
}

func (*commandHelp) Name() string {
	return "help"
}

func (*commandHelp) Help() string {
	return "help text"
}

func (cmd *commandHelp) Run([]string) string {
	str := "Commands:\r\n"
	for _, c := range commands {
		str += c.Name() + "\t\t\t" + c.Help() + "\r\n"
	}
	str += "quit - exit console"

	return str
}
