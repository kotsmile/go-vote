package main

import (
	"fmt"
	"time"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
	"github.com/kotsmile/go-vote/server"
)

func main() {
	wallet := blockchain.NewWalletFromString("8ed1d4ab8975e20a666f42783be40a345f1acffbf9660db9bd93a87883f4ff6c")

	server1 := server.NewServer(p2p.NewTcpTransport(":3001"), wallet)
	server2 := server.NewServer(p2p.NewTcpTransport(":3002"), wallet)

	go server1.Start()
	go server2.Start()

	time.Sleep(5 * time.Second)

	server1.Transport.Dial(server2.Transport.Addr())

	time.Sleep(5 * time.Second)

	var peer p2p.Peer
	for _, v := range server1.Peers {
		peer = v
		break
	}
	if peer == nil {
		panic("peer is nil")
	}

	response, err := server.SendGetLatestBlock(server1, peer, server.GetLatestBlockPayload{})
	if err != nil {
		fmt.Printf("failed to send get latest block: %v\n", err)
		return
	}

	fmt.Printf("response: %+v\n", response)

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
}
