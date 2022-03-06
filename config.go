package main

import (
	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Basic     *BasicCfg     `json:"wallet"`
	ChainInfo *ChainInfoCfg `json:"chain_info"`
	TcpInfo   *TcpInfoCfg   `json:"tcp_info"`
}
type TcpInfoCfg struct {
	PBFTBaseAddress   string `json:"PBFTBaseAddress"`
	PBFTBasePortStart int    `json:"PBFTBasePortStart"`
}

type ChainInfoCfg struct {
	DBFile     string `json:"DBFile"`
	BlockSpeed int    `json:"BlockSpeed"`
}

type BasicCfg struct {
	GroupNumber            int `json:"GroupNumber"`
	StageBlockNumber       int `json:"StageBlockNumber"`
	GameTheoryStop         int `json:"GameTheoryStop"`
	RevealStop             int `json:"RevealStop"`
	InitNodesNumberinGroup int `json:"InitNodesNumberinGroup"`
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
		},
		ChainInfo: &ChainInfoCfg{
			DBFile:     viper.GetString("ChainInfo.DBFile"),
			BlockSpeed: viper.GetInt("ChainInfo.BlockSpeed"),
		},
		TcpInfo: &TcpInfoCfg{
			PBFTBaseAddress:   viper.GetString("TcpInfo.PBFTBaseAddress"),
			PBFTBasePortStart: viper.GetInt("TcpInfo.PBFTBasePortStart"),
		},
	}
	return nil
}