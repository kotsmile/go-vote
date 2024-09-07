package p2p

type RpcMethod string

type Rpc struct {
	From    string
	Method  RpcMethod `json:"method"`
	Payload any       `json:"payload"`
}
