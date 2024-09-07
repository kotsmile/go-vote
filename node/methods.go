package node

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
)

const (
	ResponseRpcMethod       p2p.RpcMethod = "response"
	BroadcastBlockRpcMethod p2p.RpcMethod = "broadcastBlock"
	GetLatestBlockRpcMethod p2p.RpcMethod = "getLatestBlock"
)

type BroadcastBlockPayload struct {
	Block blockchain.Block `json:"block"`
}

type BroadcastBlockResponse struct{}

func (n *Node) SendBroadcastBlock(peer p2p.Peer, payload BroadcastBlockPayload) (BroadcastBlockResponse, error) {
	id := uuid.New().String()

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return BroadcastBlockResponse{}, fmt.Errorf("failed to serialize %+v: %v", payload, err)
	}

	rpc := p2p.Rpc{
		Id:      id,
		Method:  BroadcastBlockRpcMethod,
		Payload: payloadData,
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

func (n *Node) SendGetLatestBlock(peer p2p.Peer, payload GetLatestBlockPayload) (GetLatestBlockResponse, error) {
	id := uuid.New().String()

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return GetLatestBlockResponse{}, fmt.Errorf("failed to serialize %+v: %v", payload, err)
	}

	rpc := p2p.Rpc{
		Id:      id,
		Method:  GetLatestBlockRpcMethod,
		Payload: payloadData,
	}

	if err := peer.Send(rpc); err != nil {
		return GetLatestBlockResponse{}, fmt.Errorf("failed to send rpc %+v: %v", rpc, err)
	}

	rpc, err = n.GetRpcById(id)
	if err != nil {
		return GetLatestBlockResponse{}, fmt.Errorf("failed to get rpc response %s: %v", id, err)
	}

	var response GetLatestBlockResponse
	if err := json.Unmarshal(rpc.Payload, &response); err != nil {
		return GetLatestBlockResponse{}, fmt.Errorf("failed to deserialize %s: %v", string(rpc.Payload), err)
	}

	return response, nil
}
