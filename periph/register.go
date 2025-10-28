package periph

import (
	"time"

	coptonixrs232i2c "github.com/ngyewch/coptonix-rs232-i2c"
	"go.bug.st/serial"
	"periph.io/x/conn/v3/i2c/i2creg"
)

var (
	DefaultOptions = Options{
		ReadTimeout: coptonixrs232i2c.DefaultReadTimeout,
		AddrMapper:  nil,
	}
)

type Options struct {
	ReadTimeout time.Duration
	AddrMapper  func(addr uint16) (uint16, error)
}

func Register(name string, aliases []string, portName string, mode *serial.Mode, options *Options) error {
	return i2creg.Register(name, aliases, -1, NewOpener(name, portName, mode, options))
}
