package node

import (
	"encoding/json"
	"fmt"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
)

func (n *Node) Send(peer p2p.Peer, method p2p.RpcMethod, payload any) error {
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize %+v: %v", payload, err)
	}

	rpc := p2p.Rpc{
		Method:  method,
		Payload: payloadData,
	}

	if err := peer.Send(rpc); err != nil {
		return fmt.Errorf("failed to send rpc %+v: %v", rpc, err)
	}

	return nil
}

const (
	GetBlock         p2p.RpcMethod = "getBlock"
	GetBlockResponse p2p.RpcMethod = GetBlock + "Response"
	BroadcastBlock   p2p.RpcMethod = "broadcastBlock"
	GetPeers         p2p.RpcMethod = "getPeers"
	GetPeersResponse p2p.RpcMethod = GetPeers + "Response"
)

type GetBlockPayload struct {
	Nonce int `json:"nonce"`
}

type GetBlockResponsePayload struct {
	Block blockchain.Block `json:"block"`
	Nonce int              `json:"nonce"`
}

type BroadcastBlockPayload struct {
	Block blockchain.Block `json:"block"`
}

type GetPeersPayload struct{}

type GetPeersResponsePayload struct {
	Peers []string `json:"peers"`
}
