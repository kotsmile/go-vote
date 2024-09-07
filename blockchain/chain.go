package blockchain

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
				return false
			}
			continue
		}

		prevBlock := c.Blocks[i-1]
		if block.PrevBlockHash != prevBlock.BlockHash {
			return false
		}

		res, err := block.Verify()
		if err != nil {
			return false
		}
		if !res {
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

func (c *Chain) GetLatestBlock() Block {
	return c.Blocks[len(c.Blocks)-1]
}
