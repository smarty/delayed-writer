package delayed

import (
	"fmt"
	"io"
)

type Writer interface {
	Write(callback func(message fmt.Stringer))
	ListenCloser
}

type ListenCloser interface {
	Listen()
	io.Closer
}

type Monitor interface {
	Buffered()
	Discarded(fmt.Stringer)
	Written()
	WriteFailed(fmt.Stringer, error)
}
