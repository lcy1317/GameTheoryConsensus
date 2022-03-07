package main

var Conf *Config

func main() {
	InitCheck()
	SendingNewBlock(int64(Conf.ChainInfo.BlockSpeed))
	select {}
}
