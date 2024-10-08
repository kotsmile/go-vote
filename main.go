package main

import (
	"fmt"
	"time"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/node"
	"github.com/kotsmile/go-vote/p2p"
)

const MainNodeAddr = ":3001"

// TODO: if amount of bad nodes will be more that 50% it is possible to cancel blocks
// add ability for block finilization after creation to be sure that history is correct

func connectionAndBroadcasting() {
	fmt.Println("starting main node")
	mainNode := node.NewNode("main.json", p2p.NewTcpTransport(MainNodeAddr), blockchain.NewRandomWallet())
	go mainNode.Start(true)
	time.Sleep(time.Second * 1)

	fmt.Println("starting user node 1")
	wallet := blockchain.NewWalletFromString("8ed1d4ab8975e20a666f42783be40a345f1acffbf9660db9bd93a87883f4ff6c")
	node1 := node.NewNode("node1.json", p2p.NewTcpTransport(":3002"), wallet)
	go node1.Start(false)
	time.Sleep(time.Second * 1)

	fmt.Println("starting user node 2")
	node2 := node.NewNode("node2.json", p2p.NewTcpTransport(":3003"), wallet)
	go node2.Start(false)
	time.Sleep(time.Second * 1)

	fmt.Println("connecting to main node")
	node1.Connect(MainNodeAddr)
	time.Sleep(time.Second * 1)

	fmt.Println("connecting to main node")
	node2.Connect(node1.Transport.Addr())
	time.Sleep(time.Second * 1)

	if _, err := node1.SendVoting(blockchain.NewVoting("test voting")); err != nil {
		panic(fmt.Errorf("failed to send voting: %v", err))
	}

	time.Sleep(time.Second * 2)

	if _, err := node2.SendVoting(blockchain.NewVoting("test voting2")); err != nil {
		panic(fmt.Errorf("failed to send voting: %v", err))
	}
}

var (
	MainNodeWallet = blockchain.NewWalletFromString("8df93ef0a4f3200125d8d27ab1c0bd0dde92cb11774b1665b8463aa462477294")
	Node1Wallet    = blockchain.NewWalletFromString("7e6b66ecc028718f1ddecc24c2146ed3f3b625edd8deea5c4c0f59aa08e8a6dd")
	Node2Wallet    = blockchain.NewWalletFromString("65ef67a2fb269c2d190edc4fb0d488fd6030f79ed23fd08af672271d0323a2cf")
)

func main() {
	mainAddr, _ := MainNodeWallet.Address()
	fmt.Printf("main node: %s\n", mainAddr[:10])

	node1Addr, _ := Node1Wallet.Address()
	fmt.Printf("node1: %s\n", node1Addr[:10])

	node2Addr, _ := Node2Wallet.Address()
	fmt.Printf("node2: %s\n", node2Addr[:10])

	fmt.Println("starting main node")
	mainNode := node.NewNode("main.json", p2p.NewTcpTransport(MainNodeAddr), MainNodeWallet).WithName("MainNode")
	go mainNode.Start(true)
	time.Sleep(time.Second * 1)

	fmt.Println("starting user node 1")
	node1 := node.NewNode("node1.json", p2p.NewTcpTransport(":3002"), Node1Wallet).WithName("Node1")
	go node1.Start(false)
	time.Sleep(time.Second * 1)

	fmt.Println("starting user node 2")
	node2 := node.NewNode("node2.json", p2p.NewTcpTransport(":3003"), Node2Wallet).WithName("Node2")
	go node2.Start(false)

	time.Sleep(5 * time.Second)

	if err := node1.Connect(MainNodeAddr); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	if err := node2.Connect(MainNodeAddr); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// fmt.Println(mainNode.Chain.String())

	// for {
	//
	// 	fmt.Println("Last block main")
	// 	fmt.Println(mainNode.Chain.GetLastBlock().String())
	//
	// 	fmt.Println("Last block node1")
	// 	fmt.Println(node1.Chain.GetLastBlock().String())
	//
	// 	fmt.Println("Last block node2")
	// 	fmt.Println(node2.Chain.GetLastBlock().String())
	//
	// 	time.Sleep(time.Second)
	// }

	select {}
}
