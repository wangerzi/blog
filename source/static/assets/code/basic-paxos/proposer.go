package main

import (
	"errors"
	"log"
	"strconv"
)

type Proposer struct {
	id        int   // 提议者服务器id
	round     int   // 已知最大轮次
	number    int   // 提案编号（轮次，服务器id）
	acceptors []int // 接收者
}

// 生成提案号，轮次在前，服务器在后，方便比较
func (p *Proposer) generateNumber() int {
	return p.round << 16 | p.id
}

// 发起新提案
func (p *Proposer) propose(v interface{}) (interface{}, error) {
	p.round++
	p.number = p.generateNumber()
	log.Println("[info] propose ", p.number, v)

	prepareCount := 0
	maxNumber := 0
	majorAcceptorNum := len(p.acceptors) / 2 + 1

	// todo:: acceptor 打乱顺序
	// 第一阶段
	for _, aid := range p.acceptors {
		args := MsgArgs{
			Number: p.number,
			From: p.id,
			To: aid,
		}
		if p.number > maxNumber {
			maxNumber = p.number
		}

		reply := new(MsgReply)

		err := call(generateAddr(aid, AddrAccept), "Acceptor.Prepare", args, reply)
		if err != nil {
			// 允许错误
			log.Println("[error] prepare response error", aid, err)
			continue
		}

		if reply.Ok {
			prepareCount++
			// 如果 prepare 发现有提案号更高的，更新提案值
			if reply.Number > maxNumber {
				maxNumber = reply.Number
				v = reply.Value
			}
		} else {
			log.Println("reply is not valid", aid, reply)
		}

		if prepareCount >= majorAcceptorNum {
			break
		}
	}
	if prepareCount < majorAcceptorNum {
		return nil, errors.New("acceptor prepare is not enough" + strconv.Itoa(prepareCount) + "-" + strconv.Itoa(majorAcceptorNum))
	}

	// 第二阶段
	acceptCount := 0
	for _, aid := range p.acceptors {
		args := MsgArgs{
			Number: maxNumber,
			Value: v,
			From: p.id,
			To: aid,
		}

		reply := new(MsgReply)

		err := call(generateAddr(aid, AddrAccept), "Acceptor.Accept", args, reply)
		if err != nil {
			// 允许错误
			log.Println("[error] accept response error", aid, err)
			continue
		}

		if reply.Ok {
			acceptCount++
		} else {
			log.Println("[error] accept reply is not valid", aid, reply)
		}

		// 超过半数
		if acceptCount >= majorAcceptorNum {
			return v, nil
		}
	}

	return nil, errors.New("acceptor accept is not enough")
}