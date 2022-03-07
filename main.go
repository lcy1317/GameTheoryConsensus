package main

import "fmt"

var Conf *Config

func main() {

	InitCheck()
	fmt.Println(Conf.ChainInfo.BlockSpeed)
	SendingNewBlock(int64(Conf.ChainInfo.BlockSpeed))
	select {}
}
