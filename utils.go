package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
	"strings"
)

// 从地址中获得我们的主节点编号的函数，就是获得端口模100

func getGroupNodeId(addr string) int {
	slice := strings.Split(addr, ":")
	portNumber, _ := strconv.Atoi(slice[1]) // TODO: 错误处理

	return portNumber % 100
}

// 序列化PBFT信息
func IntSerialize(number int) []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(number)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 反序列化PBFT信息
func IntDeserialize(data []byte) int {
	var b int
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return b
}
