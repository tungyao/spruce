package test

import (
	"../../spruce"
	"encoding/json"
	"fmt"
	ap "git.yaop.ink/tungyao/awesome-pool"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	h := spruce.CreateHash(512)
	for i := 0; i < 100; i++ {
		h.Set([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)), 0)
	}
	//t.Log(string(h.Delete([]byte("9"))))
	////t.Log(string(h.Get([]byte("9"))))
	////h.Set([]byte(strconv.Itoa(9)), []byte(strconv.Itoa(9)),0)
	//for i := 0; i < 100; i++ {
	//	t.Log(string(h.Get([]byte(strconv.Itoa(i)))))
	//}
	//t.Log(string(h.Get([]byte(strconv.Itoa(88)))))
	h.Storage()
}
func TestReplace(t *testing.T) {
	a := "abcd\n\t\r摇动"
	fmt.Println(spruce.ReplaceTabCharacter([]byte(a)))
	spruce.ToBytes([]byte("hello world"))
	fmt.Println(json.Marshal(S{A: "asdasd"}))
}

type S struct {
	A string
}

func TestGHash(t *testing.T) {
	log.Println(85 >> 1)
	//hash := make(map[string]string, 10240)
	//for i := 0; i < 10240; i++ {
	//	hash[strconv.Itoa(i+rand.Int())] = strconv.Itoa(i)
	//}
	//for i := 0; i < 100000; i++ {
	//	t.Log(hash[strconv.Itoa(i)+"abc"])
	//}
}

func TestSet(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var Sum, miss int
	var Use time.Duration
	for i := 1; i <= runtime.NumCPU()*4; i++ {
		//	for i := 2; i <= 2; i++ {
		cnt := 1000 * 100
		if i > 9 {
			cnt = 1000 * 10
		}
		sum := i * cnt
		start := time.Now()
		miss = testQueuePutGoGet(t, i, cnt)

		end := time.Now()
		use := end.Sub(start)
		op := use / time.Duration(sum)
		fmt.Printf("%v, Grp: %3d, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
			runtime.Version(), i, sum, miss, use, op)
		Use += use
		Sum += sum
	}
	op := Use / time.Duration(Sum)
	fmt.Printf("%v, Grp: %3v, Times: %10d, miss:%6v, use: %12v, %8v/op\n",
		runtime.Version(), "Sum", Sum, 0, Use, op)
}
func testQueuePutGoGet(t *testing.T, grp, cnt int) int {
	p, _ := ap.NewPool(1, "127.0.0.1:6998")
	var wg sync.WaitGroup
	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			x := p.Get()
			x.Write(EntrySet(strconv.Itoa(g), strconv.Itoa(g), 0))
			x.Read()
			wg.Done()
		}(i)
	}
	wg.Add(grp)
	for i := 0; i < grp; i++ {
		go func(g int) {
			x := p.Get()
			x.Write(EntryGet(strconv.Itoa(g)))
			fmt.Println(string(x.Read()))
			wg.Done()
		}(i)
	}
	wg.Wait()
	return 0
}
func ParsingExpirationDate(tm interface{}) interface{} {
	switch tm.(type) {
	case []byte:
		if len(tm.([]byte)) > 2 {
			fmt.Println("input error")
		}
		var out int64 = 0
		out = int64(tm.([]byte)[1])
		out += int64(tm.([]byte)[0]) << 8
		return out
	case int:
		out := make([]byte, 2)
		out[1] = byte(tm.(int))
		out[0] = byte(tm.(int) >> 8)
		return out
	}
	return nil
}
func EntrySet(key, value string, ti int) []byte {
	out := make([]byte, 11)
	out[0] = 1
	out[1] = 2
	tm := ParsingExpirationDate(ti).([]byte)
	out[2] = tm[0]
	out[3] = tm[1]
	for _, v := range key {
		out = append(out, byte(v))
	}
	out = append(out, 0)
	for _, v := range value {
		out = append(out, byte(v))
	}
	return out
}
func EntryGet(key string) []byte {
	out := make([]byte, 11)
	out[0] = 2
	out[1] = 2
	for _, v := range key {
		out = append(out, byte(v))
	}
	return out
}
