package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const MaxKickOut = 500

type CuckooFilter struct {
	m uint64 // 有多少个桶，强制转换为 2 的次幂
	b uint64 // 每个桶能存多少个指纹
	f uint64 // 每个指纹多少位，方便实现，需要是 8 的倍数，且小于 24 * 8，不然 hash 不够用
	bucket []uint8 // 实际存放桶的空间
}

// CalcHash 计算 h1 和 finger hash 的，具体分配到哪个桶，sha input % m
func (f *CuckooFilter) CalcHash(input []byte) (h1 uint64, err error) {
	// m 最多 8 字节，返回固定 32 字节的数据，取余即可
	sum := sha256.Sum256(input)

	// 取首 8 字节，转为 uint64
	h1, err = fastRemain(binary.BigEndian.Uint64(sum[:8]), f.m)
	if err != nil {
		return
	}
	return
}

// CalcFinger 计算指纹
func (f *CuckooFilter) CalcFinger(input []byte) []byte {
	sum := sha256.Sum256(input)

	return sum[8:8+f.f>>3]
}

// AltHash 计算 h2 以及 KickOut 的新位置的，h1 ^ sha finger % m
func (f *CuckooFilter) AltHash(h uint64, finger []byte) (uint64, error) {
	fingerHash, err := f.CalcHash(finger)
	if err != nil {
		return 0, err
	}

	return h ^ fingerHash, nil
}

// CheckBucketFull 检查桶是不是满的，检查最后一个 finger 即可
func (f *CuckooFilter) CheckBucketFull(index uint64) bool {
	lastFinger := f.GetBucketFinger(index, f.b - 1)

	for i := 0; i < len(lastFinger); i++ {
		if uint8(lastFinger[i]) != 0 {
			return true
		}
	}
	return false
}

// CheckBucketHaveFinger 检查桶内是否有这个指纹
func (f *CuckooFilter) CheckBucketHaveFinger(index uint64, finger []byte) bool {
	empty := f.GetEmptyFinger()

	for i := uint64(0); i < f.b; i++ {
		currFinger := f.GetBucketFinger(index, i)
		if bytes.Equal(currFinger, finger) {
			return true
		}
		if bytes.Equal(currFinger, empty) {
			// 空了就不用比了
			return false
		}
	}
	return false
}

// InsertBucket 桶内插入数据
func (f *CuckooFilter) InsertBucket(index uint64, finger []byte) error {
	empty := f.GetEmptyFinger()
	// 找到一个空的位置就可以放进去了
	for i := uint64(0); i < f.b; i++ {
		currFinger := f.GetBucketFinger(index, i)
		if bytes.Equal(currFinger, empty) {
			f.ReplaceBucket(index, i, finger)
			return nil
		}
	}
	return errors.New("bucket fulled")
}

// ReplaceBucket 替换 Bucket 中某个 finger
func (f *CuckooFilter) ReplaceBucket(index uint64, bucketIndex uint64, finger []byte) []byte {
	origin := f.GetBucketFinger(index, bucketIndex)

	for i := 0; i < len(finger); i++ {
		origin[i] = finger[i]
	}

	return origin
}
func (f *CuckooFilter) GetBucketFinger(index uint64, bucketIndex uint64) []byte {
	fingerByte := f.f>>3
	bucketBase := index * f.b * fingerByte

	return f.bucket[bucketBase+bucketIndex*fingerByte:bucketBase+(bucketIndex+1)*fingerByte]
}

func (f *CuckooFilter) GetEmptyFinger() []byte {
	fingerByte := f.f>>3
	return make([]byte, fingerByte)
}

// RandomReplaceInFullBucket 随机替换满桶中的某个 finger
func (f *CuckooFilter) RandomReplaceInFullBucket(index uint64, finger []byte) (replaceFinger []byte, err error) {
	if !f.CheckBucketFull(index) {
		return nil, errors.New("bucket not full")
	}
	randomIndex := uint64(rand.Intn(int(f.b)))

	replaceFinger= f.ReplaceBucket(index, randomIndex, finger)

	return
}

// DeleteFingerFromBucket 从桶中删除 finger，类似数组操作
func (f *CuckooFilter) DeleteFingerFromBucket(index uint64, finger []byte) error {
	empty := f.GetEmptyFinger()

	for i := uint64(0); i < f.b; i++ {
		currFinger := f.GetBucketFinger(index, i)
		if bytes.Equal(currFinger, finger) {
			f.ReplaceBucket(index, i, empty)
			for j := i; j < f.b; j++ {
				if j == f.b-1 {
					f.ReplaceBucket(index, j, empty)
				} else {
					nextFinger := f.GetBucketFinger(index, j + 1)
					f.ReplaceBucket(index, j, nextFinger)
				}
			}
			return nil
		}
	}
	return errors.New("finger not found")
}

// GetHashFinger 获取 finger
func (f *CuckooFilter) GetHashFinger(input []byte) (h1 uint64, h2 uint64, finger []byte, err error) {
	h1, err = f.CalcHash(input)
	if err != nil {
		return
	}

	finger = f.CalcFinger(input)

	h2, err = f.AltHash(h1, finger)

	return
}

// Contain 是否包含
func (f *CuckooFilter) Contain(input []byte) (bool, error) {
	h1, h2, finger, err := f.GetHashFinger(input)
	if err != nil {
		return false, err
	}
	return f.CheckBucketHaveFinger(h1, finger) || f.CheckBucketHaveFinger(h2, finger), nil
}

// Add 加一个
func (f *CuckooFilter) Add(input []byte) error {
	h1, h2, finger, err := f.GetHashFinger(input)
	if err != nil {
		return err
	}

	h1Full := f.CheckBucketFull(h1)
	h2Full := f.CheckBucketFull(h2)

	rand.Seed(time.Now().UnixNano())
	if !h1Full || !h2Full {
		index := uint64(0)
		if !h1Full && !h2Full {
			// 随机选择 h1 和 h2，有空位就插入
			if rand.Intn(2) > 0 {
				index = h1
			} else {
				index = h2
			}
		} else if !h1Full {
			index = h1
		} else {
			index = h2
		}
		// 放到桶里
		err = f.InsertBucket(index, finger)
		if err != nil {
			return err
		}
		// 成功
		return nil
	} else {
		// kick out
		// 随机选择 h1 和 h2
		currIndex := uint64(0)
		if rand.Intn(2) > 0 {
			currIndex = h1
		} else {
			currIndex = h2
		}
		currFinger := finger
		for i:= 0; i < MaxKickOut; i++ {
			if f.CheckBucketFull(currIndex) {
				// 满了没空位，就随机替换出来
				currFinger, err = f.RandomReplaceInFullBucket(currIndex, currFinger)
				if err != nil {
					return err
				}

				// 把换出来的数据，换一个桶存
				currIndex, err = f.AltHash(currIndex, currFinger)
				if err != nil {
					return err
				}
			} else {
				// 有空位就插
				err = f.InsertBucket(currIndex, currFinger)
				if err != nil {
					return err
				}
				// 成功
				return nil
			}
		}

		return errors.New("over max kick out")
	}
}

// Delete 删除一个
func (f *CuckooFilter) Delete(input []byte) error {
	h1, h2, finger, err := f.GetHashFinger(input)
	if err != nil {
		return err
	}

	// 如果能找到的话，就删掉
	if f.CheckBucketHaveFinger(h1, finger) {
		return f.DeleteFingerFromBucket(h1, finger)
	}
	if f.CheckBucketHaveFinger(h2, finger) {
		return f.DeleteFingerFromBucket(h2, finger)
	}

	return errors.New("finger not found")
}

// checkPow2 检查是否为 2 的整次幂
func checkPow2(a uint64) bool {
	return a & (a-1) == 0
}

// fastRemain 求 a % b，前提是 b 需要是 2 的次幂, remain = a & (b - 1)
func fastRemain(a uint64, b uint64) (uint64, error) {
	if !checkPow2(b) {
		return 0, errors.New("fast remain failed")
	}
	return a & (b - 1), nil
}

// New m 为桶的数量，b 为每个桶内含指纹数量，f 为指纹长度，New(1<<10, 4, 8)
func New(m uint64, b uint64, f uint64) (*CuckooFilter, error) {
	// 确保 m 为 2 的 n 次幂，确保 f 为 8 的倍数
	if !checkPow2(m) {
		return nil, errors.New("m needs 2^n")
	}
	if remain, err := fastRemain(f, 8); err != nil || remain > 0 || f < 8 {
		return nil, errors.New("m needs 8*n and > 0")
	}
	return &CuckooFilter{
		m:      m,
		b:      b,
		f:      f,
		bucket: make([]uint8, m * b * (f >> 3)), // 分配的 bit 大小是 m * n * f/8
	}, nil
}

func main() {
	filter, err := New(1<<5, 2, 8)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	_ = filter.Add([]byte("Hello"))
	_ = filter.Add([]byte("World"))

	res1, err := filter.Contain([]byte("hello"))
	res2, err := filter.Contain([]byte("Hello"))

	if err != nil {
		fmt.Println("filter contain", err)
		return
	}

	fmt.Println("hello", res1, "Hello", res2)
}
