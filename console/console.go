// console.go
package console

import (
	"net"
	"strings"

	"github.com/gbember/gt/console/command"
	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/module"
	"github.com/gbember/gt/network"
	"github.com/gbember/gt/network/msg"
)

type console struct {
	server     *network.TCPServer
	addr       string
	maxConnNum int
	maxDataLen int
}

type agent struct {
	prompt    []byte
	conn      net.Conn
	msgParser msg.MsgParser
}

func RegisterModule(addr string, maxConnNum int, maxDataLen int) {
	c := new(console)
	c.addr = addr
	c.maxConnNum = maxConnNum
	c.maxDataLen = maxDataLen
	module.Register(c)
}

func (c *console) OnInit() {
	msgParser, err := msg.NewMsgParserLine(c.maxDataLen)
	if err != nil {
		panic(err)
	}
	server, err := network.StartTCPServer(c.addr, c.maxConnNum, msgParser, NewAgant)
	if err != nil {
		panic(err)
	}
	c.server = server
	logger.Info("console start...")
}

func (c *console) OnDestroy() {
	if c.server != nil {
		c.server.Close()
	}
}

func (c *console) Run(closeSign chan bool) {
}

func NewAgant(conn net.Conn, msgParser msg.MsgParser) network.TCPAgent {
	agent := new(agent)
	agent.conn = conn
	agent.msgParser = msgParser
	agent.prompt = []byte(">")
	return agent
}

func (a *agent) Run() {
	for {
		if a.prompt != nil {
			a.msgParser.Write(a.conn, a.prompt)
		}
		bs, err := a.msgParser.Read(a.conn)
		if err != nil {
			return
		}
		line := strings.TrimSpace(string(bs))
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		if args[0] == "quit" {
			break
		}
		str := command.Run(args[0], args[1:])
		a.msgParser.Write(a.conn, []byte(str))
	}
}
func (a *agent) Close(int8) {}
