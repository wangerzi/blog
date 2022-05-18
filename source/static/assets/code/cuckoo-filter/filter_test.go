package main

import (
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New(1<<5, 2, 8)
	if err != nil {
		t.Error(err)
	}
}

func TestCuckooFilter_Add(t *testing.T) {
	filter, err := New(1<<5, 2, 8)
	if err != nil {
		t.Error(err)
	}

	_ = filter.Add([]byte("Hello"))
	_ = filter.Add([]byte("World"))

	res1, err := filter.Contain([]byte("hello"))
	res2, err := filter.Contain([]byte("Hello"))

	if err != nil {
		t.Error(err)
		return
	}

	if res1 != false || res2 != true {
		t.Error("res is not valid", res1, res2)
	}
}

func TestCuckooFilter_AddKickOut(t *testing.T) {
	filter, err := New(1<<5, 2, 8)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 5; i++ {
		err = filter.Add([]byte("Hello"))
	}

	// must kick out
	if err == nil {
		t.Error("repeat item > 2b doesn't kick out")
	}

	res1, err := filter.Contain([]byte("hello"))
	res2, err := filter.Contain([]byte("Hello"))

	if err != nil || res1 != false || res2 != true {
		t.Error("contain should be fine", err, res1, res2)
	}
}

func TestCuckooFilter_Delete(t *testing.T) {
	filter, err := New(1<<5, 2, 8)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 4; i++ {
		err = filter.Add([]byte("Hello"))
	}

	if err != nil {
		t.Error("can't kick out", err)
	}

	res1, err := filter.Contain([]byte("hello"))
	res2, err := filter.Contain([]byte("Hello"))

	if err != nil || res1 != false || res2 != true {
		t.Error("contain should be fine", err, res1, res2)
	}

	// delete not exists
	err = filter.Delete([]byte("hello"))
	if err == nil {
		t.Error("delete not exists item should be error")
	}
	// delete exists and leave 1
	for i := 0; i < 3; i++ {
		err = filter.Delete([]byte("Hello"))
		if err != nil {
			t.Error("delete exists item should be right", err)
		}
	}

	// check contain
	res3, err := filter.Contain([]byte("Hello"))
	if err != nil {
		t.Error("already exists 1", err)
	}

	if res3 != true {
		t.Error("already exists 1, should be true")
	}

	// delete last one
	err = filter.Delete([]byte("Hello"))
	if err != nil {
		t.Error("delete exists item can't be error", err)
	}

	res3, err = filter.Contain([]byte("Hello"))
	if err != nil {
		t.Error("can't be error when not exists", err)
	}

	if res3 != false {
		t.Error("delete all should be false")
	}
}
