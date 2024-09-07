package main

import (
	"fmt"
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

func main() {
	wallet := blockchain.NewWalletFromString("8ed1d4ab8975e20a666f42783be40a345f1acffbf9660db9bd93a87883f4ff6c")

	node1 := node.NewNode(p2p.NewTcpTransport(":3001"), wallet)
	node2 := node.NewNode(p2p.NewTcpTransport(":3002"), wallet)

	go node1.Start()
	go node2.Start()

	time.Sleep(5 * time.Second)

	node1.Transport.Dial(node2.Transport.Addr())

	time.Sleep(5 * time.Second)

	var peer p2p.Peer
	for _, v := range node1.Peers {
		peer = v
		break
	}
	if peer == nil {
		panic("peer is nil")
	}

	response, err := node1.SendGetLatestBlock(peer, node.GetLatestBlockPayload{})
	if err != nil {
		fmt.Printf("failed to send get latest block: %v\n", err)
		return
	}

	fmt.Printf("response: %+v\n", response)
}
