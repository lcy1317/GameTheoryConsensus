package main

import (
	"colorout"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var testBlockMessage = Block{
	Version:   "0.0 QwQ",
	Timestamp: time.Now().Unix(),
	Transactions: []*Transaction{
		{
			TXid:      []byte("第一个交易"),
			Type:      0,
			Hash:      []byte("Hash"),
			Number:    50.0,
			Signature: []byte("Signature"),
			PubKey:    []byte("PubKey"),
		},
	},

	PrevBlockHash: []byte("PrevBlockHash"),
	Hash:          []byte("Hash"),
	StageHash:     []byte("StageHash"),
	BlockNum:      1,
	StageNum:      1,
}

var testPBFTmessage = &PBFTMessage{
	MajorNode: 0, // 定义主节点
	BlockInfo: testBlockMessage,
	PBFTStage: CRequest, // 发送给主节点的消息
}

var transactions []*Transaction

func testSendTransactions() {
	a := 0
	for {
		a++
		time.Sleep(time.Duration(rand.Intn(600)+1000) * time.Millisecond) // 设置延时
		fmt.Println(colorout.Yellow("开始随机间隔发送交易信息，正在发送消息编号" + strconv.Itoa(a)))
		testTx := new(Transaction)
		testTx.TXid = IntSerialize(a)
		testTx.Type = rand.Intn(2)
		testTx.Hash = []byte("Hash")
		testTx.Number = 50.0
		testTx.Signature = []byte("Signature")
		testTx.PubKey = []byte("PubKey")
		TcpDial(testTx.TXSerialize(), Conf.TcpInfo.ClientAddr)
	}
}

// 监听交易的一个函数
func TcpListenWrapper() {
	// 该端口监听交易。
	// 创建一个TCP监听
	var addr = Conf.TcpInfo.ClientAddr
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	consoleMessage := "交易打包节点开启监听，地址：" + addr + "\n"
	fmt.Printf(colorout.Green(consoleMessage))

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		tcpMessage, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		tx := new(Transaction)
		*tx = TXDeserialize(tcpMessage)
		transactions = append(transactions, tx) // 将收到的消息放进全局变量transactions里
		fmt.Println(colorout.Cyan(addr+"接受到来自"+conn.RemoteAddr().String()+"的事务消息:"), tx.printString())
	}
}
func SendingPBFTCRequest(duration int64) {
	messageCheck = make(map[int]PBFT)
	//首先读出当前的区块编号
	_, blockNumberByte := BoltDBView(Conf.ChainInfo.DBFile, InitBucketName, []byte(InitBucketName))
	blockNumber := IntDeserialize(blockNumberByte)
	// 定时器，定时配置时间生成区块
	ticker := time.NewTicker(time.Millisecond * time.Duration(duration))
	go func() {
		for t := range ticker.C { // 每进入一次新建一个
			blockNumber++
			// 在BoltDB中存入我们的blockNumber
			_ = BoltDBPut(Conf.ChainInfo.DBFile, InitBucketName, []byte(InitBucketName), IntSerialize(blockNumber))
			fmt.Println(colorout.Cyan("每Ns出块一个"), t, colorout.Cyan("当前区块："+strconv.Itoa(blockNumber)))
			fmt.Println(colorout.Purple("当前交易池交易数：" + strconv.Itoa(len(transactions))))
			testPBFTmessage.BlockInfo.Transactions = transactions
			testPBFTmessage.BlockNumberSet(blockNumber) // 更新区块信息以及阶段信息。
			fmt.Println(colorout.Yellow("准备发送PBFT消息，BlockNumber=" + testPBFTmessage.printString()))
			messageCheck[blockNumber] = NewPBFT(*testPBFTmessage, Conf.Basic.GroupNumber) // TODO: 发送打包好的messagePool
			TcpDial(testPBFTmessage.PBFTSerialize(), "127.0.0.1:1300"+strconv.Itoa(testPBFTmessage.MajorNode))
			transactions = []*Transaction{} //清空当前消息池
		}
	}()
	select {}
}

// 该程序用来作为主节点打包交易，然后发送交易
// TODO: 当前模拟定时2s发送一个区块。缺少交易申报过程，单纯发送模拟交易。
func SendingNewBlock(duration int64) {
	//TODO: messageCheck的创建，从区块映射到PBFT结构体。
	messageCheck = make(map[int]PBFT)
	ticker := time.NewTicker(time.Millisecond * time.Duration(duration))

	//使用time.Ticker:
	_, blockNumberByte := BoltDBView(Conf.ChainInfo.DBFile, InitBucketName, []byte(InitBucketName))
	blockNumber := IntDeserialize(blockNumberByte)
	go func() {
		for t := range ticker.C {
			blockNumber++
			// 在BoltDB中存入我们的blockNumber
			_ = BoltDBPut(Conf.ChainInfo.DBFile, InitBucketName, []byte(InitBucketName), IntSerialize(blockNumber))

			fmt.Println(colorout.Cyan("每10s出块一个"), t, colorout.Cyan("当前区块："+strconv.Itoa(blockNumber)))
			message := "测试发送第" + strconv.Itoa(blockNumber) + "区块"
			fmt.Println(colorout.Purple(message))
			testPBFTmessage.BlockInfo.BlockNum = blockNumber                              // 设置当前的blockNumber值
			messageCheck[blockNumber] = NewPBFT(*testPBFTmessage, Conf.Basic.GroupNumber) // TODO: 发送消息前设置好messagePool
			TcpDial(testPBFTmessage.PBFTSerialize(), "127.0.0.1:1300"+strconv.Itoa(testPBFTmessage.MajorNode))
		}
	}()

	select {}
}

//// 该程序用来作为主节点打包交易，然后发送交易
//func SendingNewBlock(duration int64) {
//	ticker := time.NewTicker(time.Millisecond * time.Duration(duration))
//
//	//使用time.Ticker:
//	blockNumber := 0
//	go func() {
//		for t := range ticker.C {
//			blockNumber++
//			fmt.Println(colorout.Cyan("每10s出块一个"), t, colorout.Cyan("当前区块："+strconv.Itoa(blockNumber)))
//			message := "测试发送第" + strconv.Itoa(blockNumber) + "包"
//			fmt.Println(colorout.Purple(message))
//			TcpDial(testPBFTmessage.PBFTSerialize(), "127.0.0.1:1300"+strconv.Itoa(rand.Intn(Conf.Basic.GroupNumber)))
//		}
//	}()
//
//	select {}
//}
