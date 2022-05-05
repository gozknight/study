package znet

import (
	"fmt"
	"github.com/zinx/ziface"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listening at IP : %v, Port : %v is starting\n", s.IP, s.Port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%v:%v", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen %v error: %v\n", s.IPVersion, err)
			return
		}
		fmt.Printf("start Zinx server success: %v is listening\n", s.Name)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept err: %v\n", err)
				continue
			}
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Printf("Recv buf err: %v\n", err)
						continue
					}
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Printf("Write back err: %v\n", err)
						continue
					}
				}
			}()
		}
	}()

}
func (s *Server) Stop() {
	// TODO: 释放资源
}
func (s *Server) Serve() {
	s.Start()
	// TODO: 扩展业务
	select {}
}

// NewServer ss
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
