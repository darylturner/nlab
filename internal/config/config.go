package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Parse(cfgFile string) (*Topology, error) {
	var config Topology
	var cfgData []byte
	var err error

	if cfgFile == "-" {
		cfgData, err = ioutil.ReadAll(os.Stdin)
	} else {
		cfgData, err = ioutil.ReadFile(cfgFile)
	}
	if err != nil {
		return &config, err
	}

	err = json.Unmarshal(cfgData, &config)
	return &config, err
}

type Topology struct {
	Tag              string     `json:"tag"`
	Nodes            []NodeConf `json:"nodes"`
	ManagementBridge string     `json:"management_bridge"`
}

type NodeConf struct {
	Tag       string       `json:"tag"`
	Resources ResourceConf `json:"resources"`
	Network   NetworkConf  `json:"network"`
}

type NetworkConf struct {
	Management bool     `json:"management"`
	VirtIO     bool     `json:"virtio"`
	Links      []string `json:"links"`
}

type ResourceConf struct {
	CPU    int        `json:"cpu"`
	Memory int        `json:"mem"`
	Disks  []DiskConf `json:"disks"`
	CDROM  string     `json:"cdrom"`
}

type DiskConf struct {
	File   string `json:"file"`
	Format string `json:"format"`
}
