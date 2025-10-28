package periph

import (
	"fmt"
	"time"

	"go.bug.st/serial"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func Register(name string, aliases []string, portName string, mode *serial.Mode, readTimeout time.Duration) error {
	fmt.Println("Register")
	return i2creg.Register(name, aliases, -1, NewOpener(name, portName, mode, readTimeout))
}
