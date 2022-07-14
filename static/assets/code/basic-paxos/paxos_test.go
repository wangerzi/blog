package main

import (
	"testing"
	"time"
)

func start(acceptorId []int, learnerId []int)([]*Acceptor, []*Learner, error) {
	acceptorList := make([]*Acceptor, len(acceptorId))
	for i := 0; i < len(acceptorId); i++ {
		acceptor, err := newAcceptor(acceptorId[i], learnerId)
		if err != nil {
			return nil, nil, err
		}
		acceptorList[i] = acceptor
	}
	learnerList := make([]*Learner, len(learnerId))
	for i := 0; i < len(learnerId); i++ {
		learner, err := NewLearner(learnerId[i], acceptorId)
		if err != nil {
			return nil, nil, err
		}
		learnerList[i] = learner
	}

	return acceptorList, learnerList, nil
}

func cleanup(acceptorList []*Acceptor, learnerList []*Learner) {
	for _, a := range acceptorList {
		_ = a.listener.Close()
	}
	for _, l := range learnerList {
		_ = l.listener.Close()
	}
}

func TestSingleProposer(t *testing.T) {
	acceptorId := []int {20001, 20002, 20003}
	learnerId := []int {30001}

	acceptorList, learnerList, err := start(acceptorId, learnerId)

	if err != nil {
		t.Error("start failed", err)
	}

	defer cleanup(acceptorList, learnerList)

	proposer := &Proposer{
		id:        1,
		acceptors: acceptorId,
	}

	// 发起提案
	val, err := proposer.propose("hello world")
	if err != nil {
		t.Error("propose failed", err)
	}

	// 检查返回值
	if val != "hello world" {
		t.Error("propose return value is not valid", val)
	}

	// 等一秒同步完毕再检查
	time.Sleep(1 * time.Second)
	choose := learnerList[0].choose()
	if choose != "hello world" {
		t.Error("learn value is not valid", choose)
	}
}

func TestTwoProposer(t *testing.T) {
	acceptorId := []int {20001, 20002, 20003}
	learnerId := []int {30001}

	acceptorList, learnerList, err := start(acceptorId, learnerId)

	if err != nil {
		t.Error("start failed", err)
	}

	defer cleanup(acceptorList, learnerList)

	proposer := &Proposer{
		id:        1,
		acceptors: acceptorId,
	}

	// 发起提案1
	val, err := proposer.propose("hello world")
	if err != nil {
		t.Error("propose failed", err)
	}

	// 检查返回值
	if val != "hello world" {
		t.Error("propose return value is not valid", val)
	}

	// 等一秒同步完毕再检查
	time.Sleep(1 * time.Second)
	choose := learnerList[0].choose()
	if choose != "hello world" {
		t.Error("learn value is not valid", choose)
	}

	// 发起提案2
	proposer2 := &Proposer{
		id:        2,
		acceptors:  []int{20003, 20002, 20001},
	}
	val, err = proposer2.propose("hello world2")
	if err != nil {
		t.Error("propose2 failed", err)
	}
	// 检查返回值
	if val != "hello world2" {
		t.Error("propose2 return value is not valid", val)
	}
	// 再提一次
	val, err = proposer.propose("hello world3")
	if err != nil {
		t.Error("propose3 failed", err)
	}
	// 检查返回值
	if val != "hello world3" {
		t.Error("propose3 return value is not valid", val)
	}
	// 等一秒同步完毕再检查，是提案3的值
	time.Sleep(1 * time.Second)
	choose = learnerList[0].choose()
	if choose != "hello world3" {
		t.Error("learn value is not valid", choose)
	}
}
