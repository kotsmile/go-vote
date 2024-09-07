package node

import (
	"encoding/json"
	"fmt"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
)

type Node struct {
	Signer    blockchain.Wallet
	Transport p2p.Transport
	Peers     map[string]p2p.Peer
	Chain     blockchain.Chain
}

func NewNode(transport p2p.Transport, signer blockchain.Wallet) *Node {
	server := Node{
		Transport: transport,
		Signer:    signer,
		Chain:     blockchain.NewChain([]blockchain.Block{blockchain.GenesisBlock}), // TODO: get from sqlite
		Peers:     make(map[string]p2p.Peer),
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

			block := n.Chain.GetLatestBlock()
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

func (n *Node) onPeer(peer p2p.Peer) error {
	n.Peers[peer.Addr()] = peer
	return nil
}
