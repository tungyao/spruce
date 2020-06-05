package spruce

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/forgoer/openssl"
)

func getRandomInt(start, end int) int {
	randomMutex.Lock()
	<-time.After(1 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := start + r.Intn(end-start+1)
	randomMutex.Unlock()
	return n
}
func CreateLocalPWD() []byte {
	rands := []byte("0123456789abcdefghjiklmnopqrstuvwxyz#&_+=")
	outs := make([]byte, 0)
	for i := 0; i < 256; i++ {
		outs = append(outs, rands[getRandomInt(0, len(rands)-1)])
	}
	f, err := os.OpenFile("./pass.ewm", os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	err = binary.Write(f, binary.BigEndian, outs)
	if err != nil {
		log.Println(err)
	}
	return outs
}
func Encrypt(text []byte) ([]byte, error) {
	fs, err := os.OpenFile("./pass.ewm", os.O_RDONLY|os.O_CREATE, 666)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	n1, err := ioutil.ReadAll(fs)
	return openssl.AesECBDecrypt(text, n1, openssl.PKCS7_PADDING)
}

func Decrypt(s []byte) ([]byte, error) {
	fs, err := os.OpenFile("./pass.ewm", os.O_RDONLY|os.O_CREATE, 666)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	n1, err := ioutil.ReadAll(fs)
	if err != nil {
		return nil, err
	}
	return openssl.AesECBDecrypt(s, n1, openssl.PKCS7_PADDING)
}
func MD5(b []byte) []byte {
	m := md5.New()
	m.Write(b)
	return []byte(fmt.Sprintf("%x", m.Sum(nil)))
}
