package periph

import (
	"fmt"
	"time"

	"github.com/ngyewch/coptonix-rs232-i2c"
	"go.bug.st/serial"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func NewOpener(name string, portName string, mode *serial.Mode, readTimeout time.Duration) i2creg.Opener {
	fmt.Println("NewOpener")
	return func() (i2c.BusCloser, error) {
		fmt.Printf("coptonix-rs232-i2c opening %s %s\n", name, portName)
		serialPort, err := coptonixrs232i2c.OpenSerialPort(portName, mode, readTimeout)
		if err != nil {
			return nil, err
		}
		dev := coptonixrs232i2c.New(serialPort)
		bus := NewBus(name, dev)
		return bus, nil
	}
}
