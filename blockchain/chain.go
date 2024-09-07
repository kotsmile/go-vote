package blockchain

import "fmt"

type Chain struct {
	Blocks []Block
}

func NewChain(initBlocks []Block) Chain {
	return Chain{
		Blocks: initBlocks,
	}
}

func (c *Chain) Validate() bool {
	for i, block := range c.Blocks {
		if i == 0 {
			if !block.Equal(GenesisBlock) {
				fmt.Println("incorrect genesis block")
				return false
			}
			continue
		}

		prevBlock := c.Blocks[i-1]
		if block.PrevBlockHash != prevBlock.BlockHash {
			fmt.Println("incorrect prevBlockHash")
			return false
		}

		if block.Nonce-1 != prevBlock.Nonce {
			fmt.Println("incorrect nonce")
			return false
		}

		res, err := block.Verify()
		if err != nil {
			return false
		}
		if !res {
			fmt.Println("incorrect signature")
			return false
		}
	}

	return true
}

func (c *Chain) PushBlock(b Block) bool {
	newChain := NewChain(append(c.Blocks, b))
	if !newChain.Validate() {
		return false
	}

	c.Blocks = append(c.Blocks, b)
	return true
}

func (c *Chain) GetLastBlock() Block {
	return c.GetBlock(c.Length() - 1)
}

func (c Chain) Length() int {
	return len(c.Blocks)
}

func (c Chain) GetBlock(nonce int) Block {
	return c.Blocks[nonce]
}

func (c Chain) PrintDebug() {
	for i, block := range c.Blocks {
		fmt.Printf("Block #%d\n", i)
		fmt.Printf("\tprevBlockHash: %s\n", block.PrevBlockHash)
		fmt.Printf("\tnonce: %d\n", block.Nonce)
		fmt.Printf("\tfrom: %s\n", block.From)
		fmt.Printf("\tdata: %s\n", string(block.Data))
		fmt.Printf("\tdifficulty: %d\n", block.Difficulty)
		fmt.Printf("\tsalt: %d\n", block.Salt)
		fmt.Printf("\tblockHash: %s\n", block.BlockHash)
		fmt.Println()
	}
}
