// tcp_agent.go
package network

type TCPAgent interface {
	Run()
	//0:注销或正常关闭 1:服务器人数已满 2:关服关闭
	Close(closeReason int8)
}
