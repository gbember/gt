// msg_parser.go
package msg

import (
	"io"
)

type MsgParser interface {
	Read(io.Reader) ([]byte, error)
	Write(io.Writer, []byte) error
	Clone() MsgParser
}
