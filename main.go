package main

var Conf *Config
var nowBlockNumber int // 全局的当前BlockNumber记录
var nowStageNumber int // 全局的当前StageNumber记录
var nowHash string     // 一个暂时的解决办法，不知道为什么反序列化会在空交易的时候反序列化错误。

func main() {
	InitCheck()
	//BoltDBPutByte(Conf.ChainInfo.DBFile, []byte(InitBucketNameForBlockHash), []byte("测试"), []byte("你妈的"))
	//_, caoLe := BoltDBViewByte(Conf.ChainInfo.DBFile, []byte(InitBucketNameForBlockHash), []byte("测试"))
	//fmt.Println("caole ", string(caoLe), " caole")
	//select {}

	//printBoltDBBucket(Conf.ChainInfo.DBFile)
	//printBoltDBBucket("./storage/MessagePoolNode1.db")
	/**************************************************************/
	go testSendTransactions()                                // 开启协程，Sleep随机不断发送交易
	go TcpListenWrapper()                                    // 开启协程，监听收交易的端口
	go SendingPBFTCRequest(int64(Conf.ChainInfo.BlockSpeed)) // 开启协程，定时发送PBFT消息
	select {}
	/**************************************************************/
}
