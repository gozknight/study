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

func NewServer() *Server {
	return &Server{&rpc.Server{}}
}

func (s *Server) Serve(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Print("tinyrpc.Serve: accept:", err.Error())
			return
		}
		go s.Server.ServeCodec(codec.NewServerCodec(conn)) // 使用TinyRPC的解码器
	}
}
