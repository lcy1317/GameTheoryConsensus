package main

import (
	"GameTheoryConsensus/StartTest"
	"colorout"
	"log"
	"strconv"
)

func InitCheck() {
	ConfigCheck()           // 测试是否能够读取配置文件。
	BoltDBConnectionCheck() // 测试BoltDB数据库是否存在，没有则创建。
	BoltDBViewCheck()       // 测试BoltDB是否能正确读取
	PortListeningInit()     // 根据配置文件中的群主个数，从12000端口开端口监听
}
func PortListeningInit() {
	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		portNumber := Conf.TcpInfo.PBFTBasePortStart + i
		address := Conf.TcpInfo.PBFTBaseAddress + ":" + strconv.Itoa(portNumber)
		go TcpListen(address)
	}
}
func BoltDBViewCheck() {
	// 测试BoltDB是否能正确读取
	if err := StartTest.BoltDBReadTest(Conf.ChainInfo.DBFile); err != nil {
		log.Println(colorout.Yellow("链数据库读取错误:" + err.Error()))
	}
}
func BoltDBConnectionCheck() {
	// 测试BoltDB数据库是否存在，没有则创建。
	if err := StartTest.IfBoltDBExist(Conf.ChainInfo.DBFile); err != nil {
		log.Println(colorout.Red("链数据库不存在:" + err.Error()))
		log.Println(colorout.Red("将自动创建链数据库"))
		if err := StartTest.BoltDBCreate(Conf.ChainInfo.DBFile); err != nil {
			log.Println(colorout.Red("链数据库创建错误:" + err.Error()))
		}
	}
}

func ConfigCheck() {
	// 测试是否能够读取配置文件。
	if err := configInitial(); err != nil {
		log.Fatalln(colorout.Red("ReadInConfig error:" + err.Error()))
	}
}
