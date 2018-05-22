package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

func main() {

	// 读取配置
	filename := "./conf/logcollect.conf"
	err := LoadConf("ini", filename)
	if err != nil {
		fmt.Printf("load conf failed, err: %v\n", err)
		panic("load conf failed")
	}

	//初始化配置
	err = InitLogger()
	if err != nil {
		fmt.Printf("load conf failed, err: %v\n", err)
		panic(err)
	}

	logs.Debug("initial success!")
	logs.Debug("logs filename: %v", AppConfig)

	// 初始化tailf配置
	// Refer this document： https://github.com/hpcloud/tail
	err = InitTail(AppConfig.collectConf, AppConfig.ChanSize)
	if err != nil {
		logs.Error("init tail failed, error %v", err)

	}

	logs.Debug("init all success...")

	// 真正的业务逻辑处理函数
	err = ServerRun()
	if err != nil {
		logs.Error("server run failed ...")
	}

	logs.Info("program exited...")

}
