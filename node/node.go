package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

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

func (n *Node) Start(verbose bool) error {
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
		case BroadcastBlockRpcMethod:
			var payload BroadcastBlockPayload

			if err := json.Unmarshal(rpc.Payload, &payload); err != nil {
				fmt.Printf("failed deserialize payload %v: %v", rpc.Payload, err)
				continue
			}

			ok, err := n.Chain.PushBlock(payload.Block)
			if err != nil {
				if !errors.Is(err, blockchain.ErrBlockIncluded) {
					fmt.Printf("failed to push block: %v\n", err)
				}
				continue
			}
			if !ok {
				fmt.Printf("not ok")
			}

			if verbose {
				payload.Block.Print()
			}

			if err := n.BroadcastExcept(BroadcastBlockRpcMethod, payload, rpc.From); err != nil {
				fmt.Printf("failed to broadcase block: %v", err)
				continue
			}
		}
	}
}

func (n *Node) SendVoting(voting blockchain.Voting) (string, error) {
	return n.SendData(voting.Data())
}

func (n *Node) SendVote(vote blockchain.Vote) (string, error) {
	return n.SendData(vote.Data())
}

func (n *Node) SendData(data []byte) (string, error) {
	lastBlock := n.Chain.GetLastBlock()

	newBlock, err := blockchain.NewBlock(lastBlock, n.Signer, data)
	if err != nil {
		return "", fmt.Errorf("failed to create new block: %v", err)
	}

	if err := newBlock.Mine(0, math.MaxUint64); err != nil {
		return "", fmt.Errorf("failed to mine block %+v: %v", newBlock, err)
	}

	if err := newBlock.Sign(); err != nil {
		return "", fmt.Errorf("failed to sign block %+v: %v", newBlock, err)
	}

	ok, res := n.Chain.PushBlock(newBlock)
	if !ok {
		return "", fmt.Errorf("failed to push block %+v: %s", newBlock, res)
	}

	if err := n.BroadcastExcept(BroadcastBlockRpcMethod, BroadcastBlockPayload{
		Block: newBlock,
	}, ""); err != nil {
		return "", fmt.Errorf("failed to broadcast new block")
	}

	return newBlock.BlockHash, nil
}

func (n *Node) BroadcastExcept(method p2p.RpcMethod, payload any, exceptAddress string) error {
	for addr, peer := range n.Peers {
		if peer.Addr() == exceptAddress {
			continue
		}
		err := n.Send(peer, method, payload)
		if err != nil {
			return fmt.Errorf("failed to broadcast to %s: %v", addr, err)
		}
	}

	return nil
}

func (n *Node) onPeer(peer p2p.Peer) error {
	n.Peers[peer.Addr()] = peer
	return nil
}
