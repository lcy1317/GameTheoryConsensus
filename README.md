# GameTheoryConsensus

## Todo
- 进行区块历史的数据库存储设计
- 文件名统一大小写
- Commit之后节点的监听。

## 更新日志
- 9th Commit: 理论上所有的交易信息应该可以发给任意的PBFT节点，然后节点发送PBFT消息得到足够的Commit之后再做出反馈，达到一定时间后打包所有交易。<font color = red>仿真中协程太多所以我有点混乱，现在的处理方式是在ClientAddr开一个监听，所有交易都会发送到这个打包节点（对应函数TcpListenWrapper），由一个全局变量`transactions`维护当前的交易池。每隔一定出块时间打包完成之后（由另一个协程函数SendingPBFTCRequest）打包所有交易发送一个PBFT消息。并同时清空当前的交易池，这也算是为了仿真的设定吧。</font>
- 8th Commit: 总算总算现在实现了不带任何签名认证的PBFT，中间使用了一个本地变量，然后上锁这种操作，真的绝了。
- 7th Commit: 增加了区块信息的写入读取，下一步需要实现PBFT，这里可能需要考虑的问题在于写到哪里去，应当给每个节点单独安排一个数据库专门存取他们的消息数据，避免并发问题。希望别崩了好吧。
- 6th Commit: 设计了PBFT格式，交易格式，区块格式，并测试了消息的发送读取，下一步需要实现PBFT
- 5th Commit: 验证定时发送消息，下一步需要补充client的相应代码，设计区块数据结构等操作，从数据库读取data构建一个链
- 4th Commit: 增加了Tcp的部分，开启协程进行端口的监听，便于PBFT的实现，下一步实现简单的PBFT流程，然后考虑是否增加==签名验证==部分。
- 3rd Commit: 重构一下目录，增加了InitCheck文件，进行一些乱七八糟东西的检查
- 2nd Commit: 设置BoltDB的配置文件，及创建、读写BoltDB的测试代码
- 1st Commit: 设置viper的配置读取

## 备注信息
### 1.映射关系
本地PBFT消息池由一个区块编号-->PBFT结构体-->节点编号-->int类型
就是一个区块对应一个PBFT，一个PBFT里面有两个消息池+两个是否已经回复的map，并且通过互斥锁避免了并发读写带来的问题。
### 2.节点的区块数据何时存储？
节点在收到足够的Commit消息后，我这里没有再加监听去收取已经完成的消息，后续再加，当前直接打印了一些信息。
```go
//TODO: Commit之后节点的监听。
```

### 3. 消息的打包
给指定端口发送消息，并每隔 n 秒会自动打包并由客户端发送PBFT消息。这里需要考虑的是发送的消息应当由不同的节点向主节点发送。这里先做一个简化操作：由一个统一的12999端口进行监听以及交易的打包和发送。暂定如此。


## 测试用程序

### 1.测试交易的序列化反序列化程序
```go
func test2() {
	var test = Transaction{

		TXid:      []byte("第一个交易"),
		Type:      0,
		Hash:      []byte("Hash"),
		Number:    50.0,
		Signature: []byte("Signature"),
		PubKey:    []byte("PubKey"),
	}
	fmt.Println(test.TXSerialize())
	cao := DeserializeTX(test.TXSerialize())
	fmt.Println(string(cao.TXid))
	fmt.Println(string(cao.Signature))
}

```
### 2.测试区块信息的序列化程序
```go
func test() {
	var testBlockMessage = Block{
		Version:   "0.0 QwQ",
		Timestamp: time.Now().Unix(),
		Transactions: []*Transaction{
			{
				TXid:      []byte("第一个交易"),
				Type:      0,
				Hash:      []byte("Hash"),
				Number:    50.0,
				Signature: []byte("Signature"),
				PubKey:    []byte("PubKey"),
			},
		},

		PrevBlockHash: []byte("PrevBlockHash"),
		Hash:          []byte("Hash"),
		StageHash:     []byte("StageHash"),
		BlockNum:      1,
	}
	fmt.Println(testBlockMessage.BlockSerialize())
	cao := BlockDeserialize(testBlockMessage.BlockSerialize())
	fmt.Println(cao.Version)
	fmt.Println(string(cao.Transactions[0].TXid))
}
```
### 3.测试Tcp的监听程序
```go
// 测试Tcp的Code，写在main里直接测试的。
	for i := 1; i < 100; i++ {
		message := []byte("测试" + strconv.Itoa(i))
		TcpDial(message, "127.0.0.1:1300"+strconv.Itoa(rand.Intn(Conf.Basic.GroupNumber)))
		time.Sleep(time.Second / 5)
	}
```
