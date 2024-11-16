package rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func Bytes(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func HexString(bytes []byte) string {
	key := new(big.Int).SetBytes(bytes)
	base16str := fmt.Sprintf("%X", key)
	return base16str
}

func String(length uint) (string, error) {
	var realLength uint
	if length%2 == 0 {
		realLength = length / 2 //nolint:mnd
	} else {
		realLength = (length / 2) + 1 //nolint:mnd
	}
	bytes, err := Bytes(realLength)
	if err != nil {
		return "", err
	}
	return HexString(bytes)[:length], nil
}
