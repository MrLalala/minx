package utils

import (
	"../ziface"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// GlobalObj 配置结构
type GlobalObj struct {
	TCPServer     ziface.IServer //全局Server对象
	Host          string         //IP地址
	TCPPort       int            //IP端口
	Name          string         //服务器名字
	Version       string         //服务器版本
	MaxPacketSize uint32         //单个包的最大字节数
	MaxConn       int            //当前服务器的最大连接数
}

// Reload 载入配置文件
func (g *GlobalObj) Reload() {
	root := os.Getenv("ZINX_ROOT_PATH")
	if len(strings.TrimSpace(root)) == 0 {
		panic("error root path")
	}
	file := "conf/zinx.json"
	data, err := ioutil.ReadFile(root + "/" +file)
	if err != nil {
		panic(err)
	}

	// 配置文件转对象
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// GlobalObject 全局配置文件
var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:          "Zinx Main Server",
		Version:       "V0.4",
		TCPPort:       7777,
		Host:          "127.0.0.1",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	//GlobalObject.Reload()
}
