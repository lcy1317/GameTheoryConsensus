package main

import (
	"bytes"
	"colorout"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
)

// 用来存事务中的节点编号
type nodesInfo struct {
	GeneralID int
	GroupID   int
	number    float64
}

type stageInfo struct {
	stageNumber  int
	upLayerNodes []int       //上层节点的list，长度应当与配置中的群组数量保持一致
	gameNodes    []nodesInfo // 参与游戏的节点
}

type tmpStageInfo struct {
	stage  stageInfo
	ifSort bool
	lock   sync.Mutex
}

var stagePool tmpStageInfo

func stageCheck(blockNumber int) {
	stagePool.lock.Lock()
	defer stagePool.lock.Unlock()
	stageNumber, gameStop, revealStop := getStagesByBlockNum(blockNumber)
	stageStart := Conf.Basic.StageBlockNumber * (stageNumber - 1) //该阶段开始时间
	stageEnd := Conf.Basic.StageBlockNumber * stageNumber
	if blockNumber >= stageStart && blockNumber <= gameStop {
		stagePool.ifSort = false
		stagePool.stage.stageNumber = stageNumber // 更新阶段编号
		stagePool.stage.gameNodes = []nodesInfo{} // 清空
		stagePool.stage.upLayerNodes = []int{}    // 清空
	}
	if blockNumber > revealStop && blockNumber < stageEnd {
		if stagePool.ifSort == false {
			stagePool.stage.selectNode()
			fmt.Println(colorout.Blue(stagePool.stage.printString()))
			// 存储当前阶段的信息。
			BoltDBPutByte(Conf.ChainInfo.DBFile, []byte(InitBucketNameForChainStageInfo), getNumberByte(stagePool.stage.stageNumber), stagePool.stage.StageInfoSerialize())
			stagePool.ifSort = true
		}
	}
}
func (s *stageInfo) selectNode() {
	sort.Sort(s) // 从小到大排序
	if len(s.gameNodes) == 0 {
		log.Println("没有节点上报信息")
		return
	}
	var nodes map[int]int // 保存所有群组中节点个数的临时map
	nodes = make(map[int]int)
	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		nodes[i] = 0
	}

	for i := 0; i < len(s.gameNodes); i++ {
		// 遍历所有的节点，找到每个群组的个数应该是多少。
		nodes[s.gameNodes[i].GroupID]++
	}
	nowNodes := 0

	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		// 找到平均数节点
		if nodes[i] == 0 { // 当前群组节点为空
			log.Println("存在无信息群组，进入下一轮")
			return
		}
		tot, k := 0.0, 0.0
		for j := 0; j < nodes[i]; j++ {
			tot += s.gameNodes[j].number
			k += 1
		}
		tot = tot / k
		ansk := nowNodes
		maxd := 200.0
		for j := 0; j < nodes[i]; j++ {
			if math.Abs(s.gameNodes[ansk].number-s.gameNodes[nowNodes+j].number) < maxd {
				maxd = math.Abs(s.gameNodes[ansk].number - s.gameNodes[nowNodes+j].number)
				ansk = j + nowNodes
			}
		}
		nowNodes = nowNodes + nodes[i]
		s.upLayerNodes = append(s.upLayerNodes, s.gameNodes[ansk].GeneralID)
	}
}

//func (s *stageInfo) selectNode() {
//	sort.Sort(s) // 从小到大排序
//	if len(s.gameNodes) == 0 {
//		log.Println("没有节点上报信息")
//		return
//	}
//	var nodes map[int]int // 保存所有群组中节点个数的临时map
//	nodes = make(map[int]int)
//	for i := 0; i < Conf.Basic.GroupNumber; i++ {
//		nodes[i] = 0
//	}
//
//	for i := 0; i < len(s.gameNodes); i++ {
//		// 遍历所有的节点，找到每个群组的个数应该是多少。
//		nodes[s.gameNodes[i].GroupID]++
//	}
//	nowNodes := 0
//
//	for i := 0; i < Conf.Basic.GroupNumber; i++ {
//		// 找到中位数节点
//		if nodes[i] == 0 { // 当前群组节点为空
//			log.Println("存在无信息群组，进入下一轮")
//			return
//		}
//		middle := (nowNodes + nowNodes + nodes[i]) / 2
//		nowNodes = nowNodes + nodes[i]
//		s.upLayerNodes = append(s.upLayerNodes, s.gameNodes[middle].GeneralID)
//	}
//}
func (s stageInfo) Swap(i, j int) { s.gameNodes[i], s.gameNodes[j] = s.gameNodes[j], s.gameNodes[i] }
func (s stageInfo) Len() int      { return len(s.gameNodes) }
func (s stageInfo) Less(i, j int) bool {
	if s.gameNodes[i].GroupID < s.gameNodes[j].GroupID {
		return true
	} else if s.gameNodes[i].GroupID == s.gameNodes[j].GroupID && s.gameNodes[i].number < s.gameNodes[j].number {
		return true
	}
	return false
}

var nodesGameStage map[int]bool
var nodesRevealStage map[int]bool

// 该函数用来从配置中生成并返回一个list，代表有效的节点。避免重复，大概如此。
func validNodes() {
	for groupID := 0; groupID < Conf.Basic.GroupNumber; groupID++ {
		for myID := 0; myID < Conf.Basic.InitNodesNumberinGroup; myID++ {
			generalID := myID*100 + groupID
			nodesGameStage[generalID] = false
			nodesRevealStage[generalID] = false
		}
	}
	return
}

// 序列化区块信息
func (s stageInfo) StageInfoSerialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(s)
	if err != nil {
		return []byte("")
	}
	return encoded.Bytes()
}

// 反序列化区块信息
func StageInfoDeserialize(data []byte) stageInfo {
	var b stageInfo
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return b
}
