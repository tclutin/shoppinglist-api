package hash

import (
	crypto "crypto/rand"
	"math/big"
)

func NewCryptoRand(size int64) (int64, error) {
	safeNum, err := crypto.Int(crypto.Reader, big.NewInt(size))
	if err != nil {
		return 0, err
	}
	return safeNum.Int64(), nil
}
