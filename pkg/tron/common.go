package tron

import (
	"crypto/sha256"
	"fmt"

	"github.com/shengdoushi/base58"
)

type Address []byte

const (
	AddressLength  = 21
	TronBytePrefix = byte(0x41)
	addressLength  = 20
	prefixMainnet  = 0x41
)

func Encode(input []byte) string {
	return base58.Encode(input, base58.BitcoinAlphabet)
}

func EncodeCheck(input []byte) string {
	h256h0 := sha256.New()
	h256h0.Write(input)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	inputCheck := append(append([]byte(nil), input...), h1[:4]...)

	return Encode(inputCheck)
}

func Decode(input string) ([]byte, error) {
	return base58.Decode(input, base58.BitcoinAlphabet)
}

func DecodeCheck(input string) ([]byte, error) {
	decodeCheck, err := Decode(input)
	if err != nil {
		return nil, err
	}

	if len(decodeCheck) < 4 {
		return nil, fmt.Errorf("b58 check error")
	}

	if len(decodeCheck) != addressLength+4+1 {
		return nil, fmt.Errorf("invalid address length: %d", len(decodeCheck))
	}

	if decodeCheck[0] != prefixMainnet {
		return nil, fmt.Errorf("invalid prefix")
	}

	decodeData := decodeCheck[:len(decodeCheck)-4]

	h256h0 := sha256.New()
	h256h0.Write(decodeData)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	if h1[0] == decodeCheck[len(decodeData)] &&
		h1[1] == decodeCheck[len(decodeData)+1] &&
		h1[2] == decodeCheck[len(decodeData)+2] &&
		h1[3] == decodeCheck[len(decodeData)+3] {

		return decodeData, nil
	}

	return nil, fmt.Errorf("b58 check error")
}

func Base58ToAddress(s string) (Address, error) {
	addr, err := DecodeCheck(s)
	if err != nil {
		return nil, fmt.Errorf("base58 decode %q: %w", s, err)
	}
	return addr, nil
}

func (a Address) IsValid() bool {
	if len(a) != AddressLength {
		return false
	}
	if a[0] != TronBytePrefix {
		return false
	}
	encoded := EncodeCheck(a)
	_, err := DecodeCheck(encoded)
	return err == nil
}
