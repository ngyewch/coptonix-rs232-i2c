package periph

import (
	"fmt"
	"time"

	coptonixrs232i2c "github.com/ngyewch/coptonix-rs232-i2c"
	"periph.io/x/conn/v3/physic"
)

// Bus interface for the Coptonix RS232-I2C device.
type Bus struct {
	name string
	dev  *coptonixrs232i2c.Dev
}

// NewBus constructs a new Bus instance.
func NewBus(name string, dev *coptonixrs232i2c.Dev) *Bus {
	return &Bus{
		name: name,
		dev:  dev,
	}
}

// Close closes the bus.
func (bus *Bus) Close() error {
	return bus.dev.Close()
}

// String returns the bus name.
func (bus *Bus) String() string {
	return bus.name
}

// Tx does a single transaction.
func (bus *Bus) Tx(addr uint16, w, r []byte) error {
	if addr >= 0x80 {
		return fmt.Errorf("address out of range")
	}
	if len(r) > 255 {
		return fmt.Errorf("read length out of range")
	}
	if w != nil {
		_, err := bus.dev.WriteI2C(uint8(addr), w)
		if err != nil {
			return err
		}
	}
	if (w != nil) && (r != nil) {
		time.Sleep(1 * time.Millisecond)
	}
	if r != nil {
		data, success, err := bus.dev.ReadI2C(uint8(addr), uint8(len(r)))
		if err != nil {
			return err
		}
		if !success {
			return fmt.Errorf("failed to read data")
		}
		for i, b := range data {
			if i >= len(r) {
				break
			}
			r[i] = b
		}
	}
	return nil
}

// SetSpeed changes the bus speed, if supported.
func (bus *Bus) SetSpeed(f physic.Frequency) error {
	speedInHz := int64(f) / 1_000_000
	if (speedInHz < 0) || (speedInHz > 0xffffffff) {
		return fmt.Errorf("invalid speed")
	}
	_, err := bus.dev.SetSCLFrequency(uint32(speedInHz))
	return err
}
