package network

import (
	"net"
	"sync"
	"time"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/network/msg"
)

type TCPServer struct {
	//地址
	addr string
	//最大连接数
	maxConnNum int
	ln         net.Listener
	wgLn       sync.WaitGroup
	conns      map[net.Conn]struct{}
	mutexConns sync.Mutex
	wgConns    sync.WaitGroup
	msgParser  msg.MsgParser
	newAgent   func(net.Conn, msg.MsgParser) TCPAgent
}

func StartTCPServer(addr string, maxConnNum int, msgParser msg.MsgParser, newAgent func(net.Conn, msg.MsgParser) TCPAgent) (*TCPServer, error) {
	server := new(TCPServer)
	server.addr = addr
	server.maxConnNum = maxConnNum
	server.msgParser = msgParser
	server.newAgent = newAgent
	server.conns = make(map[net.Conn]struct{})
	err := server.init()
	if err != nil {
		return nil, err
	}
	go server.run()
	return server, nil
}

func (server *TCPServer) init() error {
	ln, err := net.Listen("tcp", server.addr)
	if err != nil {
		return err
	}
	server.ln = ln
	return nil
}

func (server *TCPServer) run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()

	var tempDelay time.Duration
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				logger.Error("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		//判断是否超最高在线个数
		server.mutexConns.Lock()
		if len(server.conns) >= server.maxConnNum {
			server.mutexConns.Unlock()
			conn.Close()
			logger.Info("too many connections")
			continue
		}
		server.conns[conn] = struct{}{}
		server.mutexConns.Unlock()

		agent := server.newAgent(conn, server.msgParser)
		server.wgConns.Add(1)
		go func() {
			agent.Run()
			conn.Close()
			server.wgConns.Done()
		}()
	}
}

func (server *TCPServer) Close() {
	server.ln.Close()
	server.wgLn.Wait()

	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()
	server.wgConns.Wait()
}
