package tinyrpc

import (
	"gozknight.top/tinyrpc/codec"
	"log"
	"net"
	"net/rpc"
)

type Server struct {
	*rpc.Server
}

func (s *Server) Sereve(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("tinyrpc.Serve: accept error: %v\n", err)
			return
		}
		go s.Server.ServeCodec(codec.NewServerCodec(conn))
	}
}
