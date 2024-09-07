package main

import (
	"encoding/json"
	"fmt"

	"github.com/kotsmile/go-vote/blockchain"
)

func main() {
	wallet := blockchain.NewWalletFromString("8ed1d4ab8975e20a666f42783be40a345f1acffbf9660db9bd93a87883f4ff6c")

	block, err := blockchain.NewBlock(blockchain.GenesisBlock, wallet, []byte{})

	block.Mine(3973551, 3973551+1)
	block.Sign()

	jsonBlock, err := json.Marshal(block)
	if err != nil {
		fmt.Printf("failed to json object %+v: %v\n", block, err)
		return
	}

	fmt.Printf("block: %s\n", string(jsonBlock))

	fmt.Println(block.Verify())
}
