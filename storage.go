package spruce

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// save memory data to local , default 60s run one ,but you can advance or delay ，
// save self's data
func localStorageFile() {
	allkey := balala.Get([]byte("*"))
	//fmt.Println(allkey)
	fs, err := os.OpenFile(string(MD5([]byte(time.Now().String())))+".spb", os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Println(err)
	}
	defer fs.Close()
	_, err = fs.Write(Encrypt([]byte(allkey)))
}

// this method will read [.spb file] form current folder,and load into memory
func localStorageFileRead() {
	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if path[len(path)-4:] == ".spb" {
			fs, err := os.Open(path)
			if err != nil {
				log.Println(err)
			}
			ad, _ := ioutil.ReadAll(fs)
			ad = Decrypt(ad)
			p := 0
			for k, v := range ad {
				if v == 0xFF {
					key := make([]byte, 0)
					val := make([]byte, 0)
					tim := make([]byte, 0)
					x := 0
					for c, j := range ad[p:k] {
						if j == 0xFE {
							if len(key) == 0 {
								key = ad[p:k][x:c]
								x = c
							}
							if len(val) == 0 {
								val = ad[p:k][x:c]
								x = c
							}
							if len(tim) == 0 {
								tim = ad[p:k][x:c]
								x = c
							}
						}
					}
				}
				p = k
			}
		}
		return nil
	})
}

// 通知所有的机器都进行备份
func remoteStoregeFile() {
	// 获取所有远程机器
	oAll := AllSlot
	// 饭后依次遍历 ，让其他电脑也同事备份
	for _, v := range oAll {
		go getRemote([]byte("*"), v.IP)
	}
	// 如果不出错，那么其他掉也会同时保存
}
