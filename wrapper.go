package main

import (
	"colorout"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var testBlockMessage = Block{
	Version:       "0.0 QwQ",
	Timestamp:     time.Now().Unix(),
	PrevBlockHash: "PrevBlockHash",
	Hash:          "Hash",
	StageHash:     "StageHash",
	BlockNum:      1,
	StageNum:      1,
}

var testPBFTmessage = &PBFTMessage{
	MajorNode: 0, // 定义主节点
	BlockInfo: testBlockMessage,
	PBFTStage: CRequest, // 发送给主节点的消息
}

var transactions []*Transaction
var a = 0

func testSendTransactions() {
	for {
		a++
		// TODO: 解密和上报时候信息不一样哦！
		conn, err := net.Dial("tcp", Conf.TcpInfo.ClientAddr)
		if err != nil {
			log.Println("初始发送链接错误", err)
			return
		}
		time.Sleep(time.Duration(Conf.ChainInfo.TransactionSpeed) * time.Microsecond) // 设置延时
		if Conf.PrintControl.MessageID {
			fmt.Println(colorout.Yellow("正在发送消息编号"+strconv.Itoa(a)) + "   ")
		}
		testTx := new(Transaction)
		testTx.TXid = IntSerialize(a)
		testTx.Type = rand.Intn(2)
		testTx.MyID = rand.Intn(Conf.Basic.InitNodesNumberinGroup)
		testTx.GroupID = rand.Intn(Conf.Basic.GroupNumber)
		testTx.getGeneralID()
		testTx.Number = float64(rand.Intn(math.MaxInt)) / float64(math.MaxInt) * 100
		testTx.Hash = testTx.getHash() //TODO：上报过程这个hash是自己算的，解密时候是公布数字
		testTx.Signature = []byte("Signature")
		testTx.PubKey = []byte("PubKey")
		_, err = conn.Write(testTx.TXSerialize())
		if err != nil {
			log.Fatal(err)
		}
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
		*tx = TXDeserialize(tcpMessage) // 反序列化出来我们的事务。
		txValidating, message := tx.validating()
		if !txValidating {
			if Conf.PrintControl.ReceiveTxMessage {
				fmt.Println(colorout.Red(addr+message+"  接受到非法事务消息:"), tx.printString())
			}
			continue
		}
		transactions = append(transactions, tx) // 将收到的消息放进全局变量transactions里
		if Conf.PrintControl.ReceiveTxMessage {
			fmt.Println(colorout.Cyan(addr+message+" 接受到事务消息:"), tx.printString())
		}
	}
}
func SendingPBFTCRequest(duration int64) {
	messageCheck.message = make(map[int]PBFT) // PBFT.go For message validating
	nodesGameStage = make(map[int]bool)       // stage.go For validation checking, avoiding several tx in one single stage
	nodesRevealStage = make(map[int]bool)     // stage.go For validation checking, avoiding several tx in one single stage
	stagePool.ifSort = false
	storeBlockInfo.check = false // 当前还没有存储过区块的信息（客户端收到任何一个Reply就存）
	//首先读出当前的区块编号
	_, blockNumberByte := BoltDBView(Conf.ChainInfo.DBFile, InitBucketNameForBlockNumber, []byte(InitBucketNameForBlockNumber))
	blockNumber := IntDeserialize(blockNumberByte)
	// 定时器，定时配置时间生成区块
	ticker := time.NewTicker(time.Millisecond * time.Duration(duration))
	go func() {
		for t := range ticker.C { // 每进入一次新建一个
			delayCal.printDelay()
			storeBlockInfo.check = false // 当前还没有存储过区块的信息（客户端收到任何一个Reply就存）
			validNodes()                 // Refresh Stage Check
			blockNumber++                // Add the BlockNumber
			stageCheck(blockNumber)      // 检查当前的stage是不是结束了要排序等
			// 在BoltDB中存入我们的blockNumber
			_ = BoltDBPut(Conf.ChainInfo.DBFile, InitBucketNameForBlockNumber, []byte(InitBucketNameForBlockNumber), IntSerialize(blockNumber))
			nowBlockNumber, nowStageNumber = getBlockNumandStageNum() // 获得全局的BlockNumber以及StageNumber，后面很多地方还没有重构。
			fmt.Println(t, colorout.Purple("当前交易总数："+strconv.Itoa(a)+"当前区块："+strconv.Itoa(nowBlockNumber)+" 当前阶段："+strconv.Itoa(nowStageNumber)+" 当前交易池交易数："+strconv.Itoa(len(transactions))))
			// 将生成的交易打包
			testPBFTmessage.BlockInfo.Transactions = transactions
			testPBFTmessage.BlockInfo.PrevBlockHash = nowHash
			//testPBFTmessage.BlockInfo.PrevBlockHash = getPrevBlockHash() // TODO: 一旦用这个函数直接爆炸，不知道为什么！！！！
			testPBFTmessage.BlockInfo.newBlockInfo(blockNumber) // 更新打包的区块内的区块信息
			if Conf.PrintControl.PBFTMessagePrint {
				fmt.Println(colorout.Yellow("准备发送PBFT消息，BlockNumber=" + testPBFTmessage.printString()))
			}
			messageCheck.message[blockNumber] = NewPBFT(*testPBFTmessage, Conf.Basic.GroupNumber) // TODO: 发送打包好的messagePool
			delayCal.initDelay()
			TcpDial(testPBFTmessage.PBFTSerialize(), "127.0.0.1:1300"+strconv.Itoa(testPBFTmessage.MajorNode))
			transactions = []*Transaction{} //清空当前消息池
		}
	}()
	select {}
}
