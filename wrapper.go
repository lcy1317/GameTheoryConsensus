package main

import (
	"GameTheoryConsensus/ConsensusUtils"
	"colorout"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func SendingNewBlock(duration int64) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(duration))

	//使用time.Ticker:
	blockNumber := 0
	go func() {
		for t := range ticker.C {
			blockNumber++
			fmt.Println(colorout.Cyan("每2s出块一个"), t, colorout.Cyan("当前区块："+strconv.Itoa(blockNumber)))
			message := "测试发送第" + strconv.Itoa(blockNumber) + "包"
			fmt.Println(colorout.Purple(message))
			ConsensusUtils.TcpDial([]byte(message), "127.0.0.1:1300"+strconv.Itoa(rand.Intn(Conf.Basic.GroupNumber)))
		}
	}()

	select {}
}
