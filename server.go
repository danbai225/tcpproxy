package main

import (
	logs "github.com/danbai225/go-logs"
	"github.com/xtaci/smux"
	"io"
	"net"
)

type Server struct {
	listenAddr string
	pass       string
	dstAddr    string
}

func (Server) New(sPass, dstAddr, listenAddr string) *Server {
	return &Server{listenAddr: listenAddr, pass: sPass, dstAddr: dstAddr}
}
func (s *Server) Start() {
	//监听端口
	listen, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info("服务端启动成功", s.listenAddr, "->", s.dstAddr)
	for listen != nil {
		c, err2 := listen.Accept()
		if err2 != nil {
			logs.Err(err)
			continue
		}
		//认证连接
		go s.auth(c)
	}
}
func (s *Server) auth(conn net.Conn) {
	bytes := make([]byte, 1024)
	//封装连接io（加密）
	UserAuthS, err := createAuth(conn, s.pass)
	if err != nil {
		logs.Err(err)
		return
	}
	read, err := UserAuthS.Read(bytes)
	if err != nil {
		logs.Err(err)
		return
	}
	if s.pass == string(bytes[:read]) {
		_, _ = UserAuthS.Write([]byte("ok"))
	} else {
		_, _ = UserAuthS.Write([]byte("err"))
		logs.Info("认证失败", conn.RemoteAddr())
	}
	//对连接进行复用
	ss, err3 := smux.Server(UserAuthS, nil)
	if err3 != nil {
		logs.Err(conn)
		return
	}
	go s.handC(ss)
}
func (s *Server) handC(session *smux.Session) {
	defer func() {
		if session != nil {
			_ = session.Close()
		}
	}()
	for session != nil {
		// Accept a stream
		stream, err := session.AcceptStream()
		if err != nil {
			if err != io.EOF {
				logs.Err(err)
			}
			return
		}
		go s.handStream(stream)
	}
}
func (s *Server) handStream(stream *smux.Stream) {
	dial, err := net.Dial("tcp", s.dstAddr)
	if err != nil {
		return
	}
	defer func() {
		if stream != nil {
			_ = stream.Close()
		}
		if dial != nil {
			_ = dial.Close()
		}
	}()
	//对拷流量
	go func() {
		io.Copy(dial, stream)
		dial.Close()
	}()
	io.Copy(stream, dial)
}
