package main

import (
	bolt "bbolt"
	"colorout"
	"fmt"
	"log"
	"strconv"
)

const InitBucketNameForBlockNumber = "blockNumber"
const InitBucketNameForChainInfo = "blockInfo"
const InitBucketNameForChainStageInfo = "stageInfo"
const InitBucketNameForBlockHash = "hashInfo"

func InitCheck() {
	ConfigCheck()           // 测试是否能够读取配置文件。
	BoltDBConnectionCheck() // 测试BoltDB数据库是否存在，没有则创建。
	BoltDBViewCheck()       // 测试BoltDB是否能正确读取
	BoltDBBlockNumberInit() // 初始化boltDB的BlockNumber
	PortListeningInit()     // 根据配置文件中的群主个数，从12000端口开端口监听
}
func PortListeningInit() {
	for i := 0; i < Conf.Basic.GroupNumber; i++ {
		portNumber := Conf.TcpInfo.PBFTBasePortStart + i
		address := Conf.TcpInfo.PBFTBaseAddress + ":" + strconv.Itoa(portNumber)
		go PBFTTcpListen(address)
	}
}

func BoltDBBlockNumberInit() {
	db, err := bolt.Open(Conf.ChainInfo.DBFile, 0600, nil)
	if err != nil {
		log.Fatal("数据库打开错误", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建
		_, err := tx.CreateBucketIfNotExists([]byte(InitBucketNameForBlockNumber))
		if err != nil {
			log.Fatalf(colorout.Red("创建区块数保存的Bucket出错:")+"%s", err.Error())
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(InitBucketNameForChainInfo))
		if err != nil {
			log.Fatalf(colorout.Red("创建区块链详细数据保存的Bucket出错:")+"%s", err.Error())
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(InitBucketNameForChainStageInfo))
		if err != nil {
			log.Fatalf(colorout.Red("创建区块链阶段数据保存的Bucket出错:")+"%s", err.Error())
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(InitBucketNameForBlockHash))
		if err != nil {
			log.Fatalf(colorout.Red("创建区块链区块哈希保存的Bucket出错:")+"%s", err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf(colorout.Red("创建区块数保存Bucket错误")+"%s", err.Error())
	}
	db.Close() // 及时关闭数据库

	// 读取测试，看是否有区块数，没有就需要放入键值
	_, blockNumber := BoltDBView(Conf.ChainInfo.DBFile, InitBucketNameForBlockNumber, []byte(InitBucketNameForBlockNumber))
	if blockNumber == nil {
		_ = BoltDBPut(Conf.ChainInfo.DBFile, InitBucketNameForBlockNumber, []byte(InitBucketNameForBlockNumber), IntSerialize(0))
	} else {
		fmt.Println(colorout.Yellow(strconv.Itoa(IntDeserialize(blockNumber))))
	}
	log.Print(colorout.Green("创建Bucket成功！读取或初始化区块数成功！"))
}

func BoltDBViewCheck() {
	// 测试BoltDB是否能正确读取
	if err := BoltDBReadTest(Conf.ChainInfo.DBFile); err != nil {
		log.Println(colorout.Yellow("链数据库读取错误:" + err.Error()))
	}
}
func BoltDBConnectionCheck() {
	// 测试BoltDB数据库是否存在，没有则创建。
	if err := IfBoltDBExist(Conf.ChainInfo.DBFile); err != nil {
		log.Println(colorout.Red("链数据库不存在:" + err.Error()))
		log.Println(colorout.Red("将自动创建链数据库"))
		if err := BoltDBCreate(Conf.ChainInfo.DBFile); err != nil {
			log.Println(colorout.Red("链数据库创建错误:" + err.Error()))
		}
	}
}

func ConfigCheck() {
	// 测试是否能够读取配置文件。
	if err := configInitial(); err != nil {
		log.Fatalln(colorout.Red("ReadInConfig error:" + err.Error()))
	}
}
