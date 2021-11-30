package delayed

import (
	"io"
	"sync"
)

type writer struct {
	pool    *sync.Pool
	buffer  chan Message
	target  io.WriteCloser
	monitor Monitor
}

func newWriter(config configuration) Writer {
	pool := &sync.Pool{New: func() interface{} { return config.Source() }}
	buffer := make(chan Message, config.ChannelSize)

	for i := 0; i < config.PoolSize; i++ {
		pool.Put(config.Source())
	}

	return &writer{pool: pool, buffer: buffer, target: config.Target, monitor: config.Monitor}
}

func (this *writer) Write(callback func(message Message)) {
	message := this.pool.Get().(Message)
	callback(message)

	select {
	case this.buffer <- message:
		this.monitor.Buffered()
	default:
		this.monitor.Discarded(message)
		this.pool.Put(message)
	}
}

func (this *writer) Listen() {
	defer func() { _ = this.target.Close() }()
	for message := range this.buffer {
		this.writeMessage(message)
		this.pool.Put(message)
	}
}
func (this *writer) writeMessage(message Message) {
	if _, err := message.WriteTo(this.target); err != nil {
		this.monitor.WriteFailed(message, err)
	} else {
		this.monitor.Written()
	}
}

func (this *writer) Close() error {
	close(this.buffer)
	return nil
}
