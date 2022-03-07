package main

import (
	"colorout"
	"fmt"
	"math/rand"
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
	NodeGroupId: 100,
	NodeId:      10100,
	BlockInfo:   testBlockMessage,
	PBFTStage:   CPrePrepare,
}

func SendingNewBlock(duration int64) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(duration))

	//使用time.Ticker:
	blockNumber := 0
	go func() {
		for t := range ticker.C {
			blockNumber++
			fmt.Println(colorout.Cyan("每10s出块一个"), t, colorout.Cyan("当前区块："+strconv.Itoa(blockNumber)))
			message := "测试发送第" + strconv.Itoa(blockNumber) + "包"
			fmt.Println(colorout.Purple(message))
			TcpDial(testPBFTmessage.PBFTSerialize(), "127.0.0.1:1300"+strconv.Itoa(rand.Intn(Conf.Basic.GroupNumber)))
		}
	}()

	select {}
}
