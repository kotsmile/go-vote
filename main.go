package main

import (
	"fmt"
	"time"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/node"
	"github.com/kotsmile/go-vote/p2p"
)

const MainNodeAddr = ":3001"

func main() {
	fmt.Println("starting main node")
	mainNode := node.NewNode(p2p.NewTcpTransport(MainNodeAddr), blockchain.NewRandomWallet())
	go mainNode.Start(true)
	time.Sleep(time.Second * 1)

	fmt.Println("starting user node 1")
	wallet := blockchain.NewWalletFromString("8ed1d4ab8975e20a666f42783be40a345f1acffbf9660db9bd93a87883f4ff6c")
	node1 := node.NewNode(p2p.NewTcpTransport(":3002"), wallet)
	go node1.Start(false)
	time.Sleep(time.Second * 1)

	fmt.Println("starting user node 2")
	node2 := node.NewNode(p2p.NewTcpTransport(":3003"), wallet)
	go node2.Start(false)
	time.Sleep(time.Second * 1)

	fmt.Println("connecting to main node")
	node1.Transport.Dial(MainNodeAddr)
	time.Sleep(time.Second * 1)

	fmt.Println("connecting to main node")
	node2.Transport.Dial(node1.Transport.Addr())
	time.Sleep(time.Second * 1)

	if _, err := node1.SendVoting(blockchain.NewVoting("test voting")); err != nil {
		panic(fmt.Errorf("failed to send voting: %v", err))
	}

	time.Sleep(time.Second * 2)

	if _, err := node2.SendVoting(blockchain.NewVoting("test voting2")); err != nil {
		panic(fmt.Errorf("failed to send voting: %v", err))
	}

	select {}
}
