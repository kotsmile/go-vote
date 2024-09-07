package server

import (
	"fmt"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
)

// getBlocks

const (
	BroadcastBlockRpcMethod p2p.RpcMethod = "broadcastBlock"
	GetLatestBlockRpcMethod p2p.RpcMethod = "getLatestBlock"
)

type BroadcastBlockPayload struct {
	Block blockchain.Block `json:"block"`
}

type BroadcastBlockResponse struct{}

func SendBroadcastBlock(peer p2p.Peer, payload BroadcastBlockPayload) (BroadcastBlockResponse, error) {
	rpc := p2p.Rpc{
		Method:  BroadcastBlockRpcMethod,
		Payload: payload,
	}

	if err := peer.Send(rpc); err != nil {
		return BroadcastBlockResponse{}, fmt.Errorf("failed to send rpc %+v: %v", rpc, err)
	}

	return BroadcastBlockResponse{}, nil
}

type GetLatestBlockPayload struct{}

type GetLatestBlockResponse struct {
	Block blockchain.Block `json:"block"`
}

func SendGetBlock(peer p2p.Peer, payload GetLatestBlockPayload) (GetLatestBlockResponse, error) {
	rpc := p2p.Rpc{
		Method:  GetLatestBlockRpcMethod,
		Payload: payload,
	}

	if err := peer.Send(rpc); err != nil {
		return GetLatestBlockResponse{}, fmt.Errorf("failed to send rpc %+v: %v", rpc, err)
	}

	var response GetLatestBlockResponse
	if err := peer.Receive(&response); err != nil {
		return GetLatestBlockResponse{}, fmt.Errorf("failed to receive data: %v", err)
	}

	return response, nil
}
