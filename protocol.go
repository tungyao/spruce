package spruce

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

type protocol struct {
}

//传输协议的制定 0-> 表示操作码 1 状态码 2-3 过期时间 4-10 暂时没想到 11-N 语句
// 操作码  0 delete | 1 set | 2 get | 3 status
// 状态码 0 操作失败 1操作成功 2 操作中
// set key 和 value 以0作为分割
// 新加 2-3 存放过期时间 最大 0xFFFF or 65535
func SplitKeyValue(b []byte) ([]byte, []byte) {
	for k, v := range b {
		if v == 0 {
			return b[:k], b[k+1:]
		}
	}
	return nil, nil
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

func SendSetMessage(lang []byte, ti []byte) []byte {
	var out = make([]byte, len(lang)+11)
	out[0] = 1
	out[1] = 2
	out[2] = ti[0]
	out[3] = ti[1]
	for i := 11; i < len(out); i++ {
		out[i] = lang[i-11]
	}
	return out
}
func SendGetMessage(key []byte) []byte {
	out := make([]byte, len(key)+11)
	out[0] = 2
	out[0] = 2
	for i := 11; i < len(out); i++ {
		out[i] = key[i-11]
	}
	return out
}
func SendDeleteMessage() {

}
func SendStatusMessage(b bool) []byte {
	if b {
		return []byte{0x3, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	} else {
		return []byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	}
}
func CreateLocalPWD() []byte {
	rands := []byte("0123456789abcdefghjiklmnopqrstuvwxyz#&_+=")
	outs := make([]byte, 0)
	for i := 0; i < rand.Int()*256; i++ {
		outs = append(outs, rands[getRandomInt(0, len(rands)-1)])
	}
	f, err := os.OpenFile("./pass.ewm", os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	_, err = f.Write(outs)
	if err != nil {
		log.Println(err)
	}
	return outs
}
func Encrypt(string2 []byte) []byte {
	fs, err := os.OpenFile("./pass.ewm", os.O_RDONLY|os.O_CREATE, 666)
	if err != nil {
		log.Panicln(err)
		return nil
	}
	defer fs.Close()
	n1, err := ioutil.ReadAll(fs)
	if len(n1) == 0 || err != nil {
		log.Panicln("There is no password in the document", err)
	}

	for k, v := range string2 {
		string2[k] = v + n1[len(n1)-1%int(v)]
	}
	return string2
}
func Decrypt(s []byte) []byte {
	fs, err := os.OpenFile("./pass.ewm", os.O_RDONLY|os.O_CREATE, 666)
	if err != nil {
		log.Println(err)
	}
	defer fs.Close()
	n1, err := ioutil.ReadAll(fs)
	for k, v := range s {
		s[k] = v - n1[len(n1)-1%int(v)]
	}
	return s
}
