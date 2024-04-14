package configs

import (
	_ "embed"
)

//go:embed node_config.yml
var NodeConfigFile []byte

//go:embed cluster_config.yml
var ClusterConfigFile []byte
