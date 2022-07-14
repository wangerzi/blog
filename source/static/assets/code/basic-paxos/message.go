package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// MsgArgs 请求结构
type MsgArgs struct {
	Number int
	Value  interface{}
	From   int
	To     int
}

// MsgReply 响应结构
type MsgReply struct {
	Ok     bool
	Number int
	Value  interface{}
}

// 约定：接受者1开头，学习者端口2开头
const (
	AddrAccept = iota + 1
	AddrLearner = iota + 1
)
func generateAddr(id int, t int) string {
	return fmt.Sprintf("127.0.0.1:%d", id)
}

func call(addr string, name string, args interface{}, reply interface{}) error {
	c, err := rpc.Dial("tcp",addr)

	if err != nil {
		return err
	}

	defer func(c *rpc.Client) {
		_ = c.Close()
	}(c)

	err = c.Call(name, args, reply)
	if err != nil {
		return err
	}

	return nil
}

func server(addr string, obj interface{}) (net.Listener, error) {
	s := rpc.NewServer()
	if err := s.Register(obj); err != nil {
		return nil, err
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				// ignore error and continue
				log.Println("[error] tcp", err)
				break
			}

			go s.ServeConn(conn)
		}
	}()

	return l, nil
}