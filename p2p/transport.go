package p2p

type Peer interface {
	Send(Rpc) error
	Addr() string
}

type Transport interface {
	SetOnPeer(func(Peer) error)
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Close() error
	Consume() <-chan Rpc
}
