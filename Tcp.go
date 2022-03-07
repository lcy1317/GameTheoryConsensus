package main

import (
	"colorout"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

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

		//ifBroadCast = checkNowStage()
		// TODO: 当前收到了多少条消息，是否满足2f+1,满足之后才能够进行下一步。

		portList, nextStage := ParseNextStep(addr, tcpPBFTMessage.PBFTStage) // TODO: 回参：下一步的指令，给哪些端口发 入参：当前地址,当前阶段
		fmt.Println(colorout.Cyan(addr + "接受到来自" + conn.RemoteAddr().String() + "Tcp消息，当前阶段为:" + tcpPBFTMessage.PBFTStage))
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

func TcpListen(addr string) {
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
		fmt.Println(colorout.Cyan(addr + "接受到来自" + conn.RemoteAddr().String() + "Tcp消息" + ":" + string(tcpMessage)))
	}

}
func TcpDial(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Println("消息发送链接错误", err)
		return
	}

	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}

/*
// 测试Tcp的Code，写在main里直接测试的。
	for i := 1; i < 100; i++ {
		message := []byte("测试" + strconv.Itoa(i))
		ConsensusUtils.TcpDial(message, "127.0.0.1:1300"+strconv.Itoa(rand.Intn(Conf.Basic.GroupNumber)))
		time.Sleep(time.Second / 5)
	}
*/
