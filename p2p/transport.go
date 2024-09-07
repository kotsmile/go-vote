package p2p

type Peer interface {
	Send(RPC) error
	Addr() string
}

type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Close() error
	Consume() <-chan RPC
}
