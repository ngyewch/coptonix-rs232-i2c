package coptonixrs232i2c

import (
	"time"

	"go.bug.st/serial"
)

func OpenSerialPort(portName string, mode *serial.Mode, readTimeout time.Duration) (serial.Port, error) {
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}

	err = port.SetReadTimeout(readTimeout)
	if err != nil {
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	return port, nil
}
