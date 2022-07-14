package main

import (
	"log"
	"net"
)

type Learner struct {
	listener    net.Listener
	id          int
	acceptedMsg map[int]MsgArgs
}

// Learn 接受 acceptor 发来的学习请求并更新
func (l *Learner) Learn(args *MsgArgs, reply *MsgReply) error {
	a := l.acceptedMsg[args.From]
	log.Println("[info] learn", args)
	if a.Number < args.Number {
		l.acceptedMsg[args.From] = *args
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

// 选出来最多 acceptor 认可的值即可
func (l *Learner) choose() interface{} {
	acceptCounts := make(map[int]int)
	// number -> msg的映射
	acceptMsg := make(map[int]MsgArgs)
	majorCount := len(l.acceptedMsg) / 2

	//log.Println("[info] learner ", l.id, l.acceptedMsg)
	for _, accepted := range l.acceptedMsg {
		if accepted.Number != 0 {
			acceptCounts[accepted.Number]++
			acceptMsg[accepted.Number] = accepted
		}

		if acceptCounts[accepted.Number] > majorCount {
			return accepted.Value
		}
	}

	return nil
}

func NewLearner(id int, acceptorIds []int) (*Learner, error) {
	learner := &Learner{
		listener:    nil,
		id:          id,
		acceptedMsg: make(map[int]MsgArgs),
	}

	listener, err := server(generateAddr(id, AddrLearner), learner)

	if err != nil {
		return learner, err
	}
	learner.listener = listener
	// 初始化
	for _, aid := range acceptorIds {
		learner.acceptedMsg[aid] = MsgArgs{
			Number: 0,
			Value: nil,
		}
	}

	return learner, nil
}