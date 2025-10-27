package coptonixrs232i2c

import (
	"time"

	"go.bug.st/serial"
)

const (
	DefaultReadTimeout = 5 * time.Second
)

func DefaultMode() serial.Mode {
	return serial.Mode{
		BaudRate: 19200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
}
