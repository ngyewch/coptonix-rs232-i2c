package periph

import (
	"time"

	coptonixrs232i2c "github.com/ngyewch/coptonix-rs232-i2c"
	"go.bug.st/serial"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

var (
	busMap = make(map[string]i2c.BusCloser, 0)
)

func NewOpener(name string, portName string, mode *serial.Mode, readTimeout time.Duration) i2creg.Opener {
	return func() (i2c.BusCloser, error) {
		bus, ok := busMap[name]
		if !ok {
			serialPort, err := coptonixrs232i2c.OpenSerialPort(portName, mode, readTimeout)
			if err != nil {
				return nil, err
			}
			dev := coptonixrs232i2c.New(serialPort)
			bus = NewBus(name, dev)
			busMap[name] = bus
		}
		return bus, nil
	}
}
