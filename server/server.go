package server

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kotsmile/go-vote/blockchain"
	"github.com/kotsmile/go-vote/p2p"
)

type Server struct {
	Signer    blockchain.Wallet
	Transport p2p.Transport
	Peers     map[string]p2p.Peer
	Blocks    []blockchain.Block

	responses     map[string]p2p.Rpc
	responseMutex sync.Mutex
}

func NewServer(transport p2p.Transport, signer blockchain.Wallet) *Server {
	server := Server{
		Transport: transport,
		Signer:    signer,
		Blocks:    []blockchain.Block{blockchain.GenesisBlock}, // TODO: get from sqlite
		Peers:     make(map[string]p2p.Peer),
		responses: make(map[string]p2p.Rpc),
	}

	transport.SetOnPeer(server.onPeer)
	return &server
}

func (s *Server) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return fmt.Errorf("failed to start transport")
	}

	for {
		rpc := <-s.Transport.Consume()
		switch rpc.Method {
		case GetLatestBlockRpcMethod:
			conn, ok := s.Peers[rpc.From]
			if !ok {
				fmt.Printf("failed to find peer with %s", rpc.From)
				continue
			}

			block := s.GetLatestBlock()
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
			s.responseMutex.Lock()
			s.responses[rpc.Id] = rpc
			s.responseMutex.Unlock()
			break
		}
	}
}

func (s *Server) GetRpcById(id string) (p2p.Rpc, error) {
	retries := 0
	for {
		s.responseMutex.Lock()
		rpc, ok := s.responses[id]
		if ok {
			delete(s.responses, id)
			s.responseMutex.Unlock()
			return rpc, nil
		}

		s.responseMutex.Unlock()
		time.Sleep(time.Second * 5)
		retries++

		if retries > 3 {
			return p2p.Rpc{}, fmt.Errorf("max retries")
		}
	}
}

func (s *Server) GetLatestBlock() blockchain.Block {
	return s.Blocks[len(s.Blocks)-1]
}

func (s *Server) onPeer(peer p2p.Peer) error {
	s.Peers[peer.Addr()] = peer
	return nil
}
