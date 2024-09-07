package node

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
)

type Node struct {
	Signer    blockchain.Wallet
	Transport p2p.Transport
	Peers     map[string]p2p.Peer
	Blocks    []blockchain.Block

	responses     map[string]p2p.Rpc
	responseMutex sync.Mutex
}

func NewNode(transport p2p.Transport, signer blockchain.Wallet) *Node {
	server := Node{
		Transport: transport,
		Signer:    signer,
		Blocks:    []blockchain.Block{blockchain.GenesisBlock}, // TODO: get from sqlite
		Peers:     make(map[string]p2p.Peer),
		responses: make(map[string]p2p.Rpc),
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
		case GetLatestBlockRpcMethod:
			conn, ok := n.Peers[rpc.From]
			if !ok {
				fmt.Printf("failed to find peer with %s", rpc.From)
				continue
			}

			block := n.GetLatestBlock()
			payload, err := json.Marshal(GetLatestBlockResponse{
				Block: block,
			})
			if err != nil {
				fmt.Printf("failed to serialize response: %v\n", err)
				continue
			}

			if err := conn.Send(p2p.Rpc{
				Id:      rpc.Id,
				Method:  ResponseRpcMethod,
				Payload: payload,
			}); err != nil {
				fmt.Printf("failed to send response on %+v: %v\n", rpc, err)
			}

		case ResponseRpcMethod:
			n.responseMutex.Lock()
			n.responses[rpc.Id] = rpc
			n.responseMutex.Unlock()
			break
		}
	}
}

func (n *Node) GetRpcById(id string) (p2p.Rpc, error) {
	retries := 0
	for {
		n.responseMutex.Lock()
		rpc, ok := n.responses[id]
		if ok {
			delete(n.responses, id)
			n.responseMutex.Unlock()
			return rpc, nil
		}

		n.responseMutex.Unlock()
		time.Sleep(time.Second * 5)
		retries++

		if retries > 3 {
			return p2p.Rpc{}, fmt.Errorf("max retries")
		}
	}
}

func (n *Node) GetLatestBlock() blockchain.Block {
	return n.Blocks[len(n.Blocks)-1]
}

func (n *Node) onPeer(peer p2p.Peer) error {
	n.Peers[peer.Addr()] = peer
	return nil
}
