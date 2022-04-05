package main

import (
	"math/rand"
	"net"
	"time"
)

var Conf *Config
var nowBlockNumber int // 全局的当前BlockNumber记录
var nowStageNumber int // 全局的当前StageNumber记录
var nowHash string     // TODO: 一个暂时的解决办法，不知道为什么反序列化会在空交易的时候反序列化错误。
var TcpConn map[string]map[string]net.Conn

func main() {
	DelayInit()
	rand.Seed(int64(time.Now().Nanosecond())) // 随机数种子
	InitCheck()                               // 初始化检查
	go TcpListenWrapper()                     // 开启协程，监听收交易的端口
	time.Sleep(3 * time.Second)
	go testSendTransactions()                                // 开启协程，Sleep随机不断发送交易
	go SendingPBFTCRequest(int64(Conf.ChainInfo.BlockSpeed)) // 开启协程，定时发送PBFT消息
	select {}
}
