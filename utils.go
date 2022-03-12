package main

import (
	"bytes"
	"color"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"strconv"
	"strings"
)

func getSHA256Hash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// 将区块编号如Int类型的1157转变成byte类型”1157“，作为存入数据库的Key。
func getNumberByte(num int) []byte {
	return []byte(strconv.Itoa(num))
}

// 从数据库中获取上一个区块的Hash值。作为当前区块的prevHash

func getPrevBlockHash() string { // TODO: 这段一旦加上一定出错，操了，我吐了，怎么就读取爆炸？？？？
	if nowBlockNumber == 0 {
		return "Genesis Block!"
	}
	// 处理完commit消息的Cend之后去区块链数据库中存一下区块信息。
	_, prevBlockByte := BoltDBViewByte(Conf.ChainInfo.DBFile, []byte(InitBucketNameForBlockHash), getNumberByte(nowBlockNumber-1))
	color.Redln(prevBlockByte)
	hash := ""
	if len(prevBlockByte) < 10 {
		hash = "Get PrevHash Error!"
	} else {
		hash = string(prevBlockByte)
	}

	return hash
}

// 从配置文件中获取当前区块编号，阶段编号， GameTheoryStop， RevealStop
func getStagesByBlockNum(blockNumber int) (int, int, int) { //以12 6 9 为例，当前区块25
	stageNumber := (blockNumber / Conf.Basic.StageBlockNumber) + 1                                // 3
	gameTheoryStop := (stageNumber-1)*Conf.Basic.StageBlockNumber + Conf.Basic.GameTheoryStop - 1 //25
	revealStop := (stageNumber-1)*Conf.Basic.StageBlockNumber + Conf.Basic.RevealStop - 1         //28
	return stageNumber, gameTheoryStop, revealStop
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

// 序列化Int
func IntSerialize(number int) []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(number)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 反序列化Int
func IntDeserialize(data []byte) int {
	var b int
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return b
}
