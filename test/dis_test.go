package test

import (
	"../../spruce"
	"fmt"
	"golang.org/x/net/webdav"
	"net/http"
	"testing"
)

type User struct {
	Name string
	Pass string
}

var checkUser [2]User

func init() {
	checkUser[0] = User{
		Name: "abcd",
		Pass: "1121",
	}
	checkUser[1] = User{
		Name: "feng",
		Pass: "1121331",
	}
}
func TestTE2(t *testing.T) {
	fs := &webdav.Handler{
		FileSystem: webdav.Dir("D:\\phpstudy_pro\\WWW\\auth"),
		LockSystem: webdav.NewMemLS(),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var p bool
		for _, v := range checkUser {
			if v.Name != username && v.Pass != password {
				continue
			}
			p = true
		}
		if !p {
			http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}
		fs.ServeHTTP(w, req)
	})
	http.ListenAndServe(":4444", nil)
}
func TestBack(t *testing.T) {

	http.HandleFunc("/door", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.RemoteAddr)
		writer.Write([]byte(`{"run":"yes","ip":"101.132.172.196","port":"443","continue_time":"5"}`))
	})
	http.HandleFunc("/report", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.RawQuery)
	})
	http.ListenAndServe(":80", nil)
}
func TestDIS3(t *testing.T) {
	conf := make([]spruce.DCSConfig, 1)
	conf[0] = spruce.DCSConfig{
		Name:     "master",
		Ip:       "192.168.0.105:82",
		Weigh:    2,
		Password: "",
	}
	//conf[1] = spruce.DCSConfig{
	//	Name:     "node",
	//	Ip:       "192.168.0.114:82",
	//	Weigh:    1,
	//	Password: "",
	//}
	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.MEMORY,
		DCSConfigFile: "",
		DCSConfigs:    conf,
		Addr:          ":81",
		NowIP:         "192.168.0.105:82",
		KeepAlive:     false,
		IsBackup:      false,
	})
}
func TestUUID(t *testing.T) {
	d := []int{0, 1, 2, 3, 4, 5, 6}
	fmt.Println(d[:3])
}
func TestClock(t *testing.T) {
	n := spruce.NewClockTask(2)
	n.NewClock("test1", func() {
		fmt.Println("123")
	})
	n.NewClock("test2", func() {
		fmt.Println("456")
	})
	n.Start()
}
func TestDIS4(t *testing.T) {
	conf := make([]spruce.DCSConfig, 2)
	conf[0] = spruce.DCSConfig{
		Name:     "master",
		Ip:       "192.168.0.105:82",
		Weigh:    2,
		Password: "",
	}
	conf[1] = spruce.DCSConfig{
		Name:     "node",
		Ip:       "192.168.0.114:82",
		Weigh:    1,
		Password: "",
	}

	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.MEMORY,
		DCSConfigFile: "",
		DCSConfigs:    conf,
		Addr:          ":82",
		NowIP:         "192.168.0.114:82",
		KeepAlive:     false,
		IsBackup:      false,
	})
}
func TestSplit(t *testing.T) {
	t.Log(spruce.SplitString([]byte("set**hello**word"), []byte("**")))
}
