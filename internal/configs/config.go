package configs

import (
	"fmt"

	"distributed-kvs/configs"
	"gopkg.in/yaml.v3"
)

type StoreType string

const (
	RaftStoreType StoreType = "raft"
	ZabStoreType  StoreType = "zab"
	GLAStoreType  StoreType = "gla"
)

type Config struct {
	HTTPServer HTTPServerConfig `yaml:"httpServer"`
	GRPCServer GRPCServerConfig `yaml:"grpcServer"`
	Store      StoreConfig      `yaml:"store"`
}

type HTTPServerConfig struct {
	Port    int  `yaml:"port"`
	Logging bool `yaml:"logging"`
}

type GRPCServerConfig struct {
	Port int `yaml:"port"`
}

type StoreConfig struct {
	JoinURL string `yaml:"joinURL"`

	Type StoreType        `yaml:"type"`
	Raft *RaftStoreConfig `yaml:"raft"`
	Zab  *RaftStoreConfig `yaml:"zab"`
}

type RaftStoreConfig struct {
	NodeID       string `yaml:"nodeID"`
	InMem        bool   `yaml:"inMem"`
	Addr         string `yaml:"addr"`
	EnableSingle bool   `yaml:"enableSingle"`
}

type ZabStoreConfig struct {
}

func ReadLocal() (Config, error) {
	var config Config

	err := yaml.Unmarshal(configs.NodeConfigFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("yaml decode response: %w", err)
	}

	return config, nil
}

func ReadClusterLocal() ([]Config, error) {
	var config []Config

	err := yaml.Unmarshal(configs.ClusterConfigFile, &config)
	if err != nil {
		return nil, fmt.Errorf("yaml decode response: %w", err)
	}

	return config, nil
}
