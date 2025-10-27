package coptonixrs232i2c

import (
	"fmt"
	"strconv"
)

var (
	ErrInvalidChecksum  = fmt.Errorf("invalid checksum")
	ErrChecksumMismatch = fmt.Errorf("checksum mismatch")
)

func calcChecksumForString(s string) (uint8, error) {
	var checksum uint8
	for _, c := range s {
		if c > 127 {
			return 0, fmt.Errorf("non-ASCII character found in input")
		}
		checksum += uint8(c)
	}
	return uint8(0x100 - uint16(checksum)), nil
}

func verifyChecksumForString(s string) (string, error) {
	if len(s) < 2 {
		return "", ErrInvalidChecksum
	}
	checksumPart := s[len(s)-2:]
	actualChecksum, err := strconv.ParseUint(checksumPart, 16, 32)
	if err != nil {
		return "", ErrInvalidChecksum
	}
	dataPart := s[0 : len(s)-2]
	expectedChecksum, err := calcChecksumForString(dataPart)
	if err != nil {
		return "", err
	}
	if uint8(actualChecksum) != expectedChecksum {
		fmt.Printf("s=%s, expected=%02x, actual=%02x\n", s, expectedChecksum, actualChecksum)
		return "", ErrChecksumMismatch
	}
	return dataPart, nil
}
