package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
)

// 考虑到仅实现共识流程，所以没有其他的输入输出部分，仅博弈数字相关的交易

type Transaction struct {
	TXid      []byte  //储存该交易所引用的交易id
	Type      int     // 0 - 上报 1 - 数字解密
	Hash      string  //上报时候需要的字段 Hash(StageNumber + Number)
	Number    float64 // 解密时候的字段
	Signature []byte  //TODO:签名
	PubKey    []byte  //TODO:公钥
}

func (tx *Transaction) validating() bool {
	blockNum, stageNum, gameStop, revealStop := getBlockNumStageNumGameRevealStop()
	//color.Redln(blockNum, stageNum, (stageNum-1)*Conf.Basic.StageBlockNumber, gameStop, revealStop)
	if tx.Type == 0 { // 上报
		if blockNum > (stageNum-1)*Conf.Basic.StageBlockNumber && blockNum <= gameStop {
			return true
		} else {
			return false
		}
	}
	if tx.Type == 1 { // 解密
		if blockNum > gameStop && blockNum <= revealStop {
			return true
		} else {
			return false
		}
	}
	return true
}

func (tx *Transaction) getFloatNumString() string {
	return strconv.FormatFloat(tx.Number, 'f', Conf.Basic.NumberPrecision, 64)
}

func (tx *Transaction) getTXHash() string {
	//交易的哈希值
	var txHashes [][]byte
	bn, sn := getBlockNumandStageNum()
	txHashes = append(txHashes, tx.TXid)
	txHashes = append(txHashes, getNumberByte(tx.Type))
	txHashes = append(txHashes, getNumberByte(bn))
	txHashes = append(txHashes, getNumberByte(sn))
	return getSHA256Hash(bytes.Join(txHashes, []byte{}))
}

func (tx *Transaction) getHash() string {
	//上报时候需要的字段 Hash(StageNumber + Number)
	var txHashes [][]byte
	_, sn := getBlockNumandStageNum()
	txHashes = append(txHashes, IntSerialize(sn))
	txHashes = append(txHashes, []byte(tx.getFloatNumString()))
	return getSHA256Hash(bytes.Join(txHashes, []byte{}))
}

// 序列化交易
func (tx Transaction) TXSerialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 反序列化交易
func TXDeserialize(data []byte) Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}
