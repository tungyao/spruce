package spruce

import (
	"sync"
)

// 每一次触发都会使指针向后面移动
type MessageRing struct {
	Msg [2048]*msg `用来存放没有成功推送出去的消息,消息都会及时推送出去`
	All int
	mux sync.Mutex
	P   int
	Tag bool
}
type Ring struct {
	Rings []*MessageRing
	N     int
	mux   sync.Mutex
}
type msg struct {
	From    string
	To      string
	Time    string
	Content string
	Repeat  int `重复发送次数 ，超过3次就删除`
	Loop    int `在环上是第几圈`
}

// TODO Message
// 要 new 出一个环形的msg
func NewMessage() *Ring {
	nMessageRing = new(MessageRing)
	pMsg = make(chan msg, 1000)
	return &Ring{
		Rings: make([]*MessageRing, 1000),
	}
}
func (m *Ring) MSG(from, to, times, content string) msg {
	return msg{
		From:    from,
		To:      to,
		Time:    times,
		Content: content,
		Repeat:  0,
		Loop:    0,
	}
}
func (m *Ring) Push(sg msg) {
	if m.N == 0 {
		m.Rings[0] = &MessageRing{
			Msg: [2048]*msg{},
			All: 0,
			mux: sync.Mutex{},
			P:   0,
			Tag: false,
		}
		m.N++
	}
	for i := 0; i < m.N; i++ {
		m.mux.Lock()
		if m.Rings[i].Tag {
			continue
		}
		p := hashString([]rune(sg.To)) % len(m.Rings[i].Msg)
		m.Rings[i].mux.Lock()
		if m.Rings[i].Msg[p] == nil {
			m.Rings[i].Msg[p] = &sg
		} else {

		}
		m.Rings[i].mux.Unlock()
		m.mux.Unlock()
		break
	}

}
func (m *Ring) Pull() {

}

// 用来表示当前的在环的那个位置
func (m *MessageRing) Pointer() int {
	return m.P
}

// 继续循环或者开始循环
func (m *MessageRing) Loop() {

}
func hashString(str []rune) int {
	var (
		key        = str
		seed1 uint = 0x7FED7FED
		seed2 uint = 0xEEEEEEEE
		ch    uint = 0
	)
	for _, v := range key {
		ch = uint(v)
		seed1 = CRY[(2<<2)+ch] ^ (seed1 + seed2)
		seed2 = ch + seed1 + seed2 + (seed2 << 5) + 3
	}
	return int(seed1)
}
