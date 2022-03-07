# GameTheoryConsensus

## 更新日志
- 5th Commit: 验证定时发送消息，下一步需要补充client的相应代码，设计区块数据结构等操作，从数据库读取data构建一个链
- 4th Commit: 增加了Tcp的部分，开启协程进行端口的监听，便于PBFT的实现，下一步实现简单的PBFT流程，然后考虑是否增加==签名验证==部分。
- 3rd Commit: 重构一下目录，增加了InitCheck文件，进行一些乱七八糟东西的检查
- 2nd Commit: 设置BoltDB的配置文件，及创建、读写BoltDB的测试代码
- 1st Commit: 设置viper的配置读取

## 测试用程序
### 测试Tcp的监听程序
```go
// 测试Tcp的Code，写在main里直接测试的。
	for i := 1; i < 100; i++ {
		message := []byte("测试" + strconv.Itoa(i))
		ConsensusUtils.TcpDial(message, "127.0.0.1:1300"+strconv.Itoa(rand.Intn(Conf.Basic.GroupNumber)))
		time.Sleep(time.Second / 5)
	}
```