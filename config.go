package main

import (
	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Basic *BasicCfg `json:"wallet"`
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
	}
	return nil
}
