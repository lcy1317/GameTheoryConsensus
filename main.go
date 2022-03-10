package main

var Conf *Config

func main() {
	InitCheck()
	//printBoltDBBucket(Conf.ChainInfo.DBFile)
	//printBoltDBBucket("./storage/MessagePoolNode1.db")
	/**************************************************************/
	go testSendTransactions()                                // 开启协程，Sleep随机不断发送交易
	go TcpListenWrapper()                                    // 开启协程，监听收交易的端口
	go SendingPBFTCRequest(int64(Conf.ChainInfo.BlockSpeed)) // 开启协程，定时发送PBFT消息
	select {}
	/**************************************************************/
}
