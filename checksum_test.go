package coptonixrs232i2c

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc(t *testing.T) {
	testCalcChecksum(t, 0xdb, "wC4A11F225CB0")
	testCalcChecksum(t, 0xba, "7701C4")
	testCalcChecksum(t, 0xbb, "7700C4")
	testCalcChecksum(t, 0xbd, "C")
	testCalcChecksum(t, 0x35, "7301")
}

func TestVerify(t *testing.T) {
	testVerifyChecksum(t, "7701C4", "7701C4BA")
	testVerifyChecksum(t, "7700C4", "7700C4BB")
}

func TestHandleResponse(t *testing.T) {
	testHandleResponse(t, []uint8{0x77, 0x01, 0xc4}, 0x77, "7701C4BA")
	testHandleResponse(t, []uint8{0x77, 0x00, 0xc4}, 0x77, "7700C4BB")
}

func testCalcChecksum(t *testing.T, expected uint8, s string) {
	checksum, err := calcChecksumForString(s)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, checksum)
	}
}

func testVerifyChecksum(t *testing.T, expectedDataPart string, s string) {
	actualDataPart, err := verifyChecksumForString(s)
	if assert.NoError(t, err) {
		assert.Equal(t, expectedDataPart, actualDataPart)
	}
}

func testHandleResponse(t *testing.T, expectedData []byte, command uint8, s string) {
	actualData, err := handleResponse(command, s)
	if assert.NoError(t, err) {
		assert.EqualValues(t, expectedData, actualData)
	}
}
