package spruce

import (
	"fmt"
	"log"
	"os"
)

// save memory data to local , default 60s run one ,but you can advance or delay
func localStorageFile() {
	allkey := balala.Get([]byte("*"))
	fmt.Println(allkey)
	fs, err := os.OpenFile("./spruce.db", os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Println(err)
	}
	defer fs.Close()
	_, err = fs.Write(Encrypt([]byte(allkey)))
}

// 这个b方法怎么写哟 ，不球晓得，TMD
func remoteStoregeFile() {
	// 获取所有远程机器
	oAll := AllSlot
	// 饭后依次遍历 ，让其他电脑也同事备份
	for _, v := range oAll {
		go getRemote([]byte(""), v.IP)
	}
	// 如果不出错，那么其他掉也会同时保存
}
