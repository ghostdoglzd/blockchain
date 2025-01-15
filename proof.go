package main

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

const targetBits = 24
const MaxNonce = math.MaxInt64

type proofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *proofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &proofOfWork{b, target}
	return pow
}

func (pow *proofOfWork) prepareData(nonce int64) []byte {
	data := [][]byte{
		pow.block.PrevHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}
	return bytes.Join(data, []byte{})
}

func (pow *proofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0
	for nonce < MaxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *proofOfWork) validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
