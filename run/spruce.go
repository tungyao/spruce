package main

import "../../spruce"

func main() {
	conf := make([]spruce.DCSConfig, 1)
	conf[0] = spruce.DCSConfig{
		Name:     "master",
		Ip:       "127.0.0.1:6999",
		Weigh:    2,
		Password: "",
	}
	// conf[1] = spruce.DCSConfig{
	//	Name:     "node",
	//	Ip:       "192.168.0.114:82",
	//	Weigh:    1,
	//	Password: "",
	// }
	spruce.StartSpruceDistributed(spruce.Config{
		ConfigType:    spruce.FILE,
		DCSConfigFile: "./config.yml",
		//DCSConfigFile: "/go/src/spruce/run/config.yml",
		Addr:      "127.0.0.1:6998",
		KeepAlive: true,
		IsBackup:  false,
		NowIP:     "127.0.0.1:6999",
	})
	// spruce.StartSpruceDistributed(spruce.Config{
	// 	ConfigType:    spruce.MEMORY,
	// 	DCSConfigFile: "",
	// 	DCSConfigs:    conf,
	// 	Addr:          ":6998",
	// 	NowIP:         "127.0.0.1:6999",
	// 	KeepAlive:     true,
	// 	IsBackup:      false,
	// })
}
