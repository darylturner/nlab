package node

import (
	"errors"
	"github.com/darylturner/nlab/internal/config"
	"github.com/darylturner/nlab/internal/node/qemu"
	"runtime"
)

func New(cfg config.NodeConf) (Node, error) {
	os := runtime.GOOS
	switch os {
	case "linux":
		node := qemu.New(cfg)
		return node, nil
	default:
		return nil, errors.New("running on unsupported os")
	}
}

type Node interface {
	Run(*config.Topology, bool) error
	Stop(*config.Topology) error
}
