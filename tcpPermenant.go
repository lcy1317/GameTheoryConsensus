package main

import (
	"net"
	"strconv"
	"time"
)

var pbftTCP map[string]net.Conn

func initAllTcp() {
	pbftTCP = make(map[string]net.Conn)
	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		portNumber := Conf.TcpInfo.PBFTBasePortStart + i
		address := Conf.TcpInfo.PBFTBaseAddress + ":" + strconv.Itoa(portNumber)
		pbftTCP[address], _ = net.Dial("tcp", address)
		time.Sleep(50 * time.Millisecond)
	}
}
func sendTcpMessage(content []byte, addr string) {
	pbftTCP[addr].Write(content)
}
