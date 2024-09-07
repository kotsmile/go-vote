package blockchain

import (
	"encoding/hex"
	"fmt"
)

const (
	DefaultDifficulty = 3
	ZeroHash          = "0000000000000000000000000000000000000000000000000000000000000000"
)

type Block struct {
	signer Wallet

	SignBlockData
	Signature Signature `json:"signature"`
}

var GenesisBlock = Block{
	SignBlockData: SignBlockData{
		BlockData: BlockData{
			PrevBlockHash: ZeroHash,
			Difficulty:    DefaultDifficulty,
		},
		BlockHash: ZeroHash,
	},
}

func NewBlock(prevBlock Block, signer Wallet, data []byte) (Block, error) {
	addr, err := signer.Address()
	if err != nil {
		return Block{}, fmt.Errorf("failed to get address from wallet: %v", err)
	}

	return Block{
		signer: signer,
		SignBlockData: SignBlockData{
			BlockData: BlockData{
				PrevBlockHash: prevBlock.BlockHash,
				Nonce:         prevBlock.Nonce + 1,
				From:          addr,
				Data:          data,
				Difficulty:    DefaultDifficulty,
			},
		},
	}, nil
}

func (b *Block) Mine(start uint64, stop uint64) error {
	for salt := start; salt < stop; salt++ {
		hash, err := b.BlockDataWithSalt(salt).Hash()
		if err != nil {
			return fmt.Errorf("failed to get hash for salt %d: %v", salt, err)
		}

		difficulty := GetDifficulty(hash)

		if difficulty >= b.Difficulty {
			b.Salt = salt
			b.BlockHash = hex.EncodeToString(hash[:])
			break
		}
	}

	return nil
}

func (b *Block) Sign() error {
	signBlockData := SignBlockData{
		BlockData: BlockData{
			PrevBlockHash: b.PrevBlockHash,
			Nonce:         b.Nonce,
			From:          b.From,
			Data:          b.Data,
			Difficulty:    b.Difficulty,
			Salt:          b.Salt,
		},
		BlockHash: b.BlockHash,
	}

	hash, err := signBlockData.Hash()
	if err != nil {
		return fmt.Errorf("failed to get hash for data %v: %v", signBlockData, err)
	}

	signature, err := b.signer.Sign(hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign hash %v: %v", hash, err)
	}

	b.Signature = signature
	return nil
}

func (b *Block) BlockDataWithSalt(salt uint64) BlockData {
	return BlockData{
		PrevBlockHash: b.PrevBlockHash,
		Nonce:         b.Nonce,
		From:          b.From,
		Data:          b.Data,
		Difficulty:    b.Difficulty,
		Salt:          salt,
	}
}

type BlockData struct {
	PrevBlockHash string  `json:"prevBlockHash"`
	Nonce         uint64  `json:"nonce"`
	From          Address `json:"from"`
	Data          []byte  `json:"data"`
	Difficulty    uint64  `json:"difficulty"`
	Salt          uint64  `json:"salt"`
}

func (d BlockData) Hash() ([32]byte, error) {
	return Hash(d)
}

type SignBlockData struct {
	BlockData
	BlockHash string `json:"blockHash"`
}

func (d SignBlockData) Hash() ([32]byte, error) {
	return Hash(d)
}
