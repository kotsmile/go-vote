package node

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
)

type Node struct {
	Signer    blockchain.Wallet
	Transport p2p.Transport
	Peers     map[string]p2p.Peer
	Chain     blockchain.Chain
	wg        *sync.WaitGroup
}

func NewNode(wg *sync.WaitGroup, transport p2p.Transport, signer blockchain.Wallet) *Node {
	wg.Add(1)
	server := Node{
		Transport: transport,
		Signer:    signer,
		Chain:     blockchain.NewChain([]blockchain.Block{blockchain.GenesisBlock}), // TODO: get from sqlite
		Peers:     make(map[string]p2p.Peer),
		wg:        wg,
	}

	transport.SetOnPeer(server.onPeer)
	return &server
}

func (n *Node) Start() error {
	if err := n.Transport.ListenAndAccept(); err != nil {
		return fmt.Errorf("failed to start transport")
	}

	for {
		rpc := <-n.Transport.Consume()
		switch rpc.Method {
		case GetBlockRpcMethod:
			peer, ok := n.Peers[rpc.From]
			if !ok {
				fmt.Printf("failed to find peer with %s", rpc.From)
				continue
			}

			var payload GetBlockPayload
			if err := json.Unmarshal(rpc.Payload, &payload); err != nil {
				fmt.Printf("failed deserialize payload %v: %v", rpc.Payload, err)
				continue
			}

			block := n.Chain.GetLastBlock()
			if payload.Nonce != -1 && payload.Nonce < n.Chain.Length() {
				block = n.Chain.GetBlock(payload.Nonce)
			}

			if err := n.Send(peer, GetBlockResponseRpcMethod, GetBlockResponsePayload{
				Block: block,
			}); err != nil {
				fmt.Printf("failed to send response on %s: %v\n", GetBlockRpcMethod, err)
				continue
			}
		case GetBlockResponseRpcMethod:
			var payload GetBlockResponsePayload

			if err := json.Unmarshal(rpc.Payload, &payload); err != nil {
				fmt.Printf("failed deserialize payload %v: %v", rpc.Payload, err)
				continue
			}

			// TODO: add processor
			fmt.Printf("block %+v", payload.Block)
		}
	}
}

func (n *Node) SendVoting(voting blockchain.Voting) error {
	lastBlock := n.Chain.GetLastBlock()

	newBlock, err := blockchain.NewBlock(lastBlock, n.Signer, voting.Data())
	if err != nil {
		return fmt.Errorf("failed to create new block: %v", err)
	}

	if err := newBlock.Mine(0, math.MaxUint64); err != nil {
		return fmt.Errorf("failed to mine block %+v: %v", newBlock, err)
	}

	if err := newBlock.Sign(); err != nil {
		return fmt.Errorf("failed to sign block %+v: %v", newBlock, err)
	}

	if !n.Chain.PushBlock(newBlock) {
		return fmt.Errorf("failed to push block %+v", newBlock)
	}

	return nil
}

func (n *Node) onPeer(peer p2p.Peer) error {
	n.Peers[peer.Addr()] = peer
	return nil
}
