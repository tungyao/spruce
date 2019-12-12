package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)
import "../../spruce"

var randomMutex sync.Mutex
var str string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetRandomInt(start, end int) int {
	//访问加同步锁，是因为并发访问时容易因为时间种子相同而生成相同的随机数，那就狠不随机鸟！
	randomMutex.Lock()

	//利用定时器阻塞1纳秒，保证时间种子得以更改
	<-time.After(1 * time.Nanosecond)

	//根据时间纳秒（种子）生成随机数对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//得到[start,end]之间的随机数
	n := start + r.Intn(end-start+1)

	//释放同步锁，供其它协程调用
	randomMutex.Unlock()
	return n
}

func TestRing(t *testing.T) {
	spruce.CreateHash(1024)
	r := spruce.NewMessage()
	i := 100
	for i != 0 {
		i--
		r.Push(r.MSG("a", string(str[GetRandomInt(0, 52)]), "now", "23131"))
		if i == 0 {
			break
		}
	}
}
func TestJump(t *testing.T) {
	i := 100
end:
	for i != 0 {
		i--
		if i == 0 {
			break end
		}
	}
	fmt.Println(i)
}
func TestTime(t *testing.T) {
	fs, _ := os.Open("./config.yml")
	str, _ := ioutil.ReadAll(fs)
	fmt.Println(strings.Split(string(str), "\r\n"))
}
