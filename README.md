# GameTheoryConsensus

## Todo
- 设计区块数据结构（包括交易的序列化反序列化）
- 构建区块链
- PBFT的实现，地址池设置，消息往哪里发等等
- 进行区块历史的数据库存储设计
- 消息锁
- 文件名统一大小写

## 更新日志
- 7th Commit: 增加了区块信息的写入读取，下一步需要实现PBFT，这里可能需要考虑的问题在于写到哪里去，应当给每个节点单独安排一个数据库专门存取他们的消息数据，避免并发问题。希望别崩了好吧。
- 6th Commit: 设计了PBFT格式，交易格式，区块格式，并测试了消息的发送读取，下一步需要实现PBFT
- 5th Commit: 验证定时发送消息，下一步需要补充client的相应代码，设计区块数据结构等操作，从数据库读取data构建一个链
- 4th Commit: 增加了Tcp的部分，开启协程进行端口的监听，便于PBFT的实现，下一步实现简单的PBFT流程，然后考虑是否增加==签名验证==部分。
- 3rd Commit: 重构一下目录，增加了InitCheck文件，进行一些乱七八糟东西的检查
- 2nd Commit: 设置BoltDB的配置文件，及创建、读写BoltDB的测试代码
- 1st Commit: 设置viper的配置读取

## 测试用程序
### 测试交易的序列化反序列化程序
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
### 测试区块信息的序列化程序
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
### 测试Tcp的监听程序
```go
// 测试Tcp的Code，写在main里直接测试的。
	for i := 1; i < 100; i++ {
		message := []byte("测试" + strconv.Itoa(i))
		TcpDial(message, "127.0.0.1:1300"+strconv.Itoa(rand.Intn(Conf.Basic.GroupNumber)))
		time.Sleep(time.Second / 5)
	}
```