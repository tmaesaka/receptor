package index

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const (
	SHA224DexIDLen    = sha256.Size224
	SHA224DexIDHexLen = sha256.Size224 * 2
)

type SHA224DexID struct {
	value [SHA224DexIDLen]byte
}

// Algo implements the ID interface.
func (id SHA224DexID) Algo() uint8 {
	return DexIDTypeSHA224
}

// Value implements the ID interface.
func (id SHA224DexID) Value() []byte {
	return id.value[:]
}

// HexString implements the ID interface.
func (id SHA224DexID) HexString() string {
	return hex.EncodeToString(id.value[:])
}

// SetValue implements the ID interface.
func (id *SHA224DexID) SetValue(val []byte) {
	copy(id.value[:], val)
}

// SetValueFromString implements the ID interface.
func (id *SHA224DexID) SetValueFromString(val string) {
	id.value = sha256.Sum224([]byte(val))
}

// SetValueFromHex implements the ID interface.
func (id *SHA224DexID) SetValueFromHexString(val string) error {
	if len(val) != SHA224DexIDHexLen {
		return errors.New("invalid hex string")
	}
	raw, err := hex.DecodeString(val)
	if err != nil {
		return err
	}
	copy(id.value[:], raw)
	return nil
}
