package tcpproxy

import (
	"errors"
	logs "github.com/danbai225/go-logs"
	"github.com/xtaci/smux"
	"io"
	"net"
)

type Client struct {
	pass       string
	lAddr      string
	serverAddr string
	session    *smux.Session
	listen     net.Listener
}

func (Client) New(pass, serverAddr, lAddr string) *Client {
	return &Client{
		pass:       pass,
		lAddr:      lAddr,
		serverAddr: serverAddr,
	}
}
func (c *Client) Start() {
	var err error
	c.listen, err = net.Listen("tcp", c.lAddr)
	if err != nil {
		logs.Err(err)
		return
	}
	err = c.connServer()
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info("客户端启动成功", c.lAddr, "->", c.serverAddr)
	for c.listen != nil {
		accept, err2 := c.listen.Accept()
		if err2 == nil {
			go c.hanC(accept)
		}
	}
}
func (c *Client) Stop() error {
	return c.listen.Close()
}
func (c *Client) connServer() error {
	conn2, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		logs.Err(err)
		return err
	}
	s, _ := createAuth(conn2, c.pass)
	_, err = s.Write([]byte(c.pass))
	if err != nil {
		return err
	}
	bytes := make([]byte, len(c.pass))
	read, err := s.Read(bytes)
	if string(bytes[:read]) != "ok" {
		return errors.New("认证失败")
	}
	c.session, err = smux.Client(s, nil)
	if err != nil {
		logs.Err(err)
		return err
	}
	return nil
}
func (c *Client) hanC(con net.Conn) {
	defer func() {
		if con != nil {
			_ = con.Close()
		}
	}()
	if c.session == nil {
		err := c.connServer()
		if err != nil {
			logs.Err(err)
			return
		}
	}
	// Open a new stream
	stream, err := c.session.OpenStream()
	defer func() {
		if stream != nil {
			stream.Close()
		}
	}()
	if err != nil {
		logs.Err(err)
		_ = c.session.Close()
		c.session = nil
		return
	}
	go func() {
		io.Copy(con, stream)
		con.Close()
	}()
	io.Copy(stream, con)
}
