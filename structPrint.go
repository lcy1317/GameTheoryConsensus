package main

import "strconv"

func (tx Transaction) printString() string {
	str := ""
	str = str + "TXid:" + strconv.Itoa(IntDeserialize(tx.TXid))
	if tx.Type == 0 {
		str = str + " Type:" + strconv.Itoa(tx.Type) + "上报"
	} else {
		str = str + " Type:" + strconv.Itoa(tx.Type) + "解密"
	}

	str = str + " Hash:" + string(tx.Hash)
	str = str + " Number:" + tx.getFloatNumString()
	str = str + " Signature:" + string(tx.Signature)
	str = str + " PubKey:" + string(tx.PubKey)
	return str
}

func (b Block) printString() string {
	str := ""
	str = str + "Version:" + b.Version
	str = str + " Timestamp:" + strconv.FormatInt(b.Timestamp, 10)
	for i := 0; i < len(b.Transactions); i++ {
		str = str + " \n  Transactions" + strconv.Itoa(i) + ":" + b.Transactions[i].printString()
	}
	str = str + " \nPrevBlockHash:" + string(b.PrevBlockHash)
	str = str + " Hash:" + string(b.Hash)
	str = str + " StageHash:" + string(b.StageHash)
	str = str + " BlockNum:" + strconv.Itoa(b.BlockNum)
	str = str + " StageNum:" + strconv.Itoa(b.StageNum)
	return str
}

func (p PBFTMessage) printString() string {
	str := ""
	str = str + "MajorNode:" + strconv.Itoa(p.MajorNode)
	str = str + " GroupNodeId:" + strconv.Itoa(p.GroupNodeId)
	str = str + " BlockInfo:" + p.BlockInfo.printString()
	str = str + " PBFTStage:" + p.PBFTStage
	return str
}
