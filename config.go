package delayed

import "io"

func New(options ...option) Writer {
	var config configuration
	Options.apply(options...)(&config)

	if config.Target == nil {
		return &nopWriter{}
	} else if _, ok := config.Target.(*nop); ok {
		return &nopWriter{}
	}

	return newWriter(config)
}

func (singleton) Source(value func() Message) option {
	return func(this *configuration) { this.Source = value }
}
func (singleton) Target(value io.WriteCloser) option {
	return func(this *configuration) { this.Target = value }
}
func (singleton) PoolSize(value int) option {
	return func(this *configuration) { this.PoolSize = value }
}
func (singleton) ChannelSize(value int) option {
	return func(this *configuration) { this.ChannelSize = value }
}
func (singleton) Sequence(value uint64) option {
	return func(this *configuration) { this.Sequence = value }
}
func (singleton) Monitor(value Monitor) option {
	return func(this *configuration) { this.Monitor = value }
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, item := range Options.defaults(options...) {
			item(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	empty := &nop{}

	return append([]option{
		Options.Source(func() Message { return empty }),
		Options.Target(empty),
		Options.PoolSize(1024),
		Options.ChannelSize(128),
		Options.Sequence(0),
		Options.Monitor(empty),
	}, options...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type configuration struct {
	Source      func() Message
	Target      io.WriteCloser
	PoolSize    int
	ChannelSize int
	Sequence    uint64
	Monitor     Monitor
}
type option func(*configuration)
type singleton struct{}

var Options singleton

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type nop struct{}

func (*nop) Buffered()                      {}
func (*nop) Discarded(_ Message)            {}
func (*nop) Written()                       {}
func (*nop) WriteFailed(_ Message, _ error) {}

func (*nop) Sequence(uint64)                  {}
func (*nop) WriteTo(io.Writer) (int64, error) { return 0, nil }

func (*nop) Write([]byte) (int, error) { return 0, nil }
func (*nop) Close() error              { return nil }
