package blockchain

import (
	"bytes"
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
		blockHash, err := b.BlockDataWithSalt(salt).Hash()
		if err != nil {
			return fmt.Errorf("failed to get hash for salt %d: %v", salt, err)
		}

		difficulty := GetDifficulty(blockHash)

		if difficulty >= b.Difficulty {
			b.Salt = salt
			b.BlockHash = hex.EncodeToString(blockHash[:])
			return nil
		}
	}

	return fmt.Errorf("failed to mine block")
}

func (b Block) Equal(other Block) bool {
	return b.SignBlockData.Equal(other.SignBlockData) &&
		b.Signature == other.Signature
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

func (b *Block) Verify() (bool, error) {
	signerAddress := b.From
	blockHash, err := b.SignBlockData.Hash()
	if err != nil {
		return false, fmt.Errorf("failed to hash block: %v", err)
	}

	blockDataHash, _ := b.BlockData.Hash()
	if b.BlockHash != hex.EncodeToString(blockDataHash[:]) {
		return false, nil
	}

	ok, err := b.Signature.Verify(signerAddress, blockHash[:])
	if err != nil {
		return false, fmt.Errorf("failed to verify signagure: %v", err)
	}

	return ok, nil
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

func (d BlockData) Equal(other BlockData) bool {
	return d.PrevBlockHash == other.PrevBlockHash &&
		d.Nonce == other.Nonce &&
		d.From == other.From &&
		bytes.Equal(d.Data, other.Data) &&
		d.Difficulty == other.Difficulty &&
		d.Salt == other.Salt
}

type SignBlockData struct {
	BlockData
	BlockHash string `json:"blockHash"`
}

func (d SignBlockData) Hash() ([32]byte, error) {
	return Hash(d)
}

func (d SignBlockData) Equal(other SignBlockData) bool {
	return d.BlockData.Equal(other.BlockData) &&
		d.BlockHash == other.BlockHash
}
