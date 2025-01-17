package blockchain

import (
	"time"
)

type Block struct {
	Index     int64
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
	Nonce     int64
}

func NewBlock(index int64, data []byte, prevHash []byte) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Data:      data,
		PrevHash:  prevHash,
		Hash:      []byte{},
		Nonce:     0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock(0, []byte("Genesis Block"), []byte{})
}

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) AddBlock(data []byte) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
