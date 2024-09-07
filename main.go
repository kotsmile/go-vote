package main

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/kotsmile/go-vote/blockchain"
)

func main() {
	wallet := blockchain.NewRandomWallet()

	block, err := blockchain.NewBlock(blockchain.GenesisBlock, wallet, []byte{})

	block.Mine(0, math.MaxUint64)
	block.Sign()

	jsonBlock, err := json.Marshal(block)
	if err != nil {
		fmt.Printf("failed to json object %+v: %v\n", block, err)
		return
	}

	fmt.Printf("block: %s\n", string(jsonBlock))
}
