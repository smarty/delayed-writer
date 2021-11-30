package delayed

import "io"

type Message interface {
	Sequence(uint64)
	io.WriterTo
}

type Writer interface {
	Write(callback func(Message))
	ListenCloser
}

type ListenCloser interface {
	Listen()
	io.Closer
}

type Monitor interface {
	Buffered()
	Discarded(Message)
	Written()
	WriteFailed(Message, error)
}
