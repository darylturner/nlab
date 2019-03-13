package network

import (
	"errors"
	"fmt"
	"github.com/darylturner/nlab/internal/config"
	"github.com/darylturner/nlab/internal/network/netlinux"
	"math/rand"
	"runtime"
	"time"
)

type Network interface {
	Create() error
	Destroy() error
	AddNode(string)
}

type PseudoWire struct {
	Type  string
	Port  int
	Nodes []string
}

func newNet(link string, nd config.NodeConf, allLinks map[string]Network) error {
	if net, ok := allLinks[link]; ok {
		net.AddNode(nd.Tag)
		return nil
	}

	switch runtime.GOOS {
	case "linux":
		net := netlinux.New(link)
		net.AddNode(nd.Tag)
		allLinks[link] = &net
	default:
		return errors.New("running on unsupported os")
	}

	return nil
}

func GetMap(cfg *config.Topology) (map[string]Network, error) {
	allLinks := make(map[string]Network)
	for _, nd := range cfg.Nodes {
		if nd.Network.Management == true {
			link := "_" + cfg.ManagementBridge
			if err := newNet(link, nd, allLinks); err != nil {
				return allLinks, err
			}
		}

		if !cfg.PseudoWire {
			for _, link := range nd.Network.Links {
				if err := newNet(link, nd, allLinks); err != nil {
					return allLinks, err
				}
			}
		}
	}
	return allLinks, nil
}

func GetPseudoWireMap(cfg *config.Topology) (map[string]*PseudoWire, error) {
	allLinks := make(map[string]*PseudoWire)
	portBase := 32768 + (cfg.LabID * 1024)
	os := runtime.GOOS
	for _, nd := range cfg.Nodes {
		for _, link := range nd.Network.Links {
			if pw, ok := allLinks[link]; ok {
				pw.Nodes = append(pw.Nodes, nd.Tag)
			} else {
				switch os {
				case "linux":
					allLinks[link] = &PseudoWire{
						Type:  "qemu-unicast-udp",
						Port:  portBase,
						Nodes: []string{nd.Tag},
					}
				default:
					return allLinks, errors.New("running on unsupported os")
				}
			}

			portBase += 2
		}
	}
	return allLinks, nil
}

func randomByte() int {
	return rand.Intn(256)
}

func RandomMAC() string {
	baseMAC := "52:54:00"
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(
		baseMAC+":%02x:%02x:%02x",
		randomByte(),
		randomByte(),
		randomByte(),
	)
}
