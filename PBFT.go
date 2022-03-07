package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

const CRequest = "Request"
const CPrePrepare = "PrePrepare"
const CPrepare = "Prepare"
const CCommit = "Commit"

type PBFTMessage struct {
	NodeGroupId int // 所属的群组id，0默认为主节点，其余为非主节点
	NodeId      int // 第一阶段不考虑这个NodeId，这是群组中的节点编号，我们以NodeGroupId*10000+NodeId做他的编号好了
	BlockInfo   Block
	PBFTStage   string // 参考ConsensusUtils里面的定义，阶段是不同的。
	//临时消息池，消息摘要对应消息本体
	MessagePool map[string][]byte
	//存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
	PrePareConfirmCount map[string]map[string]bool
	//存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	CommitConfirmCount map[string]map[string]bool
	//该笔消息是否已进行Commit广播
	//isCommitBordcast map[string]bool
	//该笔消息是否已对客户端进行Reply
	//isReply map[string]bool
}

// 序列化PBFT信息
func (p PBFTMessage) PBFTSerialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(p)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 反序列化PBFT信息
func PBFTDeserialize(data []byte) PBFTMessage {
	var b PBFTMessage
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return b
}
