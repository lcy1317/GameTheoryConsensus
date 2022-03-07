package main

import (
	"colorout"
	"fmt"
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
}

var testPBFTmessage = PBFTMessage{
	MajorNode: 0,   // 定义主节点
	NodeId:    100, // TODO:深入考虑这个字段，PBFT的消息只需要主节点参数应该就够了。
	BlockInfo: testBlockMessage,
	PBFTStage: CRequest, // 发送给主节点的消息
}

// 该程序用来作为主节点打包交易，然后发送交易
// TODO: 当前模拟定时2s发送一个区块。缺少交易申报过程，单纯发送模拟交易。
func SendingNewBlock(duration int64) {
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
			testPBFTmessage.BlockInfo.BlockNum = blockNumber // 设置当前的blockNumber值
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
