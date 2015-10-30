// command.go
package command

import "fmt"

type Command interface {
	//命令名字
	Name() string
	//命令帮助
	Help() string

	Run(args []string) string
}

var commands map[string]Command = make(map[string]Command)

//注册命令
func Register(cmd Command) {
	name := cmd.Name()
	if _, ok := commands[name]; ok {
		panic(fmt.Sprintf("command %v is already registered", name))
	}
	commands[name] = cmd
}

//运行命令
func Run(name string, args []string) string {
	cmd, ok := commands[name]
	if ok {
		return cmd.Run(args) + "\r\n"
	}
	return "command not found, try `help` for help\r\n"
}
