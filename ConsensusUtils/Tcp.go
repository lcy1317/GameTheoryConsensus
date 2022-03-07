package ConsensusUtils

import (
	"colorout"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

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
		fmt.Println(colorout.Cyan(addr + "接受到Tcp消息" + ":" + string(tcpMessage)))
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
