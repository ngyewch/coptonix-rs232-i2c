package coptonixrs232i2c

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"go.bug.st/serial"
)

var (
	ErrInvalidResponse   = fmt.Errorf("invalid response")
	ErrInterruptDetected = fmt.Errorf("interrupt detected")
	ErrChecksumError     = fmt.Errorf("checksum error")
	ErrUnknownCommand    = fmt.Errorf("unknown command")
	ErrCommandMismatch   = fmt.Errorf("command mismatch")
)

type Dev struct {
	port  serial.Port
	mutex sync.Mutex
}

func New(port serial.Port) *Dev {
	return &Dev{
		port: port,
	}
}

func (dev *Dev) Close() error {
	return dev.port.Close()
}

func (dev *Dev) ReadI2C(addr uint8, bytesToRead uint8) ([]uint8, bool, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.writeRequestString(fmt.Sprintf("R%02X%02X", addr, bytesToRead))
	if err != nil {
		return nil, false, err
	}
	response, err := dev.readResponseString()
	if err != nil {
		return nil, false, err
	}
	responseData, err := handleResponse('R', response)
	if err != nil {
		return nil, false, err
	}
	if len(responseData) < 3 {
		return nil, false, fmt.Errorf("invalid response, too short")
	}
	responseCode := responseData[1]
	responseSlaveAddr := responseData[2]
	responseCount := responseData[(len(responseData) - 1)]
	if responseSlaveAddr != addr {
		return nil, false, fmt.Errorf("invalid response, slave address mismatch")
	}
	if responseCount != bytesToRead {
		return nil, false, fmt.Errorf("invalid response, count mismatch")
	}
	switch responseCode {
	case 0x00:
		return nil, false, nil
	case 0x01:
		data := responseData[3 : len(responseData)-1]
		if len(data) != int(bytesToRead) {
			return data, false, fmt.Errorf("invalid response, read count mismatch")
		}
		return data, true, nil
	default:
		return nil, false, fmt.Errorf("invalid response, unknown result code")
	}
}

func (dev *Dev) WriteI2C(addr uint8, data []uint8) (bool, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.writeRequestString(fmt.Sprintf("w%02X%s", addr, strings.ToUpper(hex.EncodeToString(data))))
	if err != nil {
		return false, err
	}
	response, err := dev.readResponseString()
	if err != nil {
		return false, err
	}
	responseData, err := handleResponse('w', response)
	if err != nil {
		return false, err
	}
	if len(responseData) < 3 {
		return false, fmt.Errorf("invalid response, too short")
	}
	responseCode := responseData[1]
	responseSlaveAddr := responseData[2]
	if responseSlaveAddr != addr {
		return false, fmt.Errorf("invalid response, slave address mismatch")
	}
	switch responseCode {
	case 0x00:
		return false, nil
	case 0x01:
		return true, nil
	default:
		return false, fmt.Errorf("invalid response, unknown result code")
	}
}

func (dev *Dev) CheckSlaveAddress(addr uint8) (bool, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.writeRequestString(fmt.Sprintf("c%02X", addr))
	if err != nil {
		return false, err
	}
	response, err := dev.readResponseString()
	if err != nil {
		return false, err
	}
	responseData, err := handleResponse('c', response)
	if err != nil {
		return false, err
	}
	if len(responseData) != 3 {
		return false, fmt.Errorf("invalid response, too short/long")
	}
	responseCode := responseData[2]
	responseSlaveAddr := responseData[1]
	if responseSlaveAddr != addr {
		return false, fmt.Errorf("invalid response, slave address mismatch")
	}
	switch responseCode {
	case 0x00:
		return false, nil
	case 0x01:
		return true, nil
	default:
		return false, fmt.Errorf("invalid response, unknown result code")
	}
}

func (dev *Dev) ScanI2C() ([]uint8, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.writeRequestString("C")
	if err != nil {
		return nil, err
	}
	response, err := dev.readResponseString()
	if err != nil {
		return nil, err
	}
	responseData, err := handleResponse('C', response)
	if err != nil {
		return nil, err
	}
	if len(responseData) < 2 {
		return nil, fmt.Errorf("invalid response, too short")
	}
	count := int(responseData[1])
	if len(responseData)-2 != count {
		return nil, fmt.Errorf("invalid response, count mismatch")
	}
	return responseData[2:], nil
}

func (dev *Dev) GetSCLFrequency() (uint32, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.writeRequestString("I")
	if err != nil {
		return 0, err
	}
	response, err := dev.readResponseString()
	if err != nil {
		return 0, err
	}
	responseData, err := handleResponse('I', response)
	if err != nil {
		return 0, err
	}
	if len(responseData) != 5 {
		return 0, fmt.Errorf("invalid response, too short/long")
	}
	freq := binary.LittleEndian.Uint32(responseData[1:5])
	return freq, nil
}

func (dev *Dev) SetSCLFrequency(freq uint32) (uint32, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.writeRequestString(fmt.Sprintf("E%s", hex.EncodeToString(binary.LittleEndian.AppendUint32(nil, freq))))
	if err != nil {
		return 0, err
	}
	response, err := dev.readResponseString()
	if err != nil {
		return 0, err
	}
	responseData, err := handleResponse('E', response)
	if err != nil {
		return 0, err
	}
	if len(responseData) != 5 {
		return 0, fmt.Errorf("invalid response, too short/long")
	}
	newFreq := binary.LittleEndian.Uint32(responseData[1:5])
	return newFreq, nil
}

func (dev *Dev) writeRequestString(s string) error {
	checksum, err := calcChecksumForString(s)
	if err != nil {
		return err
	}
	line := s + fmt.Sprintf("%02X\r", checksum)
	data := []byte(line)
	slog.Debug("write request string",
		slog.String("data", line),
	)
	_, err = dev.port.Write(data)
	return nil
}

func (dev *Dev) readResponseString() (string, error) {
	buf := bytes.NewBuffer(nil)
	b := make([]byte, 1)
	for {
		n, err := dev.port.Read(b)
		if err != nil {
			return "", err
		}
		if n == 0 {
			break
		}
		if b[0] == '\r' {
			break
		}
		_, err = buf.Write(b[:n])
		if err != nil {
			return "", err
		}
	}
	s := buf.String()
	slog.Debug("read response string",
		slog.String("data", s),
	)
	return s, nil
}

func handleResponse(command uint8, s string) ([]byte, error) {
	dataString, err := verifyChecksumForString(s)
	if err != nil {
		return nil, err
	}
	data, err := hex.DecodeString(dataString)
	if err != nil {
		return nil, ErrInvalidResponse
	}
	if len(data) < 1 {
		return nil, ErrInvalidResponse
	}
	switch data[0] {
	case 0x70:
		if (len(data) == 2) && (data[1] == 0x01) {
			return nil, ErrInterruptDetected
		}
	case 0x73:
		if (len(data) == 2) && (data[1] == 0x01) {
			return nil, ErrChecksumError
		}
	case 0xff:
		if (len(data) == 2) && (data[1] == 0x00) {
			return nil, ErrUnknownCommand
		}
	}
	if data[0] != command {
		return nil, ErrCommandMismatch
	}
	return data, nil
}
