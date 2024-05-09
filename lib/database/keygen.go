package database

import (
	"encoding/binary"
)

func Keygen(typeByte byte, injectionIDs ...uint64) []byte {
	keyLength := 1 + (len(injectionIDs) * binary.MaxVarintLen64)
	key := make([]byte, keyLength)
	key[0] = typeByte
	for offset, id := range injectionIDs {
		index := 1 + (offset * binary.MaxVarintLen64)
		binary.PutUvarint(key[index:], id)
	}

	return key
}
