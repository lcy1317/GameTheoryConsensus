package main

import (
	"bytes"
	"colorout"
	"encoding/gob"
	"fmt"
	"log"
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
			fmt.Println(colorout.Purple(stagePool.stage.printString()))
			// 存储当前阶段的信息。
			BoltDBPutByte(Conf.ChainInfo.DBFile, []byte(InitBucketNameForChainStageInfo), getNumberByte(stagePool.stage.stageNumber), stagePool.stage.StageInfoSerialize())
			stagePool.ifSort = true
		}
	}
}

func (s *stageInfo) selectNode() { // TODO： 重新处理，当前不能有节点不报。
	sort.Sort(s) // 从小到大排序
	if len(s.gameNodes) < Conf.Basic.GroupNumber*Conf.Basic.InitNodesNumberinGroup {
		log.Println("节点为全部上报信息")
		return
	}
	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		// 找到中位数节点
		middle := i * Conf.Basic.InitNodesNumberinGroup
		middle = middle + Conf.Basic.InitNodesNumberinGroup/2
		s.upLayerNodes = append(s.upLayerNodes, s.gameNodes[middle].GeneralID)
	}
}
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
