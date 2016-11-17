package main

import (
	"fmt"
	"os"

	"github.com/bykovme/goconfig"
)

type config struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	ServerIP string `json:"ip_address"`
}

const filename = "test.conf"
const folder = "etc/testconfig"

func main() {

	var conf config
	conf.Username = "TestUser"
	conf.Password = "testPassword"
	conf.ServerIP = "127.0.0.1"

	userHome, err := goconfig.GetUserHomePath()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	fullPath := userHome + folder + filename
	err = goconfig.SaveConfig(fullPath, conf)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	foundPath, err := goconfig.SearchForConfigFile(filename, folder, true, true, false)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	var loadedConf config
	err = goconfig.LoadConfig(foundPath, &loadedConf)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	if loadedConf.Username == conf.Username &&
		loadedConf.Password == conf.Password &&
		loadedConf.ServerIP == conf.ServerIP {
		fmt.Println("Config saved and loaded correctly")
	} else {
		fmt.Println("Config is not loaded")
		os.Exit(1)
	}

}
