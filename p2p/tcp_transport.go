package p2p

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type TcpTransport struct {
	listenAddr string

	rpcCh    chan Rpc
	listener net.Listener
	encoder  Encoder

	OnPeer func(Peer) error
}

var _ Transport = (*TcpTransport)(nil)

func NewTcpTransport(listenAddr string) *TcpTransport {
	return &TcpTransport{
		listenAddr: listenAddr,
		rpcCh:      make(chan Rpc, 1024),
		encoder:    Encoder{},
	}
}

func (t *TcpTransport) Addr() string {
	return t.listenAddr
}

func (t *TcpTransport) SetOnPeer(onPeer func(Peer) error) {
	t.OnPeer = onPeer
}

func (t *TcpTransport) Close() error {
	if t.listener != nil {
		t.listener.Close()
	}

	return nil
}

func (t *TcpTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to dial %s: %v", addr, err)
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TcpTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", t.listenAddr, err)
	}

	go t.acceptLoop()

	return nil
}

func (t *TcpTransport) Consume() <-chan Rpc {
	return t.rpcCh
}

func (t *TcpTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("failed to accept conn on %s: %v\n", t.listenAddr, err)
			continue
		}

		go t.handleConn(conn, false)
	}
}

func (t *TcpTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTcpPeer(conn, outbound)
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			fmt.Printf("failed to call 'onPeer': %v", err)
			return
		}
	}

	for {
		rpc := Rpc{}
		err = t.encoder.Decode(conn, &rpc)
		if err != nil {
			return
		}

		rpc.From = peer.Addr()
		t.rpcCh <- rpc
	}
}

type TcpPeer struct {
	net.Conn
	outbound bool

	encoder Encoder
}

var _ Peer = (*TcpPeer)(nil)

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		Conn:     conn,
		outbound: outbound,
		encoder:  Encoder{},
	}
}

func (p *TcpPeer) Send(data Rpc) error {
	encoder := json.NewEncoder(p.Conn)

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode and send data %+v: %v", data, err)
	}

	return nil
}

func (p *TcpPeer) Addr() string {
	return p.Conn.RemoteAddr().String()
}
