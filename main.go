package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/node"
	"github.com/kotsmile/go-vote/p2p"
)

// block, err := blockchain.NewBlock(blockchain.GenesisBlock, wallet, []byte{})
//
// block.Mine(3973551, 3973551+1)
// block.Sign()
//
// jsonBlock, err := json.Marshal(block)
// if err != nil {
// 	fmt.Printf("failed to json object %+v: %v\n", block, err)
// 	return
// }
//
// fmt.Printf("block: %s\n", string(jsonBlock))
//
// fmt.Println(block.Verify())

const MainNodeAddr = ":3001"

func main() {
	var wg sync.WaitGroup

	fmt.Println("starting main node")
	mainNode := node.NewNode(&wg, p2p.NewTcpTransport(MainNodeAddr), blockchain.NewRandomWallet())
	go mainNode.Start()
	time.Sleep(time.Second * 1)

	fmt.Println("starting user node 1")
	wallet := blockchain.NewWalletFromString("8ed1d4ab8975e20a666f42783be40a345f1acffbf9660db9bd93a87883f4ff6c")
	node1 := node.NewNode(&wg, p2p.NewTcpTransport(":3002"), wallet)
	go node1.Start()
	time.Sleep(time.Second * 1)

	fmt.Println("connecting to main node")
	node1.Transport.Dial(MainNodeAddr)
	time.Sleep(time.Second * 1)

	if err := node1.SendVoting(blockchain.NewVoting("test voting")); err != nil {
		panic(fmt.Errorf("failed to send voting: %v", err))
	}
	if err := node1.SendVoting(blockchain.NewVoting("test voting2")); err != nil {
		panic(fmt.Errorf("failed to send voting: %v", err))
	}

	node1.Chain.PrintDebug()

	wg.Wait()
}
