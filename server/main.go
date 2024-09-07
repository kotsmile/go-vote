package server

type Server struct {
	listenAddr string
}

func NewServer(addr string) Server {
	return Server{
		listenAddr: addr,
	}
}
