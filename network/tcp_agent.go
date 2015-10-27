// tcp_agent.go
package network

type TCPAgent interface {
	Run()
	Close()
}
