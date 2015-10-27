// command_cpuprof.go
package command

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

type commandCPUProf struct{}

func init() {
	cmd := new(commandCPUProf)
	Register(cmd)
}

func (c *commandCPUProf) Name() string {
	return "cpuprof"
}

func (c *commandCPUProf) Help() string {
	return "CPU profiling for the current process"
}

func (c *commandCPUProf) Run(args []string) string {
	if len(args) == 0 {
		return c.usage()
	}

	switch args[0] {
	case "start":
		fn := profileName() + ".cpuprof"
		f, err := os.Create(fn)
		if err != nil {
			return err.Error()
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			f.Close()
			return err.Error()
		}
		return fn
	case "stop":
		pprof.StopCPUProfile()
		return "success"
	default:
		return c.usage()
	}
}

func (*commandCPUProf) usage() string {
	return "cpuprof writes runtime profiling data in the format expected by \r\n" +
		"the pprof visualization tool\r\n\r\n" +
		"Usage: cpuprof start|stop\r\n" +
		"  start - enables CPU profiling\r\n" +
		"  stop  - stops the current CPU profile"
}

func profileName() string {
	now := time.Now()
	return fmt.Sprintf("./%d%02d%02d_%02d_%02d_%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
}
