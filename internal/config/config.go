package config

import (
	"github.com/hashicorp/hcl"
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

	err = hcl.Unmarshal(cfgData, &config)
	return &config, err
}

type Topology struct {
	Tag              string     `hcl:"tag"`
	Nodes            []NodeConf `hcl:"node"`
	Links            []string   `hcl:"link_tags"`
	ManagementBridge string     `hcl:"management_bridge"`
}

type NodeConf struct {
	Tag       string       `hcl:",key"`
	Resources ResourceConf `hcl:"resources"`
	Network   NetworkConf  `hcl:"network"`
}

type NetworkConf struct {
	Management bool     `hcl:"management"`
	VirtIO     bool     `hcl:"virtio"`
	Links      []string `hcl:"links"`
}

type ResourceConf struct {
	CPU    string     `hcl:"cpu"`
	Memory string     `hcl:"mem"`
	Disks  []DiskConf `hcl:"disk"`
	CDROM  string     `hcl:"cdrom"`
}

type DiskConf struct {
	File   string `hcl:",key"`
	Format string `hcl:"format"`
}
