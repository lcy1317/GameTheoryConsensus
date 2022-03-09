package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Version       string // 是不是String再说吧， 我准备用String，事实上可以简化一下不要这个字段
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	StageHash     []byte // 阶段哈希，由什么构成呢？ 例如由Miner+StageNum进行哈希。其实也可以用数字编号标识进行简化。
	BlockNum      int    // 创世纪块至今的区块数，模配置文件里的信息可得到对应的阶段内区块编号。
	StageNum      int    // 创世纪块至今的阶段编号，模配置文件里的信息可得到对应的阶段编号。
}

// 序列化区块信息
func (b Block) BlockSerialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 反序列化区块信息
func BlockDeserialize(data []byte) Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return b
}
