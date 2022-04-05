package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type timeDelay struct {
	lock       sync.Mutex
	delayStart int64
	delayTot   int64
	delayCnt   int64
	allTot     float64
	allCnt     float64
}

var delayCal *timeDelay
var delayStart int64

func DelayInit() {
	delayCal = &timeDelay{
		delayStart: time.Now().UnixMicro(),
		delayCnt:   0,
		delayTot:   0,
		allCnt:     0,
		allTot:     0,
	}
}
func (d *timeDelay) saveDelay(nodeID int) {
	tmp := strconv.FormatInt(time.Now().UnixMicro()-delayStart, 10) + "\n"
	filePath := "./delay/" + strconv.Itoa(nodeID) + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		f, _ := os.Create(filePath)
		f.Close()
		file, _ = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(tmp)
	write.Flush()
}
func (d *timeDelay) initDelay() {
	d.lock.Lock()
	d.delayStart = time.Now().UnixMicro()
	d.lock.Unlock()
}

func (d *timeDelay) setDelay() {
	d.lock.Lock()
	d.delayTot = d.delayTot + (time.Now().UnixMicro() - d.delayStart)
	d.delayCnt = d.delayCnt + 1
	d.lock.Unlock()
}

func (d *timeDelay) printDelay() {
	d.lock.Lock()
	defer d.lock.Unlock()
	tot := float64(d.delayTot)
	cnt := float64(d.delayCnt)
	tmp := tot / cnt / 1000
	d.delayCnt = 0
	d.delayTot = 0
	if math.IsNaN(tmp) == false {
		d.allCnt = d.allCnt + 1
		d.allTot = d.allTot + tmp
		go d.saveDelayToFile(tmp, d.allTot/d.allCnt)
	}

	go fmt.Println("Consensus Delay:", tmp, "ms", " Average Consensus Delay:", d.allTot/d.allCnt, "ms")
}
func (d *timeDelay) saveDelayToFile(round float64, average float64) {
	filePath := "./delay/" + strconv.Itoa(Conf.Basic.GroupNumber) + "_" + strconv.Itoa(Conf.Basic.InitNodesNumberinGroup) + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		f, _ := os.Create(filePath)
		f.Close()
		file, _ = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(strconv.FormatFloat(round, 'f', Conf.Basic.NumberPrecision, 64) + " " + strconv.FormatFloat(average, 'f', Conf.Basic.NumberPrecision, 64) + "\n")
	write.Flush()
}
