package main

import (
	bolt "bbolt"
	"bytes"
	"colorout"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const CRequest = "Request"
const CPrePrepare = "PrePrepare"
const CPrepare = "Prepare"
const CCommit = "Commit"
const CEnded = "Ended"
const ifPrepare = "ifPrepare"
const ifCommit = "ifCommit"
const prePrepareNum = "prePrepareNum"
const prepareNum = "prepareNum"
const commitNum = "commitNum"
const blockInformation = "blockInfo"

type PBFTMessage struct {
	MajorNode   int // 主节点
	GroupNodeId int // 所属的群组id，0默认为主节点，其余为非主节点
	NodeId      int // 第一阶段不考虑这个NodeId，这是群组中的节点编号，我们以NodeGroupId*10000+NodeId做他的编号好了
	BlockInfo   Block
	PBFTStage   string // 参考ConsensusUtils里面的定义，阶段是不同的。
	lock        sync.Mutex
	//临时消息池，消息摘要对应消息本体
	//MessagePool map[string][]byte
	//存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
	//PrePareConfirmCount map[string]map[string]bool
	//存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	//CommitConfirmCount map[string]map[string]bool
}

var mutexListen sync.Mutex
var mutexTcp sync.RWMutex
var CommitNumbers = 0

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

		tcpPBFTMessage.handleStore() // TODO: 不太应该在这全部创建，单纯为了简单。 需要修改，如果有空
		//ifBroadCast = checkNowStage()
		// TODO: 当前收到了多少条消息，是否满足2f+1,满足之后才能够进行下一步。
		ifNext := false
		switch tcpPBFTMessage.PBFTStage {
		case CRequest:
			ifNext = tcpPBFTMessage.handleRequest(getGroupNodeId(addr))
		case CPrePrepare:
			ifNext = tcpPBFTMessage.handlePrePrepare(getGroupNodeId(addr))
		case CPrepare:
			ifNext = tcpPBFTMessage.handlePrepare(getGroupNodeId(addr))
			fmt.Println(colorout.Yellow("收到Prepare消息"), ifNext)
		case CCommit:
			ifNext = tcpPBFTMessage.handleCommit(getGroupNodeId(addr))
			fmt.Println(colorout.Yellow("收到Commit消息"), ifNext)
			CommitNumbers++
			if CommitNumbers >= 27 {
				println("Fucking Commit", CommitNumbers)
			}
			// TODO: CEnded
		}
		fmt.Println(colorout.Cyan(addr + "接受到来自" + conn.RemoteAddr().String() + "Tcp消息，当前阶段为:" + tcpPBFTMessage.PBFTStage))
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
				}
			}
		}
	}
}
func (p PBFTMessage) handleStore() {
	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		dbFileName := Conf.ChainInfo.NodeDBFile + "MessagePoolNode" + strconv.Itoa(i) + ".db"
		db, err := bolt.Open(dbFileName, 0600, nil)
		if err != nil {
			log.Println(colorout.Red("节点数据库打开出错:")+"%s", err.Error())
		}

		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(IntSerialize(p.BlockInfo.BlockNum))
			if err != nil {
				log.Fatalf(colorout.Red("创建Bucket出错:")+"%s", err.Error())
				return err
			}
			if err = bucket.Put([]byte(ifPrepare), IntSerialize(0)); err != nil {
				log.Fatalf(colorout.Red("节点Bucket存放数据错误:")+"%s", err.Error())
				return err
			}
			if err = bucket.Put([]byte(ifCommit), IntSerialize(0)); err != nil {
				log.Fatalf(colorout.Red("节点Bucket存放数据错误:")+"%s", err.Error())
				return err
			}
			if err = bucket.Put([]byte(prepareNum), IntSerialize(0)); err != nil {
				log.Fatalf(colorout.Red("节点Bucket存放数据错误:")+"%s", err.Error())
				return err
			}
			if err = bucket.Put([]byte(commitNum), IntSerialize(0)); err != nil {
				log.Fatalf(colorout.Red("节点Bucket存放数据错误:")+"%s", err.Error())
				return err
			}
			// 放入区块体的数据
			if err = bucket.Put([]byte(blockInformation), p.PBFTSerialize()); err != nil {
				log.Fatalf(colorout.Red("节点Bucket存放区块错误:")+"%s", err.Error())
				return err
			}
			return nil
		})
		if err != nil {
			log.Fatalf(colorout.Red("更新数据库错误")+"%s", err.Error())
		}
		_ = db.Close() // 及时关闭数据库
	}
}
func (p PBFTMessage) handleRequest(nodeID int) bool {

	dbFileName := Conf.ChainInfo.NodeDBFile + "MessagePoolNode" + strconv.Itoa(nodeID) + ".db"
	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		log.Println(colorout.Red("主节点数据库打开出错:")+"%s", err.Error())
	}
	defer db.Close() // 及时关闭数据库

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(IntSerialize(p.BlockInfo.BlockNum))
		if err != nil {
			log.Fatalf(colorout.Red("创建Bucket出错:")+"%s", err.Error())
			return err
		}
		// 初始化放入了一些信息就是是否确认，以及确认数等。
		if err = bucket.Put([]byte(ifPrepare), IntSerialize(0)); err != nil {
			log.Fatalf(colorout.Red("主节点Bucket存放数据错误:")+"%s", err.Error())
			return err
		}
		if err = bucket.Put([]byte(ifCommit), IntSerialize(0)); err != nil {
			log.Fatalf(colorout.Red("主节点Bucket存放数据错误:")+"%s", err.Error())
			return err
		}
		if err = bucket.Put([]byte(prepareNum), IntSerialize(0)); err != nil {
			log.Fatalf(colorout.Red("主节点Bucket存放数据错误:")+"%s", err.Error())
			return err
		}
		if err = bucket.Put([]byte(commitNum), IntSerialize(0)); err != nil {
			log.Fatalf(colorout.Red("主节点Bucket存放数据错误:")+"%s", err.Error())
			return err
		}
		// 放入区块体的数据
		if err = bucket.Put([]byte(blockInformation), p.PBFTSerialize()); err != nil {
			log.Fatalf(colorout.Red("主节点Bucket存放区块错误:")+"%s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf(colorout.Red("更新数据库错误")+"%s", err.Error())
	}
	return true
}

// 收到preprepare肯定也必须发消息，所以这里不需要初始化任何东西，但是肯定要去发Prepare消息的。
func (p PBFTMessage) handlePrePrepare(nodeID int) bool {
	return true
}
func (p PBFTMessage) handlePrepare(nodeID int) bool {

	dbFileName := Conf.ChainInfo.NodeDBFile + "MessagePoolNode" + strconv.Itoa(nodeID) + ".db"

	err, prepareNumsByte := BoltDBByte("View", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(prepareNum), []byte(""))

	// 包括自己要收到2f条，我这边自己不发所以是2f-1条
	numbersNeeded := 2*(Conf.Basic.GroupNumber/3) - 1

	// Whatever 先更新消息总数
	number := IntDeserialize(prepareNumsByte) + 1
	err, _ = BoltDBByte("Put", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(prepareNum), IntSerialize(number))
	if err != nil {
		log.Fatalf(colorout.Red("更新数据库PrePrepare消息数量错误")+"%s", err.Error())
	}
	fmt.Println("PrePare Stage", number, numbersNeeded)

	if IntDeserialize(prepareNumsByte)+1 < numbersNeeded {
		// 数量不够，那就不能下一步传播，但是要更新数据库里的数量
		return false
	} else {
		// 如果数量达到要求
		err, ifPrepareByte := BoltDBByte("View", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(ifPrepare), []byte(""))
		if err != nil {
			log.Fatalf(colorout.Red("读取数据库是否已经传递Prepare消息错误")+"%s", err.Error())
		}
		fmt.Println("IfPrepare:", IntDeserialize(ifPrepareByte), "blockNumber = ", p.BlockInfo.BlockNum, time.Nanosecond)
		if IntDeserialize(ifPrepareByte) == 0 {
			// 如果是0代表还没有开始传播commit消息，因此可以return true进入下一阶段
			// 同时要更新成 1
			err, _ = BoltDBByte("Put", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(ifPrepare), IntSerialize(1))
			if err != nil {
				log.Fatalf(colorout.Red("更新数据库是否已经传递Prepare消息错误")+"%s", err.Error())
			}

			return true
		}
		return false
	}
	return false
}
func (p PBFTMessage) handleCommit(nodeID int) bool {

	dbFileName := Conf.ChainInfo.NodeDBFile + "MessagePoolNode" + strconv.Itoa(nodeID) + ".db"

	err, commitNumsByte := BoltDBByte("View", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(commitNum), []byte(""))

	numbersNeeded := 2 * (Conf.Basic.GroupNumber / 3)

	number := IntDeserialize(commitNumsByte) + 1
	err, _ = BoltDBByte("Put", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(commitNum), IntSerialize(number))
	if err != nil {
		log.Fatalf(colorout.Red("更新数据库Commit消息数量错误")+"%s", err.Error())
	}
	fmt.Println(number, numbersNeeded)
	if IntDeserialize(commitNumsByte)+1 < numbersNeeded {
		// 数量不够，那就不能下一步传播，但是要更新数据库里的数量
		return false
	} else {
		// 如果数量达到要求
		err, ifCommitByte := BoltDBByte("View", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(ifCommit), []byte(""))
		if err != nil {
			log.Fatalf(colorout.Red("读取数据库是否已经传递Commit消息错误")+"%s", err.Error())
		}
		if IntDeserialize(ifCommitByte) == 0 {
			// 如果是0代表还没有开始传播commit消息，因此可以return true进入下一阶段
			// 同时要更新成 1
			err, _ = BoltDBByte("Put", dbFileName, IntSerialize(p.BlockInfo.BlockNum), []byte(ifCommit), IntSerialize(1))
			if err != nil {
				log.Fatalf(colorout.Red("更新数据库是否已经传递Commit消息错误")+"%s", err.Error())
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
