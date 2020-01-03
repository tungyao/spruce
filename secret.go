package spruce

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

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
func MD5(b []byte) []byte {
	m := md5.New()
	m.Write(b)
	return []byte(fmt.Sprintf("%x", m.Sum(nil)))
}
