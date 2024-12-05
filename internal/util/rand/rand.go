package rand

import (
	"crypto/rand"
	"encoding/base64"
)

func Bytes(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HexString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
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
