package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
)

const CRequest = "Request"
const CPrePrepare = "PrePrepare"
const CPrepare = "Prepare"
const CCommit = "Commit"
const CEnded = "Ended"

type PBFTMessage struct {
	MajorNode   int // 主节点
	GroupNodeId int // 所属的群组id，0默认为主节点，其余为非主节点
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

func ParseNextStep(addr string, PBFTstage string) ([]string, string) {
	// 获得下一个阶段应该访问的地址。
	// 当前地址是127.0.0.1:13000 形式
	var portList []string
	for i := 0; i < Conf.Basic.GroupNumber; i++ { // 循环所有地址
		portNumber := Conf.TcpInfo.PBFTBasePortStart + i
		listenAddr := Conf.TcpInfo.PBFTBaseAddress + ":" + strconv.Itoa(portNumber) // 首先构造对应的监听地址
		if listenAddr == addr {
			continue
		}
		portList = append(portList, listenAddr) // 发给除了自己之外的全部其他节点
	}
	// 获得下一个应进行的阶段
	var stage = CRequest
	switch PBFTstage {
	case CRequest:
		stage = CPrePrepare
	case CPrePrepare:
		stage = CPrepare
	case CPrepare:
		stage = CCommit
	case CCommit:
		stage = CEnded
	default:
		stage = CEnded
	}

	return portList, stage
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
