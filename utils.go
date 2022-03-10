package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
	"strings"
)

// 将区块编号如Int类型的1157转变成byte类型”1157“，作为存入数据库的Key。
func getBlockNumberByte(blockNumber int) []byte {
	return []byte(strconv.Itoa(blockNumber))
}

// 从数据库中获取上一个区块的Hash值。作为当前区块的prevHash

func getPrevBlockHash() []byte {
	_, blockNumberByte := BoltDBView(Conf.ChainInfo.DBFile, InitBucketNameForBlockNumber, []byte(InitBucketNameForBlockNumber))
	blockNumber := IntDeserialize(blockNumberByte)
	if blockNumber == 0 {
		return []byte("Genesis Block!")
	}
	blockNumber--
	//TODO: 处理完commit消息的Cend之后去区块链数据库中存一下区块信息。
	return []byte("")
}

// 从配置文件中获取当前区块编号，阶段编号， GameTheoryStop， RevealStop
func getBlockNumStageNumGameRevealStop() (int, int, int, int) {
	blockNumber, stageNumber := getBlockNumandStageNum()
	gameTheoryStop := (stageNumber-1)*Conf.Basic.StageBlockNumber + Conf.Basic.GameTheoryStop
	revealStop := (stageNumber-1)*Conf.Basic.StageBlockNumber + Conf.Basic.RevealStop
	return blockNumber, stageNumber, gameTheoryStop, revealStop
}

// 获得区块数以及我们的阶段数
func getBlockNumandStageNum() (int, int) {
	_, blockNumberByte := BoltDBView(Conf.ChainInfo.DBFile, InitBucketNameForBlockNumber, []byte(InitBucketNameForBlockNumber))
	blockNumber := IntDeserialize(blockNumberByte)
	stageNumber := (blockNumber / Conf.Basic.StageBlockNumber) + 1
	return blockNumber, stageNumber
}

func getGroupNodeId(addr string) int {
	slice := strings.Split(addr, ":")
	portNumber, _ := strconv.Atoi(slice[1]) // TODO: 错误处理

	return portNumber % 100
}

// 序列化PBFT信息
func IntSerialize(number int) []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(number)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 反序列化PBFT信息
func IntDeserialize(data []byte) int {
	var b int
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return b
}
