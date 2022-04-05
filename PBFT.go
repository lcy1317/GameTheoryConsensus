package main

import (
	"bytes"
	"colorout"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"sync"
)

const CRequest = "Request"
const CPrePrepare = "PrePrepare"
const CPrepare = "Prepare"
const CCommit = "Commit"
const CEnded = "Ended"

type PBFTMessage struct {
	MajorNode   int // 主节点
	GroupNodeId int // 所属的群组id，0默认为主节点，其余为非主节点
	BlockInfo   Block
	PBFTStage   string // 参考ConsensusUtils里面的定义，阶段是不同的。
}

type PBFT struct {
	Message PBFTMessage // TODO:PBFT对应的Message
	//lock    sync.Mutex
	//存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
	prePareConfirmCount map[int]int
	//存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	commitConfirmCount map[int]int
	//该笔消息是否已进行Commit广播
	isCommitBordcast map[int]bool
	//该笔消息是否已对客户端进行Reply
	isReply map[int]bool
}

func NewPBFT(message PBFTMessage, nodesnum int) PBFT {
	var p PBFT
	//p.lock.Lock()
	p.Message = message
	p.prePareConfirmCount = make(map[int]int)
	p.commitConfirmCount = make(map[int]int)
	p.isCommitBordcast = make(map[int]bool)
	p.isReply = make(map[int]bool)
	// 初始化
	for i := 0; i < nodesnum; i++ {
		p.prePareConfirmCount[i] = 0
		p.commitConfirmCount[i] = 0
		p.isCommitBordcast[i] = false
		p.isReply[i] = false
	}
	//p.lock.Unlock()
	return p
}

type tmpPBFT struct {
	message map[int]PBFT
	lock    sync.Mutex
}

var messageCheck tmpPBFT

func PBFTTcpListen(addr string) {
	// 创建一个TCP监听
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	consoleMessage := "节点开启监听，地址：" + addr + "\n"
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
		var tcpPBFTMessage = PBFTDeserialize(tcpMessage) // 获取反序列化的消息

		// TODO: 当前收到了多少条消息，是否满足2f+1,满足之后才能够进行下一步。
		ifNext := false
		switch tcpPBFTMessage.PBFTStage {
		case CRequest:
			ifNext = tcpPBFTMessage.handleRequest(getGroupNodeId(addr))
		case CPrePrepare:
			ifNext = tcpPBFTMessage.handlePrePrepare(getGroupNodeId(addr))
		case CPrepare:
			ifNext = tcpPBFTMessage.handlePrepare(getGroupNodeId(addr))
		case CCommit:
			ifNext = tcpPBFTMessage.handleCommit(getGroupNodeId(addr))
			// TODO: CEnded
		}
		// PBFT消息打印
		// fmt.Println(colorout.Cyan(addr + "接受到来自" + conn.RemoteAddr().String() + "Tcp消息，当前阶段为:" + tcpPBFTMessage.PBFTStage))
		if ifNext {
			// 回参：下一步的指令，给哪些端口发 入参：当前地址,当前阶段
			portList, nextStage := ParseNextStep(addr, tcpPBFTMessage.PBFTStage)
			//fmt.Println(portList, nextStage)
			if nextStage != CEnded {
				var nextTcpMessage = tcpPBFTMessage
				nextTcpMessage.PBFTStage = nextStage              // 切换stage
				nextTcpMessage.GroupNodeId = getGroupNodeId(addr) // 根据地址确定当前发消息来的节点是那个，传下去消息的时候节点会变
				for _, port := range portList {
					TcpDial(nextTcpMessage.PBFTSerialize(), port)
					//_, err := TcpConn[addr][port].Write(nextTcpMessage.PBFTSerialize())
					//if err != nil {
					//	log.Fatal(err)
					//}
				}
			}
		}
	}
}

// 收到request不论如何向下传递
func (p PBFTMessage) handleRequest(nodeID int) bool {
	delayCal.initDelay()
	return true
}

// 收到preprepare肯定也必须发消息，所以这里不需要初始化任何东西，但是肯定要去发Prepare消息的。
func (p PBFTMessage) handlePrePrepare(nodeID int) bool {
	return true
}

// 收到prepare消息肯定要进行加加减减的操作
func (p PBFTMessage) handlePrepare(nodeID int) bool {
	blockNumber := p.BlockInfo.BlockNum
	// 互斥锁，保障同时只有一个线程
	messageCheck.lock.Lock()
	defer messageCheck.lock.Unlock()
	messageCheck.message[blockNumber].prePareConfirmCount[nodeID] += 1 // 增加一下prePare的数量
	if messageCheck.message[blockNumber].prePareConfirmCount[nodeID] < 2*(Conf.Basic.GroupNumber/3)-1 {
		// 说明消息不够
		return false
	} else {
		// 如果当前还没有进行Commit广播
		if messageCheck.message[blockNumber].isCommitBordcast[nodeID] == false {
			// 首先更改确认传递的状态
			messageCheck.message[blockNumber].isCommitBordcast[nodeID] = true
			return true
		}
		return false
	}
	return false
}
func (p PBFTMessage) handleCommit(nodeID int) bool {
	dbFileName := Conf.ChainInfo.NodeDBFile + "MessagePoolNode" + strconv.Itoa(nodeID) + ".db" //TODO: 收到足够的Commit消息之后再存入区块数据，中途用变量维护收到的消息数。
	blockNumber := p.BlockInfo.BlockNum
	// 互斥锁，保障同时只有一个线程
	messageCheck.lock.Lock()
	defer messageCheck.lock.Unlock()
	messageCheck.message[blockNumber].commitConfirmCount[nodeID] += 1 // 增加一下prePare的数量
	if messageCheck.message[blockNumber].commitConfirmCount[nodeID] < 2*(Conf.Basic.GroupNumber/3) {
		// 说明Commit消息数量不够 需要等待其他消息
		return false
	} else {
		// 如果当前还没有进行Commit广播
		if messageCheck.message[blockNumber].isReply[nodeID] == false {
			// 首先更改确认传递的状态
			messageCheck.message[blockNumber].isReply[nodeID] = true
			// 在这里收到足够的Commit，所以要在数据库中存下数据。
			BoltDBPutByte(dbFileName, []byte(strconv.Itoa(blockNumber)), []byte(strconv.Itoa(blockNumber)), p.PBFTSerialize())
			// 在主链上存储区块信息。
			delayCal.setDelay() // 计算延时
			// p.BlockInfo.storeBlockInfo()
			if Conf.PrintControl.Commit {
				fmt.Println(colorout.Purple("节点" + strconv.Itoa(nodeID) + "已完成Commit"))
			}

			return true
		}
		return false
	}
	return false
}

// 解析下一阶段发送的地址以及下一阶段的阶段名
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
