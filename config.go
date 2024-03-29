package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Basic        *BasicCfg     `json:"wallet"`
	ChainInfo    *ChainInfoCfg `json:"chain_info"`
	TcpInfo      *TcpInfoCfg   `json:"tcp_info"`
	PrintControl *PrintControl `json:"print_control"`
}
type PrintControl struct {
	Commit           bool `json:"Commit"`
	MessageID        bool `json:"MessageID"`
	PBFTMessagePrint bool `json:"PBFTMessagePrint"`
	ReceiveTxMessage bool `json:"ReceiveTxMessage"`
}
type TcpInfoCfg struct {
	PBFTBaseAddress   string `json:"PBFTBaseAddress"`
	PBFTBasePortStart int    `json:"PBFTBasePortStart"`
	ClientAddr        string `json:"ClientAddr"`
}

type ChainInfoCfg struct {
	DBFile           string `json:"DBFile"`
	NodeDBFile       string `json:"NodeDBFile"`
	BlockSpeed       int    `json:"BlockSpeed"`
	TransactionSpeed int    `json:"TransactionSpeed"`
}

type BasicCfg struct {
	GroupNumber            int `json:"GroupNumber"`
	StageBlockNumber       int `json:"StageBlockNumber"`
	GameTheoryStop         int `json:"GameTheoryStop"`
	RevealStop             int `json:"RevealStop"`
	InitNodesNumberinGroup int `json:"InitNodesNumberinGroup"`
	NumberPrecision        int `json:"NumberPrecision"`
}

func configInitial() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	Conf = &Config{
		Basic: &BasicCfg{
			GroupNumber:            viper.GetInt("BasicCfg.GroupNumber"),
			StageBlockNumber:       viper.GetInt("BasicCfg.StageBlockNumber"),
			GameTheoryStop:         viper.GetInt("BasicCfg.GameTheoryStop"),
			RevealStop:             viper.GetInt("BasicCfg.RevealStop"),
			InitNodesNumberinGroup: viper.GetInt("BasicCfg.InitNodesNumberinGroup"),
			NumberPrecision:        viper.GetInt("BasicCfg.NumberPrecision"),
		},
		ChainInfo: &ChainInfoCfg{
			DBFile:           viper.GetString("ChainInfo.DBFile"),
			NodeDBFile:       viper.GetString("ChainInfo.NodeDBFile"),
			BlockSpeed:       viper.GetInt("ChainInfo.BlockSpeed"),
			TransactionSpeed: viper.GetInt("ChainInfo.TransactionSpeed"),
		},
		TcpInfo: &TcpInfoCfg{
			PBFTBaseAddress:   viper.GetString("TcpInfo.PBFTBaseAddress"),
			PBFTBasePortStart: viper.GetInt("TcpInfo.PBFTBasePortStart"),
			ClientAddr:        viper.GetString("TcpInfo.ClientAddr"),
		},
		PrintControl: &PrintControl{
			Commit:           viper.GetBool("PrintControl.Commit"),
			MessageID:        viper.GetBool("PrintControl.MessageID"),
			PBFTMessagePrint: viper.GetBool("PrintControl.PBFTMessagePrint"),
			ReceiveTxMessage: viper.GetBool("PrintControl.ReceiveTxMessage"),
		},
	}
	return nil
}
