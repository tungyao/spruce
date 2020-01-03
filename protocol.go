package spruce

import (
	"fmt"
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
