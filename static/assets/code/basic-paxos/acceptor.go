package main

import (
	"log"
	"net"
)

type Acceptor struct {
	listener    net.Listener
	id          int         // 服务器id
	minProposal int         // 最小提案号
	acceptedN   int         // 已接受的值
	acceptedV   interface{} // 已接受的提案值
	learner     []int       // 学习者id
}

func (a *Acceptor) Prepare(args *MsgArgs, reply *MsgReply) error {
	if args.Number > a.minProposal {
		a.minProposal = args.Number
		reply.Number = a.acceptedN
		reply.Value = a.acceptedV
		reply.Ok = true
	} else {
		reply.Ok = false
	}

	return nil
}

func (a *Acceptor) Accept(args *MsgArgs, reply *MsgReply) error {
	if args.Number >= a.minProposal {
		a.minProposal = args.Number
		a.acceptedN = args.Number
		a.acceptedV = args.Value
		reply.Ok = true

		// 转发给学习者
		for _, lid := range a.learner {
			go func(learner int) {
				addr := generateAddr(learner, AddrLearner)
				args.From = a.id
				args.To = learner

				resp := new(MsgReply)

				err := call(addr, "Learner.Learn", args, resp)
				if err != nil {
					log.Println("[error] learn failed", err)
				}
				if !resp.Ok {
					log.Println("[error] learn response is not valid")
				}
			}(lid)
		}
	} else {
		log.Println("[warn] can't accept", args, a)
		reply.Ok = false
	}
	return nil
}

func newAcceptor(id int, learner []int) (*Acceptor, error) {
	acceptor := &Acceptor{
		listener:    nil,
		id:          id,
		minProposal: 0,
		acceptedN:   0,
		acceptedV:   nil,
		learner:      learner,
	}
	listener, err := server(generateAddr(id, AddrAccept), acceptor)
	if err != nil {
		return nil, err
	}

	acceptor.listener = listener

	return acceptor, nil
}