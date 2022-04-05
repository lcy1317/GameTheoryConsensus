package main

import (
	"color"
	"fmt"
	"sync"
	"time"
)

type timeDelay struct {
	lock       sync.Mutex
	delayStart int64
	delayTot   int64
	delayCnt   int64
}

var delayCal timeDelay

func DelayInit() {
	delayCal.lock.Lock()
	defer delayCal.lock.Unlock()
	delayCal.delayStart = time.Now().UnixMicro()
	delayCal.delayCnt = 0
	delayCal.delayTot = 0
}
func (d timeDelay) initDelay() {
	d.lock.Lock()
	d.delayStart = time.Now().UnixMicro()
	d.lock.Unlock()
}

func (d timeDelay) setDelay() {
	d.lock.Lock()
	d.delayTot = d.delayTot + (time.Now().UnixMicro() - d.delayStart)
	d.delayCnt = d.delayCnt + 1
	color.Redln(d)
	d.lock.Unlock()
}

func (d timeDelay) printDelay() {
	tot := float64(d.delayTot)
	cnt := float64(d.delayCnt)
	tmp := tot / cnt / 1000
	color.Redln(d)
	fmt.Println("Consensus Delay:", tmp, "ms")
}
