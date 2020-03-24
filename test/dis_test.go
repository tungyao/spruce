package test

import (
	"../../spruce"
	"fmt"
	"net/http"
	"testing"
)

func TestDIS2(t *testing.T) {
	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.FILE,
		DCSConfigFile: "./config.yml",
		Addr:          "127.0.0.1:9102",
		KeepAlive:     false,
		IsBackup:      false,
		NowIP:         "192.168.0.102:9102",
	})
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
	conf:=make([]spruce.DCSConfig,2)
	conf[0] = spruce.DCSConfig{
		Name:     "master",
		Ip:       "127.0.0.1:81",
		Weigh:    2,
		Password: "",
	}
	conf[1] = spruce.DCSConfig{
		Name:     "node",
		Ip:       "127.0.0.1:82",
		Weigh:    1,
		Password: "",
	}

	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.MEMORY,
		DCSConfigFile: "",
		DCSConfigs:    conf,
		Addr:          ":81",
		NowIP:         "127.0.0.1:81",
		KeepAlive:     false,
		IsBackup:      false,
	})
}
func TestDIS4(t *testing.T) {
	conf:=make([]spruce.DCSConfig,2)
	conf[0] = spruce.DCSConfig{
		Name:     "master",
		Ip:       "127.0.0.1:81",
		Weigh:    2,
		Password: "",
	}
	conf[1] = spruce.DCSConfig{
		Name:     "node",
		Ip:       "127.0.0.1:82",
		Weigh:    1,
		Password: "",
	}

	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.MEMORY,
		DCSConfigFile: "",
		DCSConfigs:    conf,
		Addr:          ":82",
		NowIP:         "127.0.0.1:82",
		KeepAlive:     false,
		IsBackup:      false,
	})
}
func TestSplit(t *testing.T) {
	t.Log(spruce.SplitString([]byte("set**hello**word"), []byte("**")))
}
