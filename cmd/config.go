package cmd

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func parseConfig(cfgFile string) (*Topology, error) {
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

	err = yaml.Unmarshal(cfgData, &config)
	return &config, err
}

type Topology struct {
	Nodes            []Node `yaml:"nodes"`
	Links            []Link `yaml:"links"`
	ManagementBridge string `yaml:"management_bridge"`
}

type Node struct {
	Tag       string       `yaml:"tag"`
	Network   NetworkConf  `yaml:"network"`
	Resources ResourceConf `yaml:"resources"`
}

type NetworkConf struct {
	Management bool     `yaml:"management"`
	Links      []string `yaml:"links"`
}

type ResourceConf struct {
	CPU    string   `yaml:"cpu"`
	Memory string   `yaml:"mem"`
	Disk   DiskConf `yaml:"disk"`
	CDROM  string   `yaml:"cdrom"`
}

type DiskConf struct {
	File   string `yaml:"file"`
	Format string `yaml:"format"`
}

type Link struct {
	Tag string `yaml:"tag"`
}
