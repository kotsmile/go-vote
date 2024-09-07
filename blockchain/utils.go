package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func GetDifficulty(data [32]byte) uint64 {
	var difficulty uint64 = 0

	for i := 0; i < len(data); i++ {
		if data[i] == 0 {
			difficulty++
		} else {
			break
		}
	}

	return difficulty
}

func Hash(payload any) ([32]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal %v: %v", payload, err)
	}

	return sha256.Sum256(data), nil
}
