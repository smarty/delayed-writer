package delayed

import (
	"fmt"
	"io"
	"sync"
)

type writer struct {
	pool    *sync.Pool
	buffer  chan fmt.Stringer
	target  io.WriteCloser
	monitor Monitor
}

func newWriter(config configuration) Writer {
	pool := &sync.Pool{New: func() interface{} { return config.Source() }}
	buffer := make(chan fmt.Stringer, config.ChannelSize)

	for i := 0; i < config.PoolSize; i++ {
		pool.Put(config.Source())
	}

	return &writer{pool: pool, buffer: buffer, target: config.Target, monitor: config.Monitor}
}

func (this *writer) Write(callback func(message fmt.Stringer)) {
	item := this.pool.Get().(fmt.Stringer)
	callback(item)

	select {
	case this.buffer <- item:
		this.monitor.Buffered()
	default:
		this.monitor.Discarded(item)
		this.pool.Put(item)
	}
}

func (this *writer) Listen() {
	defer func() { _ = this.target.Close() }()
	for item := range this.buffer {
		this.writeMessage(item)
		this.pool.Put(item)
	}
}
func (this *writer) writeMessage(item fmt.Stringer) {
	message := item.String()
	if _, err := io.WriteString(this.target, message); err != nil {
		this.monitor.WriteFailed(item, err)
	} else {
		this.monitor.Written()
	}
}

func (this *writer) Close() error {
	close(this.buffer)
	return nil
}
