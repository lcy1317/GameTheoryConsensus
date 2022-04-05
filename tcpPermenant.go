package main

import (
	"color"
	"fmt"
	"net"
	"strconv"
	"time"
)

type tcpAll struct {
	addr string
	conn net.Conn
}

var pbftTCP []*tcpAll

func initAllTcp() {
	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		portNumber := Conf.TcpInfo.PBFTBasePortStart + i
		address := Conf.TcpInfo.PBFTBaseAddress + ":" + strconv.Itoa(portNumber)
		tmpConn, err := net.Dial("tcp", address)
		if err != nil {
			color.Redln(err)
		}
		tmp := &tcpAll{
			addr: address,
			conn: tmpConn,
		}
		pbftTCP = append(pbftTCP, tmp)
		time.Sleep(50 * time.Millisecond)
	}
	for _, v := range pbftTCP {
		color.Redln(v.addr)
		v.conn.Close()
	}
}
func sendTcpMessage(content []byte, addr string) {
	for _, v := range pbftTCP {
		if v.addr == addr {
			_, err := v.conn.Write(content)
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
}
