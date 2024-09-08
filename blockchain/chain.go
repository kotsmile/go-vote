package blockchain

import (
	"errors"
	"fmt"
)

type Chain struct {
	Blocks []Block
}

func NewChain(initBlocks []Block) Chain {
	return Chain{
		Blocks: initBlocks,
	}
}

var (
	ErrIncorrectGenesisBlock  = errors.New("incorrect genesis block")
	ErrIncorrectPrevBlockHash = errors.New("incorrect prevBlockHash")
	ErrIncorrectNonce         = errors.New("incorrect nonce")
	ErrIncorrectSignature     = errors.New("incorrect signature")
	ErrBlockIncluded          = errors.New("block has been included")
)

func (c *Chain) Validate() (bool, error) {
	for i, block := range c.Blocks {
		if i == 0 {
			if !block.Equal(GenesisBlock) {
				return false, ErrIncorrectGenesisBlock
			}
			continue
		}

		prevBlock := c.Blocks[i-1]
		if block.PrevBlockHash != prevBlock.BlockHash {
			return false, ErrIncorrectPrevBlockHash
		}

		if block.Nonce-1 != prevBlock.Nonce {
			return false, ErrIncorrectNonce
		}

		res, err := block.Verify()
		if err != nil {
			return false, fmt.Errorf("failed to verify: %v", err)
		}
		if !res {
			return false, ErrIncorrectSignature
		}
	}

	return true, nil
}

func (c *Chain) PushBlock(b Block) (bool, error) {
	for _, block := range c.Blocks {
		if block.Equal(b) {
			return false, ErrBlockIncluded
		}
	}

	newChain := NewChain(append(c.Blocks, b))
	ok, err := newChain.Validate()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("not ok")
	}

	c.Blocks = append(c.Blocks, b)
	return true, nil
}

func (c *Chain) GetLastBlock() Block {
	block, _ := c.GetBlock(c.Length() - 1)
	return block
}

func (c Chain) Length() int {
	return len(c.Blocks)
}

func (c Chain) GetBlock(nonce int) (Block, bool) {
	if nonce >= len(c.Blocks) {
		return Block{}, false
	}

	return c.Blocks[nonce], true
}

func (c Chain) Print() {
	for _, block := range c.Blocks {
		block.Print()
		fmt.Println()
	}
}
